package services

import (
	"context"
	"net/http"

	"github.com/tradingiq/projectx-client/client"
	"github.com/tradingiq/projectx-client/models"
)

type PositionService struct {
	client *client.Client
}

func NewPositionService(c *client.Client) *PositionService {
	return &PositionService{client: c}
}

func (s *PositionService) SearchOpenPositions(ctx context.Context, req *models.SearchPositionRequest) (*models.SearchPositionResponse, error) {
	var resp models.SearchPositionResponse
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodPost,
		Path:   "/api/Position/searchOpen",
		Body:   req,
	}, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *PositionService) CloseContractPosition(ctx context.Context, req *models.CloseContractPositionRequest) (*models.ClosePositionResponse, error) {
	var resp models.ClosePositionResponse
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodPost,
		Path:   "/api/Position/closeContract",
		Body:   req,
	}, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *PositionService) PartialCloseContractPosition(ctx context.Context, req *models.PartialCloseContractPositionRequest) (*models.PartialClosePositionResponse, error) {
	var resp models.PartialClosePositionResponse
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodPost,
		Path:   "/api/Position/partialCloseContract",
		Body:   req,
	}, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
