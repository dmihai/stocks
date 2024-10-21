package stocks

import (
	"golang.org/x/sync/errgroup"

	"github.com/dmihai/stocks/pkg/data"
)

type hybrid struct {
	fmp     *fmp
	polygon *polygon
}

func NewHybridClient(fmp *fmp, polygon *polygon) *hybrid {
	return &hybrid{
		fmp:     fmp,
		polygon: polygon,
	}
}

func (h *hybrid) GetAvailableSymbolsByExchange(exchange string) ([]data.Symbol, error) {
	return h.fmp.GetAvailableSymbolsByExchange(exchange)
}

func (h *hybrid) GetAllRealtimePrices() ([]data.Intraday, error) {
	return h.fmp.GetAllRealtimePrices()
}

func (h *hybrid) GetSymbolDetails(symbol string) (*SymbolDetails, error) {
	g := errgroup.Group{}

	var fmpSymbol *SymbolDetails
	g.Go(func() error {
		var err error
		fmpSymbol, err = h.fmp.GetSymbolDetails(symbol)
		return err
	})

	var polygonSymbol *SymbolDetails
	g.Go(func() error {
		var err error
		polygonSymbol, err = h.polygon.GetSymbolDetails(symbol)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &SymbolDetails{
		Symbol:            fmpSymbol.Symbol,
		Name:              fmpSymbol.Name,
		Industry:          fmpSymbol.Industry,
		Sector:            fmpSymbol.Sector,
		IpoDate:           fmpSymbol.IpoDate,
		SharesOutstanding: polygonSymbol.SharesOutstanding,
	}, nil
}
