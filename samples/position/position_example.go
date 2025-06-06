package main

import (
	"context"
	"fmt"
	"log"
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
	err := tsxClient.LoginAndConnect(ctx, username, apiKey)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	account, err := tsxClient.GetFirstActiveAccount(ctx)
	if err != nil {
		log.Fatalf("GetFistActiveAccount failed: %v", err)
	}
	accountID := account.ID

	fmt.Printf("Using account ID: %d\n", accountID)

	fmt.Println("\n=== Position Service Examples ===")

	fmt.Println("\n1. Search Current Positions")
	positionsResp, err := tsxClient.Position.SearchOpenPositions(ctx, &models.SearchPositionRequest{
		AccountID: int32(accountID),
	})
	if err != nil {
		log.Printf("Search positions failed: %v", err)
	} else if positionsResp.Success {
		fmt.Printf("Found %d open positions\n", len(positionsResp.Positions))

		if len(positionsResp.Positions) > 0 {
			fmt.Println("\nPosition Details:")
			for _, position := range positionsResp.Positions {
				fmt.Printf("\n  Position ID: %d\n", position.ID)
				fmt.Printf("  Contract: %s\n", position.ContractID)
				fmt.Printf("  Type: %s\n", getPositionTypeString(position.Type))
				fmt.Printf("  Size: %d\n", position.Size)
				fmt.Printf("  Average Price: %.2f\n", position.AveragePrice)
				fmt.Printf("  Created: %s\n", position.CreationTimestamp.Format("2006-01-02 15:04:05"))

				currentPL := calculatePL(position)
				fmt.Printf("  Estimated P&L: %.2f (example calculation)\n", currentPL)
			}
		} else {
			fmt.Println("No open positions found.")
		}
	}

	fmt.Println("\n2. Position Management Examples")
	fmt.Println("NOTE: These examples show API usage. Actual execution requires open positions.")

	searchText := "ES"
	contractsResp, err := tsxClient.Contract.SearchContracts(ctx, &models.SearchContractRequest{
		SearchText: &searchText,
		Live:       false,
	})
	if err != nil || !contractsResp.Success || len(contractsResp.Contracts) == 0 {
		log.Fatal("Failed to find contracts")
	}

	contractID := contractsResp.Contracts[0].ID
	fmt.Printf("\nUsing contract: %s (%s) for examples\n", contractID, contractsResp.Contracts[0].Name)

	fmt.Println("\n2.1 Close Position Example")
	closeResp, err := tsxClient.Position.CloseContractPosition(ctx, &models.CloseContractPositionRequest{
		AccountID:  int32(accountID),
		ContractID: contractID,
	})
	if err != nil {
		fmt.Printf("Close position error (expected if no position): %v\n", err)
	} else if closeResp.Success {
		fmt.Println("Position closed successfully!")
	} else {
		fmt.Printf("Close position failed with error code: %v\n", closeResp.ErrorCode)
		if closeResp.ErrorMessage != nil {
			fmt.Printf("Error message: %s\n", *closeResp.ErrorMessage)
		}
	}

	fmt.Println("\n2.2 Partial Close Position Example")
	partialSize := int32(5)
	partialCloseResp, err := tsxClient.Position.PartialCloseContractPosition(ctx, &models.PartialCloseContractPositionRequest{
		AccountID:  int32(accountID),
		ContractID: contractID,
		Size:       partialSize,
	})
	if err != nil {
		fmt.Printf("Partial close position error (expected if no position): %v\n", err)
	} else if partialCloseResp.Success {
		fmt.Printf("Partially closed %d contracts successfully!\n", partialSize)
	} else {
		fmt.Printf("Partial close failed with error code: %v\n", partialCloseResp.ErrorCode)
		if partialCloseResp.ErrorMessage != nil {
			fmt.Printf("Error message: %s\n", *partialCloseResp.ErrorMessage)
		}
	}
	fmt.Println("\n4. Position Summary")
	finalPositionsResp, err := tsxClient.Position.SearchOpenPositions(ctx, &models.SearchPositionRequest{
		AccountID: int32(accountID),
	})
	if err != nil {
		log.Printf("Search final positions failed: %v", err)
	} else if finalPositionsResp.Success {
		fmt.Printf("\nFinal position count: %d\n", len(finalPositionsResp.Positions))

		if len(finalPositionsResp.Positions) > 0 {
			var totalLongSize, totalShortSize int32
			var totalLongValue, totalShortValue float64

			for _, position := range finalPositionsResp.Positions {
				if position.Type == models.PositionTypeLong {
					totalLongSize += position.Size
					totalLongValue += float64(position.Size) * position.AveragePrice
				} else if position.Type == models.PositionTypeShort {
					totalShortSize += position.Size
					totalShortValue += float64(position.Size) * position.AveragePrice
				}
			}

			fmt.Printf("\nPosition Summary:\n")
			fmt.Printf("  Total Long: %d contracts (Value: %.2f)\n", totalLongSize, totalLongValue)
			fmt.Printf("  Total Short: %d contracts (Value: %.2f)\n", totalShortSize, totalShortValue)
			fmt.Printf("  Net Position: %d contracts\n", totalLongSize-totalShortSize)
		}
	}
}

func getPositionTypeString(posType models.PositionType) string {
	switch posType {
	case models.PositionTypeLong:
		return "LONG"
	case models.PositionTypeShort:
		return "SHORT"
	default:
		return "UNDEFINED"
	}
}

func calculatePL(position models.PositionModel) float64 {
	currentPrice := position.AveragePrice * 1.001

	if position.Type == models.PositionTypeLong {
		return float64(position.Size) * (currentPrice - position.AveragePrice)
	} else if position.Type == models.PositionTypeShort {
		return float64(position.Size) * (position.AveragePrice - currentPrice)
	}

	return 0.0
}
