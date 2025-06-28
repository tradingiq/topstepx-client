# TopStepX Go Client Library

A comprehensive Go SDK for interacting with the TopStepX trading platform API. This library provides complete access to trading operations, real-time market data, account management, and historical data retrieval.

## Features

- **Complete Trading Operations**: Place, modify, and cancel orders with all order types supported
- **Real-time Data Streaming**: WebSocket integration for live market data and user data updates
- **Account Management**: Multi-account support with balance tracking and permissions
- **Position Management**: Monitor and manage open positions with partial closing capabilities
- **Market Data**: Access historical price data and real-time quotes, trades, and market depth
- **Automatic Authentication**: Built-in token management and session handling
- **Type-Safe**: Structured data models with type-safe enums for all trading operations
- **Robust WebSocket**: Automatic reconnection, subscription recovery, and health monitoring

## Installation

```bash
go get github.com/tradingiq/topstepx-client
```

## Requirements

- Go 1.24 or higher
- Active TopStepX trading account
- API credentials (username and API key)

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
    "github.com/tradingiq/topstepx-client"
    "github.com/tradingiq/topstepx-client/models"
)

func main() {
    // Create client
    client := topstepx.NewClient()
    
    ctx := context.Background()
    
    // Login with API credentials
    resp, err := client.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
        UserName: os.Getenv("TOPSTEPX_USERNAME"),
        APIKey:   os.Getenv("TOPSTEPX_API_KEY"),
    })
    if err != nil {
        log.Fatal("Failed to login:", err)
    }
    if !resp.Success {
        log.Fatal("Login failed:", resp.ErrorMessage)
    }
    
    // Get active accounts
    accountsResp, err := client.Account.SearchAccounts(ctx, &models.SearchAccountRequest{
        OnlyActiveAccounts: true,
    })
    if err != nil || !accountsResp.Success || len(accountsResp.Accounts) == 0 {
        log.Fatal("Failed to get accounts")
    }
    account := accountsResp.Accounts[0]
    
    // Place a market order
    orderResp, err := client.Order.PlaceOrder(ctx, &models.PlaceOrderRequest{
        AccountID:  int32(account.ID),
        ContractID: "MES",
        Type:       models.OrderTypeMarket,
        Side:       models.OrderSideBid,
        Size:       1,
    })
    if err != nil {
        log.Fatal("Failed to place order:", err)
    }
    
    fmt.Printf("Order placed: %d\n", *orderResp.OrderID)
}
```

## Core Services

### Authentication Service
```go
// Login with API key
resp, err := client.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
    UserName: username,
    APIKey:   apiKey,
})

// Check authentication status
isValid := client.Auth.IsTokenValid()

// Logout
logoutResp, err := client.Auth.Logout(ctx)
```

### Account Service
```go
// Search accounts with filters
resp, err := client.Account.SearchAccounts(ctx, &models.SearchAccountRequest{
    OnlyActiveAccounts: true,
})

// Get specific account
accountResp, err := client.Account.GetAccount(ctx, &models.GetAccountRequest{
    AccountID: accountID,
})
```

### Order Service
```go
// Place order (supports all order types)
orderResp, err := client.Order.PlaceOrder(ctx, &models.PlaceOrderRequest{
    AccountID:  int32(accountID),
    ContractID: contractID,
    Type:       models.OrderTypeMarket,  // or Limit, Stop, StopLimit, etc.
    Side:       models.OrderSideBid,     // or Ask
    Size:       1,
    LimitPrice: &limitPrice,  // for limit orders
    StopPrice:  &stopPrice,   // for stop orders
})

// Modify order
modifyResp, err := client.Order.ModifyOrder(ctx, &models.ModifyOrderRequest{
    AccountID:  int32(accountID),
    OrderID:    orderID,
    LimitPrice: &newPrice,
    Size:       newQuantity,
})

// Cancel order
cancelResp, err := client.Order.CancelOrder(ctx, &models.CancelOrderRequest{
    AccountID: int32(accountID),
    OrderID:   orderID,
})

// Search orders
ordersResp, err := client.Order.SearchOrders(ctx, &models.SearchOrderRequest{
    AccountID:      int32(accountID),
    StartTimestamp: startTime,
    EndTimestamp:   &endTime,
})
```

### Position Service
```go
// Search open positions
positionsResp, err := client.Position.SearchOpenPositions(ctx, &models.SearchPositionRequest{
    AccountID: int32(accountID),
})

// Close entire position by contract
closeResp, err := client.Position.CloseContractPosition(ctx, &models.CloseContractPositionRequest{
    AccountID:  int32(accountID),
    ContractID: contractID,
})

// Close partial position
closeResp, err := client.Position.CloseContractPosition(ctx, &models.CloseContractPositionRequest{
    AccountID:  int32(accountID),
    ContractID: contractID,
    Size:       2,  // partial close
})
```

### Market Data Service
```go
// Search contracts
searchText := "MES"
contractsResp, err := client.Contract.SearchContracts(ctx, &models.SearchContractRequest{
    SearchText: &searchText,
    Live:       false,  // false for demo, true for live
})

// Get historical data
historyResp, err := client.History.GetBarHistory(ctx, &models.GetBarHistoryRequest{
    ContractID:    contractID,
    BarSizeUnit:   models.BarSizeUnitMinute,
    BarSizeValue:  1,
    StartDateTime: startTime,
    EndDateTime:   &endTime,
})
```

## WebSocket Streaming

### User Data WebSocket
```go
// Connect to user data stream
err := client.UserData.Connect(ctx)
if err != nil {
    log.Fatal("Failed to connect:", err)
}

// Subscribe to all events for an account
err = client.UserData.SubscribeAll(accountID)

// Or subscribe to specific events
err = client.UserData.SubscribeOrders(accountID)
err = client.UserData.SubscribePositions(accountID)
err = client.UserData.SubscribeTrades(accountID)

// Set handlers for order updates
client.UserData.SetOrderHandler(func(data *models.OrderUpdateData) {
    fmt.Printf("Order update: %d - Status: %s\n", data.Data.ID, data.Data.Status)
})

// Set handlers for position updates
client.UserData.SetPositionHandler(func(data *models.PositionUpdateData) {
    fmt.Printf("Position update: %s - Size: %d\n", data.Data.ContractID, data.Data.Size)
})

// Set handlers for trade updates
client.UserData.SetTradeHandler(func(data interface{}) {
    fmt.Printf("Trade update: %+v\n", data)
})
```

### Market Data WebSocket
```go
// Connect to market data stream
err := client.MarketData.Connect(ctx)
if err != nil {
    log.Fatal("Failed to connect:", err)
}

// Subscribe to all data types for a contract
err = client.MarketData.SubscribeAll(contractID)

// Or subscribe to specific data types
err = client.MarketData.SubscribeQuotes(contractID)
err = client.MarketData.SubscribeTrades(contractID)
err = client.MarketData.SubscribeDepth(contractID)

// Set quote handler
client.MarketData.SetQuoteHandler(func(contractID string, quote models.Quote) {
    fmt.Printf("Quote for %s - Bid: %.2f, Ask: %.2f\n", 
        contractID, quote.BestBid, quote.BestAsk)
})

// Set trade handler
client.MarketData.SetTradeHandler(func(contractID string, trades models.TradeData) {
    for _, trade := range trades {
        fmt.Printf("Trade for %s - Price: %.2f, Volume: %.0f\n", 
            contractID, trade.Price, trade.Volume)
    }
})

// Set depth handler
client.MarketData.SetDepthHandler(func(contractID string, depth models.MarketDepthData) {
    fmt.Printf("Market depth update for %s\n", contractID)
})
```

## Examples

The library includes comprehensive examples demonstrating all major functionalities:

### Basic Trading Examples
- **market_order**: Place market orders and manage positions
- **limit_order**: Place and cancel limit orders
- **order**: Comprehensive order management (search, place, modify, cancel)
- **position**: Position management with partial closing

### Real-time Data Examples
- **marketdata**: Stream real-time quotes, trades, and market depth
- **userdata**: Monitor account updates, orders, positions, and trades
- **order_handler**: Handle real-time order lifecycle events
- **position_handler**: Track position changes in real-time

### Account & Information Examples
- **account**: Search and filter trading accounts
- **contract**: Search trading instruments and get contract specifications
- **history**: Retrieve historical price data with various intervals
- **trade**: Analyze trade history and calculate statistics
- **status**: Monitor API health and connection status

To run an example:
```bash
cd samples/market_order
go run market_order_example.go
```

## Configuration

Set your credentials using environment variables:
```bash
export TOPSTEPX_USERNAME="your_username"
export TOPSTEPX_API_KEY="your_api_key"
```

Or use a `.env` file:
```env
TOPSTEPX_USERNAME=your_username
TOPSTEPX_API_KEY=your_api_key
```

## Architecture

### Service-Oriented Design
- Each service handles a specific domain (orders, positions, accounts, etc.)
- Services are accessible through the main client instance
- All operations support context for cancellation and timeouts

### WebSocket Implementation
- Built on SignalR protocol for reliable real-time communication
- Automatic reconnection with exponential backoff
- Connection state monitoring and recovery
- Type-safe event handlers with automatic data parsing

### Error Handling
- Typed error codes for different failure scenarios
- Context-aware error propagation
- Detailed error messages for debugging

## API Endpoints

- Base API: `https://api.topstepx.com`
- User Data WebSocket: `https://rtc.topstepx.com/hubs/user`
- Market Data WebSocket: `https://rtc.topstepx.com/hubs/market`

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This project is licensed under the MIT License with Attribution - see the LICENSE file for details.

## Author

- Victor Geyer
- Trading IQ

## Support

For issues, questions, or contributions, please visit the [GitHub repository](https://github.com/tradingiq/topstepx-client).