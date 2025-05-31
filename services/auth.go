package services

import (
	"context"
	"net/http"

	"github.com/tradingiq/topstepx-client/client"
	"github.com/tradingiq/topstepx-client/models"
)

type AuthService struct {
	client *client.Client
}

func NewAuthService(c *client.Client) *AuthService {
	return &AuthService{client: c}
}

func (s *AuthService) LoginApp(ctx context.Context, req *models.LoginAppRequest) (*models.LoginResponse, error) {
	var resp models.LoginResponse
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodPost,
		Path:   "/api/Auth/loginApp",
		Body:   req,
	}, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Success && resp.Token != nil {
		s.client.SetToken(*resp.Token)
	}

	return &resp, nil
}

func (s *AuthService) LoginKey(ctx context.Context, req *models.LoginApiKeyRequest) (*models.LoginResponse, error) {
	var resp models.LoginResponse
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodPost,
		Path:   "/api/Auth/loginKey",
		Body:   req,
	}, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Success && resp.Token != nil {
		s.client.SetToken(*resp.Token)
	}

	return &resp, nil
}

func (s *AuthService) Logout(ctx context.Context) (*models.LogoutResponse, error) {
	var resp models.LogoutResponse
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodPost,
		Path:   "/api/Auth/logout",
	}, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Success {
		s.client.SetToken("")
	}

	return &resp, nil
}

func (s *AuthService) Validate(ctx context.Context) (*models.ValidateResponse, error) {
	var resp models.ValidateResponse
	err := s.client.DoJSON(ctx, &client.Request{
		Method: http.MethodPost,
		Path:   "/api/Auth/validate",
	}, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Success && resp.NewToken != nil {
		s.client.SetToken(*resp.NewToken)
	}

	return &resp, nil
}
