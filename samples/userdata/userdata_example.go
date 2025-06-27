package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/tradingiq/topstepx-client/samples"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/tradingiq/topstepx-client"
	"github.com/tradingiq/topstepx-client/models"
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
		if loginResp.ErrorMessage != nil {
			log.Fatalf("Login failed: %s", *loginResp.ErrorMessage)
		}
		log.Fatalf("Login failed with error code: %v", loginResp.ErrorCode)
	}

	fmt.Println("Login successful!")

	fmt.Println("\nRetrieving active accounts...")
	accountsResp, err := client.Account.SearchAccounts(ctx, &models.SearchAccountRequest{
		OnlyActiveAccounts: true,
	})
	if err != nil {
		log.Fatalf("Failed to retrieve accounts: %v", err)
	}

	if !accountsResp.Success {
		if accountsResp.ErrorMessage != nil {
			log.Fatalf("Failed to retrieve accounts: %s", *accountsResp.ErrorMessage)
		}
		log.Fatalf("Failed to retrieve accounts with error code: %v", accountsResp.ErrorCode)
	}

	if len(accountsResp.Accounts) == 0 {
		log.Fatal("No active accounts found")
	}

	fmt.Println("\n=== Active Accounts ===")
	for i, account := range accountsResp.Accounts {
		fmt.Printf("%d. Name: %s\n", i+1, account.Name)
		fmt.Printf("   ID: %d\n", account.ID)
		fmt.Printf("   Balance: $%.2f\n", account.Balance)
		fmt.Println()
	}

	var selectedAccount models.TradingAccountModel

	if len(accountsResp.Accounts) == 1 {
		selectedAccount = accountsResp.Accounts[0]
		fmt.Printf("Using the only available account: %s (ID: %d)\n", selectedAccount.Name, selectedAccount.ID)
	} else {

		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Printf("Select an account (1-%d): ", len(accountsResp.Accounts))
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			choice, err := strconv.Atoi(input)
			if err == nil && choice >= 1 && choice <= len(accountsResp.Accounts) {
				selectedAccount = accountsResp.Accounts[choice-1]
				break
			}
			fmt.Println("Invalid selection, please try again.")
		}
	}

	fmt.Printf("\nSelected account: %s (ID: %d)\n", selectedAccount.Name, selectedAccount.ID)

	fmt.Println("\nSetting up WebSocket handlers...")
	fmt.Println("Note: Order handler now supports both structured (type-safe) and raw (interface{}) data access")

	client.UserData.SetConnectionHandler(func(state services.ConnectionState) {
		switch state {
		case services.StateDisconnected:
			fmt.Println("\nWebSocket disconnected")
		case services.StateConnecting:
			fmt.Println("\nConnecting to WebSocket...")
		case services.StateConnected:
			fmt.Println("\nWebSocket connected!")
		case services.StateReconnecting:
			fmt.Println("\nReconnecting to WebSocket...")
		}
	})

	client.UserData.SetAccountHandler(func(data interface{}) {
		fmt.Printf("\n[ACCOUNT UPDATE] %+v\n", data)
	})

	client.UserData.SetPositionHandler(func(data interface{}) {
		fmt.Printf("\n[POSITION UPDATE] %+v\n", data)
	})

	client.UserData.SetTradeHandler(func(data interface{}) {
		fmt.Printf("\n[TRADE UPDATE] %+v\n", data)
	})

	fmt.Println("\nConnecting to WebSocket...")
	if err := client.UserData.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}

	fmt.Printf("\nSubscribing to all events for account %d...\n", selectedAccount.ID)
	if err := client.UserData.SubscribeAll(int(selectedAccount.ID)); err != nil {
		log.Fatalf("Failed to subscribe to events: %v", err)
	}
	fmt.Println("Successfully subscribed to all events!")

	fmt.Println("\n=== WebSocket Subscriptions Active ===")
	fmt.Println("Listening for real-time updates on:")
	fmt.Println("- Account updates")
	fmt.Println("- Order updates")
	fmt.Println("- Position updates")
	fmt.Println("- Trade updates")
	fmt.Println("\nPress Ctrl+C to exit...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n\nShutting down...")

	fmt.Println("Unsubscribing from events...")
	if err := client.UserData.UnsubscribeAll(); err != nil {
		log.Printf("Error unsubscribing: %v", err)
	}

	fmt.Println("Disconnecting WebSocket...")
	if err := client.UserData.Disconnect(); err != nil {
		log.Printf("Error disconnecting WebSocket: %v", err)
	}

	fmt.Println("Logging out...")
	if _, err := client.Auth.Logout(ctx); err != nil {
		log.Printf("Error logging out: %v", err)
	}

	fmt.Println("Shutdown complete!")
}
