package data

import "time"

type Candle struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int
}

type Price struct {
	Price  float64
	Volume int
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
	Symbol        string
	PercentChange float64
}

type TopGainer struct {
	Gainer
	Yesterday Candle
	Current   Price
}
