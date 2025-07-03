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

	// Get user data websocket service
	userDataWS := client.UserData

	// Set up account update handler with detailed printf debugging
	userDataWS.SetAccountHandler(func(data interface{}) {
		fmt.Println("\n========== ACCOUNT UPDATE RECEIVED ==========")
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
		if accountMap, ok := data.(map[string]interface{}); ok {
			fmt.Println("\n--- Extracted Fields ---")

			// Common account fields to look for
			fields := []string{
				"accountId", "AccountId", "id", "Id",
				"name", "Name", "accountName", "AccountName",
				"balance", "Balance", "equity", "Equity",
				"marginUsed", "MarginUsed", "marginAvailable", "MarginAvailable",
				"currency", "Currency", "status", "Status",
				"type", "Type", "accountType", "AccountType",
			}

			for _, field := range fields {
				if value, exists := accountMap[field]; exists {
					fmt.Printf("%s: %v\n", field, value)
				}
			}
		}

		fmt.Println("=============================================\n")
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

	// Subscribe to accounts
	fmt.Println("Subscribing to account updates...")
	if err := userDataWS.SubscribeAccounts(); err != nil {
		log.Fatalf("Failed to subscribe to accounts: %v", err)
	}
	fmt.Println("Successfully subscribed to account updates!")
	fmt.Println("\nListening for account updates... (Press Ctrl+C to exit)")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Keep the program running
	<-sigChan

	// Cleanup
	fmt.Println("\nShutting down...")

	// Unsubscribe from accounts
	if err := userDataWS.UnsubscribeAccounts(); err != nil {
		fmt.Printf("Error unsubscribing from accounts: %v\n", err)
	}

	// Disconnect
	if err := userDataWS.Disconnect(); err != nil {
		fmt.Printf("Error disconnecting: %v\n", err)
	}

	fmt.Println("Shutdown complete.")
}
