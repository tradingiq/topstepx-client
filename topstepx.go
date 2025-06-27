package topstepx

import (
	"context"
	"fmt"

	"github.com/tradingiq/topstepx-client/client"
	"github.com/tradingiq/topstepx-client/models"
	"github.com/tradingiq/topstepx-client/services"
)

type Client struct {
	client *client.Client

	Auth       *services.AuthService
	Account    *services.AccountService
	Contract   *services.ContractService
	Order      *services.OrderService
	Position   *services.PositionService
	History    *services.HistoryService
	Trade      *services.TradeService
	Status     *services.StatusService
	UserData   *services.UserDataWebSocketService
	MarketData *services.MarketDataWebSocketService
}

func NewClient(httpOpts ...client.Option) *Client {
	c := client.NewClient(httpOpts...)

	return &Client{
		client:     c,
		Auth:       services.NewAuthService(c),
		Account:    services.NewAccountService(c),
		Contract:   services.NewContractService(c),
		Order:      services.NewOrderService(c),
		Position:   services.NewPositionService(c),
		History:    services.NewHistoryService(c),
		Trade:      services.NewTradeService(c),
		Status:     services.NewStatusService(c),
		UserData:   services.NewUserDataWebSocketService(c),
		MarketData: services.NewMarketDataWebSocketService(c),
	}
}

func (c *Client) SetToken(token string) {
	c.client.SetToken(token)
}

func (c *Client) GetToken() string {
	return c.client.GetToken()
}

func (c *Client) LoginAndConnect(ctx context.Context, username, apiKey string) error {

	resp, err := c.Auth.LoginKey(ctx, &models.LoginApiKeyRequest{
		UserName: username,
		APIKey:   apiKey,
	})
	if err != nil {
		return err
	}

	if !resp.Success {
		if resp.ErrorMessage != nil {
			return fmt.Errorf("login failed: %s", *resp.ErrorMessage)
		}
		return fmt.Errorf("login failed with error code: %v", resp.ErrorCode)
	}
	return nil
}

func (c *Client) ConnectWebSocketWithAccount(ctx context.Context, accountID int) error {

	if c.GetToken() == "" {
		return fmt.Errorf("no authentication token available, please login first")
	}

	if err := c.UserData.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect WebSocket: %w", err)
	}

	if err := c.UserData.SubscribeAll(accountID); err != nil {
		return fmt.Errorf("failed to subscribe to WebSocket events: %w", err)
	}

	return nil
}

func (c *Client) GetFirstActiveAccount(ctx context.Context) (*models.TradingAccountModel, error) {
	resp, err := c.Account.SearchAccounts(ctx, &models.SearchAccountRequest{
		OnlyActiveAccounts: true,
	})
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		if resp.ErrorMessage != nil {
			return nil, fmt.Errorf("failed to search accounts: %s", *resp.ErrorMessage)
		}
		return nil, fmt.Errorf("failed to search accounts with error code: %v", resp.ErrorCode)
	}

	if len(resp.Accounts) == 0 {
		return nil, fmt.Errorf("no active accounts found")
	}

	return &resp.Accounts[0], nil
}

func (c *Client) GetActiveAccounts(ctx context.Context) ([]models.TradingAccountModel, error) {
	resp, err := c.Account.SearchAccounts(ctx, &models.SearchAccountRequest{
		OnlyActiveAccounts: true,
	})
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		if resp.ErrorMessage != nil {
			return nil, fmt.Errorf("failed to search accounts: %s", *resp.ErrorMessage)
		}
		return nil, fmt.Errorf("failed to search accounts with error code: %v", resp.ErrorCode)
	}

	return resp.Accounts, nil
}

func (c *Client) ConnectMarketData(ctx context.Context) error {
	if c.GetToken() == "" {
		return fmt.Errorf("no authentication token available, please login first")
	}

	if err := c.MarketData.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect market data WebSocket: %w", err)
	}

	return nil
}

func (c *Client) Disconnect(ctx context.Context) error {

	if c.UserData.IsConnected() {
		c.UserData.UnsubscribeAll()
		c.UserData.Disconnect()
	}

	if c.MarketData.IsConnected() {
		c.MarketData.UnsubscribeAllContracts()
		c.MarketData.Disconnect()
	}

	_, err := c.Auth.Logout(ctx)
	return err
}
