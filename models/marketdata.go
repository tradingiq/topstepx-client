package models

import "time"

// Quote represents market quote data
type Quote struct {
	BestAsk       float64   `json:"bestAsk"`
	BestBid       float64   `json:"bestBid"`
	Change        float64   `json:"change"`
	ChangePercent float64   `json:"changePercent"`
	LastPrice     float64   `json:"lastPrice"`
	LastUpdated   time.Time `json:"lastUpdated"`
	Symbol        string    `json:"symbol"`
	Timestamp     time.Time `json:"timestamp"`
	Volume        float64   `json:"volume"`
}

// Trade represents a single market trade
type Trade struct {
	Price     float64   `json:"price"`
	SymbolID  string    `json:"symbolId"`
	Timestamp time.Time `json:"timestamp"`
	Type      int       `json:"type"`
	Volume    float64   `json:"volume"`
}

// TradeData represents an array of trades received from websocket
type TradeData []Trade

// MarketDepth represents a single market depth entry
type MarketDepth struct {
	Price         float64   `json:"price"`
	Volume        float64   `json:"volume"`
	CurrentVolume float64   `json:"currentVolume"`
	Type          int       `json:"type"`
	Timestamp     time.Time `json:"timestamp"`
}

// MarketDepthData represents an array of market depth entries received from websocket
type MarketDepthData []MarketDepth