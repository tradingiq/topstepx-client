package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tradingiq/projectx-client"
	"github.com/tradingiq/projectx-client/models"
	"github.com/tradingiq/projectx-client/samples"
)

// Demonstrates: Placing a market order
func main() {
	client := projectx.NewClient()
	ctx := context.Background()

	// Authenticate
	loginResp, _ := client.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
		UserName: samples.Config.Username,
		APIKey:   samples.Config.ApiKey,
	})
	if !loginResp.Success {
		log.Fatal("Login failed")
	}

	// Get first active account
	accountsResp, _ := client.Account.SearchAccounts(ctx, &models.SearchAccountRequest{
		OnlyActiveAccounts: true,
	})
	if len(accountsResp.Accounts) == 0 {
		log.Fatal("No active accounts")
	}
	account := accountsResp.Accounts[0]

	// Find MES contract
	searchText := "MES"
	contractsResp, _ := client.Contract.SearchContracts(ctx, &models.SearchContractRequest{
		SearchText: &searchText,
		Live:       false,
	})
	if len(contractsResp.Contracts) == 0 {
		log.Fatal("MES contract not found")
	}
	contract := contractsResp.Contracts[0]

	// Place market buy order
	orderResp, err := client.Order.PlaceOrder(ctx, &models.PlaceOrderRequest{
		AccountID:  int32(account.ID),
		ContractID: contract.ID,
		Type:       models.OrderTypeMarket,
		Side:       models.OrderSideBid,
		Size:       1,
	})

	if err != nil {
		log.Fatal("Failed to place order:", err)
	}

	if orderResp.Success && orderResp.OrderID != nil {
		fmt.Printf("✅ Market order placed!\n")
		fmt.Printf("Order ID: %d\n", *orderResp.OrderID)
	} else {
		fmt.Printf("❌ Order failed: %s\n", *orderResp.ErrorMessage)
	}
}
