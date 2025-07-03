package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tradingiq/topstepx-client"
	"github.com/tradingiq/topstepx-client/models"
	"github.com/tradingiq/topstepx-client/samples"
	"github.com/tradingiq/topstepx-client/services"
)

func main() {
	// Initialize the client
	client := topstepx.NewClient()
	ctx := context.Background()

	// Authenticate using API key
	fmt.Println("Authenticating...")
	resp, err := client.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
		UserName: samples.Config.Username,
		APIKey:   samples.Config.ApiKey,
	})
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}
	if !resp.Success {
		log.Fatalf("Authentication failed: %v", resp.ErrorMessage)
	}
	fmt.Println("Authentication successful!")

	// Get accounts to select one for trade monitoring
	fmt.Println("\nFetching accounts...")
	accounts, err := client.GetActiveAccounts(ctx)
	if err != nil {
		log.Fatalf("Failed to get accounts: %v", err)
	}

	if len(accounts) == 0 {
		log.Fatal("No accounts found")
	}

	// Use the first account
	account := accounts[0]
	fmt.Printf("Using account: %s (ID: %d)\n", account.Name, account.ID)

	// Get user data websocket service
	userDataWS := client.UserData

	// Set up trade update handler with detailed printf debugging
	userDataWS.SetTradeHandler(func(data interface{}) {
		fmt.Println("\n========== TRADE UPDATE RECEIVED ==========")
		fmt.Printf("Timestamp: %s\n", time.Now().Format("2006-01-02 15:04:05.000"))

		// Print raw data type
		fmt.Printf("Data Type: %T\n", data)

		// Pretty print the raw JSON data
		jsonBytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling data: %v\n", err)
			fmt.Printf("Raw data: %+v\n", data)
		} else {
			fmt.Printf("JSON Data:\n%s\n", string(jsonBytes))
		}

		// Try to extract specific fields if it's a map
		if tradeMap, ok := data.(map[string]interface{}); ok {
			fmt.Println("\n--- Extracted Trade Fields ---")

			// Common trade fields to look for
			fields := []string{
				"tradeId", "TradeId", "id", "Id",
				"orderId", "OrderId", "orderID", "OrderID",
				"accountId", "AccountId", "accountID", "AccountID",
				"symbol", "Symbol", "contractId", "ContractId",
				"side", "Side", "direction", "Direction",
				"quantity", "Quantity", "qty", "Qty",
				"price", "Price", "fillPrice", "FillPrice",
				"commission", "Commission", "fees", "Fees",
				"executionTime", "ExecutionTime", "timestamp", "Timestamp",
				"tradeType", "TradeType", "type", "Type",
				"status", "Status", "state", "State",
			}

			for _, field := range fields {
				if value, exists := tradeMap[field]; exists {
					fmt.Printf("%s: %v", field, value)

					// Special formatting for certain fields
					switch field {
					case "side", "Side", "direction", "Direction":
						if v, ok := value.(float64); ok {
							if v == 1 {
								fmt.Printf(" (BUY)")
							} else if v == -1 || v == 2 {
								fmt.Printf(" (SELL)")
							}
						}
					case "executionTime", "ExecutionTime", "timestamp", "Timestamp":
						// Try to parse as timestamp if it's a string
						if timeStr, ok := value.(string); ok {
							fmt.Printf(" (parsed: %s)", timeStr)
						}
					}
					fmt.Println()
				}
			}

			// Calculate P&L if available
			if qty, hasQty := tradeMap["quantity"]; hasQty {
				if price, hasPrice := tradeMap["price"]; hasPrice {
					if qtyFloat, ok1 := qty.(float64); ok1 {
						if priceFloat, ok2 := price.(float64); ok2 {
							fmt.Printf("Trade Value: %.2f\n", qtyFloat*priceFloat)
						}
					}
				}
			}
		}

		fmt.Println("==========================================\n")
	})

	// Set connection state handler
	userDataWS.SetConnectionHandler(func(state services.ConnectionState) {
		fmt.Printf("[Connection State Changed] %s at %s\n",
			state, time.Now().Format("15:04:05.000"))
	})

	// Connect to websocket
	fmt.Println("\nConnecting to user data websocket...")
	if err := userDataWS.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Wait a moment for connection to establish
	time.Sleep(2 * time.Second)

	// Subscribe to trades
	fmt.Printf("Subscribing to trade updates for account %d...\n", account.ID)
	if err := userDataWS.SubscribeTrades(int(account.ID)); err != nil {
		log.Fatalf("Failed to subscribe to trades: %v", err)
	}
	fmt.Println("Successfully subscribed to trade updates!")
	fmt.Println("\nListening for trade updates... (Press Ctrl+C to exit)")
	fmt.Println("NOTE: Trade updates will appear when trades are executed on this account")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Keep the program running
	<-sigChan

	// Cleanup
	fmt.Println("\nShutting down...")

	// Unsubscribe from trades
	if err := userDataWS.UnsubscribeTrades(int(account.ID)); err != nil {
		fmt.Printf("Error unsubscribing from trades: %v\n", err)
	}

	// Disconnect
	if err := userDataWS.Disconnect(); err != nil {
		fmt.Printf("Error disconnecting: %v\n", err)
	}

	fmt.Println("Shutdown complete.")
}
