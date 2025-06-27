# TopStepX Go Client Library

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**This library is currently under active development and may undergo breaking changes.**

A comprehensive Go client library for the TopStepX trading platform API, providing full access to trading operations,
real-time data streaming, and account management functionality.

## Features

### Core Trading Operations

- **Order Management**: Place, modify, and cancel orders with support for all order types
- **Position Management**: Monitor and manage trading positions
- **Account Management**: Multi-account support with balance tracking
- **Trade History**: Access historical trade data and execution details

### Real-time Data Streaming

- **WebSocket Integration**: Real-time updates via SignalR
- **Event Handlers**: Custom handlers for account, order, position, and trade updates
- **Automatic Reconnection**: Built-in connection recovery and resubscription

### Market Data & Analysis

- **Historical Data**: Multi-timeframe price bars and candlestick data
- **Contract Information**: Search and retrieve trading instrument details
- **Live Market Data**: Real-time quotes, trades, and market depth via WebSocket
- **Market Data Streaming**: Real-time price updates, volume, and order book data

### Advanced Features

- **Authentication Management**: Automatic token handling and refresh
- **Error Handling**: Comprehensive error codes and structured responses
- **Context Support**: Full context.Context support for timeouts and cancellation
- **Concurrent Operations**: Thread-safe operations with proper synchronization
- **Connection Management**: Automatic reconnection with exponential backoff and health monitoring
- **Ping Timeout Protection**: Built-in connection validation to prevent indefinite blocking

## Installation

```bash
go get github.com/tradingiq/topstepx-client
```

### Prerequisites

- Go 1.24 or higher
- Active TopStepX trading account
- API credentials (username and API key)

## Quick Start

### Basic Setup

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tradingiq/topstepx-client"
	"github.com/tradingiq/topstepx-client/models"
)

func main() {
	// Create new client
	client := topstepx.NewClient()
	ctx := context.Background()

	// Login with API credentials
	err := client.LoginAndConnect(ctx, "your-username", "your-api-key")
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	// Get active trading accounts
	account, err := client.GetFirstActiveAccount(ctx)
	if err != nil {
		log.Fatalf("Failed to get account: %v", err)
	}

	fmt.Printf("Account: %s, Balance: $%.2f\n", account.Name, account.Balance)

	// Clean shutdown
	defer client.Disconnect(ctx)
}
```

### Environment Configuration

Create a `.env` file in your project root:

```env
TOPSTEPX_USERNAME=your-username
TOPSTEPX_API_KEY=your-api-key
```

## API Reference

### High-Level Convenience Methods

The client provides several convenience methods that simplify common workflows:

```go
// Login and connect to both API and WebSocket services in one call
err := client.LoginAndConnect(ctx, "username", "api-key")
if err != nil {
    log.Fatalf("Login and connection failed: %v", err)
}

// Connect WebSocket services with account context
accountID := 12345
err := client.ConnectWebSocketWithAccount(ctx, accountID)
if err != nil {
    log.Fatalf("WebSocket connection failed: %v", err)
}

// Get the first active trading account (most common use case)
account, err := client.GetFirstActiveAccount(ctx)
if err != nil {
    log.Fatalf("Failed to get active account: %v", err)
}

// Get all active accounts
accounts, err := client.GetActiveAccounts(ctx)
if err != nil {
    log.Fatalf("Failed to get accounts: %v", err)
}

// Connect to market data WebSocket
err = client.ConnectMarketData(ctx)
if err != nil {
    log.Fatalf("Market data connection failed: %v", err)
}

// Graceful shutdown of all connections
defer client.Disconnect(ctx)
```

### Client Initialization

```go
// Basic client
client := topstepx.NewClient()

// Client with custom HTTP options
client := topstepx.NewClient(
client.WithBaseURL("https://custom-api-url.com"),
client.WithTimeout(30 * time.Second),
)
```

### Authentication

```go
// Login with API key
resp, err := client.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
UserName: "username",
APIKey:   "api-key",
})

// Login with app credentials
resp, err := client.Auth.LoginApp(ctx, &models.LoginAppRequest{
UserName: "username",
Password: "password",
})

// Validate current session
resp, err := client.Auth.Validate(ctx)

// Logout
resp, err := client.Auth.Logout(ctx)
```

### Account Management

```go
// Search all accounts
resp, err := client.Account.SearchAccounts(ctx, &models.SearchAccountRequest{
OnlyActiveAccounts: true,
})

// Get active accounts (convenience method)
accounts, err := client.GetActiveAccounts(ctx)

// Get first active account (convenience method)
account, err := client.GetFirstActiveAccount(ctx)
```

### Order Management

```go
// Place a market order
resp, err := client.Order.PlaceOrder(ctx, &models.PlaceOrderRequest{
AccountID:  accountID,
ContractID: contractID,
Type:       models.OrderTypeMarket,
Side:       models.OrderSideBid,
Size:       1,
})

// Place a limit order
resp, err := client.Order.PlaceOrder(ctx, &models.PlaceOrderRequest{
AccountID:  accountID,
ContractID: contractID,
Type:       models.OrderTypeLimit,
Side:       models.OrderSideBid,
Size:       1,
LimitPrice: &limitPrice,
})

// Search open orders
resp, err := client.Order.SearchOpenOrders(ctx, &models.SearchOpenOrderRequest{
AccountID: accountID,
})

// Modify existing order
resp, err := client.Order.ModifyOrder(ctx, &models.ModifyOrderRequest{
AccountID:  accountID,
OrderID:    orderID,
LimitPrice: &newPrice,
Size:       &newSize,
})

// Cancel order
resp, err := client.Order.CancelOrder(ctx, &models.CancelOrderRequest{
AccountID: accountID,
OrderID:   orderID,
})
```

### Position Management

```go
// Search current positions
resp, err := client.Position.SearchPositions(ctx, &models.SearchPositionRequest{
AccountID: accountID,
})

// Close entire position
resp, err := client.Position.CloseContractPosition(ctx, &models.CloseContractPositionRequest{
AccountID:  accountID,
ContractID: contractID,
})

// Partial position close
resp, err := client.Position.PartialCloseContractPosition(ctx, &models.PartialCloseContractPositionRequest{
AccountID:  accountID,
ContractID: contractID,
Size:       partialSize,
})
```

### Historical Data

```go
// Get price bars
resp, err := client.History.RetrieveBars(ctx, &models.RetrieveBarRequest{
ContractID:        contractID,
Live:              false,
Unit:              models.AggregateBarUnitMinute,
UnitNumber:        5, // 5-minute bars
StartTime:         startTime,
EndTime:           endTime,
Limit:             1000,
IncludePartialBar: false,
})
```

### Contract Information

```go
// Search contracts
searchText := "ES" // E-mini S&P 500
resp, err := client.Contract.SearchContracts(ctx, &models.SearchContractRequest{
SearchText: &searchText,
Live:       true,
})

// Get contract by ID
resp, err := client.Contract.SearchContractByID(ctx, &models.SearchContractByIdRequest{
ContractID: contractID,
})
```

### User Data Real-time Updates

#### Type-Safe Event Handlers

```go
// Type-safe order handler with structured data
client.UserData.SetOrderHandler(func(data *models.OrderUpdateData) {
    fmt.Printf("Order %d: %s - %s %d at $%.2f\n", 
        data.Data.ID, data.Data.Status, data.Data.Side, 
        data.Data.Size, data.Data.LimitPrice)
})

// Type-safe position handler with structured data
client.UserData.SetPositionHandler(func(data *models.PositionUpdateData) {
    fmt.Printf("Position %d: %d shares at avg $%.2f\n", 
        data.Data.ID, data.Data.Size, data.Data.AveragePrice)
})

// Generic handlers for account and trade data
client.UserData.SetAccountHandler(func(data interface{}) {
    fmt.Printf("Account update: %+v\n", data)
})

client.UserData.SetTradeHandler(func(data interface{}) {
    fmt.Printf("Trade update: %+v\n", data)
})

// Connection state handler
client.UserData.SetConnectionHandler(func(state services.ConnectionState) {
    fmt.Printf("Connection state: %v\n", state)
})
```

#### Connection and Subscription Management

```go
// Connect and subscribe
err := client.UserData.Connect(ctx)
if err != nil {
    log.Fatalf("UserData connection failed: %v", err)
}

// Subscribe to all events for an account
err = client.UserData.SubscribeAll(accountID)
if err != nil {
    log.Fatalf("Subscription failed: %v", err)
}

// Or subscribe to specific events
err = client.UserData.SubscribeAccounts()
err = client.UserData.SubscribeOrders(accountID)
err = client.UserData.SubscribePositions(accountID)
err = client.UserData.SubscribeTrades(accountID)
```

### Market Data Real-time Streaming

#### Setting Up Market Data Handlers

```go
// Real-time quote handler (bid/ask prices, volume, etc.)
client.MarketData.SetQuoteHandler(func(contractID string, quote models.Quote) {
    fmt.Printf("Quote for %s: Bid $%.2f, Ask $%.2f, Last $%.2f, Volume %.0f\n",
        contractID, quote.BestBid, quote.BestAsk, quote.LastPrice, quote.Volume)
})

// Real-time trade handler
client.MarketData.SetTradeHandler(func(contractID string, trades models.TradeData) {
    for _, trade := range trades {
        fmt.Printf("Trade for %s: $%.2f, Volume %.0f at %v\n",
            contractID, trade.Price, trade.Volume, trade.Timestamp)
    }
})

// Market depth/order book handler
client.MarketData.SetDepthHandler(func(contractID string, depth models.MarketDepthData) {
    fmt.Printf("Market depth for %s: %d levels\n", contractID, len(depth))
    for _, level := range depth {
        side := "Bid"
        if level.Type == 3 { // Ask side
            side = "Ask"
        }
        fmt.Printf("  %s: $%.2f x %.0f\n", side, level.Price, level.Volume)
    }
})

// Connection state monitoring
client.MarketData.SetConnectionHandler(func(state services.ConnectionState) {
    fmt.Printf("MarketData connection state: %v\n", state)
})
```

#### Market Data Connection and Subscriptions

```go
// Connect to market data WebSocket
err := client.MarketData.Connect(ctx)
if err != nil {
    log.Fatalf("MarketData connection failed: %v", err)
}

// Subscribe to all data types for a contract
contractID := "12345" // Contract ID from search
err = client.MarketData.SubscribeAll(contractID)
if err != nil {
    log.Fatalf("Market data subscription failed: %v", err)
}

// Or subscribe to specific data types
err = client.MarketData.SubscribeContractQuotes(contractID)
err = client.MarketData.SubscribeContractTrades(contractID)
err = client.MarketData.SubscribeContractMarketDepth(contractID)

// Unsubscribe when done
defer func() {
    client.MarketData.UnsubscribeAll(contractID)
    client.MarketData.Disconnect()
}()
```

### Connection Management

Both UserData and MarketData WebSocket services include sophisticated connection management:

#### Connection States

```go
// Connection states
const (
    StateDisconnected ConnectionState = iota
    StateConnecting
    StateConnected
    StateReconnecting
)

// Check current connection state
if client.UserData.IsConnected() {
    fmt.Println("UserData is connected")
}

state := client.MarketData.GetConnectionState()
fmt.Printf("MarketData state: %v\n", state)
```

#### Automatic Reconnection Features

- **Exponential Backoff**: Automatic reconnection with increasing delays (up to 30 seconds)
- **Health Monitoring**: Ping/pong health checks every 5 seconds
- **Subscription Recovery**: Automatically resubscribes to previous subscriptions after reconnection
- **Timeout Protection**: 10-second ping timeout prevents indefinite blocking
- **Connection Validation**: Built-in validation to ensure stable connections

#### Best Practices

```go
// Set up connection handlers to monitor state changes
client.UserData.SetConnectionHandler(func(state services.ConnectionState) {
    switch state {
    case services.StateConnected:
        log.Println("UserData connected - ready for subscriptions")
    case services.StateReconnecting:
        log.Println("UserData reconnecting - please wait")
    case services.StateDisconnected:
        log.Println("UserData disconnected")
    }
})

// Always use proper cleanup
defer func() {
    client.UserData.Disconnect()
    client.MarketData.Disconnect()
    client.Disconnect(ctx)
}()
```

## Order Types

The library supports all TopStepX order types:

```go
// Market Orders
models.OrderTypeMarket

// Limit Orders
models.OrderTypeLimit

// Stop Orders
models.OrderTypeStop

// Stop-Limit Orders
models.OrderTypeStopLimit

// Trailing Stop Orders
models.OrderTypeTrailingStop

// Join Bid/Offer Orders
models.OrderTypeJoinBid
models.OrderTypeJoinAsk
```

## Bar Types and Timeframes

Historical data supports multiple timeframes:

```go
models.AggregateBarUnitSecond // Second bars
models.AggregateBarUnitMinute // Minute bars (1, 5, 15, 30, etc.)
models.AggregateBarUnitHour    // Hourly bars
models.AggregateBarUnitDay     // Daily bars
models.AggregateBarUnitWeek    // Weekly bars
models.AggregateBarUnitMonth // Monthly bars
```

## Error Handling

The library provides structured error handling:

```go
resp, err := client.Order.PlaceOrder(ctx, orderRequest)
if err != nil {
// Network or client error
log.Printf("Request failed: %v", err)
return
}

if !resp.Success {
// API error
if resp.ErrorMessage != nil {
log.Printf("API error: %s", *resp.ErrorMessage)
}
log.Printf("Error code: %v", resp.ErrorCode)
return
}

// Success - use resp.Order
```

## Performance Considerations

### Connection Management

- Reuse client instances across operations
- Use context with appropriate timeouts
- Handle user data reconnections gracefully

### Rate Limiting

- The API may have rate limits - implement appropriate backoff strategies
- Consider batching operations when possible
- Monitor for rate limit error codes

### Memory Management

- Close user data connections when done
- Use context cancellation for long-running operations
- Properly handle large result sets

## Testing

Run the example programs to test functionality:

```bash
# Account management example
cd samples/account && go run account_example.go

# Basic order operations
cd samples/order && go run order_example.go

# Market order with automatic position closing
cd samples/market_order && go run market_order_example.go

# Limit order example
cd samples/limit_order && go run main.go

# Type-safe order event handling
cd samples/order_handler && go run main.go

# Position management example
cd samples/position && go run position_example.go

# Type-safe position event handling
cd samples/position_handler && go run main.go

# User data WebSocket integration
cd samples/userdata && go run userdata_example.go

# Real-time market data streaming
cd samples/marketdata && go run marketdata_example.go

# Historical data example
cd samples/history && go run history_example.go

# Contract search example  
cd samples/contract && go run contract_example.go

# Trade history example
cd samples/trade && go run trade_example.go

# Status monitoring example
cd samples/status && go run status_example.go
```

## Examples

The `samples/` directory contains comprehensive examples:

- **Account Management**: `samples/account/account_example.go`
- **Order Operations**: `samples/order/order_example.go`
- **Market Orders**: `samples/market_order/market_order_example.go` - Complete workflow with position closing
- **Limit Orders**: `samples/limit_order/main.go`
- **Type-Safe Order Handlers**: `samples/order_handler/main.go` - Real-time order event handling
- **Position Management**: `samples/position/position_example.go`
- **Type-Safe Position Handlers**: `samples/position_handler/main.go` - Real-time position event handling
- **User Data Integration**: `samples/userdata/userdata_example.go` - WebSocket user data streaming
- **Market Data Streaming**: `samples/marketdata/marketdata_example.go` - Real-time market data
- **Historical Data**: `samples/history/history_example.go`
- **Contract Search**: `samples/contract/contract_example.go`
- **Trade History**: `samples/trade/trade_example.go`
- **Status Monitoring**: `samples/status/status_example.go`

## Security

- **Never commit API credentials** to version control
- Use environment variables for sensitive configuration
- Implement proper token management and refresh
- Monitor for unauthorized access patterns
- Use HTTPS for all API communications

## Status Codes

Common response codes and their meanings:

- `Success = true`: Operation completed successfully
- `Success = false`: Check `ErrorCode` and `ErrorMessage` for details
- Network errors: Connection, timeout, or DNS issues
- Authentication errors: Invalid credentials or expired tokens

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Documentation**: Check the `samples/` directory for usage examples
- **Issues**: Report bugs and feature requests via GitHub Issues
- **API Reference**: Refer to TopStepX official API documentation

## Changelog

### Latest Updates

- **Type-Safe Event Handlers**: Structured data models for order and position updates with automatic parsing
- **Market Data Streaming**: Full WebSocket support for real-time quotes, trades, and market depth
- **Enhanced Connection Management**: Automatic reconnection with exponential backoff and health monitoring
- **Ping Timeout Protection**: Built-in connection validation prevents indefinite blocking
- **Comprehensive Sample Applications**: New examples for market orders, limit orders, and event handling
- **Convenience Methods**: High-level methods like `LoginAndConnect` and `GetFirstActiveAccount`
- **Connection State Monitoring**: Real-time connection state tracking with event handlers
- **Subscription Recovery**: Automatic resubscription after WebSocket reconnection
- **Multi-account Management**: Enhanced account search and selection capabilities
- **Advanced Order Types**: Support for all TopStepX order types including trailing stops

---

**Note**: This library is for educational and development purposes. Always test thoroughly in a demo environment before
using with live trading accounts. Trading involves significant risk of loss.