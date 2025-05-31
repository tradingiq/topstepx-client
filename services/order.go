package services

import (
	"context"
	"net/http"

	"github.com/tradingiq/topstepx-client/client"
	"github.com/tradingiq/topstepx-client/models"
)

type OrderService struct {
	client *client.Client
}

func NewOrderService(c *client.Client) *OrderService {
	return &OrderService{client: c}
}

func (s *OrderService) SearchOrders(ctx context.Context, req *models.SearchOrderRequest) (*models.SearchOrderResponse, error) {
	var resp models.SearchOrderResponse
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodPost,
		Path:   "/api/Order/search",
		Body:   req,
	}, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *OrderService) SearchOpenOrders(ctx context.Context, req *models.SearchOpenOrderRequest) (*models.SearchOrderResponse, error) {
	var resp models.SearchOrderResponse
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodPost,
		Path:   "/api/Order/searchOpen",
		Body:   req,
	}, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *OrderService) PlaceOrder(ctx context.Context, req *models.PlaceOrderRequest) (*models.PlaceOrderResponse, error) {
	var resp models.PlaceOrderResponse
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodPost,
		Path:   "/api/Order/place",
		Body:   req,
	}, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, req *models.CancelOrderRequest) (*models.CancelOrderResponse, error) {
	var resp models.CancelOrderResponse
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodPost,
		Path:   "/api/Order/cancel",
		Body:   req,
	}, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *OrderService) ModifyOrder(ctx context.Context, req *models.ModifyOrderRequest) (*models.ModifyOrderResponse, error) {
	var resp models.ModifyOrderResponse
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodPost,
		Path:   "/api/Order/modify",
		Body:   req,
	}, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
