package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/philippseith/signalr"
	"github.com/tradingiq/projectx-client/client"
	"github.com/tradingiq/projectx-client/models"
)

const (
	MarketHubURL = "https://rtc.projectx.com/hubs/market?access_token=%s"
)

type ConnectionState int

const (
	StateDisconnected ConnectionState = iota
	StateConnecting
	StateConnected
	StateReconnecting
)

type MarketDataWebSocketService struct {
	client            *client.Client
	conn              signalr.Client
	receiver          *MarketDataReceiver
	mu                sync.Mutex
	state             ConnectionState
	subscriptions     map[string]map[string]bool
	ctx               context.Context
	cancel            context.CancelFunc
	reconnectChan     chan struct{}
	connectionHandler func(ConnectionState)
	maxReconnectDelay time.Duration
	reconnectAttempts int
}

type MarketDataReceiver struct {
	quoteHandler func(string, models.Quote)
	tradeHandler func(string, models.TradeData)
	depthHandler func(string, models.MarketDepthData)
	mu           sync.RWMutex
	service      *MarketDataWebSocketService
}

func NewMarketDataReceiver(service *MarketDataWebSocketService) *MarketDataReceiver {
	return &MarketDataReceiver{
		service: service,
	}
}

func (r *MarketDataReceiver) SetQuoteHandler(handler func(string, models.Quote)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.quoteHandler = handler
}

func (r *MarketDataReceiver) SetTradeHandler(handler func(string, models.TradeData)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tradeHandler = handler
}

func (r *MarketDataReceiver) SetDepthHandler(handler func(string, models.MarketDepthData)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.depthHandler = handler
}

func (r *MarketDataReceiver) ConnectionClosed() {
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

func (r *MarketDataReceiver) GatewayQuote(contractID string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	var quote models.Quote
	if err := json.Unmarshal(jsonData, &quote); err != nil {
		return
	}

	r.mu.RLock()
	handler := r.quoteHandler
	r.mu.RUnlock()

	if handler != nil {
		handler(contractID, quote)
	}
}

func (r *MarketDataReceiver) GatewayTrade(contractID string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	var trades models.TradeData
	if err := json.Unmarshal(jsonData, &trades); err != nil {
		return
	}

	r.mu.RLock()
	handler := r.tradeHandler
	r.mu.RUnlock()

	if handler != nil {
		handler(contractID, trades)
	}
}

func (r *MarketDataReceiver) GatewayDepth(contractID string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	var depth models.MarketDepthData
	if err := json.Unmarshal(jsonData, &depth); err != nil {
		return
	}

	r.mu.RLock()
	handler := r.depthHandler
	r.mu.RUnlock()

	if handler != nil {
		handler(contractID, depth)
	}
}

func NewMarketDataWebSocketService(c *client.Client) *MarketDataWebSocketService {
	s := &MarketDataWebSocketService{
		client:            c,
		subscriptions:     make(map[string]map[string]bool),
		state:             StateDisconnected,
		maxReconnectDelay: 30 * time.Second,
		reconnectChan:     make(chan struct{}, 1),
	}
	s.receiver = NewMarketDataReceiver(s)
	return s
}

func (s *MarketDataWebSocketService) Connect(ctx context.Context) error {
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
		signalr.WithHttpConnection(s.ctx, fmt.Sprintf(MarketHubURL, token),
			signalr.WithHTTPHeaders(func() http.Header {
				headers := http.Header{}
				headers.Set("Authorization", "Bearer "+token)
				return headers
			}),
		),
		signalr.WithReceiver(s.receiver),
		signalr.MaximumReceiveMessageSize(1024*1024),
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

func (s *MarketDataWebSocketService) Disconnect() error {
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
	s.subscriptions = make(map[string]map[string]bool)

	return nil
}

func (s *MarketDataWebSocketService) IsConnected() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.state == StateConnected
}

func (s *MarketDataWebSocketService) GetConnectionState() ConnectionState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.state
}

func (s *MarketDataWebSocketService) SetConnectionHandler(handler func(ConnectionState)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connectionHandler = handler
}

func (s *MarketDataWebSocketService) setState(state ConnectionState) {
	s.state = state
	if s.connectionHandler != nil {

		go s.connectionHandler(state)
	}
}

func (s *MarketDataWebSocketService) SubscribeContractQuotes(contractID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("SubscribeContractQuotes", contractID)
	if result.Error != nil {
		return fmt.Errorf("failed to subscribe to quotes for %s: %w", contractID, result.Error)
	}

	if s.subscriptions[contractID] == nil {
		s.subscriptions[contractID] = make(map[string]bool)
	}
	s.subscriptions[contractID]["quotes"] = true
	return nil
}

func (s *MarketDataWebSocketService) UnsubscribeContractQuotes(contractID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("UnsubscribeContractQuotes", contractID)
	if result.Error != nil {
		return fmt.Errorf("failed to unsubscribe from quotes for %s: %w", contractID, result.Error)
	}

	if s.subscriptions[contractID] != nil {
		delete(s.subscriptions[contractID], "quotes")
		if len(s.subscriptions[contractID]) == 0 {
			delete(s.subscriptions, contractID)
		}
	}
	return nil
}

func (s *MarketDataWebSocketService) SubscribeContractTrades(contractID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("SubscribeContractTrades", contractID)
	if result.Error != nil {
		return fmt.Errorf("failed to subscribe to trades for %s: %w", contractID, result.Error)
	}

	if s.subscriptions[contractID] == nil {
		s.subscriptions[contractID] = make(map[string]bool)
	}
	s.subscriptions[contractID]["trades"] = true
	return nil
}

func (s *MarketDataWebSocketService) UnsubscribeContractTrades(contractID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("UnsubscribeContractTrades", contractID)
	if result.Error != nil {
		return fmt.Errorf("failed to unsubscribe from trades for %s: %w", contractID, result.Error)
	}

	if s.subscriptions[contractID] != nil {
		delete(s.subscriptions[contractID], "trades")
		if len(s.subscriptions[contractID]) == 0 {
			delete(s.subscriptions, contractID)
		}
	}
	return nil
}

func (s *MarketDataWebSocketService) SubscribeContractMarketDepth(contractID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("SubscribeContractMarketDepth", contractID)
	if result.Error != nil {
		return fmt.Errorf("failed to subscribe to market depth for %s: %w", contractID, result.Error)
	}

	if s.subscriptions[contractID] == nil {
		s.subscriptions[contractID] = make(map[string]bool)
	}
	s.subscriptions[contractID]["depth"] = true
	return nil
}

func (s *MarketDataWebSocketService) UnsubscribeContractMarketDepth(contractID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateConnected {
		return fmt.Errorf("websocket not connected")
	}

	result := <-s.conn.Invoke("UnsubscribeContractMarketDepth", contractID)
	if result.Error != nil {
		return fmt.Errorf("failed to unsubscribe from market depth for %s: %w", contractID, result.Error)
	}

	if s.subscriptions[contractID] != nil {
		delete(s.subscriptions[contractID], "depth")
		if len(s.subscriptions[contractID]) == 0 {
			delete(s.subscriptions, contractID)
		}
	}
	return nil
}

func (s *MarketDataWebSocketService) SubscribeAll(contractID string) error {
	if err := s.SubscribeContractQuotes(contractID); err != nil {
		return err
	}
	if err := s.SubscribeContractTrades(contractID); err != nil {
		return err
	}
	if err := s.SubscribeContractMarketDepth(contractID); err != nil {
		return err
	}
	return nil
}

func (s *MarketDataWebSocketService) UnsubscribeAll(contractID string) error {
	if err := s.UnsubscribeContractQuotes(contractID); err != nil {
		return err
	}
	if err := s.UnsubscribeContractTrades(contractID); err != nil {
		return err
	}
	if err := s.UnsubscribeContractMarketDepth(contractID); err != nil {
		return err
	}
	return nil
}

func (s *MarketDataWebSocketService) UnsubscribeAllContracts() error {
	s.mu.Lock()
	contracts := make([]string, 0, len(s.subscriptions))
	for contractID := range s.subscriptions {
		contracts = append(contracts, contractID)
	}
	s.mu.Unlock()

	for _, contractID := range contracts {
		if err := s.UnsubscribeAll(contractID); err != nil {
			return err
		}
	}
	return nil
}

func (s *MarketDataWebSocketService) SetQuoteHandler(handler func(string, models.Quote)) {
	s.receiver.SetQuoteHandler(handler)
}

func (s *MarketDataWebSocketService) SetTradeHandler(handler func(string, models.TradeData)) {
	s.receiver.SetTradeHandler(handler)
}

func (s *MarketDataWebSocketService) SetDepthHandler(handler func(string, models.MarketDepthData)) {
	s.receiver.SetDepthHandler(handler)
}

func (s *MarketDataWebSocketService) handleReconnection() {

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

func (s *MarketDataWebSocketService) checkConnectionHealth() {
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

func (s *MarketDataWebSocketService) reconnect() {
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
			signalr.WithHttpConnection(s.ctx, fmt.Sprintf(MarketHubURL, token),
				signalr.WithHTTPHeaders(func() http.Header {
					headers := http.Header{}
					headers.Set("Authorization", "Bearer "+token)
					return headers
				}),
			),
			signalr.WithReceiver(s.receiver),
			signalr.MaximumReceiveMessageSize(1024*1024),
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
			// Connection test failed
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

		// Resubscribe to all previous subscriptions
		s.resubscribe()

		// Successfully reconnected
		return
	}
}

func (s *MarketDataWebSocketService) resubscribe() {
	s.mu.Lock()
	// Create a copy of subscriptions to avoid holding lock during invocations
	subs := make(map[string]map[string]bool)
	for contractID, dataTypes := range s.subscriptions {
		subs[contractID] = make(map[string]bool)
		for dataType, subscribed := range dataTypes {
			subs[contractID][dataType] = subscribed
		}
	}
	s.mu.Unlock()

	// Resubscribe to all previously subscribed data
	for contractID, dataTypes := range subs {
		if dataTypes["quotes"] {
			<-s.conn.Send("SubscribeContractQuotes", contractID)
		}
		if dataTypes["trades"] {
			<-s.conn.Send("SubscribeContractTrades", contractID)
		}
		if dataTypes["depth"] {
			<-s.conn.Send("SubscribeContractMarketDepth", contractID)
		}
	}
}

func (s *MarketDataWebSocketService) GetSubscriptions() map[string]map[string]bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Return a copy to prevent external modifications
	result := make(map[string]map[string]bool)
	for contractID, dataTypes := range s.subscriptions {
		result[contractID] = make(map[string]bool)
		for dataType, subscribed := range dataTypes {
			result[contractID][dataType] = subscribed
		}
	}
	return result
}
