package services

import (
	"context"
	"net/http"

	"github.com/tradingiq/projectx-client/client"
	"github.com/tradingiq/projectx-client/models"
)

type HistoryService struct {
	client *client.Client
}

func NewHistoryService(c *client.Client) *HistoryService {
	return &HistoryService{client: c}
}

func (s *HistoryService) GetBars(ctx context.Context, req *models.RetrieveBarRequest) (*models.RetrieveBarResponse, error) {
	var resp models.RetrieveBarResponse
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodPost,
		Path:   "/api/History/retrieveBars",
		Body:   req,
	}, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
