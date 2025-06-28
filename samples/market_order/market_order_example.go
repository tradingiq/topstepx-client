package main

import (
	"context"
	"fmt"
	"github.com/tradingiq/topstepx-client"
	"github.com/tradingiq/topstepx-client/models"
	"github.com/tradingiq/topstepx-client/samples"
	"log"
	"time"
)

func main() {
	username := samples.Config.Username
	apiKey := samples.Config.ApiKey

	if username == "" || apiKey == "" {
		log.Fatal("Please set TOPSTEPX_USERNAME and TOPSTEPX_API_KEY environment variables")
	}

	ctx := context.Background()

	client := topstepx.NewClient()

	fmt.Println("Logging in...")
	// Login
	resp, err := client.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
		UserName: username,
		APIKey:   apiKey,
	})
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}
	if !resp.Success {
		if resp.ErrorMessage != nil {
			log.Fatalf("Login failed: %s", *resp.ErrorMessage)
		}
		log.Fatalf("Login failed with error code: %v", resp.ErrorCode)
	}

	// Get active accounts
	accountsResp, err := client.Account.SearchAccounts(ctx, &models.SearchAccountRequest{
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

	searchText := "MES"
	contractsResp, err := client.Contract.SearchContracts(ctx, &models.SearchContractRequest{
		SearchText: &searchText,
		Live:       false,
	})
	if err != nil {
		log.Fatalf("Failed to search contracts: %v", err)
	}

	if len(contractsResp.Contracts) == 0 {
		log.Fatal("No ES contracts found")
	}

	contractID := contractsResp.Contracts[0].ID
	fmt.Printf("Using contract: %s (ID: %d)\n", contractsResp.Contracts[0].Name, contractID)

	fmt.Println("\n=== Placing Market Buy Order ===")
	placeOrderResp, err := client.Order.PlaceOrder(ctx, &models.PlaceOrderRequest{
		AccountID:  int32(accountID),
		ContractID: contractID,
		Type:       models.OrderTypeMarket,
		Side:       models.OrderSideBid,
		Size:       1,
	})

	if err != nil {
		log.Fatalf("Failed to place market order: %v", err)
	}

	if placeOrderResp.ErrorCode != 0 {
		log.Fatalf("Order placement failed with error code %d: %s", placeOrderResp.ErrorCode, *placeOrderResp.ErrorMessage)
	}

	orderID := placeOrderResp.OrderID
	fmt.Printf("Market order placed successfully! Order ID: %d\n", orderID)

	fmt.Println("\n=== Waiting 5 seconds before closing position ===")
	time.Sleep(5 * time.Second)

	fmt.Println("\n=== Checking for open positions ===")
	positionsResp, err := client.Position.SearchOpenPositions(ctx, &models.SearchPositionRequest{
		AccountID: int32(accountID),
	})
	if err != nil {
		log.Fatalf("Failed to search open positions: %v", err)
	}

	var positionToClose models.PositionModel
	for _, position := range positionsResp.Positions {
		if position.ContractID == contractID {
			positionToClose = position
			break
		}
	}

	fmt.Println("\n=== Closing Position ===")
	closeResp, err := client.Position.CloseContractPosition(ctx, &models.CloseContractPositionRequest{
		AccountID:  int32(accountID),
		ContractID: contractID,
	})

	if err != nil {
		log.Fatalf("Failed to close position: %v", err)
	}

	if closeResp.ErrorCode != 0 {
		log.Fatalf("Position close failed with error code %d: %s", closeResp.ErrorCode, closeResp.ErrorMessage)
	}

	fmt.Printf("Position closed successfully! Close position ID: %d\n", positionToClose.ID)
	fmt.Println("\n=== Market Order Example Complete ===")
}
