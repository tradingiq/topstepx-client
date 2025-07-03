package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	searchText := "MES"
	contractsResp, _ := client.Contract.SearchContracts(ctx, &models.SearchContractRequest{
		SearchText: &searchText,
		Live:       false,
	})
	if len(contractsResp.Contracts) == 0 {
		log.Fatal("Contract not found")
	}
	contractID := contractsResp.Contracts[0].ID

	client.MarketData.SetQuoteHandler(func(contractID string, quote models.Quote) {
		fmt.Printf("[%s] %s - Bid: %.2f, Ask: %.2f, Last: %.2f, Volume: %d\n",
			quote.Timestamp.Format("15:04:05"), contractID,
			quote.BestBid, quote.BestAsk, quote.LastPrice, quote.Volume)
	})

	fmt.Println("Connecting to market data...")
	if err := client.MarketData.Connect(ctx); err != nil {
		log.Fatal("Failed to connect:", err)
	}

	fmt.Printf("Subscribing to quotes for %s...\n", contractID)
	if err := client.MarketData.SubscribeContractQuotes(contractID); err != nil {
		log.Fatal("Failed to subscribe:", err)
	}

	fmt.Println("✅ Streaming quotes (Press Ctrl+Close to stop)...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down...")
	client.MarketData.UnsubscribeAllContracts()
	client.MarketData.Disconnect()
}
