package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type SearchAccountErrorCode int

const (
	SearchAccountErrorCodeSuccess SearchAccountErrorCode = 0
)

type LoginErrorCode int

const (
	LoginErrorCodeSuccess                    LoginErrorCode = 0
	LoginErrorCodeUserNotFound               LoginErrorCode = 1
	LoginErrorCodePasswordVerificationFailed LoginErrorCode = 2
	LoginErrorCodeInvalidCredentials         LoginErrorCode = 3
	LoginErrorCodeAppNotFound                LoginErrorCode = 4
	LoginErrorCodeAppVerificationFailed      LoginErrorCode = 5
	LoginErrorCodeInvalidDevice              LoginErrorCode = 6
	LoginErrorCodeAgreementsNotSigned        LoginErrorCode = 7
	LoginErrorCodeUnknownError               LoginErrorCode = 8
	LoginErrorCodeApiSubscriptionNotFound    LoginErrorCode = 9
)

type LogoutErrorCode int

const (
	LogoutErrorCodeSuccess        LogoutErrorCode = 0
	LogoutErrorCodeInvalidSession LogoutErrorCode = 1
	LogoutErrorCodeUnknownError   LogoutErrorCode = 2
)

type ValidateErrorCode int

const (
	ValidateErrorCodeSuccess         ValidateErrorCode = 0
	ValidateErrorCodeInvalidSession  ValidateErrorCode = 1
	ValidateErrorCodeSessionNotFound ValidateErrorCode = 2
	ValidateErrorCodeExpiredToken    ValidateErrorCode = 3
	ValidateErrorCodeUnknownError    ValidateErrorCode = 4
)

type SearchContractErrorCode int

const (
	SearchContractErrorCodeSuccess SearchContractErrorCode = 0
)

type SearchContractByIdErrorCode int

const (
	SearchContractByIdErrorCodeSuccess          SearchContractByIdErrorCode = 0
	SearchContractByIdErrorCodeContractNotFound SearchContractByIdErrorCode = 1
)

type RetrieveBarErrorCode int

const (
	RetrieveBarErrorCodeSuccess          RetrieveBarErrorCode = 0
	RetrieveBarErrorCodeContractNotFound RetrieveBarErrorCode = 1
)

type AggregateBarUnit int

const (
	AggregateBarUnitUnspecified AggregateBarUnit = 0
	AggregateBarUnitSecond      AggregateBarUnit = 1
	AggregateBarUnitMinute      AggregateBarUnit = 2
	AggregateBarUnitHour        AggregateBarUnit = 3
	AggregateBarUnitDay         AggregateBarUnit = 4
	AggregateBarUnitWeek        AggregateBarUnit = 5
	AggregateBarUnitMonth       AggregateBarUnit = 6
)

type SearchOrderErrorCode int

const (
	SearchOrderErrorCodeSuccess         SearchOrderErrorCode = 0
	SearchOrderErrorCodeAccountNotFound SearchOrderErrorCode = 1
)

type OrderStatus int

const (
	OrderStatusNone      OrderStatus = 0
	OrderStatusOpen      OrderStatus = 1
	OrderStatusFilled    OrderStatus = 2
	OrderStatusCancelled OrderStatus = 3
	OrderStatusExpired   OrderStatus = 4
	OrderStatusRejected  OrderStatus = 5
	OrderStatusPending   OrderStatus = 6
)

func (o OrderStatus) String() string {
	switch o {
	case OrderStatusNone:
		return "NONE"
	case OrderStatusOpen:
		return "OPEN"
	case OrderStatusFilled:
		return "FILLED"
	case OrderStatusCancelled:
		return "CANCELLED"
	case OrderStatusExpired:
		return "EXPIRED"
	case OrderStatusRejected:
		return "REJECTED"
	case OrderStatusPending:
		return "PENDING"
	default:
		return "UNKNOWN"
	}
}

type OrderType int

const (
	OrderTypeUnknown      OrderType = 0
	OrderTypeLimit        OrderType = 1
	OrderTypeMarket       OrderType = 2
	OrderTypeStopLimit    OrderType = 3
	OrderTypeStop         OrderType = 4
	OrderTypeTrailingStop OrderType = 5
	OrderTypeJoinBid      OrderType = 6
	OrderTypeJoinAsk      OrderType = 7
)

func (o OrderType) String() string {
	switch o {
	case OrderTypeUnknown:
		return "UNKNOWN"
	case OrderTypeLimit:
		return "LIMIT"
	case OrderTypeMarket:
		return "MARKET"
	case OrderTypeStopLimit:
		return "STOP_LIMIT"
	case OrderTypeStop:
		return "STOP"
	case OrderTypeTrailingStop:
		return "TRAILING_STOP"
	case OrderTypeJoinBid:
		return "JOIN_BID"
	case OrderTypeJoinAsk:
		return "JOIN_ASK"
	default:
		return "UNKNOWN"
	}
}

type OrderSide int

const (
	OrderSideBid OrderSide = 0
	OrderSideAsk OrderSide = 1
)

func (o OrderSide) String() string {
	switch o {
	case OrderSideBid:
		return "BUY"
	case OrderSideAsk:
		return "SELL"
	default:
		return "UNKNOWN"
	}
}

type PlaceOrderErrorCode int

const (
	PlaceOrderErrorCodeSuccess             PlaceOrderErrorCode = 0
	PlaceOrderErrorCodeAccountNotFound     PlaceOrderErrorCode = 1
	PlaceOrderErrorCodeOrderRejected       PlaceOrderErrorCode = 2
	PlaceOrderErrorCodeInsufficientFunds   PlaceOrderErrorCode = 3
	PlaceOrderErrorCodeAccountViolation    PlaceOrderErrorCode = 4
	PlaceOrderErrorCodeOutsideTradingHours PlaceOrderErrorCode = 5
	PlaceOrderErrorCodeOrderPending        PlaceOrderErrorCode = 6
	PlaceOrderErrorCodeUnknownError        PlaceOrderErrorCode = 7
	PlaceOrderErrorCodeContractNotFound    PlaceOrderErrorCode = 8
	PlaceOrderErrorCodeContractNotActive   PlaceOrderErrorCode = 9
	PlaceOrderErrorCodeAccountRejected     PlaceOrderErrorCode = 10
)

type CancelOrderErrorCode int

const (
	CancelOrderErrorCodeSuccess         CancelOrderErrorCode = 0
	CancelOrderErrorCodeAccountNotFound CancelOrderErrorCode = 1
	CancelOrderErrorCodeOrderNotFound   CancelOrderErrorCode = 2
	CancelOrderErrorCodeRejected        CancelOrderErrorCode = 3
	CancelOrderErrorCodePending         CancelOrderErrorCode = 4
	CancelOrderErrorCodeUnknownError    CancelOrderErrorCode = 5
	CancelOrderErrorCodeAccountRejected CancelOrderErrorCode = 6
)

type ModifyOrderErrorCode int

const (
	ModifyOrderErrorCodeSuccess          ModifyOrderErrorCode = 0
	ModifyOrderErrorCodeAccountNotFound  ModifyOrderErrorCode = 1
	ModifyOrderErrorCodeOrderNotFound    ModifyOrderErrorCode = 2
	ModifyOrderErrorCodeRejected         ModifyOrderErrorCode = 3
	ModifyOrderErrorCodePending          ModifyOrderErrorCode = 4
	ModifyOrderErrorCodeUnknownError     ModifyOrderErrorCode = 5
	ModifyOrderErrorCodeAccountRejected  ModifyOrderErrorCode = 6
	ModifyOrderErrorCodeContractNotFound ModifyOrderErrorCode = 7
)

type SearchPositionErrorCode int

const (
	SearchPositionErrorCodeSuccess         SearchPositionErrorCode = 0
	SearchPositionErrorCodeAccountNotFound SearchPositionErrorCode = 1
)

type PositionType int

const (
	PositionTypeFlat  PositionType = 0
	PositionTypeLong  PositionType = 1
	PositionTypeShort PositionType = 2
)

func (p PositionType) String() string {
	switch p {
	case PositionTypeFlat:
		return "FLAT"
	case PositionTypeLong:
		return "LONG"
	case PositionTypeShort:
		return "SHORT"
	default:
		return "UNKNOWN"
	}
}

type ClosePositionErrorCode int

const (
	ClosePositionErrorCodeSuccess           ClosePositionErrorCode = 0
	ClosePositionErrorCodeAccountNotFound   ClosePositionErrorCode = 1
	ClosePositionErrorCodePositionNotFound  ClosePositionErrorCode = 2
	ClosePositionErrorCodeContractNotFound  ClosePositionErrorCode = 3
	ClosePositionErrorCodeContractNotActive ClosePositionErrorCode = 4
	ClosePositionErrorCodeOrderRejected     ClosePositionErrorCode = 5
	ClosePositionErrorCodeOrderPending      ClosePositionErrorCode = 6
	ClosePositionErrorCodeUnknownError      ClosePositionErrorCode = 7
	ClosePositionErrorCodeAccountRejected   ClosePositionErrorCode = 8
)

type PartialClosePositionErrorCode int

const (
	PartialClosePositionErrorCodeSuccess           PartialClosePositionErrorCode = 0
	PartialClosePositionErrorCodeAccountNotFound   PartialClosePositionErrorCode = 1
	PartialClosePositionErrorCodePositionNotFound  PartialClosePositionErrorCode = 2
	PartialClosePositionErrorCodeContractNotFound  PartialClosePositionErrorCode = 3
	PartialClosePositionErrorCodeContractNotActive PartialClosePositionErrorCode = 4
	PartialClosePositionErrorCodeInvalidCloseSize  PartialClosePositionErrorCode = 5
	PartialClosePositionErrorCodeOrderRejected     PartialClosePositionErrorCode = 6
	PartialClosePositionErrorCodeOrderPending      PartialClosePositionErrorCode = 7
	PartialClosePositionErrorCodeUnknownError      PartialClosePositionErrorCode = 8
	PartialClosePositionErrorCodeAccountRejected   PartialClosePositionErrorCode = 9
)

type SearchTradeErrorCode int

const (
	SearchTradeErrorCodeSuccess         SearchTradeErrorCode = 0
	SearchTradeErrorCodeAccountNotFound SearchTradeErrorCode = 1
)

type SearchAccountRequest struct {
	OnlyActiveAccounts bool `json:"onlyActiveAccounts"`
}

type LoginAppRequest struct {
	UserName  string `json:"userName"`
	Password  string `json:"password"`
	DeviceID  string `json:"deviceId"`
	AppID     string `json:"appId"`
	VerifyKey string `json:"verifyKey"`
}

type LoginApiKeyRequest struct {
	UserName string `json:"userName"`
	APIKey   string `json:"apiKey"`
}

type SearchContractRequest struct {
	SearchText *string `json:"searchText,omitempty"`
	Live       bool    `json:"live"`
}

type SearchContractByIdRequest struct {
	ContractID string `json:"contractId"`
}

type RetrieveBarRequest struct {
	ContractID        string           `json:"contractId"`
	Live              bool             `json:"live"`
	StartTime         time.Time        `json:"startTime"`
	EndTime           time.Time        `json:"endTime"`
	Unit              AggregateBarUnit `json:"unit"`
	UnitNumber        int32            `json:"unitNumber"`
	Limit             int32            `json:"limit"`
	IncludePartialBar bool             `json:"includePartialBar"`
}

type SearchOrderRequest struct {
	AccountID      int32      `json:"accountId"`
	StartTimestamp time.Time  `json:"startTimestamp"`
	EndTimestamp   *time.Time `json:"endTimestamp,omitempty"`
}

type SearchOpenOrderRequest struct {
	AccountID int32 `json:"accountId"`
}

type PlaceOrderRequest struct {
	AccountID     int32     `json:"accountId"`
	ContractID    string    `json:"contractId"`
	Type          OrderType `json:"type"`
	Side          OrderSide `json:"side"`
	Size          int32     `json:"size"`
	LimitPrice    *float64  `json:"limitPrice,omitempty"`
	StopPrice     *float64  `json:"stopPrice,omitempty"`
	TrailPrice    *float64  `json:"trailPrice,omitempty"`
	CustomTag     *string   `json:"customTag,omitempty"`
	LinkedOrderID *int32    `json:"linkedOrderId,omitempty"`
}

type CancelOrderRequest struct {
	AccountID int32 `json:"accountId"`
	OrderID   int32 `json:"orderId"`
}

type ModifyOrderRequest struct {
	AccountID  int32    `json:"accountId"`
	OrderID    int32    `json:"orderId"`
	Size       *int32   `json:"size,omitempty"`
	LimitPrice *float64 `json:"limitPrice,omitempty"`
	StopPrice  *float64 `json:"stopPrice,omitempty"`
	TrailPrice *float64 `json:"trailPrice,omitempty"`
}

type SearchPositionRequest struct {
	AccountID int32 `json:"accountId"`
}

type CloseContractPositionRequest struct {
	AccountID  int32  `json:"accountId"`
	ContractID string `json:"contractId"`
}

type PartialCloseContractPositionRequest struct {
	AccountID  int32  `json:"accountId"`
	ContractID string `json:"contractId"`
	Size       int32  `json:"size"`
}

type SearchTradeRequest struct {
	AccountID      int32      `json:"accountId"`
	StartTimestamp *time.Time `json:"startTimestamp,omitempty"`
	EndTimestamp   *time.Time `json:"endTimestamp,omitempty"`
}

type SearchAccountResponse struct {
	Success      bool                   `json:"success"`
	ErrorCode    SearchAccountErrorCode `json:"errorCode"`
	ErrorMessage *string                `json:"errorMessage,omitempty"`
	Accounts     []TradingAccountModel  `json:"accounts,omitempty"`
}

type LoginResponse struct {
	Success      bool           `json:"success"`
	ErrorCode    LoginErrorCode `json:"errorCode"`
	ErrorMessage *string        `json:"errorMessage,omitempty"`
	Token        *string        `json:"token,omitempty"`
}

type LogoutResponse struct {
	Success      bool            `json:"success"`
	ErrorCode    LogoutErrorCode `json:"errorCode"`
	ErrorMessage *string         `json:"errorMessage,omitempty"`
}

type ValidateResponse struct {
	Success      bool              `json:"success"`
	ErrorCode    ValidateErrorCode `json:"errorCode"`
	ErrorMessage *string           `json:"errorMessage,omitempty"`
	NewToken     *string           `json:"newToken,omitempty"`
}

type SearchContractResponse struct {
	Success      bool                    `json:"success"`
	ErrorCode    SearchContractErrorCode `json:"errorCode"`
	ErrorMessage *string                 `json:"errorMessage,omitempty"`
	Contracts    []ContractModel         `json:"contracts,omitempty"`
}

type SearchContractByIdResponse struct {
	Success      bool                        `json:"success"`
	ErrorCode    SearchContractByIdErrorCode `json:"errorCode"`
	ErrorMessage *string                     `json:"errorMessage,omitempty"`
	Contract     *ContractModel              `json:"contract,omitempty"`
}

type RetrieveBarResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    RetrieveBarErrorCode `json:"errorCode"`
	ErrorMessage *string              `json:"errorMessage,omitempty"`
	Bars         []AggregateBarModel  `json:"bars,omitempty"`
}

type SearchOrderResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    SearchOrderErrorCode `json:"errorCode"`
	ErrorMessage *string              `json:"errorMessage,omitempty"`
	Orders       []OrderModel         `json:"orders,omitempty"`
}

type PlaceOrderResponse struct {
	Success      bool                `json:"success"`
	ErrorCode    PlaceOrderErrorCode `json:"errorCode"`
	ErrorMessage *string             `json:"errorMessage,omitempty"`
	OrderID      *int32              `json:"orderId,omitempty"`
}

type CancelOrderResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    CancelOrderErrorCode `json:"errorCode"`
	ErrorMessage *string              `json:"errorMessage,omitempty"`
}

type ModifyOrderResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    ModifyOrderErrorCode `json:"errorCode"`
	ErrorMessage *string              `json:"errorMessage,omitempty"`
}

type SearchPositionResponse struct {
	Success      bool                    `json:"success"`
	ErrorCode    SearchPositionErrorCode `json:"errorCode"`
	ErrorMessage *string                 `json:"errorMessage,omitempty"`
	Positions    []PositionModel         `json:"positions,omitempty"`
}

type ClosePositionResponse struct {
	Success      bool                   `json:"success"`
	ErrorCode    ClosePositionErrorCode `json:"errorCode"`
	ErrorMessage *string                `json:"errorMessage,omitempty"`
}

type PartialClosePositionResponse struct {
	Success      bool                          `json:"success"`
	ErrorCode    PartialClosePositionErrorCode `json:"errorCode"`
	ErrorMessage *string                       `json:"errorMessage,omitempty"`
}

type SearchHalfTradeResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    SearchTradeErrorCode `json:"errorCode"`
	ErrorMessage *string              `json:"errorMessage,omitempty"`
	Trades       []HalfTradeModel     `json:"trades,omitempty"`
}

type TradingAccountModel struct {
	ID        int32   `json:"id"`
	Name      string  `json:"name"`
	Balance   float64 `json:"balance"`
	CanTrade  bool    `json:"canTrade"`
	IsVisible bool    `json:"isVisible"`
}

type ContractModel struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	TickSize       float64 `json:"tickSize"`
	TickValue      float64 `json:"tickValue"`
	ActiveContract bool    `json:"activeContract"`
}

type AggregateBarModel struct {
	T time.Time `json:"t"`
	O float64   `json:"o"`
	H float64   `json:"h"`
	L float64   `json:"l"`
	C float64   `json:"c"`
	V int64     `json:"v"`
}

type OrderModel struct {
	ID                int32       `json:"id"`
	AccountID         int32       `json:"accountId"`
	ContractID        string      `json:"contractId"`
	CreationTimestamp time.Time   `json:"creationTimestamp"`
	UpdateTimestamp   *time.Time  `json:"updateTimestamp,omitempty"`
	Status            OrderStatus `json:"status"`
	Type              OrderType   `json:"type"`
	Side              OrderSide   `json:"side"`
	Size              int32       `json:"size"`
	LimitPrice        *float64    `json:"limitPrice,omitempty"`
	StopPrice         *float64    `json:"stopPrice,omitempty"`
	FillVolume        int32       `json:"fillVolume"`
}

type OrderUpdateData struct {
	Action int                `json:"action"`
	Data   OrderUpdatePayload `json:"data"`
}

type OrderUpdatePayload struct {
	AccountID         int32       `json:"accountId"`
	ContractID        string      `json:"contractId"`
	CreationTimestamp time.Time   `json:"creationTimestamp"`
	FillVolume        int32       `json:"fillVolume"`
	ID                int32       `json:"id"`
	LimitPrice        float64     `json:"limitPrice"`
	Side              OrderSide   `json:"side"`
	Size              int32       `json:"size"`
	Status            OrderStatus `json:"status"`
	Type              OrderType   `json:"type"`
	UpdateTimestamp   time.Time   `json:"updateTimestamp"`
}

func (o *OrderUpdatePayload) UnmarshalJSON(data []byte) error {

	type rawOrderUpdatePayload struct {
		AccountID         float64 `json:"accountId"`
		ContractID        string  `json:"contractId"`
		CreationTimestamp string  `json:"creationTimestamp"`
		FillVolume        int32   `json:"fillVolume"`
		ID                float64 `json:"id"`
		LimitPrice        float64 `json:"limitPrice"`
		Side              int32   `json:"side"`
		Size              int32   `json:"size"`
		Status            int32   `json:"status"`
		Type              int32   `json:"type"`
		UpdateTimestamp   string  `json:"updateTimestamp"`
	}

	var raw rawOrderUpdatePayload
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	o.AccountID = int32(raw.AccountID)
	o.ContractID = raw.ContractID
	o.FillVolume = raw.FillVolume
	o.ID = int32(raw.ID)
	o.LimitPrice = raw.LimitPrice
	o.Side = OrderSide(raw.Side)
	o.Size = raw.Size
	o.Status = OrderStatus(raw.Status)
	o.Type = OrderType(raw.Type)

	var err error
	o.CreationTimestamp, err = parseTimestamp(raw.CreationTimestamp)
	if err != nil {
		return fmt.Errorf("failed to parse CreationTimestamp: %w", err)
	}

	o.UpdateTimestamp, err = parseTimestamp(raw.UpdateTimestamp)
	if err != nil {
		return fmt.Errorf("failed to parse UpdateTimestamp: %w", err)
	}

	return nil
}

func parseTimestamp(timestamp string) (time.Time, error) {
	if timestamp == "" {
		return time.Time{}, nil
	}

	t, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {

		t, err = time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return time.Time{}, fmt.Errorf("unable to parse timestamp %q: %w", timestamp, err)
		}
	}

	return t, nil
}

type PositionModel struct {
	ID                int32        `json:"id"`
	AccountID         int32        `json:"accountId"`
	ContractID        string       `json:"contractId"`
	CreationTimestamp time.Time    `json:"creationTimestamp"`
	Type              PositionType `json:"type"`
	Size              int32        `json:"size"`
	AveragePrice      float64      `json:"averagePrice"`
}

type PositionUpdateData struct {
	Action int                   `json:"action"`
	Data   PositionUpdatePayload `json:"data"`
}

type PositionUpdatePayload struct {
	AccountID         int32        `json:"accountId"`
	AveragePrice      float64      `json:"averagePrice"`
	ContractID        string       `json:"contractId"`
	CreationTimestamp time.Time    `json:"creationTimestamp"`
	ID                int32        `json:"id"`
	Size              int32        `json:"size"`
	Type              PositionType `json:"type"`
}

func (p *PositionUpdatePayload) UnmarshalJSON(data []byte) error {
	type rawPositionUpdatePayload struct {
		AccountID         float64 `json:"accountId"`
		AveragePrice      float64 `json:"averagePrice"`
		ContractID        string  `json:"contractId"`
		CreationTimestamp string  `json:"creationTimestamp"`
		ID                float64 `json:"id"`
		Size              int32   `json:"size"`
		Type              int32   `json:"type"`
	}

	var raw rawPositionUpdatePayload
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	p.AccountID = int32(raw.AccountID)
	p.AveragePrice = raw.AveragePrice
	p.ContractID = raw.ContractID
	p.ID = int32(raw.ID)
	p.Size = raw.Size
	p.Type = PositionType(raw.Type)

	var err error
	p.CreationTimestamp, err = parseTimestamp(raw.CreationTimestamp)
	if err != nil {
		return fmt.Errorf("failed to parse CreationTimestamp: %w", err)
	}

	return nil
}

type HalfTradeModel struct {
	ID                int32     `json:"id"`
	AccountID         int32     `json:"accountId"`
	ContractID        string    `json:"contractId"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
	Price             float64   `json:"price"`
	ProfitAndLoss     *float64  `json:"profitAndLoss,omitempty"`
	Fees              float64   `json:"fees"`
	Side              OrderSide `json:"side"`
	Size              int32     `json:"size"`
	Voided            bool      `json:"voided"`
	OrderID           int32     `json:"orderId"`
}
