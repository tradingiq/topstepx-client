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

	"github.com/tradingiq/projectx-client"
	"github.com/tradingiq/projectx-client/models"
	"github.com/tradingiq/projectx-client/samples"
	"github.com/tradingiq/projectx-client/services"
)

func main() {

	client := projectx.NewClient()
	ctx := context.Background()

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

	userDataWS := client.UserData

	userDataWS.SetAccountHandler(func(update *models.AccountUpdateData) {
		fmt.Println("\n========== ACCOUNT UPDATE RECEIVED ==========")
		fmt.Printf("Timestamp: %s\n", time.Now().Format("2006-01-02 15:04:05.000"))

		fmt.Printf("Action: %d\n", update.Action)

		fmt.Println("\n--- Account Details ---")
		fmt.Printf("ID:        %d\n", update.Data.ID)
		fmt.Printf("Name:      %s\n", update.Data.Name)
		fmt.Printf("Balance:   $%.2f\n", update.Data.Balance)
		fmt.Printf("Can Trade: %v\n", update.Data.CanTrade)
		fmt.Printf("Visible:   %v\n", update.Data.IsVisible)
		fmt.Printf("Simulated: %v\n", update.Data.Simulated)

		fmt.Println("\n--- Full JSON Structure ---")
		jsonBytes, err := json.MarshalIndent(update, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling data: %v\n", err)
		} else {
			fmt.Printf("%s\n", string(jsonBytes))
		}

		fmt.Println("=============================================\n")
	})

	userDataWS.SetConnectionHandler(func(state services.ConnectionState) {
		fmt.Printf("[Connection State Changed] %s at %s\n",
			state, time.Now().Format("15:04:05.000"))
	})

	fmt.Println("\nConnecting to user data websocket...")
	if err := userDataWS.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Subscribing to account updates...")
	if err := userDataWS.SubscribeAccounts(); err != nil {
		log.Fatalf("Failed to subscribe to accounts: %v", err)
	}
	fmt.Println("Successfully subscribed to account updates!")
	fmt.Println("\nListening for account updates... (Press Ctrl+C to exit)")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	fmt.Println("\nShutting down...")

	if err := userDataWS.UnsubscribeAccounts(); err != nil {
		fmt.Printf("Error unsubscribing from accounts: %v\n", err)
	}

	if err := userDataWS.Disconnect(); err != nil {
		fmt.Printf("Error disconnecting: %v\n", err)
	}

	fmt.Println("Shutdown complete.")
}
