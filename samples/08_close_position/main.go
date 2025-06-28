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
	account := accountsResp.Accounts[0]

	positionsResp, err := client.Position.SearchOpenPositions(ctx, &models.SearchPositionRequest{
		AccountID: int32(account.ID),
	})

	if err != nil {
		log.Fatal("Failed to get positions:", err)
	}

	if len(positionsResp.Positions) == 0 {
		fmt.Println("No open positions to close")
		return
	}

	// Close the first position found
	position := positionsResp.Positions[0]
	fmt.Printf("Closing position: %s (Size: %d)\n", position.ContractID, position.Size)

	closeResp, err := client.Position.CloseContractPosition(ctx, &models.CloseContractPositionRequest{
		AccountID:  int32(account.ID),
		ContractID: position.ContractID,
	})

	if err != nil {
		log.Fatal("Failed to close position:", err)
	}

	if closeResp.Success {
		fmt.Printf("✅ Position closed successfully\n")
	} else {
		fmt.Printf("❌ Close failed: %s\n", closeResp.ErrorMessage)
	}
}
