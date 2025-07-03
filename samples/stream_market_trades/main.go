package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tradingiq/projectx-client"
	"github.com/tradingiq/projectx-client/models"
	"github.com/tradingiq/projectx-client/samples"
)

func main() {
	client := projectx.NewClient()
	ctx := context.Background()

	loginResp, _ := client.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
		UserName: samples.Config.Username,
		APIKey:   samples.Config.ApiKey,
	})
	if !loginResp.Success {
		log.Fatal("Login failed")
	}

	searchText := "ES"
	contractsResp, _ := client.Contract.SearchContracts(ctx, &models.SearchContractRequest{
		SearchText: &searchText,
		Live:       false,
	})
	if len(contractsResp.Contracts) == 0 {
		log.Fatal("Contract not found")
	}
	contract := contractsResp.Contracts[0]

	// Set up trade handler
	client.MarketData.SetTradeHandler(func(contractID string, trades models.TradeData) {
		for _, trade := range trades {
			fmt.Printf("\n📊 TRADE - %s\n", contractID)
			fmt.Printf("  Price: $%.2f\n", trade.Price)
			fmt.Printf("  Volume: %d\n", trade.Volume)
			fmt.Printf("  Time: %s\n", trade.Timestamp.Format("15:04:05.000"))
		}
	})

	fmt.Println("Connecting to market data WebSocket...")
	if err := client.MarketData.Connect(ctx); err != nil {
		log.Fatal("Failed to connect:", err)
	}

	// Subscribe to trades
	fmt.Printf("Subscribing to trades for %s (%s)...\n", contract.Name, contract.ID)
	if err := client.MarketData.SubscribeContractTrades(contract.ID); err != nil {
		log.Fatal("Failed to subscribe:", err)
	}

	fmt.Println("✅ Monitoring trades...")
	fmt.Println("Press Ctrl+C to stop...")

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down...")
	client.MarketData.UnsubscribeContractTrades(contract.ID)
	client.MarketData.Disconnect()
}

func getAggressorSide(aggressor int) string {
	switch aggressor {
	case 1:
		return "BUY"
	case 2:
		return "SELL"
	default:
		return "UNKNOWN"
	}
}
