package stocks

import "github.com/dmihai/stocks/pkg/data"

type Client interface {
	GetAvailableSymbolsByExchange(exchange string) ([]data.Symbol, error)
	GetAllRealtimePrices() ([]data.Intraday, error)
	GetSymbolDetails(symbol string) (*SymbolDetails, error)
}
