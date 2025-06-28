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
	account := accountsResp.Accounts[0]

	searchText := "ES"
	contractsResp, _ := client.Contract.SearchContracts(ctx, &models.SearchContractRequest{
		SearchText: &searchText,
		Live:       false,
	})
	contract := contractsResp.Contracts[0]

	limitPrice := 4000.0
	orderResp, _ := client.Order.PlaceOrder(ctx, &models.PlaceOrderRequest{
		AccountID:  int32(account.ID),
		ContractID: contract.ID,
		Type:       models.OrderTypeLimit,
		Side:       models.OrderSideBid,
		Size:       1,
		LimitPrice: &limitPrice,
	})

	if !orderResp.Success {
		log.Fatal("Failed to place order")
	}

	fmt.Printf("✅ Order placed: %d @ $%.2f\n", *orderResp.OrderID, limitPrice)

	time.Sleep(2 * time.Second)

	cancelResp, err := client.Order.CancelOrder(ctx, &models.CancelOrderRequest{
		AccountID: int32(account.ID),
		OrderID:   *orderResp.OrderID,
	})

	if err != nil {
		log.Fatal("Failed to cancel order:", err)
	}

	if cancelResp.Success {
		fmt.Printf("✅ Order %d canceled successfully\n", *orderResp.OrderID)
	} else {
		fmt.Printf("❌ Cancel failed: %s\n", *cancelResp.ErrorMessage)
	}
}
