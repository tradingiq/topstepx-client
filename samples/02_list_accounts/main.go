package main

import (
	"context"
	"fmt"
	"log"

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

	resp, err := client.Account.SearchAccounts(ctx, &models.SearchAccountRequest{
		OnlyActiveAccounts: true,
	})
	if err != nil {
		log.Fatal("Failed to get accounts:", err)
	}

	fmt.Printf("Found %d active accounts:\n\n", len(resp.Accounts))
	for _, account := range resp.Accounts {
		fmt.Printf("Account: %s\n", account.Name)
		fmt.Printf("  ID: %d\n", account.ID)
		fmt.Printf("  Balance: $%.2f\n", account.Balance)
		fmt.Printf("  Can Trade: %v\n", account.CanTrade)
		fmt.Println()
	}
}
