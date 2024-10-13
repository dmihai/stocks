package fmp

const (
	ExchangeNASDAQ = "NASDAQ"
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
