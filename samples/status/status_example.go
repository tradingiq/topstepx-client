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

	fmt.Println("=== Status Service Examples ===")

	fmt.Println("\n1. Check API Status")
	statusResp, err := tsxClient.Status.Ping(ctx)
	if err != nil {
		log.Printf("Status check failed: %v", err)
	} else {
		fmt.Printf("API Status Response: %+v\n", statusResp)
	}

	fmt.Println("\n2. Periodic Status Monitoring")
	fmt.Println("Checking API status every 5 seconds for 20 seconds...")

	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)
	count := 0

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				count++
				status, err := tsxClient.Status.Ping(ctx)
				if err != nil {
					fmt.Printf("[%s] Status check failed: %v\n", t.Format("15:04:05"), err)
				} else {
					fmt.Printf("[%s] API is responding, Status %s\n", t.Format("15:04:05"), status)
				}

				if count >= 4 {
					ticker.Stop()
					done <- true
				}
			}
		}
	}()

	<-done
	fmt.Println("\nStatus monitoring completed")

	fmt.Println("\n3. Status Check with Context Timeout")
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	statusResp, err = tsxClient.Status.Ping(ctxWithTimeout)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("Status check timed out (expected behavior for demonstration)")
		} else {
			fmt.Printf("Status check with timeout failed: %v\n", err)
		}
	} else {
		fmt.Println("Status check with timeout succeeded")
	}

	fmt.Println("\n4. Health Check Pattern")
	fmt.Println("Implementing a health check pattern...")

	healthCheck := func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := tsxClient.Status.Ping(ctx)
		return err == nil
	}

	isHealthy := healthCheck()
	if isHealthy {
		fmt.Println("✓ API is healthy and responding")
	} else {
		fmt.Println("✗ API health check failed")
	}

	fmt.Println("\n5. Multiple Status Checks")
	fmt.Println("Performing rapid status checks...")

	successCount := 0
	totalChecks := 10

	start := time.Now()
	for i := 0; i < totalChecks; i++ {
		_, err := tsxClient.Status.Ping(ctx)
		if err == nil {
			successCount++
		}
		time.Sleep(100 * time.Millisecond)
	}
	duration := time.Since(start)

	fmt.Printf("\nStatus Check Summary:\n")
	fmt.Printf("  Total checks: %d\n", totalChecks)
	fmt.Printf("  Successful: %d\n", successCount)
	fmt.Printf("  Failed: %d\n", totalChecks-successCount)
	fmt.Printf("  Success rate: %.1f%%\n", float64(successCount)/float64(totalChecks)*100)
	fmt.Printf("  Total duration: %v\n", duration)
	fmt.Printf("  Average response time: %v\n", duration/time.Duration(totalChecks))
}
