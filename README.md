# ProjectX Go Client Library

A comprehensive Go SDK for interacting with the ProjectX trading platform API. This library provides complete access to trading operations, real-time market data, account management, and historical data retrieval.

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
go get github.com/tradingiq/projectx-client
```

## Requirements

- Go 1.24 or higher
- Active ProjectX trading account
- API credentials (username and API key)

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
    "github.com/tradingiq/projectx-client"
    "github.com/tradingiq/projectx-client/models"
)

func main() {
    // Create client
    client := projectx.NewClient()
    
    ctx := context.Background()
    
    // Login with API credentials
    resp, err := client.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
        UserName: os.Getenv("PROJECTX_USERNAME"),
        APIKey:   os.Getenv("PROJECTX_API_KEY"),
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

The library includes focused examples demonstrating specific features. Each example is self-contained and demonstrates a single topic.

### Running Examples

To run any example, navigate to its directory:

```bash
cd samples/01_authentication
go run main.go
```

Or run directly from the project root:

```bash
go run ./samples/01_authentication/main.go
```

### Available Examples

#### Authentication & Account Management
- **01_authentication/** - Basic API authentication
- **02_list_accounts/** - List trading accounts and view balances

#### Order Management
- **03_place_market_order/** - Place a market order
- **04_place_limit_order/** - Place a limit order
- **05_modify_order/** - Modify an existing order
- **06_cancel_order/** - Cancel an order

#### Position Management
- **07_list_positions/** - View open positions with P&L
- **08_close_position/** - Close an open position

#### Market Data
- **09_search_contracts/** - Search for trading contracts
- **10_get_historical_data/** - Retrieve historical price bars

#### Real-time Data Streaming
- **11_stream_quotes/** - Stream real-time market quotes
- **12_monitor_orders/** - Monitor order updates via WebSocket

### Example Structure

Each example follows this pattern:
1. **Authentication** - Login with API credentials
2. **Setup** - Get accounts, find contracts, etc.
3. **Main Action** - Demonstrate the specific feature
4. **Output** - Show results with clear formatting

### Key Concepts

- **Error Handling**: All examples include proper error handling and check response success flags
- **Resource Cleanup**: WebSocket examples properly disconnect and unsubscribe on exit
- **Demo vs Live Markets**: Examples use demo markets (`Live: false`) by default for safety
- **Type Safety**: Examples demonstrate the library's type-safe approach with structured data models

## Configuration

Set your credentials using environment variables:
```bash
export PROJECTX_USERNAME="your_username"
export PROJECTX_API_KEY="your_api_key"
```

Or use a `.env` file:
```env
PROJECTX_USERNAME=your_username
PROJECTX_API_KEY=your_api_key
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

- Base API: `https://api.projectx.com`
- User Data WebSocket: `https://rtc.projectx.com/hubs/user`
- Market Data WebSocket: `https://rtc.projectx.com/hubs/market`

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This project is licensed under the MIT License with Attribution - see the LICENSE file for details.

## Author

- Victor Geyer
- Trading IQ

## Support

For issues, questions, or contributions, please visit the [GitHub repository](https://github.com/tradingiq/projectx-client).