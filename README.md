# TopStepX Go Client Library

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

⚠️ **This library is currently under active development and may undergo breaking changes.**

A comprehensive Go client library for the TopStepX trading platform API, providing full access to trading operations, real-time data streaming, and account management functionality.

## 🚀 Features

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
- **Live Market Data**: Real-time price updates and market information

### Advanced Features
- **Authentication Management**: Automatic token handling and refresh
- **Error Handling**: Comprehensive error codes and structured responses
- **Context Support**: Full context.Context support for timeouts and cancellation
- **Concurrent Operations**: Thread-safe operations with proper synchronization

## 📦 Installation

```bash
go get github.com/tradingiq/topstepx-client
```

### Prerequisites

- Go 1.24 or higher
- Active TopStepX trading account
- API credentials (username and API key)

## 🔧 Quick Start

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

## 📖 API Reference

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
    OrderType:  models.OrderTypeMarket,
    Side:       models.OrderSideBuy,
    Quantity:   1,
})

// Place a limit order
resp, err := client.Order.PlaceOrder(ctx, &models.PlaceOrderRequest{
    AccountID:  accountID,
    ContractID: contractID,
    OrderType:  models.OrderTypeLimit,
    Side:       models.OrderSideBuy,
    Quantity:   1,
    Price:      &limitPrice,
})

// Search open orders
resp, err := client.Order.SearchOpenOrders(ctx, &models.SearchOpenOrdersRequest{
    AccountID: accountID,
})

// Modify existing order
resp, err := client.Order.ModifyOrder(ctx, &models.ModifyOrderRequest{
    OrderID:  orderID,
    Price:    &newPrice,
    Quantity: &newQuantity,
})

// Cancel order
resp, err := client.Order.CancelOrder(ctx, &models.CancelOrderRequest{
    OrderID: orderID,
})
```

### Position Management

```go
// Search current positions
resp, err := client.Position.SearchPositions(ctx, &models.SearchPositionsRequest{
    AccountID: accountID,
})

// Close entire position
resp, err := client.Position.ClosePosition(ctx, &models.ClosePositionRequest{
    AccountID:  accountID,
    ContractID: contractID,
})

// Partial position close
resp, err := client.Position.ClosePosition(ctx, &models.ClosePositionRequest{
    AccountID:  accountID,
    ContractID: contractID,
    Quantity:   &partialQuantity,
})
```

### Historical Data

```go
// Get price bars
resp, err := client.History.GetBars(ctx, &models.GetBarsRequest{
    ContractID: contractID,
    BarType:    models.BarTypeMinute,
    BarSize:    5, // 5-minute bars
    FromDate:   &startTime,
    ToDate:     &endTime,
    Limit:      1000,
})
```

### Contract Information

```go
// Search contracts
resp, err := client.Contract.SearchContracts(ctx, &models.SearchContractsRequest{
    SearchText: "ES", // E-mini S&P 500
    Live:       true,
})

// Get contract by ID
resp, err := client.Contract.GetContract(ctx, contractID)
```

### WebSocket Real-time Data

```go
// Set up event handlers
client.WebSocket.SetAccountHandler(func(data interface{}) {
    fmt.Printf("Account update: %+v\n", data)
})

client.WebSocket.SetOrderHandler(func(data interface{}) {
    fmt.Printf("Order update: %+v\n", data)
})

client.WebSocket.SetPositionHandler(func(data interface{}) {
    fmt.Printf("Position update: %+v\n", data)
})

client.WebSocket.SetTradeHandler(func(data interface{}) {
    fmt.Printf("Trade update: %+v\n", data)
})

// Connect and subscribe
err := client.WebSocket.Connect(ctx)
if err != nil {
    log.Fatalf("WebSocket connection failed: %v", err)
}

// Subscribe to all events for an account
err = client.WebSocket.SubscribeAll(accountID)
if err != nil {
    log.Fatalf("Subscription failed: %v", err)
}

// Or subscribe to specific events
err = client.WebSocket.SubscribeToAccount(accountID)
err = client.WebSocket.SubscribeToOrders(accountID)
err = client.WebSocket.SubscribeToPositions(accountID)
err = client.WebSocket.SubscribeToTrades(accountID)
```

## 🔧 Order Types

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

## 📊 Bar Types and Timeframes

Historical data supports multiple timeframes:

```go
models.BarTypeSecond  // Second bars
models.BarTypeMinute  // Minute bars (1, 5, 15, 30, etc.)
models.BarTypeHour    // Hourly bars
models.BarTypeDay     // Daily bars
models.BarTypeWeek    // Weekly bars
models.BarTypeMonth   // Monthly bars
```

## 🔐 Error Handling

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

## ⚡ Performance Considerations

### Connection Management
- Reuse client instances across operations
- Use context with appropriate timeouts
- Handle WebSocket reconnections gracefully

### Rate Limiting
- The API may have rate limits - implement appropriate backoff strategies
- Consider batching operations when possible
- Monitor for rate limit error codes

### Memory Management
- Close WebSocket connections when done
- Use context cancellation for long-running operations
- Properly handle large result sets

## 🧪 Testing

Run the example programs to test functionality:

```bash
# Account example
cd samples/account && go run account_example.go

# Order example
cd samples/order && go run order_example.go

# WebSocket integration example
cd samples/websocket && go run integrated_example.go

# Position management example
cd samples/position && go run position_example.go

# Historical data example
cd samples/history && go run history_example.go
```

## 📝 Examples

The `samples/` directory contains comprehensive examples:

- **Account Management**: `samples/account/account_example.go`
- **Order Operations**: `samples/order/order_example.go`
- **Position Management**: `samples/position/position_example.go`
- **Historical Data**: `samples/history/history_example.go`
- **Contract Search**: `samples/contract/contract_example.go`
- **Trade History**: `samples/trade/trade_example.go`
- **WebSocket Integration**: `samples/websocket/integrated_example.go`
- **Status Monitoring**: `samples/status/status_example.go`

## 🔒 Security

- **Never commit API credentials** to version control
- Use environment variables for sensitive configuration
- Implement proper token management and refresh
- Monitor for unauthorized access patterns
- Use HTTPS for all API communications

## 🚦 Status Codes

Common response codes and their meanings:

- `Success = true`: Operation completed successfully
- `Success = false`: Check `ErrorCode` and `ErrorMessage` for details
- Network errors: Connection, timeout, or DNS issues
- Authentication errors: Invalid credentials or expired tokens

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Documentation**: Check the `samples/` directory for usage examples
- **Issues**: Report bugs and feature requests via GitHub Issues
- **API Reference**: Refer to TopStepX official API documentation

## 🔄 Changelog

### Latest Updates
- Full WebSocket support with automatic reconnection
- Comprehensive error handling and logging
- Multi-account management capabilities
- Historical data access with multiple timeframes
- Advanced order types including trailing stops

---

**Note**: This library is for educational and development purposes. Always test thoroughly in a demo environment before using with live trading accounts. Trading involves significant risk of loss.