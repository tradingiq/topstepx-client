package main

import (
	"context"
	"fmt"
	"log"

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

	searchTerms := []string{"ES", "NQ", "CL", "GC", "6E"}

	for _, term := range searchTerms {
		fmt.Printf("\n=== Searching for '%s' contracts ===\n", term)

		searchText := term
		resp, err := client.Contract.SearchContracts(ctx, &models.SearchContractRequest{
			SearchText: &searchText,
			Live:       false,
		})

		if err != nil {
			fmt.Printf("Search failed: %v\n", err)
			continue
		}

		if !resp.Success {
			fmt.Printf("Search failed: %s\n", *resp.ErrorMessage)
			continue
		}

		fmt.Printf("Found %d contracts:\n", len(resp.Contracts))
		for i, contract := range resp.Contracts {
			if i >= 3 {
				break
			}
			fmt.Printf("  %s - %s\n", contract.ID, contract.Name)
			fmt.Printf("    Description: %s\n", contract.Description)
			fmt.Printf("    Tick Size: $%.2f\n", contract.TickSize)
			fmt.Printf("    Tick Value: $%.2f\n", contract.TickValue)
			fmt.Println()
		}
	}
}
