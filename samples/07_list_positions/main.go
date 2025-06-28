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

	if !positionsResp.Success {
		log.Fatal("Position search failed:", positionsResp.ErrorMessage)
	}

	fmt.Printf("Found %d open positions:\n\n", len(positionsResp.Positions))

	if len(positionsResp.Positions) == 0 {
		fmt.Println("No open positions")
		return
	}

	for _, position := range positionsResp.Positions {
		direction := "LONG"
		if position.Size < 0 {
			direction = "SHORT"
		}

		fmt.Printf("Position: %s\n", position.ContractID)
		fmt.Printf("  Direction: %s\n", direction)
		fmt.Printf("  Size: %d\n", position.Size)
		fmt.Printf("  Average Price: $%.2f\n", position.AveragePrice)
	}
}
