package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tradingiq/topstepx-client"
	"github.com/tradingiq/topstepx-client/client"
	"github.com/tradingiq/topstepx-client/models"
	"github.com/tradingiq/topstepx-client/samples"
)

func main() {
	username := samples.Config.Username
	apiKey := samples.Config.ApiKey

	if username == "" || apiKey == "" {
		log.Fatal("Please set TOPSTEPX_USERNAME and TOPSTEPX_API_KEY environment variables")
	}

	tsxClient := topstepx.NewClient(
		client.WithBaseURL("https://api.topstepx.com"),
	)

	ctx := context.Background()

	// Login
	resp, err := tsxClient.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
		UserName: username,
		APIKey:   apiKey,
	})
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	if !resp.Success {
		if resp.ErrorMessage != nil {
			log.Fatalf("Login failed: %s", *resp.ErrorMessage)
		}
		log.Fatalf("Login failed with error code: %v", resp.ErrorCode)
	}

	// Get active accounts
	accountsResp, err := tsxClient.Account.SearchAccounts(ctx, &models.SearchAccountRequest{
		OnlyActiveAccounts: true,
	})
	if err != nil {
		log.Fatalf("Failed to search accounts: %v", err)
	}
	if !accountsResp.Success {
		if accountsResp.ErrorMessage != nil {
			log.Fatalf("Failed to search accounts: %s", *accountsResp.ErrorMessage)
		}
		log.Fatalf("Failed to search accounts with error code: %v", accountsResp.ErrorCode)
	}
	if len(accountsResp.Accounts) == 0 {
		log.Fatal("No active accounts found")
	}
	accountID := accountsResp.Accounts[0].ID

	fmt.Printf("Using account ID: %d\n", accountID)

	searchText := "ES"
	contractsResp, err := tsxClient.Contract.SearchContracts(ctx, &models.SearchContractRequest{
		SearchText: &searchText,
		Live:       false,
	})
	if err != nil || !contractsResp.Success || len(contractsResp.Contracts) == 0 {
		log.Fatal("Failed to find contracts")
	}

	contractID := contractsResp.Contracts[0].ID
	fmt.Printf("Using contract: %s (%s)\n", contractID, contractsResp.Contracts[0].Name)

	fmt.Println("\n=== Limit Order Example ===")

	fmt.Println("\n1. Place Buy Limit Order")
	limitPrice := 4500.0
	customTag := "buy-limit-example"

	buyLimitResp, err := tsxClient.Order.PlaceOrder(ctx, &models.PlaceOrderRequest{
		AccountID:  int32(accountID),
		ContractID: contractID,
		Type:       models.OrderTypeLimit,
		Side:       models.OrderSideBid,
		Size:       1,
		LimitPrice: &limitPrice,
		CustomTag:  &customTag,
	})

	if err != nil {
		fmt.Printf("Place buy limit order error: %v\n", err)
	} else if buyLimitResp.Success && buyLimitResp.OrderID != nil {
		fmt.Printf("Buy limit order placed successfully! Order ID: %d\n", *buyLimitResp.OrderID)
		fmt.Printf("Order details: BUY 1 contract @ $%.2f\n", limitPrice)

		time.Sleep(2 * time.Second)
		cancelResp, _ := tsxClient.Order.CancelOrder(ctx, &models.CancelOrderRequest{
			AccountID: int32(accountID),
			OrderID:   *buyLimitResp.OrderID,
		})
		if cancelResp.Success {
			fmt.Println("Buy limit order cancelled")
		}
	} else {
		fmt.Printf("Buy limit order failed with error code: %v\n", buyLimitResp.ErrorCode)
	}

	fmt.Println("\n2. Place Sell Limit Order")
	sellLimitPrice := 8000.0
	sellCustomTag := "sell-limit-example"

	sellLimitResp, err := tsxClient.Order.PlaceOrder(ctx, &models.PlaceOrderRequest{
		AccountID:  int32(accountID),
		ContractID: contractID,
		Type:       models.OrderTypeLimit,
		Side:       models.OrderSideAsk,
		Size:       1,
		LimitPrice: &sellLimitPrice,
		CustomTag:  &sellCustomTag,
	})

	if err != nil {
		fmt.Printf("Place sell limit order error: %v\n", err)
	} else if sellLimitResp.Success && sellLimitResp.OrderID != nil {
		fmt.Printf("Sell limit order placed successfully! Order ID: %d\n", *sellLimitResp.OrderID)
		fmt.Printf("Order details: SELL 1 contract @ $%.2f\n", sellLimitPrice)

		time.Sleep(2 * time.Second)
		cancelResp, _ := tsxClient.Order.CancelOrder(ctx, &models.CancelOrderRequest{
			AccountID: int32(accountID),
			OrderID:   *sellLimitResp.OrderID,
		})
		if cancelResp.Success {
			fmt.Println("Sell limit order cancelled")
		}
	} else {
		fmt.Printf("Sell limit order failed with error code: %v\n", sellLimitResp.ErrorCode)
	}

	fmt.Println("\nLimit order example completed!")
}
