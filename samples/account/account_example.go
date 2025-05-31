package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tradingiq/topstepx-client"
	"github.com/tradingiq/topstepx-client/client"
	"github.com/tradingiq/topstepx-client/models"
)

func main() {
	username := os.Getenv("TOPSTEPX_USERNAME")
	apiKey := os.Getenv("TOPSTEPX_API_KEY")

	if username == "" || apiKey == "" {
		log.Fatal("Please set TOPSTEPX_USERNAME and TOPSTEPX_API_KEY environment variables")
	}

	tsxClient := topstepx.NewClient(
		client.WithBaseURL("https://api.topstepx.com"),
	)

	ctx := context.Background()

	loginResp, err := tsxClient.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
		UserName: username,
		APIKey:   apiKey,
	})
	if err != nil || !loginResp.Success {
		log.Fatal("Login failed")
	}

	fmt.Println("=== Account Service Examples ===")

	fmt.Println("\n1. Search All Accounts")
	allAccountsResp, err := tsxClient.Account.SearchAccounts(ctx, &models.SearchAccountRequest{
		OnlyActiveAccounts: false,
	})
	if err != nil {
		log.Fatalf("Search accounts failed: %v", err)
	}

	if allAccountsResp.Success {
		fmt.Printf("Found %d accounts:\n", len(allAccountsResp.Accounts))
		for _, account := range allAccountsResp.Accounts {
			fmt.Printf("  - ID: %d, Name: %s, Balance: %.2f, CanTrade: %v, IsVisible: %v\n",
				account.ID, account.Name, account.Balance, account.CanTrade, account.IsVisible)
		}
	}

	fmt.Println("\n2. Search Active Accounts Only")
	activeAccountsResp, err := tsxClient.Account.SearchAccounts(ctx, &models.SearchAccountRequest{
		OnlyActiveAccounts: true,
	})
	if err != nil {
		log.Fatalf("Search active accounts failed: %v", err)
	}

	if activeAccountsResp.Success {
		fmt.Printf("Found %d active accounts:\n", len(activeAccountsResp.Accounts))
		for _, account := range activeAccountsResp.Accounts {
			fmt.Printf("  - ID: %d, Name: %s, Balance: %.2f\n",
				account.ID, account.Name, account.Balance)
		}
	}
}
