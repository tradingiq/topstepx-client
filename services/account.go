package services

import (
	"context"
	"net/http"

	"github.com/tradingiq/projectx-client/client"
	"github.com/tradingiq/projectx-client/models"
)

type AccountService struct {
	client *client.Client
}

func NewAccountService(c *client.Client) *AccountService {
	return &AccountService{client: c}
}

func (s *AccountService) SearchAccounts(ctx context.Context, req *models.SearchAccountRequest) (*models.SearchAccountResponse, error) {
	var resp models.SearchAccountResponse
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodPost,
		Path:   "/api/Account/search",
		Body:   req,
	}, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
