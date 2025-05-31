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

	err := tsxClient.LoginAndConnect(ctx, username, apiKey)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	account, err := tsxClient.GetFirstActiveAccount(ctx)
	if err != nil {
		log.Fatalf("GetFistActiveAccount failed: %v", err)
	}
	accountID := account.ID

	fmt.Printf("Using account ID: %d\n", accountID)

	searchText := "ES"
	contractsResp, err := tsxClient.Contract.SearchContracts(ctx, &models.SearchContractRequest{
		SearchText: &searchText,
		Live:       false,
	})
	if err != nil || !contractsResp.Success || len(contractsResp.Contracts) == 0 {
		log.Fatal("Failed to find contracts")
	}

	contractID := contractsResp.Contracts[0].ID
	fmt.Printf("Using contract: %s (%s)\n", contractID, contractsResp.Contracts[0].Name)

	fmt.Println("\n=== Order Service Examples ===")

	fmt.Println("\n1. Search Historical Orders")
	endTime := time.Now()
	startTime := endTime.Add(-7 * 24 * time.Hour)
	histOrdersResp, err := tsxClient.Order.SearchOrders(ctx, &models.SearchOrderRequest{
		AccountID:      int32(accountID),
		StartTimestamp: startTime,
		EndTimestamp:   &endTime,
	})
	if err != nil {
		log.Printf("Search historical orders failed: %v", err)
	} else if histOrdersResp.Success {
		fmt.Printf("Found %d historical orders in the last 7 days\n", len(histOrdersResp.Orders))
		for _, order := range histOrdersResp.Orders {
			fmt.Printf("  - Order %d: %s %s %d @ %.2f, Status: %v\n",
				order.ID, getOrderTypeString(order.Type), getOrderSideString(order.Side),
				order.Size, getLimitPrice(order), getOrderStatusString(order.Status))
		}
	}

	fmt.Println("\n2. Search Open Orders")
	openOrdersResp, err := tsxClient.Order.SearchOpenOrders(ctx, &models.SearchOpenOrderRequest{
		AccountID: int32(accountID),
	})
	if err != nil {
		log.Printf("Search open orders failed: %v", err)
	} else if openOrdersResp.Success {
		fmt.Printf("Found %d open orders\n", len(openOrdersResp.Orders))
		for _, order := range openOrdersResp.Orders {
			fmt.Printf("  - Order %d: %s %s %d @ %.2f\n",
				order.ID, getOrderTypeString(order.Type), getOrderSideString(order.Side),
				order.Size, getLimitPrice(order))
		}
	}

	fmt.Println("\n3. Place Different Order Types (Demo)")
	fmt.Println("NOTE: These are examples. Actual execution depends on market conditions and account status.")

	fmt.Println("\n3.1 Place Limit Order")
	limitPrice := 4500.0
	customTag := "example-limit-order"
	limitOrderResp, err := tsxClient.Order.PlaceOrder(ctx, &models.PlaceOrderRequest{
		AccountID:  int32(accountID),
		ContractID: contractID,
		Type:       models.OrderTypeLimit,
		Side:       models.OrderSideBid,
		Size:       1,
		LimitPrice: &limitPrice,
		CustomTag:  &customTag,
	})
	if err != nil {
		fmt.Printf("Place limit order error: %v\n", err)
	} else if limitOrderResp.Success && limitOrderResp.OrderID != nil {
		fmt.Printf("Limit order placed successfully! Order ID: %d\n", *limitOrderResp.OrderID)

		time.Sleep(2 * time.Second)

		fmt.Println("\n4. Modify Order")
		newLimitPrice := 4505.0
		newSize := int32(2)
		modifyResp, err := tsxClient.Order.ModifyOrder(ctx, &models.ModifyOrderRequest{
			AccountID:  int32(accountID),
			OrderID:    *limitOrderResp.OrderID,
			Size:       &newSize,
			LimitPrice: &newLimitPrice,
		})
		if err != nil {
			fmt.Printf("Modify order error: %v\n", err)
		} else if modifyResp.Success {
			fmt.Println("Order modified successfully!")
		}

		time.Sleep(2 * time.Second)

		fmt.Println("\n5. Cancel Order")
		cancelResp, err := tsxClient.Order.CancelOrder(ctx, &models.CancelOrderRequest{
			AccountID: int32(accountID),
			OrderID:   *limitOrderResp.OrderID,
		})
		if err != nil {
			fmt.Printf("Cancel order error: %v\n", err)
		} else if cancelResp.Success {
			fmt.Println("Order cancelled successfully!")
		}
	} else {
		fmt.Printf("Limit order failed with error code: %v\n", limitOrderResp.ErrorCode)
	}

	fmt.Println("\n3.2 Place Stop Order Example")
	stopPrice := 4600.0
	stopOrderResp, err := tsxClient.Order.PlaceOrder(ctx, &models.PlaceOrderRequest{
		AccountID:  int32(accountID),
		ContractID: contractID,
		Type:       models.OrderTypeStop,
		Side:       models.OrderSideAsk,
		Size:       1,
		StopPrice:  &stopPrice,
	})
	if err != nil {
		fmt.Printf("Place stop order error: %v\n", err)
	} else if stopOrderResp.Success && stopOrderResp.OrderID != nil {
		fmt.Printf("Stop order placed successfully! Order ID: %d\n", *stopOrderResp.OrderID)

		cancelResp, _ := tsxClient.Order.CancelOrder(ctx, &models.CancelOrderRequest{
			AccountID: int32(accountID),
			OrderID:   *stopOrderResp.OrderID,
		})
		if cancelResp.Success {
			fmt.Println("Stop order cancelled")
		}
	}

	fmt.Println("\n3.3 Place Stop Limit Order Example")
	stopLimitPrice := 4595.0
	stopLimitOrderResp, err := tsxClient.Order.PlaceOrder(ctx, &models.PlaceOrderRequest{
		AccountID:  int32(accountID),
		ContractID: contractID,
		Type:       models.OrderTypeStopLimit,
		Side:       models.OrderSideAsk,
		Size:       1,
		LimitPrice: &stopLimitPrice,
		StopPrice:  &stopPrice,
	})
	if err != nil {
		fmt.Printf("Place stop limit order error: %v\n", err)
	} else if stopLimitOrderResp.Success && stopLimitOrderResp.OrderID != nil {
		fmt.Printf("Stop limit order placed successfully! Order ID: %d\n", *stopLimitOrderResp.OrderID)

		cancelResp, _ := tsxClient.Order.CancelOrder(ctx, &models.CancelOrderRequest{
			AccountID: int32(accountID),
			OrderID:   *stopLimitOrderResp.OrderID,
		})
		if cancelResp.Success {
			fmt.Println("Stop limit order cancelled")
		}
	}

	fmt.Println("\n6. Final Open Orders Check")
	finalOpenOrdersResp, err := tsxClient.Order.SearchOpenOrders(ctx, &models.SearchOpenOrderRequest{
		AccountID: int32(accountID),
	})
	if err != nil {
		log.Printf("Search final open orders failed: %v", err)
	} else if finalOpenOrdersResp.Success {
		fmt.Printf("Final open orders count: %d\n", len(finalOpenOrdersResp.Orders))
	}
}

func getOrderTypeString(orderType models.OrderType) string {
	switch orderType {
	case models.OrderTypeLimit:
		return "LIMIT"
	case models.OrderTypeMarket:
		return "MARKET"
	case models.OrderTypeStop:
		return "STOP"
	case models.OrderTypeStopLimit:
		return "STOP_LIMIT"
	case models.OrderTypeTrailingStop:
		return "TRAILING_STOP"
	case models.OrderTypeJoinBid:
		return "JOIN_BID"
	case models.OrderTypeJoinAsk:
		return "JOIN_ASK"
	default:
		return "UNKNOWN"
	}
}

func getOrderSideString(side models.OrderSide) string {
	if side == models.OrderSideBid {
		return "BUY"
	}
	return "SELL"
}

func getOrderStatusString(status models.OrderStatus) string {
	switch status {
	case models.OrderStatusOpen:
		return "OPEN"
	case models.OrderStatusFilled:
		return "FILLED"
	case models.OrderStatusCancelled:
		return "CANCELLED"
	case models.OrderStatusExpired:
		return "EXPIRED"
	case models.OrderStatusRejected:
		return "REJECTED"
	case models.OrderStatusPending:
		return "PENDING"
	default:
		return "NONE"
	}
}

func getLimitPrice(order models.OrderModel) float64 {
	if order.LimitPrice != nil {
		return *order.LimitPrice
	}
	if order.StopPrice != nil {
		return *order.StopPrice
	}
	return 0.0
}
