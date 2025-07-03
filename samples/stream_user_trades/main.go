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

	userDataWS.SetTradeHandler(func(update *models.TradeUpdateData) {
		fmt.Println("\n========== TRADE UPDATE RECEIVED ==========")
		fmt.Printf("Timestamp: %s\n", time.Now().Format("2006-01-02 15:04:05.000"))

		fmt.Printf("Action: %d\n", update.Action)

		fmt.Println("\n--- Trade Details ---")
		fmt.Printf("Trade ID:     %d\n", update.Data.ID)
		fmt.Printf("Account ID:   %d\n", update.Data.AccountID)
		fmt.Printf("Contract ID:  %s\n", update.Data.ContractID)
		fmt.Printf("Order ID:     %d\n", update.Data.OrderID)
		fmt.Printf("Side:         %s\n", update.Data.Side.String())
		fmt.Printf("Size:         %d\n", update.Data.Size)
		fmt.Printf("Price:        %.2f\n", update.Data.Price)
		fmt.Printf("Fees:         $%.2f\n", update.Data.Fees)
		fmt.Printf("Voided:       %v\n", update.Data.Voided)
		fmt.Printf("Executed At:  %s\n", update.Data.CreationTimestamp.Format("2006-01-02 15:04:05"))

		tradeValue := float64(update.Data.Size) * update.Data.Price
		fmt.Printf("Trade Value:  $%.2f\n", tradeValue)

		fmt.Println("\n--- Full JSON Structure ---")
		jsonBytes, err := json.MarshalIndent(update, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling data: %v\n", err)
		} else {
			fmt.Printf("%s\n", string(jsonBytes))
		}

		fmt.Println("==========================================\n")
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

	fmt.Printf("Subscribing to trade updates for account %d...\n", account.ID)
	if err := userDataWS.SubscribeTrades(int(account.ID)); err != nil {
		log.Fatalf("Failed to subscribe to trades: %v", err)
	}
	fmt.Println("Successfully subscribed to trade updates!")
	fmt.Println("\nListening for trade updates... (Press Ctrl+C to exit)")
	fmt.Println("NOTE: Trade updates will appear when trades are executed on this account")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	fmt.Println("\nShutting down...")

	if err := userDataWS.UnsubscribeTrades(int(account.ID)); err != nil {
		fmt.Printf("Error unsubscribing from trades: %v\n", err)
	}

	if err := userDataWS.Disconnect(); err != nil {
		fmt.Printf("Error disconnecting: %v\n", err)
	}

	fmt.Println("Shutdown complete.")
}
