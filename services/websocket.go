package services

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/philippseith/signalr"
	"github.com/tradingiq/topstepx-client/client"
)

const (
	UserHubURL = "https://rtc.topstepx.com/hubs/user?access_token=%s"
)

type WebSocketService struct {
	client        *client.Client
	conn          signalr.Client
	receiver      *WebSocketReceiver
	mu            sync.Mutex
	isConnected   bool
	accountID     int
	subscriptions map[string]bool
}

type WebSocketReceiver struct {
	handlers map[string]func(interface{})
	mu       sync.RWMutex
}

func NewWebSocketReceiver() *WebSocketReceiver {
	return &WebSocketReceiver{
		handlers: make(map[string]func(interface{})),
	}
}

func (r *WebSocketReceiver) SetHandler(event string, handler func(interface{})) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[event] = handler
}

func (r *WebSocketReceiver) RemoveHandler(event string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.handlers, event)
}

func (r *WebSocketReceiver) GatewayUserAccount(data interface{}) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if handler, ok := r.handlers["account"]; ok {
		handler(data)
	}
}

func (r *WebSocketReceiver) GatewayUserOrder(data interface{}) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if handler, ok := r.handlers["order"]; ok {
		handler(data)
	}
}

func (r *WebSocketReceiver) GatewayUserPosition(data interface{}) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if handler, ok := r.handlers["position"]; ok {
		handler(data)
	}
}

func (r *WebSocketReceiver) GatewayUserTrade(data interface{}) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if handler, ok := r.handlers["trade"]; ok {
		handler(data)
	}
}

func NewWebSocketService(c *client.Client) *WebSocketService {
	return &WebSocketService{
		client:        c,
		receiver:      NewWebSocketReceiver(),
		subscriptions: make(map[string]bool),
	}
}

func (s *WebSocketService) Connect(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isConnected {
		return nil
	}

	token := s.client.GetToken()
	if token == "" {
		return fmt.Errorf("authentication token not set")
	}

	conn, err := signalr.NewClient(ctx,
		signalr.WithHttpConnection(ctx, fmt.Sprintf(UserHubURL, token),
			signalr.WithHTTPHeaders(func() http.Header {
				headers := http.Header{}
				headers.Set("Authorization", "Bearer "+token)
				return headers
			}),
		),
		signalr.WithReceiver(s.receiver),
	)

	if err != nil {
		return fmt.Errorf("failed to create SignalR client: %w", err)
	}

	conn.Start()

	s.conn = conn
	s.isConnected = true

	return nil
}

func (s *WebSocketService) Disconnect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isConnected || s.conn == nil {
		return nil
	}

	s.conn.Stop()
	s.isConnected = false
	s.subscriptions = make(map[string]bool)

	return nil
}

func (s *WebSocketService) IsConnected() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isConnected
}

func (s *WebSocketService) SetAccountID(accountID int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.accountID = accountID
}

func (s *WebSocketService) SubscribeAccounts() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("SubscribeAccounts")
	if result.Error != nil {
		return fmt.Errorf("failed to subscribe to accounts: %w", result.Error)
	}

	s.subscriptions["accounts"] = true
	return nil
}

func (s *WebSocketService) UnsubscribeAccounts() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("UnsubscribeAccounts")
	if result.Error != nil {
		return fmt.Errorf("failed to unsubscribe from accounts: %w", result.Error)
	}

	delete(s.subscriptions, "accounts")
	return nil
}

func (s *WebSocketService) SubscribeOrders(accountID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isConnected {
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

func (s *WebSocketService) UnsubscribeOrders(accountID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("UnsubscribeOrders", accountID)
	if result.Error != nil {
		return fmt.Errorf("failed to unsubscribe from orders: %w", result.Error)
	}

	delete(s.subscriptions, "orders")
	return nil
}

func (s *WebSocketService) SubscribePositions(accountID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isConnected {
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

func (s *WebSocketService) UnsubscribePositions(accountID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("UnsubscribePositions", accountID)
	if result.Error != nil {
		return fmt.Errorf("failed to unsubscribe from positions: %w", result.Error)
	}

	delete(s.subscriptions, "positions")
	return nil
}

func (s *WebSocketService) SubscribeTrades(accountID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isConnected {
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

func (s *WebSocketService) UnsubscribeTrades(accountID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("UnsubscribeTrades", accountID)
	if result.Error != nil {
		return fmt.Errorf("failed to unsubscribe from trades: %w", result.Error)
	}

	delete(s.subscriptions, "trades")
	return nil
}

func (s *WebSocketService) SubscribeAll(accountID int) error {
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

func (s *WebSocketService) UnsubscribeAll() error {
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

func (s *WebSocketService) SetAccountHandler(handler func(interface{})) {
	s.receiver.SetHandler("account", handler)
}

func (s *WebSocketService) SetOrderHandler(handler func(interface{})) {
	s.receiver.SetHandler("order", handler)
}

func (s *WebSocketService) SetPositionHandler(handler func(interface{})) {
	s.receiver.SetHandler("position", handler)
}

func (s *WebSocketService) SetTradeHandler(handler func(interface{})) {
	s.receiver.SetHandler("trade", handler)
}

func (s *WebSocketService) resubscribe() {
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
