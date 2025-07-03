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
	client := topstepx.NewClient()
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

	fmt.Println("\nFetching accounts...")
	accounts, err := client.GetActiveAccounts(ctx)
	if err != nil {
		log.Fatalf("Failed to get accounts: %v", err)
	}

	if len(accounts) == 0 {
		log.Fatal("No accounts found")
	}

	account := accounts[0]
	fmt.Printf("Using account: %s (ID: %d)\n", account.Name, account.ID)

	userDataWS := client.UserData

	userDataWS.SetPositionHandler(func(update *models.PositionUpdateData) {
		fmt.Println("\n========== POSITION UPDATE RECEIVED ==========")
		fmt.Printf("Timestamp: %s\n", time.Now().Format("2006-01-02 15:04:05.000"))
		
		fmt.Printf("Action: %d\n", update.Action)
		
		fmt.Println("\n--- Position Details ---")
		fmt.Printf("Position ID:  %d\n", update.Data.ID)
		fmt.Printf("Account ID:   %d\n", update.Data.AccountID)
		fmt.Printf("Contract ID:  %s\n", update.Data.ContractID)
		fmt.Printf("Type:         %s\n", update.Data.Type.String())
		fmt.Printf("Size:         %d\n", update.Data.Size)
		fmt.Printf("Avg Price:    %.2f\n", update.Data.AveragePrice)
		fmt.Printf("Created At:   %s\n", update.Data.CreationTimestamp.Format("2006-01-02 15:04:05"))
		
		positionValue := float64(update.Data.Size) * update.Data.AveragePrice
		fmt.Printf("Position Value: $%.2f\n", positionValue)
		
		fmt.Println("\n--- Full JSON Structure ---")
		jsonBytes, err := json.MarshalIndent(update, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling data: %v\n", err)
		} else {
			fmt.Printf("%s\n", string(jsonBytes))
		}
		
		fmt.Println("===============================================\n")
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

	fmt.Printf("Subscribing to position updates for account %d...\n", account.ID)
	if err := userDataWS.SubscribePositions(int(account.ID)); err != nil {
		log.Fatalf("Failed to subscribe to positions: %v", err)
	}
	fmt.Println("Successfully subscribed to position updates!")
	fmt.Println("\nListening for position updates... (Press Ctrl+C to exit)")
	fmt.Println("NOTE: Position updates will appear when positions change on this account")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	fmt.Println("\nShutting down...")

	if err := userDataWS.UnsubscribePositions(int(account.ID)); err != nil {
		fmt.Printf("Error unsubscribing from positions: %v\n", err)
	}

	if err := userDataWS.Disconnect(); err != nil {
		fmt.Printf("Error disconnecting: %v\n", err)
	}

	fmt.Println("Shutdown complete.")
}