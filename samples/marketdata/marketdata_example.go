package main

import (
	"context"
	"fmt"
	"github.com/tradingiq/topstepx-client/samples"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	topstepx "github.com/tradingiq/topstepx-client"
	"github.com/tradingiq/topstepx-client/models"
	"github.com/tradingiq/topstepx-client/services"
)

func main() {

	client := topstepx.NewClient()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Logging in...")
	if err := client.LoginAndConnect(ctx,
		samples.Config.Username,
		samples.Config.ApiKey,
	); err != nil {
		log.Fatal("Failed to login: ", err)
	}
	fmt.Println("Login successful!")

	client.MarketData.SetConnectionHandler(func(state services.ConnectionState) {
		switch state {
		case services.StateDisconnected:
			fmt.Println("Market data WebSocket disconnected")
		case services.StateConnecting:
			fmt.Println("Connecting to market data WebSocket...")
		case services.StateConnected:
			fmt.Println("Market data WebSocket connected!")
		case services.StateReconnecting:
			fmt.Println("Reconnecting to market data WebSocket...")
		}
	})

	fmt.Println("Connecting to market data WebSocket...")
	if err := client.MarketData.Connect(ctx); err != nil {
		log.Fatal("Failed to connect to market data WebSocket: ", err)
	}

	client.MarketData.SetQuoteHandler(func(contractID string, quote models.Quote) {
		fmt.Printf("\nQUOTE for %s:\n", contractID)
		fmt.Printf("  Bid: %.2f\n", quote.BestBid)
		fmt.Printf("  Ask: %.2f\n", quote.BestAsk)
		fmt.Printf("  Last: %.2f\n", quote.LastPrice)
		fmt.Printf("  Volume: %.0f\n", quote.Volume)
		fmt.Printf("  Change: %.2f (%.2f%%)\n", quote.Change, quote.ChangePercent*100)
		fmt.Printf("  Time: %s\n", quote.Timestamp.Format("15:04:05"))
	})

	client.MarketData.SetTradeHandler(func(contractID string, trades models.TradeData) {
		for _, trade := range trades {
			fmt.Printf("\nTRADE for %s:\n", contractID)
			fmt.Printf("  Price: %.2f\n", trade.Price)
			fmt.Printf("  Volume: %.0f\n", trade.Volume)
			fmt.Printf("  Type: %d\n", trade.Type)
			fmt.Printf("  Time: %s\n", trade.Timestamp.Format("15:04:05"))
		}
	})

	client.MarketData.SetDepthHandler(func(contractID string, depth models.MarketDepthData) {
		fmt.Printf("\nMARKET DEPTH for %s:\n", contractID)

		bids := []models.MarketDepth{}
		asks := []models.MarketDepth{}

		for _, entry := range depth {
			switch entry.Type {
			case 4:
				bids = append(bids, entry)
			case 3:
				asks = append(asks, entry)
			}
		}

		if len(bids) > 0 {
			fmt.Println("  Bids:")
			for i, bid := range bids {
				if i >= 3 {
					break
				}
				fmt.Printf("    Level %d: %.2f @ %.0f\n", i+1, bid.Price, bid.Volume)
			}
		}

		if len(asks) > 0 {
			fmt.Println("  Asks:")
			for i, ask := range asks {
				if i >= 3 {
					break
				}
				fmt.Printf("    Level %d: %.2f @ %.0f\n", i+1, ask.Price, ask.Volume)
			}
		}
	})

	fmt.Println("\nSearching for E-mini Russell 2000 futures...")
	searchText := "RTY"
	searchResp, err := client.Contract.SearchContracts(ctx, &models.SearchContractRequest{
		SearchText: &searchText,
		Live:       false,
	})
	if err != nil {
		log.Fatal("Failed to search contracts: ", err)
	}

	if !searchResp.Success || len(searchResp.Contracts) == 0 {
		log.Fatal("No contracts found")
	}

	contract := searchResp.Contracts[0]
	contractID := contract.ID
	fmt.Printf("Found contract: %s (%s)\n", contract.Name, contract.Description)

	fmt.Printf("\nSubscribing to market data for %s...\n", contractID)
	if err := client.MarketData.SubscribeAll(contractID); err != nil {
		log.Fatal("Failed to subscribe to market data: ", err)
	}
	fmt.Println("Successfully subscribed to quotes, trades, and market depth!")

	if len(searchResp.Contracts) > 1 {
		contract2 := searchResp.Contracts[1]
		contractID2 := contract2.ID
		fmt.Printf("\nAlso subscribing to %s (%s)...\n", contract2.Name, contract2.Description)
		if err := client.MarketData.SubscribeAll(contractID2); err != nil {
			fmt.Printf("Failed to subscribe to second contract: %v\n", err)
		}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("\n Market data streaming started!")
	fmt.Println("Press Ctrl+C to stop...")

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	cleanup := func() {
		fmt.Println("Unsubscribing from all market data...")
		if err := client.MarketData.UnsubscribeAllContracts(); err != nil {
			fmt.Printf("Error unsubscribing: %v\n", err)
		}

		fmt.Println("Disconnecting...")
		if err := client.Disconnect(ctx); err != nil {
			fmt.Printf("Error disconnecting: %v\n", err)
		}

		fmt.Println("Goodbye!")
	}

	running := true
	for running {
		select {
		case <-sigChan:
			fmt.Println("\n\nReceived interrupt signal, shutting down...")
			running = false
		case <-ticker.C:
			subs := client.MarketData.GetSubscriptions()
			fmt.Printf("\n Active subscriptions: %d contracts\n", len(subs))
			for cid, types := range subs {
				fmt.Printf("  %s: ", cid)
				for dtype := range types {
					fmt.Printf("%s ", dtype)
				}
				fmt.Println()
			}
		}
	}

	cleanup()
}
