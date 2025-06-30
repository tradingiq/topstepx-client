package models

import "time"

type Quote struct {
	BestAsk       float64   `json:"bestAsk"`
	BestBid       float64   `json:"bestBid"`
	Change        float64   `json:"change"`
	ChangePercent float64   `json:"changePercent"`
	LastPrice     float64   `json:"lastPrice"`
	LastUpdated   time.Time `json:"lastUpdated"`
	Symbol        string    `json:"symbol"`
	Timestamp     time.Time `json:"timestamp"`
	Volume        int       `json:"volume"`
}

type Trade struct {
	Price     float64   `json:"price"`
	SymbolID  string    `json:"symbolId"`
	Timestamp time.Time `json:"timestamp"`
	Type      int       `json:"type"`
	Volume    int       `json:"volume"`
}

type TradeData []Trade

type MarketDepth struct {
	Price         float64   `json:"price"`
	Volume        float64   `json:"volume"`
	CurrentVolume int       `json:"currentVolume"`
	Type          int       `json:"type"`
	Timestamp     time.Time `json:"timestamp"`
}

type MarketDepthData []MarketDepth
