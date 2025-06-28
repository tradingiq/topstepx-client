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

// Demonstrates: Getting historical price data
func main() {
	client := topstepx.NewClient()
	ctx := context.Background()

	// Authenticate
	loginResp, _ := client.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
		UserName: samples.Config.Username,
		APIKey:   samples.Config.ApiKey,
	})
	if !loginResp.Success {
		log.Fatal("Login failed")
	}

	// Find ES contract
	searchText := "MES"
	contractsResp, _ := client.Contract.SearchContracts(ctx, &models.SearchContractRequest{
		SearchText: &searchText,
		Live:       false,
	})
	if len(contractsResp.Contracts) == 0 {
		log.Fatal("Contract not found")
	}
	contractID := contractsResp.Contracts[0].ID

	// Get 1-minute bars for the last 24 hours
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)

	fmt.Printf("Getting 1-minute bars for %s from %s to %s\n\n",
		contractID, startTime.Format("2006-01-02 15:04"), endTime.Format("2006-01-02 15:04"))

	resp, err := client.History.GetBars(ctx, &models.RetrieveBarRequest{
		ContractID: contractID,
		Unit:       models.AggregateBarUnitMinute,
		UnitNumber: 1,
		StartTime:  startTime,
		EndTime:    endTime,
	})

	if err != nil {
		log.Fatal("Failed to get historical data:", err)
	}

	if !resp.Success {
		log.Fatal("History request failed:", resp.ErrorMessage)
	}

	fmt.Printf("Retrieved %d bars:\n\n", len(resp.Bars))

	// Show last 5 bars
	start := len(resp.Bars) - 5
	if start < 0 {
		start = 0
	}

	for i := start; i < len(resp.Bars); i++ {
		bar := resp.Bars[i]
		fmt.Printf("%s: Open=%.2f High=%.2f Low=%.2f Close=%.2f Volume=%.0f\n",
			bar.T.Format("15:04"), bar.Open, bar.High, bar.Low, bar.Close, bar.Volume)
	}
}
