# TopStepX Client Samples

This directory contains comprehensive examples demonstrating how to use the TopStepX Go client library.

## Prerequisites

Before running any samples, you need to set the following environment variables:

```bash
export TOPSTEPX_USERNAME="your_username"
export TOPSTEPX_API_KEY="your_api_key"
```

## Available Samples

### 1. Authentication (`auth/`)

Demonstrates login, validation, and logout operations.

```bash
go run auth/auth_example.go
```

### 2. Account Management (`account/`)

Shows how to search for accounts and use convenience methods.

```bash
go run account/account_example.go
```

### 3. Contract Search (`contract/`)

Examples of searching for contracts by symbol and retrieving contract details.

```bash
go run contract/contract_example.go
```

### 4. Historical Data (`history/`)

Demonstrates retrieving historical bar data with different time intervals.

```bash
go run history/history_example.go
```

### 5. Order Management (`order/`)

Shows placing, modifying, canceling orders, and searching order history.

```bash
go run order/order_example.go
```

### 6. Order Handler (`order_handler/`)

Demonstrates real-time order updates using the structured order handler with type-safe enum access and automatic
timestamp parsing.

```bash
go run order_handler/main.go
```

### 7. Position Management (`position/`)

Examples of searching positions and closing positions.

```bash
go run position/position_example.go
```

### 8. Trade History (`trade/`)

Demonstrates searching trade history and calculating trading statistics.

```bash
go run trade/trade_example.go
```

### 9. Status Check (`status/`)

Shows how to check API status and implement health checks.

```bash
go run status/status_example.go
```

### 10. User Data Streaming (`userdata/`)

#### Basic User Data Example (`userdata/userdata_example.go`)

Comprehensive example showing real-time account, order, position, and trade updates with structured order handling.

```bash
go run userdata/userdata_example.go
```

#### Order Handler Demo (`userdata/order_handler_demo.go`)

Detailed demo showcasing the structured order handler with formatted output and order interpretation.

```bash
go run userdata/order_handler_demo.go
```

### 11. Client Configuration (`client_options/`)

Demonstrates various client configuration options.

```bash
go run client_options/client_options_example.go
```

### 12. Complete Trading Example (`complete_example/`)

A full workflow example combining multiple services.

```bash
go run complete_example/complete_example.go
```

## Additional Examples

### Integrated User Data Example

The `userdata/userdata_example.go` file shows a production-ready user data integration with automatic reconnection.

### Configuration Helper

The `config.go` file demonstrates loading configuration from environment variables using `.env` files.

## Important Notes

1. **Demo vs Live**: Most examples use demo mode (`Live: false`). Be careful when switching to live trading.

2. **Error Handling**: Examples include basic error handling. Production code should implement more robust error
   handling.

3. **Rate Limiting**: Be aware of API rate limits when running examples repeatedly.

4. **User Data Connections**: User data examples maintain persistent connections. Use Ctrl+C to gracefully shutdown.

5. **Account IDs**: Examples automatically use the first active account. Modify as needed for your use case.

6. **Order Handlers**: The WebSocket service now supports both structured (`SetOrderStructHandler`) and raw (
   `SetOrderHandler`) handlers for type-safe order data access with automatic timestamp parsing and enum types.

## Building Examples

To build a specific example:

```bash
cd samples/auth
go build auth_example.go
```

To build all examples:

```bash
cd samples
for dir in */; do
    if [ -f "$dir"/*.go ]; then
        echo "Building $dir"
        cd "$dir" && go build *.go && cd ..
    fi
done
```

## License

These samples are provided as-is for educational purposes. See the main repository LICENSE for details.