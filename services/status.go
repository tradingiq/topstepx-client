package services

import (
	"context"
	"net/http"

	"github.com/tradingiq/topstepx-client/client"
)

type StatusService struct {
	client *client.Client
}

func NewStatusService(c *client.Client) *StatusService {
	return &StatusService{client: c}
}

func (s *StatusService) Ping(ctx context.Context) (string, error) {
	var resp string
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodGet,
		Path:   "/api/Status/ping",
	}, &resp)
	if err != nil {
		return "", err
	}

	return resp, nil
}
