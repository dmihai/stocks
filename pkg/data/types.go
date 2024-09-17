package data

import "time"

type Candle struct {
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int     `json:"volume"`
}

type Price struct {
	Price  float64 `json:"price"`
	Volume int     `json:"volume"`
}

type Daily struct {
	Symbol string
	Day    string
	Candle
}

type Intraday struct {
	Symbol    string
	Timestamp time.Time
	Price
}

type Gainer struct {
	Symbol         string  `json:"symbol"`
	PercentChanged float64 `json:"percentChanged"`
	intradayIndex  int
}

type TopGainer struct {
	Gainer
	Yesterday   Candle    `json:"yesterday"`
	Current     Price     `json:"current"`
	LastUpdated time.Time `json:"lastUpdated"`
}
