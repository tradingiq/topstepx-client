package main

import (
	"context"
	"fmt"
	"log"
	"time"

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
	contract := contractsResp.Contracts[0]

	limitPrice := 4500.0
	orderResp, _ := client.Order.PlaceOrder(ctx, &models.PlaceOrderRequest{
		AccountID:  int32(account.ID),
		ContractID: contract.ID,
		Type:       models.OrderTypeLimit,
		Side:       models.OrderSideBid,
		Size:       1,
		LimitPrice: &limitPrice,
	})

	if !orderResp.Success {
		log.Fatal("Failed to place initial order")
	}

	fmt.Printf("✅ Order placed: %d @ $%.2f\n", *orderResp.OrderID, limitPrice)

	time.Sleep(2 * time.Second)

	// Modify the order price
	newPrice := 4505.0
	modifyResp, err := client.Order.ModifyOrder(ctx, &models.ModifyOrderRequest{
		AccountID:  int32(account.ID),
		OrderID:    *orderResp.OrderID,
		LimitPrice: &newPrice,
	})

	if err != nil {
		log.Fatal("Failed to modify order:", err)
	}

	if modifyResp.Success {
		fmt.Printf("✅ Order modified: new price $%.2f\n", newPrice)
	} else {
		fmt.Printf("❌ Modify failed: %s\n", *modifyResp.ErrorMessage)
	}
}
