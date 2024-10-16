package stocks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dmihai/stocks/pkg/data"
)

type Client interface {
	GetAvailableSymbolsByExchange(exchange string) ([]data.Symbol, error)
	GetAllRealtimePrices() ([]data.Intraday, error)
}

type fmp struct {
	apiURL string
	apiKey string

	client *http.Client
}

func NewFMPClient(apiURL string, apiKey string) Client {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	return &fmp{
		apiURL: apiURL,
		apiKey: apiKey,
		client: client,
	}
}

func (s *fmp) GetAvailableSymbolsByExchange(exchange string) ([]data.Symbol, error) {
	url := fmt.Sprintf("%s/v3/symbol/%s?apikey=%s", s.apiURL, exchange, s.apiKey)

	response, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var symbolList []StockSymbol
	err = json.Unmarshal(responseBody, &symbolList)
	if err != nil {
		return nil, err
	}

	symbols := make([]data.Symbol, len(symbolList))
	for i, symbol := range symbolList {
		symbols[i] = data.Symbol{
			Name:         symbol.Symbol,
			PrevDayClose: symbol.PreviousClose,
			Shares:       symbol.SharesOutstanding,
		}
	}

	return symbols, nil
}

func (s *fmp) GetAllRealtimePrices() ([]data.Intraday, error) {
	url := fmt.Sprintf("%s/v3/stock/full/real-time-price?apikey=%s", s.apiURL, s.apiKey)

	response, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var priceList []StockPrice
	err = json.Unmarshal(responseBody, &priceList)
	if err != nil {
		return nil, err
	}

	intradayList := make([]data.Intraday, len(priceList))
	for i, price := range priceList {
		intradayList[i] = data.Intraday{
			Symbol:    price.Symbol,
			Timestamp: price.LastUpdated.Time,
			Price: data.Price{
				Price:  price.Price,
				Volume: int(price.Volume),
			},
		}
	}

	return intradayList, nil
}
