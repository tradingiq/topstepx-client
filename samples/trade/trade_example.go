package main

import (
	"context"
	"fmt"
	"github.com/tradingiq/topstepx-client/samples"
	"log"
	"time"

	"github.com/tradingiq/topstepx-client"
	"github.com/tradingiq/topstepx-client/client"
	"github.com/tradingiq/topstepx-client/models"
)

func main() {

	tsxClient := topstepx.NewClient(
		client.WithBaseURL("https://api.topstepx.com"),
	)

	ctx := context.Background()
	// Login
	resp, err := tsxClient.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
		UserName: samples.Config.Username,
		APIKey:   samples.Config.ApiKey,
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

	response, err := tsxClient.Account.SearchAccounts(ctx, &models.SearchAccountRequest{
		OnlyActiveAccounts: true,
	})
	if err != nil {
		log.Fatalf("Failed to search accoutns: %v", err)
	}

	for _, account := range response.Accounts {
		fmt.Println("\n2. Search Trades - Last 24 Hours")
		now := time.Now()
		end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		trades24hResp, err := tsxClient.Trade.SearchHalfTurnTrades(ctx, &models.SearchTradeRequest{
			AccountID:      int32(account.ID),
			StartTimestamp: &start,
			EndTimestamp:   &end,
		})
		if err != nil {
			log.Printf("Search 24h trades failed: %v", err)
		} else if trades24hResp.Success {
			fmt.Printf("Found %d trades in the last 24 hours\n", len(trades24hResp.Trades))

			if len(trades24hResp.Trades) > 0 {
				calculateTradeStats(trades24hResp.Trades, "24 Hour")
			}
		}
	}
}

func calculateTradeStats(trades []models.HalfTradeModel, period string) {
	var totalPL float64
	var totalFees float64
	var winningTrades, losingTrades int
	var totalVolume int32

	for _, trade := range trades {
		if !trade.Voided {
			totalVolume += trade.Size
			totalFees += trade.Fees

			if trade.ProfitAndLoss != nil {
				totalPL += *trade.ProfitAndLoss
				if *trade.ProfitAndLoss > 0 {
					winningTrades++
				} else if *trade.ProfitAndLoss < 0 {
					losingTrades++
				}
			}
		}
	}

	fmt.Printf("\n%s Trading Statistics:\n", period)
	fmt.Printf("  Total Trades: %d\n", len(trades))
	fmt.Printf("  Total Volume: %d contracts\n", totalVolume)
	fmt.Printf("  Winning Trades: %d\n", winningTrades)
	fmt.Printf("  Losing Trades: %d\n", losingTrades)
	if winningTrades+losingTrades > 0 {
		winRate := float64(winningTrades) / float64(winningTrades+losingTrades) * 100
		fmt.Printf("  Win Rate: %.1f%%\n", winRate)
	}
	fmt.Printf("  Total P&L: %.2f\n", totalPL)
	fmt.Printf("  Total Fees: %.2f\n", totalFees)
	fmt.Printf("  Net P&L: %.2f\n", totalPL-totalFees)
}
