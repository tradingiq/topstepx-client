package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tradingiq/topstepx-client"
	"github.com/tradingiq/topstepx-client/client"
	"github.com/tradingiq/topstepx-client/models"
)

func main() {
	username := os.Getenv("TOPSTEPX_USERNAME")
	apiKey := os.Getenv("TOPSTEPX_API_KEY")

	if username == "" || apiKey == "" {
		log.Fatal("Please set TOPSTEPX_USERNAME and TOPSTEPX_API_KEY environment variables")
	}

	tsxClient := topstepx.NewClient(
		client.WithBaseURL("https://api.topstepx.com"),
	)

	ctx := context.Background()

	loginResp, err := tsxClient.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
		UserName: username,
		APIKey:   apiKey,
	})
	if err != nil || !loginResp.Success {
		log.Fatal("Login failed")
	}

	fmt.Println("=== History Service Examples ===")

	searchText := "ES"
	contractsResp, err := tsxClient.Contract.SearchContracts(ctx, &models.SearchContractRequest{
		SearchText: &searchText,
		Live:       true,
	})
	if err != nil || !contractsResp.Success || len(contractsResp.Contracts) == 0 {
		log.Fatal("Failed to find contracts")
	}

	contractID := contractsResp.Contracts[0].ID
	fmt.Printf("Using contract: %s (%s)\n", contractID, contractsResp.Contracts[0].Name)

	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)

	fmt.Println("\n1. Retrieve 1-Minute Bars")
	oneMinBarsResp, err := tsxClient.History.GetBars(ctx, &models.RetrieveBarRequest{
		ContractID:        contractID,
		Live:              true,
		StartTime:         startTime,
		EndTime:           endTime,
		Unit:              models.AggregateBarUnitMinute,
		UnitNumber:        1,
		Limit:             100,
		IncludePartialBar: true,
	})
	if err != nil {
		log.Printf("Retrieve 1-minute bars failed: %v", err)
	} else if oneMinBarsResp.Success {
		fmt.Printf("Retrieved %d 1-minute bars\n", len(oneMinBarsResp.Bars))
		if len(oneMinBarsResp.Bars) > 0 {
			fmt.Println("Last 5 bars:")
			start := len(oneMinBarsResp.Bars) - 5
			if start < 0 {
				start = 0
			}
			for i := start; i < len(oneMinBarsResp.Bars); i++ {
				bar := oneMinBarsResp.Bars[i]
				fmt.Printf("  %s - O: %.2f, H: %.2f, L: %.2f, C: %.2f, V: %d\n",
					bar.T.Format("15:04:05"), bar.O, bar.H, bar.L, bar.C, bar.V)
			}
		}
	}

	fmt.Println("\n2. Retrieve 5-Minute Bars")
	fiveMinBarsResp, err := tsxClient.History.GetBars(ctx, &models.RetrieveBarRequest{
		ContractID:        contractID,
		Live:              true,
		StartTime:         startTime,
		EndTime:           endTime,
		Unit:              models.AggregateBarUnitMinute,
		UnitNumber:        5,
		Limit:             50,
		IncludePartialBar: false,
	})
	if err != nil {
		log.Printf("Retrieve 5-minute bars failed: %v", err)
	} else if fiveMinBarsResp.Success {
		fmt.Printf("Retrieved %d 5-minute bars\n", len(fiveMinBarsResp.Bars))
	}

	fmt.Println("\n3. Retrieve Hourly Bars")
	hourlyStartTime := endTime.Add(-7 * 24 * time.Hour)
	hourlyBarsResp, err := tsxClient.History.GetBars(ctx, &models.RetrieveBarRequest{
		ContractID: contractID,
		Live:       true,
		StartTime:  hourlyStartTime,
		EndTime:    endTime,
		Unit:       models.AggregateBarUnitHour,
		UnitNumber: 1,
		Limit:      168,
	})
	if err != nil {
		log.Printf("Retrieve hourly bars failed: %v", err)
	} else if hourlyBarsResp.Success {
		fmt.Printf("Retrieved %d hourly bars\n", len(hourlyBarsResp.Bars))

		if len(hourlyBarsResp.Bars) > 0 {
			var totalVolume int64
			var highPrice, lowPrice float64 = hourlyBarsResp.Bars[0].H, hourlyBarsResp.Bars[0].L

			for _, bar := range hourlyBarsResp.Bars {
				totalVolume += bar.V
				if bar.H > highPrice {
					highPrice = bar.H
				}
				if bar.L < lowPrice {
					lowPrice = bar.L
				}
			}

			fmt.Printf("Summary - High: %.2f, Low: %.2f, Total Volume: %d\n",
				highPrice, lowPrice, totalVolume)
		}
	}

	fmt.Println("\n4. Retrieve Daily Bars")
	dailyStartTime := endTime.Add(-30 * 24 * time.Hour)
	dailyBarsResp, err := tsxClient.History.GetBars(ctx, &models.RetrieveBarRequest{
		ContractID: contractID,
		Live:       true,
		StartTime:  dailyStartTime,
		EndTime:    endTime,
		Unit:       models.AggregateBarUnitDay,
		UnitNumber: 1,
		Limit:      30,
	})
	if err != nil {
		log.Printf("Retrieve daily bars failed: %v", err)
	} else if dailyBarsResp.Success {
		fmt.Printf("Retrieved %d daily bars\n", len(dailyBarsResp.Bars))
		if len(dailyBarsResp.Bars) > 0 {
			fmt.Println("Last 5 daily bars:")
			start := len(dailyBarsResp.Bars) - 5
			if start < 0 {
				start = 0
			}
			for i := start; i < len(dailyBarsResp.Bars); i++ {
				bar := dailyBarsResp.Bars[i]
				fmt.Printf("  %s - O: %.2f, H: %.2f, L: %.2f, C: %.2f, V: %d\n",
					bar.T.Format("2006-01-02"), bar.O, bar.H, bar.L, bar.C, bar.V)
			}
		}
	}

	fmt.Println("\n5. Different Time Units Examples")
	units := []struct {
		unit       models.AggregateBarUnit
		unitNumber int32
		name       string
	}{
		{models.AggregateBarUnitSecond, 30, "30-second"},
		{models.AggregateBarUnitMinute, 15, "15-minute"},
		{models.AggregateBarUnitHour, 4, "4-hour"},
		{models.AggregateBarUnitDay, 1, "daily"},
		{models.AggregateBarUnitWeek, 1, "weekly"},
	}

	recentEndTime := time.Now()
	recentStartTime := recentEndTime.Add(-48 * time.Hour)

	for _, u := range units {
		resp, err := tsxClient.History.GetBars(ctx, &models.RetrieveBarRequest{
			ContractID: contractID,
			Live:       true,
			StartTime:  recentStartTime,
			EndTime:    recentEndTime,
			Unit:       u.unit,
			UnitNumber: u.unitNumber,
			Limit:      10,
		})
		if err != nil {
			log.Printf("Retrieve %s bars failed: %v", u.name, err)
		} else if resp.Success {
			fmt.Printf("  %s bars: Retrieved %d bars\n", u.name, len(resp.Bars))
		}
	}
}
