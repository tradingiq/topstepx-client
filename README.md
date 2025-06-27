# TopStepX Go Client Library

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**This library is currently under active development and may undergo breaking changes.**

A comprehensive Go client library for the TopStepX trading platform API, providing full access to trading operations, real-time data streaming, and account management functionality.

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
- **Live Market Data**: Real-time price updates and market information

### Advanced Features
- **Authentication Management**: Automatic token handling and refresh
- **Error Handling**: Comprehensive error codes and structured responses
- **Context Support**: Full context.Context support for timeouts and cancellation
- **Concurrent Operations**: Thread-safe operations with proper synchronization

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

```go
// Set up event handlers
client.UserData.SetAccountHandler(func(data interface{}) {
    fmt.Printf("Account update: %+v\n", data)
})

client.UserData.SetOrderHandler(func(data interface{}) {
    fmt.Printf("Order update: %+v\n", data)
})

client.UserData.SetPositionHandler(func(data interface{}) {
    fmt.Printf("Position update: %+v\n", data)
})

client.UserData.SetTradeHandler(func(data interface{}) {
    fmt.Printf("Trade update: %+v\n", data)
})

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
models.AggregateBarUnitSecond  // Second bars
models.AggregateBarUnitMinute  // Minute bars (1, 5, 15, 30, etc.)
models.AggregateBarUnitHour    // Hourly bars
models.AggregateBarUnitDay     // Daily bars
models.AggregateBarUnitWeek    // Weekly bars
models.AggregateBarUnitMonth   // Monthly bars
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
# Account example
cd samples/account && go run account_example.go

# Order example
cd samples/order && go run order_example.go

# User data integration example
cd samples/userdata && go run userdata_example.go

# Position management example
cd samples/position && go run position_example.go

# Historical data example
cd samples/history && go run history_example.go
```

## Examples

The `samples/` directory contains comprehensive examples:

- **Account Management**: `samples/account/account_example.go`
- **Order Operations**: `samples/order/order_example.go`
- **Position Management**: `samples/position/position_example.go`
- **Historical Data**: `samples/history/history_example.go`
- **Contract Search**: `samples/contract/contract_example.go`
- **Trade History**: `samples/trade/trade_example.go`
- **User Data Integration**: `samples/userdata/userdata_example.go`
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
- Full user data support with automatic reconnection
- Comprehensive error handling and logging
- Multi-account management capabilities
- Historical data access with multiple timeframes
- Advanced order types including trailing stops

---

**Note**: This library is for educational and development purposes. Always test thoroughly in a demo environment before using with live trading accounts. Trading involves significant risk of loss.