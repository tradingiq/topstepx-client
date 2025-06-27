package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tradingiq/topstepx-client"
	"github.com/tradingiq/topstepx-client/models"
	"github.com/tradingiq/topstepx-client/samples"
	"github.com/tradingiq/topstepx-client/services"
)

func main() {
	client := topstepx.NewClient()
	ctx := context.Background()

	fmt.Println("Logging in...")
	loginResp, err := client.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
		UserName: samples.Config.Username,
		APIKey:   samples.Config.ApiKey,
	})
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	if !loginResp.Success {
		log.Fatalf("Login failed: %v", loginResp.ErrorMessage)
	}
	fmt.Println("Login successful!")

	accountsResp, err := client.Account.SearchAccounts(ctx, &models.SearchAccountRequest{
		OnlyActiveAccounts: true,
	})
	if err != nil || !accountsResp.Success || len(accountsResp.Accounts) == 0 {
		log.Fatal("Failed to get active accounts")
	}

	selectedAccount := accountsResp.Accounts[0]
	fmt.Printf("Using account: %s (ID: %d)\n", selectedAccount.Name, selectedAccount.ID)

	handler(client)

	client.UserData.SetConnectionHandler(func(state services.ConnectionState) {
		fmt.Printf("Connection state: %v\n", state)
	})

	fmt.Println("\nConnecting to WebSocket...")
	if err := client.UserData.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	fmt.Printf("Subscribing to positions for account %d...\n", selectedAccount.ID)
	if err := client.UserData.SubscribePositions(int(selectedAccount.ID)); err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	fmt.Println("\nPosition handler ready!")
	fmt.Println("The handler will display updates when positions are:")
	fmt.Println("  • Opened")
	fmt.Println("  • Modified")
	fmt.Println("  • Closed")
	fmt.Println("  • Updated")
	fmt.Println("\nPress Ctrl+C to exit...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down...")
	client.UserData.UnsubscribeAll()
	client.UserData.Disconnect()
	client.Auth.Logout(ctx)
	fmt.Println("Shutdown complete!")
}

func handler(client *topstepx.Client) {
	client.UserData.SetPositionHandler(func(data *models.PositionUpdateData) {
		fmt.Printf("\n=== POSITION UPDATE ===\n")
		fmt.Printf("Action: %d\n", data.Action)
		fmt.Printf("Position ID: %d\n", data.Data.ID)
		fmt.Printf("Account ID: %d\n", data.Data.AccountID)
		fmt.Printf("Contract: %s\n", data.Data.ContractID)
		fmt.Printf("Size: %d\n", data.Data.Size)
		fmt.Printf("Average Price: %.2f\n", data.Data.AveragePrice)
		fmt.Printf("Type: %s (%d)\n", data.Data.Type, int(data.Data.Type))
		fmt.Printf("Created: %s\n", data.Data.CreationTimestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("=======================\n")
	})
}
