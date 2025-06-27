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

	fmt.Printf("Subscribing to orders for account %d...\n", selectedAccount.ID)
	if err := client.UserData.SubscribeOrders(int(selectedAccount.ID)); err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	fmt.Println("\n✅ Order handler ready!")
	fmt.Println("The handler will display updates when orders are:")
	fmt.Println("  • Placed")
	fmt.Println("  • Modified")
	fmt.Println("  • Filled")
	fmt.Println("  • Cancelled")
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
	client.UserData.SetOrderHandler(func(data *models.OrderUpdateData) {
		fmt.Printf("\n=== STRUCTURED ORDER UPDATE ===\n")
		fmt.Printf("Action: %s\n", getActionName(data.Action))
		fmt.Printf("Order Details:\n")
		fmt.Printf("  ID: %d\n", data.Data.ID)
		fmt.Printf("  Account: %d\n", data.Data.AccountID)
		fmt.Printf("  Contract: %s\n", data.Data.ContractID)
		fmt.Printf("  Status: %s (value: %d)\n", data.Data.Status, data.Data.Status)
		fmt.Printf("  Type: %s (value: %d)\n", data.Data.Type, data.Data.Type)
		fmt.Printf("  Side: %s (value: %d)\n", data.Data.Side, data.Data.Side)
		fmt.Printf("  Size: %d\n", data.Data.Size)
		fmt.Printf("  Limit Price: $%.2f\n", data.Data.LimitPrice)
		fmt.Printf("  Fill Volume: %d\n", data.Data.FillVolume)
		fmt.Printf("  Created: %s\n", data.Data.CreationTimestamp.Format("2006-01-02 15:04:05.000 MST"))
		fmt.Printf("  Updated: %s\n", data.Data.UpdateTimestamp.Format("2006-01-02 15:04:05.000 MST"))
		fmt.Printf("================================\n")
	})
}

func getActionName(action int) string {
	switch action {
	case 1:
		return "UPDATE"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", action)
	}
}
