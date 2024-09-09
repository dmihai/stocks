package data

type Candle struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int
}

type Gainer struct {
	Symbol        string
	PercentChange float64
}
