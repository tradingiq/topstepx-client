package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tradingiq/topstepx-client"
	"github.com/tradingiq/topstepx-client/models"
	"github.com/tradingiq/topstepx-client/samples"
)

func main() {
	client := topstepx.NewClient()
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
	if len(accountsResp.Accounts) == 0 {
		log.Fatal("No active accounts")
	}
	account := accountsResp.Accounts[0]

	searchText := "ES"
	contractsResp, _ := client.Contract.SearchContracts(ctx, &models.SearchContractRequest{
		SearchText: &searchText,
		Live:       false,
	})
	if len(contractsResp.Contracts) == 0 {
		log.Fatal("ES contract not found")
	}
	contract := contractsResp.Contracts[0]

	limitPrice := 4500.0
	orderResp, err := client.Order.PlaceOrder(ctx, &models.PlaceOrderRequest{
		AccountID:  int32(account.ID),
		ContractID: contract.ID,
		Type:       models.OrderTypeLimit,
		Side:       models.OrderSideBid,
		Size:       1,
		LimitPrice: &limitPrice,
	})

	if err != nil {
		log.Fatal("Failed to place order:", err)
	}

	if orderResp.Success && orderResp.OrderID != nil {
		fmt.Printf("✅ Limit order placed!\n")
		fmt.Printf("Order ID: %d\n", *orderResp.OrderID)
		fmt.Printf("BUY 1 contract @ $%.2f\n", limitPrice)
	} else {
		fmt.Printf("❌ Order failed: %s\n", *orderResp.ErrorMessage)
	}
}
