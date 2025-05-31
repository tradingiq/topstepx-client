package services

import (
	"context"
	"net/http"

	"github.com/tradingiq/topstepx-client/client"
	"github.com/tradingiq/topstepx-client/models"
)

type TradeService struct {
	client *client.Client
}

func NewTradeService(c *client.Client) *TradeService {
	return &TradeService{client: c}
}

func (s *TradeService) SearchHalfTurnTrades(ctx context.Context, req *models.SearchTradeRequest) (*models.SearchHalfTradeResponse, error) {
	var resp models.SearchHalfTradeResponse
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodPost,
		Path:   "/api/Trade/search",
		Body:   req,
	}, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
