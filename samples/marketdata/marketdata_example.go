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
)

func main() {
	// Create a new client
	client := topstepx.NewClient()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Login
	fmt.Println("Logging in...")
	if err := client.LoginAndConnect(ctx,
		samples.Config.Username,
		samples.Config.ApiKey,
	); err != nil {
		log.Fatal("Failed to login: ", err)
	}
	fmt.Println("Login successful!")

	// Connect to market data WebSocket
	fmt.Println("Connecting to market data WebSocket...")
	if err := client.MarketData.Connect(ctx); err != nil {
		log.Fatal("Failed to connect to market data WebSocket: ", err)
	}
	fmt.Println("Market data WebSocket connected!")

	// Set up handlers for market data events
	client.MarketData.SetQuoteHandler(func(contractID string, data interface{}) {
		fmt.Printf("\n📊 QUOTE for %s:\n", contractID)
		if quote, ok := data.(map[string]interface{}); ok {
			if bid, ok := quote["bid"].(float64); ok {
				fmt.Printf("  Bid: %.2f", bid)
			}
			if bidSize, ok := quote["bidSize"].(float64); ok {
				fmt.Printf(" (Size: %.0f)\n", bidSize)
			}
			if ask, ok := quote["ask"].(float64); ok {
				fmt.Printf("  Ask: %.2f", ask)
			}
			if askSize, ok := quote["askSize"].(float64); ok {
				fmt.Printf(" (Size: %.0f)\n", askSize)
			}
			if timestamp, ok := quote["timestamp"].(string); ok {
				fmt.Printf("  Time: %s\n", timestamp)
			}
		}
	})

	client.MarketData.SetTradeHandler(func(contractID string, data interface{}) {
		fmt.Printf("\n💹 TRADE for %s:\n", contractID)
		if trade, ok := data.(map[string]interface{}); ok {
			if price, ok := trade["price"].(float64); ok {
				fmt.Printf("  Price: %.2f", price)
			}
			if size, ok := trade["size"].(float64); ok {
				fmt.Printf(" (Size: %.0f)\n", size)
			}
			if timestamp, ok := trade["timestamp"].(string); ok {
				fmt.Printf("  Time: %s\n", timestamp)
			}
		}
	})

	client.MarketData.SetDepthHandler(func(contractID string, data interface{}) {
		fmt.Printf("\n📈 MARKET DEPTH for %s:\n", contractID)
		if depth, ok := data.(map[string]interface{}); ok {
			if bids, ok := depth["bids"].([]interface{}); ok {
				fmt.Println("  Bids:")
				for i, bid := range bids {
					if i >= 3 {
						break // Show only top 3 levels
					}
					if bidMap, ok := bid.(map[string]interface{}); ok {
						price := bidMap["price"].(float64)
						size := bidMap["size"].(float64)
						fmt.Printf("    Level %d: %.2f @ %.0f\n", i+1, price, size)
					}
				}
			}
			if asks, ok := depth["asks"].([]interface{}); ok {
				fmt.Println("  Asks:")
				for i, ask := range asks {
					if i >= 3 {
						break // Show only top 3 levels
					}
					if askMap, ok := ask.(map[string]interface{}); ok {
						price := askMap["price"].(float64)
						size := askMap["size"].(float64)
						fmt.Printf("    Level %d: %.2f @ %.0f\n", i+1, price, size)
					}
				}
			}
		}
	})

	// Search for a contract to subscribe to
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

	// Use the first active contract
	contract := searchResp.Contracts[0]
	contractID := contract.ID
	fmt.Printf("Found contract: %s (%s)\n", contract.Name, contract.Description)

	// Subscribe to all market data for this contract
	fmt.Printf("\nSubscribing to market data for %s...\n", contractID)
	if err := client.MarketData.SubscribeAll(contractID); err != nil {
		log.Fatal("Failed to subscribe to market data: ", err)
	}
	fmt.Println("Successfully subscribed to quotes, trades, and market depth!")

	// Also subscribe to another contract if available
	if len(searchResp.Contracts) > 1 {
		contract2 := searchResp.Contracts[1]
		contractID2 := contract2.ID
		fmt.Printf("\nAlso subscribing to %s (%s)...\n", contract2.Name, contract2.Description)
		if err := client.MarketData.SubscribeAll(contractID2); err != nil {
			fmt.Printf("Failed to subscribe to second contract: %v\n", err)
		}
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Print instructions
	fmt.Println("\n✅ Market data streaming started!")
	fmt.Println("Press Ctrl+C to stop...")

	// Create a ticker to show we're still running
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Wait for interrupt signal
	for {
		select {
		case <-sigChan:
			fmt.Println("\n\nReceived interrupt signal, shutting down...")
			goto cleanup
		case <-ticker.C:
			// Show current subscriptions periodically
			subs := client.MarketData.GetSubscriptions()
			fmt.Printf("\n📍 Active subscriptions: %d contracts\n", len(subs))
			for cid, types := range subs {
				fmt.Printf("  %s: ", cid)
				for dtype := range types {
					fmt.Printf("%s ", dtype)
				}
				fmt.Println()
			}
		}
	}

cleanup:
	// Unsubscribe from all contracts
	fmt.Println("Unsubscribing from all market data...")
	if err := client.MarketData.UnsubscribeAllContracts(); err != nil {
		fmt.Printf("Error unsubscribing: %v\n", err)
	}

	// Disconnect from everything
	fmt.Println("Disconnecting...")
	if err := client.Disconnect(ctx); err != nil {
		fmt.Printf("Error disconnecting: %v\n", err)
	}

	fmt.Println("Goodbye!")
}
