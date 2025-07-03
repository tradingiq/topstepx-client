package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/philippseith/signalr"
	"github.com/tradingiq/topstepx-client/client"
	"github.com/tradingiq/topstepx-client/models"
)

const (
	UserHubURL = "https://rtc.topstepx.com/hubs/user?access_token=%s"
)

type UserDataWebSocketService struct {
	client            *client.Client
	conn              signalr.Client
	receiver          *UserDataReceiver
	mu                sync.Mutex
	state             ConnectionState
	accountID         int
	subscriptions     map[string]bool
	ctx               context.Context
	cancel            context.CancelFunc
	reconnectChan     chan struct{}
	connectionHandler func(ConnectionState)
	maxReconnectDelay time.Duration
	reconnectAttempts int
}

type UserDataReceiver struct {
	handlers        map[string]func(interface{})
	accountHandler  func(*models.AccountUpdateData)
	orderHandler    func(*models.OrderUpdateData)
	positionHandler func(*models.PositionUpdateData)
	tradeHandler    func(*models.TradeUpdateData)
	mu              sync.RWMutex
	service         *UserDataWebSocketService
}

func NewUserDataReceiver(service *UserDataWebSocketService) *UserDataReceiver {
	return &UserDataReceiver{
		handlers: make(map[string]func(interface{})),
		service:  service,
	}
}

func (r *UserDataReceiver) SetHandler(event string, handler func(interface{})) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[event] = handler
}

func (r *UserDataReceiver) RemoveHandler(event string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.handlers, event)
}

func (r *UserDataReceiver) GatewayUserAccount(data interface{}) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.accountHandler != nil {
		var accountData models.AccountUpdateData
		if jsonBytes, err := json.Marshal(data); err == nil {
			if err := json.Unmarshal(jsonBytes, &accountData); err == nil {
				r.accountHandler(&accountData)
				return
			}
		}
	}

	if handler, ok := r.handlers["account"]; ok {
		handler(data)
	}
}

func (r *UserDataReceiver) GatewayUserOrder(data interface{}) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.orderHandler != nil {

		var orderData models.OrderUpdateData
		if jsonBytes, err := json.Marshal(data); err == nil {
			if err := json.Unmarshal(jsonBytes, &orderData); err == nil {
				r.orderHandler(&orderData)
				return
			}
		}
	}

	if handler, ok := r.handlers["order"]; ok {
		handler(data)
	}
}

func (r *UserDataReceiver) GatewayUserPosition(data interface{}) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.positionHandler != nil {
		var positionData models.PositionUpdateData
		if jsonBytes, err := json.Marshal(data); err == nil {
			if err := json.Unmarshal(jsonBytes, &positionData); err == nil {
				r.positionHandler(&positionData)
				return
			}
		}
	}

	if handler, ok := r.handlers["position"]; ok {
		handler(data)
	}
}

func (r *UserDataReceiver) GatewayUserTrade(data interface{}) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.tradeHandler != nil {
		var tradeData models.TradeUpdateData
		if jsonBytes, err := json.Marshal(data); err == nil {
			if err := json.Unmarshal(jsonBytes, &tradeData); err == nil {
				r.tradeHandler(&tradeData)
				return
			}
		}
	}

	if handler, ok := r.handlers["trade"]; ok {
		handler(data)
	}
}

func (r *UserDataReceiver) ConnectionClosed() {
	if r.service != nil {
		r.service.mu.Lock()
		if r.service.state == StateConnected {
			r.service.setState(StateReconnecting)
			r.service.mu.Unlock()

			select {
			case r.service.reconnectChan <- struct{}{}:
			default:
			}
		} else {
			r.service.mu.Unlock()
		}
	}
}

func NewUserDataWebSocketService(c *client.Client) *UserDataWebSocketService {
	s := &UserDataWebSocketService{
		client:            c,
		subscriptions:     make(map[string]bool),
		state:             StateDisconnected,
		maxReconnectDelay: 30 * time.Second,
		reconnectChan:     make(chan struct{}, 1),
	}
	s.receiver = NewUserDataReceiver(s)
	return s
}

func (s *UserDataWebSocketService) Connect(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state == StateConnected || s.state == StateConnecting {
		return nil
	}

	s.ctx, s.cancel = context.WithCancel(ctx)
	s.setState(StateConnecting)

	token := s.client.GetToken()
	if token == "" {
		s.setState(StateDisconnected)
		return fmt.Errorf("authentication token not set")
	}

	conn, err := signalr.NewClient(s.ctx,
		signalr.WithHttpConnection(s.ctx, fmt.Sprintf(UserHubURL, token),
			signalr.WithHTTPHeaders(func() http.Header {
				headers := http.Header{}
				headers.Set("Authorization", "Bearer "+token)
				return headers
			}),
		),
		signalr.WithReceiver(s.receiver),
		signalr.Logger(newNoopLogger(), false),
	)

	if err != nil {
		s.setState(StateDisconnected)
		return fmt.Errorf("failed to create SignalR client: %w", err)
	}

	conn.Start()

	s.conn = conn
	s.setState(StateConnected)
	s.reconnectAttempts = 0

	go s.handleReconnection()

	return nil
}

func (s *UserDataWebSocketService) Disconnect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state == StateDisconnected {
		return nil
	}

	if s.cancel != nil {
		s.cancel()
	}

	if s.conn != nil {
		s.conn.Stop()
	}

	s.setState(StateDisconnected)
	s.subscriptions = make(map[string]bool)
	s.accountID = 0

	return nil
}

func (s *UserDataWebSocketService) IsConnected() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.state == StateConnected
}

func (s *UserDataWebSocketService) GetConnectionState() ConnectionState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.state
}

func (s *UserDataWebSocketService) SetConnectionHandler(handler func(ConnectionState)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connectionHandler = handler
}

func (s *UserDataWebSocketService) setState(state ConnectionState) {
	s.state = state
	if s.connectionHandler != nil {

		go s.connectionHandler(state)
	}
}

func (s *UserDataWebSocketService) SetAccountID(accountID int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.accountID = accountID
}

func (s *UserDataWebSocketService) SubscribeAccounts() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("SubscribeAccounts")
	if result.Error != nil {
		return fmt.Errorf("failed to subscribe to accounts: %w", result.Error)
	}

	s.subscriptions["accounts"] = true
	return nil
}

func (s *UserDataWebSocketService) UnsubscribeAccounts() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("UnsubscribeAccounts")
	if result.Error != nil {
		return fmt.Errorf("failed to unsubscribe from accounts: %w", result.Error)
	}

	delete(s.subscriptions, "accounts")
	return nil
}

func (s *UserDataWebSocketService) SubscribeOrders(accountID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("SubscribeOrders", accountID)
	if result.Error != nil {
		return fmt.Errorf("failed to subscribe to orders: %w", result.Error)
	}

	s.subscriptions["orders"] = true
	s.accountID = accountID
	return nil
}

func (s *UserDataWebSocketService) UnsubscribeOrders(accountID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("UnsubscribeOrders", accountID)
	if result.Error != nil {
		return fmt.Errorf("failed to unsubscribe from orders: %w", result.Error)
	}

	delete(s.subscriptions, "orders")
	return nil
}

func (s *UserDataWebSocketService) SubscribePositions(accountID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("SubscribePositions", accountID)
	if result.Error != nil {
		return fmt.Errorf("failed to subscribe to positions: %w", result.Error)
	}

	s.subscriptions["positions"] = true
	s.accountID = accountID
	return nil
}

func (s *UserDataWebSocketService) UnsubscribePositions(accountID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("UnsubscribePositions", accountID)
	if result.Error != nil {
		return fmt.Errorf("failed to unsubscribe from positions: %w", result.Error)
	}

	delete(s.subscriptions, "positions")
	return nil
}

func (s *UserDataWebSocketService) SubscribeTrades(accountID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("SubscribeTrades", accountID)
	if result.Error != nil {
		return fmt.Errorf("failed to subscribe to trades: %w", result.Error)
	}

	s.subscriptions["trades"] = true
	s.accountID = accountID
	return nil
}

func (s *UserDataWebSocketService) UnsubscribeTrades(accountID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("UnsubscribeTrades", accountID)
	if result.Error != nil {
		return fmt.Errorf("failed to unsubscribe from trades: %w", result.Error)
	}

	delete(s.subscriptions, "trades")
	return nil
}

func (s *UserDataWebSocketService) SubscribeAll(accountID int) error {
	if err := s.SubscribeAccounts(); err != nil {
		return err
	}
	if err := s.SubscribeOrders(accountID); err != nil {
		return err
	}
	if err := s.SubscribePositions(accountID); err != nil {
		return err
	}
	if err := s.SubscribeTrades(accountID); err != nil {
		return err
	}
	return nil
}

func (s *UserDataWebSocketService) UnsubscribeAll() error {
	s.mu.Lock()
	accountID := s.accountID
	s.mu.Unlock()

	if err := s.UnsubscribeAccounts(); err != nil {
		return err
	}
	if accountID > 0 {
		if err := s.UnsubscribeOrders(accountID); err != nil {
			return err
		}
		if err := s.UnsubscribePositions(accountID); err != nil {
			return err
		}
		if err := s.UnsubscribeTrades(accountID); err != nil {
			return err
		}
	}
	return nil
}

func (s *UserDataWebSocketService) SetAccountHandler(handler func(*models.AccountUpdateData)) {
	s.receiver.mu.Lock()
	defer s.receiver.mu.Unlock()
	s.receiver.accountHandler = handler
}

func (s *UserDataWebSocketService) SetOrderHandler(handler func(*models.OrderUpdateData)) {
	s.receiver.mu.Lock()
	defer s.receiver.mu.Unlock()
	s.receiver.orderHandler = handler
}

func (s *UserDataWebSocketService) SetPositionHandler(handler func(*models.PositionUpdateData)) {
	s.receiver.mu.Lock()
	defer s.receiver.mu.Unlock()
	s.receiver.positionHandler = handler
}

func (s *UserDataWebSocketService) SetTradeHandler(handler func(*models.TradeUpdateData)) {
	s.receiver.mu.Lock()
	defer s.receiver.mu.Unlock()
	s.receiver.tradeHandler = handler
}

func (s *UserDataWebSocketService) handleReconnection() {

	for {
		select {
		case <-s.ctx.Done():

			return
		case <-time.After(5 * time.Second):

			s.checkConnectionHealth()
		case <-s.reconnectChan:

			s.reconnect()
		}
	}
}

func (s *UserDataWebSocketService) checkConnectionHealth() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateConnected || s.conn == nil {
		return
	}

	pingTimeout := time.NewTimer(10 * time.Second)
	defer pingTimeout.Stop()

	select {
	case result := <-s.conn.Invoke("ping"):
		if result.Error != nil {
			s.setState(StateReconnecting)
			select {
			case s.reconnectChan <- struct{}{}:
			default:
			}
		}
	case <-pingTimeout.C:

		s.setState(StateReconnecting)
		select {
		case s.reconnectChan <- struct{}{}:
		default:
		}
	}
}

func (s *UserDataWebSocketService) reconnect() {
	s.mu.Lock()
	if s.state == StateDisconnected {
		s.mu.Unlock()
		return
	}
	s.setState(StateReconnecting)
	s.mu.Unlock()

	baseDelay := time.Second
	maxDelay := s.maxReconnectDelay

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		s.reconnectAttempts++
		delay := baseDelay * time.Duration(1<<uint(s.reconnectAttempts-1))
		if delay > maxDelay {
			delay = maxDelay
		}

		select {
		case <-s.ctx.Done():
			return
		case <-time.After(delay):
		}

		s.mu.Lock()
		if s.conn != nil {
			s.conn.Stop()
		}
		s.mu.Unlock()

		token := s.client.GetToken()
		if token == "" {
			continue
		}

		conn, err := signalr.NewClient(s.ctx,
			signalr.WithHttpConnection(s.ctx, fmt.Sprintf(UserHubURL, token),
				signalr.WithHTTPHeaders(func() http.Header {
					headers := http.Header{}
					headers.Set("Authorization", "Bearer "+token)
					return headers
				}),
			),
			signalr.WithReceiver(s.receiver),
			signalr.Logger(newNoopLogger(), false),
		)

		if err != nil {
			continue
		}

		conn.Start()

		testTimeout := time.NewTimer(5 * time.Second)
		defer testTimeout.Stop()

		select {
		case result := <-conn.Invoke("ping"):
			if result.Error != nil {
				conn.Stop()
				continue
			}
		case <-testTimeout.C:

			conn.Stop()
			continue
		case <-s.ctx.Done():
			conn.Stop()
			return
		}

		s.mu.Lock()
		s.conn = conn
		s.setState(StateConnected)
		s.reconnectAttempts = 0
		s.mu.Unlock()

		s.resubscribe()

		return
	}
}

func (s *UserDataWebSocketService) resubscribe() {
	s.mu.Lock()
	subs := make(map[string]bool)
	for k, v := range s.subscriptions {
		subs[k] = v
	}
	accountID := s.accountID
	s.mu.Unlock()

	if subs["accounts"] {
		<-s.conn.Send("SubscribeAccounts")
	}
	if subs["orders"] && accountID > 0 {
		<-s.conn.Send("SubscribeOrders", accountID)
	}
	if subs["positions"] && accountID > 0 {
		<-s.conn.Send("SubscribePositions", accountID)
	}
	if subs["trades"] && accountID > 0 {
		<-s.conn.Send("SubscribeTrades", accountID)
	}
}
