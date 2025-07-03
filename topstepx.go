package projectx

import (
	"context"
	"fmt"

	"github.com/tradingiq/projectx-client/client"
	"github.com/tradingiq/projectx-client/models"
	"github.com/tradingiq/projectx-client/services"
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
