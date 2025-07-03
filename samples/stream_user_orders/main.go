package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tradingiq/projectx-client"
	"github.com/tradingiq/projectx-client/models"
	"github.com/tradingiq/projectx-client/samples"
)

func main() {
	client := projectx.NewClient()
	ctx := context.Background()

	loginResp, _ := client.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
		UserName: samples.Config.Username,
		APIKey:   samples.Config.ApiKey,
	})
	if !loginResp.Success {
		log.Fatal("Login failed")
	}

	accountsResp, _ := client.Account.SearchAccounts(ctx, &models.SearchAccountRequest{
		OnlyActiveAccounts: true,
	})
	account := accountsResp.Accounts[0]

	client.UserData.SetOrderHandler(func(data *models.OrderUpdateData) {
		action := getActionName(data.Action)
		order := data.Data

		fmt.Printf("\n🔔 ORDER %s\n", action)
		fmt.Printf("  ID: %d\n", order.ID)
		fmt.Printf("  Contract: %s\n", order.ContractID)
		fmt.Printf("  Status: %s\n", order.Status)
		fmt.Printf("  Type: %s\n", order.Type)
		fmt.Printf("  Side: %s\n", order.Side)
		fmt.Printf("  Size: %d\n", order.Size)
		if order.LimitPrice > 0 {
			fmt.Printf("  Limit Price: $%.2f\n", order.LimitPrice)
		}
		fmt.Printf("  Fill Volume: %d\n", order.FillVolume)
		fmt.Printf("  Time: %s\n", order.UpdateTimestamp.Format("15:04:05"))
	})

	fmt.Println("Connecting to user data WebSocket...")
	if err := client.UserData.Connect(ctx); err != nil {
		log.Fatal("Failed to connect:", err)
	}

	// Subscribe to order updates
	fmt.Printf("Subscribing to order updates for account %d...\n", account.ID)
	if err := client.UserData.SubscribeOrders(int(account.ID)); err != nil {
		log.Fatal("Failed to subscribe:", err)
	}

	fmt.Println("✅ Monitoring order updates...")
	fmt.Println("Place orders from another application to see updates here.")
	fmt.Println("Press Ctrl+Close to stop...")

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down...")
	client.UserData.UnsubscribeAll()
	client.UserData.Disconnect()
}

func getActionName(action int) string {
	switch action {
	case 1:
		return "CREATED"
	case 2:
		return "UPDATED"
	case 3:
		return "DELETED"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", action)
	}
}
