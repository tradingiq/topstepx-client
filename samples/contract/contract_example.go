package main

import (
	"context"
	"fmt"
	"github.com/tradingiq/topstepx-client"
	"github.com/tradingiq/topstepx-client/client"
	"github.com/tradingiq/topstepx-client/models"
	"github.com/tradingiq/topstepx-client/samples"
	"log"
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

	loginResp, err := tsxClient.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
		UserName: username,
		APIKey:   apiKey,
	})
	if err != nil || !loginResp.Success {
		log.Fatal("Login failed")
	}

	fmt.Println("=== Contract Service Examples ===")

	fmt.Println("\n1. Search Contracts - Live Markets")
	searchText := "MES"
	liveContractsResp, err := tsxClient.Contract.SearchContracts(ctx, &models.SearchContractRequest{
		SearchText: &searchText,
		Live:       true,
	})
	if err != nil {
		log.Fatalf("Search live contracts failed: %v", err)
	}

	if liveContractsResp.Success {
		fmt.Printf("Found %d live contracts matching '%s':\n", len(liveContractsResp.Contracts), searchText)
		for _, contract := range liveContractsResp.Contracts {
			fmt.Printf("  - ID: %s, Name: %s, Active: %v\n",
				contract.ID, contract.Name, contract.ActiveContract)
			fmt.Printf("    Description: %s\n", contract.Description)
			fmt.Printf("    Tick Size: %.4f, Tick Value: %.2f\n",
				contract.TickSize, contract.TickValue)
		}
	}

	fmt.Println("\n2. Search Contracts - Demo Markets")
	demoContractsResp, err := tsxClient.Contract.SearchContracts(ctx, &models.SearchContractRequest{
		SearchText: &searchText,
		Live:       false,
	})
	if err != nil {
		log.Fatalf("Search demo contracts failed: %v", err)
	}

	if demoContractsResp.Success {
		fmt.Printf("Found %d demo contracts matching '%s'\n", len(demoContractsResp.Contracts), searchText)
		for _, contract := range demoContractsResp.Contracts {
			fmt.Printf("  - ID: %s, Name: %s, Active: %v\n", contract.ID, contract.Name, contract.ActiveContract)
		}
	}

	fmt.Println("\n3. Search Contract by ID")
	if len(demoContractsResp.Contracts) > 0 {
		contractID := demoContractsResp.Contracts[0].ID
		contractByIdResp, err := tsxClient.Contract.SearchContractByID(ctx, &models.SearchContractByIdRequest{
			ContractID: contractID,
		})
		if err != nil {
			log.Printf("Search contract by ID failed: %v", err)
		} else if contractByIdResp.Success && contractByIdResp.Contract != nil {
			fmt.Printf("Found contract by ID '%s':\n", contractID)
			contract := contractByIdResp.Contract
			fmt.Printf("  Name: %s\n", contract.Name)
			fmt.Printf("  Description: %s\n", contract.Description)
			fmt.Printf("  Tick Size: %.4f\n", contract.TickSize)
			fmt.Printf("  Tick Value: %.2f\n", contract.TickValue)
			fmt.Printf("  Active: %v\n", contract.ActiveContract)
		}
	}

	fmt.Println("\n5. Search with Different Symbols")
	symbols := []string{"NQ", "CL", "GC", "6E"}
	for _, symbol := range symbols {
		resp, err := tsxClient.Contract.SearchContracts(ctx, &models.SearchContractRequest{
			SearchText: &symbol,
			Live:       false,
		})
		if err != nil {
			log.Printf("Search for %s failed: %v", symbol, err)
			continue
		}
		if resp.Success {
			fmt.Printf("  %s: Found %d contracts\n", symbol, len(resp.Contracts))
			for _, con := range resp.Contracts {
				fmt.Printf("  - ID: %s, Name: %s, Active: %v\n", con.ID, con.Name, con.ActiveContract)
			}
		}
	}
}
