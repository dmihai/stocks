package stocks

import (
	"encoding/json"
	"time"
)

const (
	ExchangeNASDAQ = "NASDAQ"
	ExchangeNYSE   = "NYSE"
	ExchangeAMEX   = "AMEX"
)

type StockSymbol struct {
	Symbol            string  `json:"symbol"`
	Name              string  `json:"name"`
	Price             float64 `json:"price"`
	DayHigh           float64 `json:"dayHigh"`
	DayLow            float64 `json:"dayLow"`
	YearHigh          float64 `json:"yearHigh"`
	YearLow           float64 `json:"yearLow"`
	MarketCap         int     `json:"marketCap"`
	Volume            int     `json:"volume"`
	AvgVolume         int     `json:"avgVolume"`
	Open              float64 `json:"open"`
	PreviousClose     float64 `json:"previousClose"`
	Exchange          string  `json:"exchange"`
	SharesOutstanding int     `json:"sharesOutstanding"`
}

type StockPrice struct {
	Symbol      string    `json:"symbol"`
	Price       float64   `json:"fmpLast"`
	Volume      float64   `json:"volume"`
	LastUpdated Timestamp `json:"lastUpdated"`
}

type Timestamp struct {
	time.Time
}

func (p *Timestamp) UnmarshalJSON(bytes []byte) error {
	var raw int64

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	p.Time = time.Unix(raw/1000, raw%1000)

	return nil
}
