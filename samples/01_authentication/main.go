package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tradingiq/topstepx-client"
	"github.com/tradingiq/topstepx-client/models"
	"github.com/tradingiq/topstepx-client/samples"
)

func main() {
	client := topstepx.NewClient()
	ctx := context.Background()

	resp, err := client.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
		UserName: samples.Config.Username,
		APIKey:   samples.Config.ApiKey,
	})
	if err != nil {
		log.Fatal("Login failed:", err)
	}

	if !resp.Success {
		log.Fatal("Login failed:", resp.ErrorMessage)
	}

	fmt.Println("✅ Authentication successful!")
	fmt.Printf("Token: %s...\n", client.GetToken()[:20])
}
