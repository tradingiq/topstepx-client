package services

import (
	"context"
	"net/http"

	"github.com/tradingiq/topstepx-client/client"
	"github.com/tradingiq/topstepx-client/models"
)

type ContractService struct {
	client *client.Client
}

func NewContractService(c *client.Client) *ContractService {
	return &ContractService{client: c}
}

func (s *ContractService) SearchContracts(ctx context.Context, req *models.SearchContractRequest) (*models.SearchContractResponse, error) {
	var resp models.SearchContractResponse
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodPost,
		Path:   "/api/Contract/search",
		Body:   req,
	}, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *ContractService) SearchContractByID(ctx context.Context, req *models.SearchContractByIdRequest) (*models.SearchContractByIdResponse, error) {
	var resp models.SearchContractByIdResponse
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodPost,
		Path:   "/api/Contract/searchById",
		Body:   req,
	}, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
