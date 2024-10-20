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

type SymbolDetails struct {
	Symbol            string
	Name              string
	Industry          string
	Sector            string
	IpoDate           string
	SharesOutstanding int
}

type FMPStockSymbol struct {
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

type FMPStockPrice struct {
	Symbol      string        `json:"symbol"`
	Price       float64       `json:"fmpLast"`
	Volume      float64       `json:"volume"`
	LastUpdated UnixTimestamp `json:"lastUpdated"`
}

type FMPSymbolDetails struct {
	Symbol      string `json:"symbol"`
	CompanyName string `json:"companyName"`
	Industry    string `json:"industry"`
	Sector      string `json:"sector"`
	IpoDate     string `json:"ipoDate"`
}

type PolygonResponse[T any] struct {
	RequestID string `json:"request_id"`
	Results   T      `json:"results"`
	Status    string `json:"status"`
}

type PolygonSymbolDetails struct {
	Symbol            string `json:"ticker"`
	Name              string `json:"name"`
	IpoDate           string `json:"list_date"`
	SharesOutstanding int    `json:"weighted_shares_outstanding"`
}

type UnixTimestamp struct {
	time.Time
}

func (p *UnixTimestamp) UnmarshalJSON(bytes []byte) error {
	var raw int64

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	p.Time = time.Unix(raw/1000, raw%1000)

	return nil
}
