package stocks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dmihai/stocks/pkg/data"
)

type fmp struct {
	apiURL string
	apiKey string

	client *http.Client
}

func NewFMPClient(apiURL string, apiKey string) *fmp {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	return &fmp{
		apiURL: apiURL,
		apiKey: apiKey,
		client: client,
	}
}

func (f *fmp) GetAvailableSymbolsByExchange(exchange string) ([]data.Symbol, error) {
	endpoint := fmt.Sprintf("v3/symbol/%s", exchange)

	symbolList, err := getFMPResponse[[]FMPStockSymbol](f, endpoint)
	if err != nil {
		return nil, err
	}

	// TODO not found response
	if symbolList == nil || len(*symbolList) == 0 {
		return nil, err
	}

	symbols := make([]data.Symbol, len(*symbolList))
	for i, symbol := range *symbolList {
		symbols[i] = data.Symbol{
			Name:         symbol.Symbol,
			PrevDayClose: symbol.PreviousClose,
			Shares:       symbol.SharesOutstanding,
		}
	}

	return symbols, nil
}

func (f *fmp) GetAllRealtimePrices() ([]data.Intraday, error) {
	endpoint := "v3/stock/full/real-time-price"

	priceList, err := getFMPResponse[[]FMPStockPrice](f, endpoint)
	if err != nil {
		return nil, err
	}

	// TODO not found response
	if priceList == nil || len(*priceList) == 0 {
		return nil, err
	}

	intradayList := make([]data.Intraday, len(*priceList))
	for i, price := range *priceList {
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

func (f *fmp) GetSymbolDetails(symbol string) (*SymbolDetails, error) {
	endpoint := fmt.Sprintf("v3/profile/%s", symbol)

	details, err := getFMPResponse[[]FMPSymbolDetails](f, endpoint)
	if err != nil {
		return nil, err
	}

	// TODO not found response
	if details == nil || len(*details) == 0 {
		return nil, err
	}

	return &SymbolDetails{
		Symbol:   (*details)[0].Symbol,
		Name:     (*details)[0].CompanyName,
		Industry: (*details)[0].Industry,
		Sector:   (*details)[0].Sector,
		IpoDate:  (*details)[0].IpoDate,
	}, nil
}

func getFMPResponse[T any](f *fmp, endpoint string) (*T, error) {
	url := fmt.Sprintf("%s/%s", f.apiURL, endpoint)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	query := request.URL.Query()
	query.Add("apikey", f.apiKey)

	request.URL.RawQuery = query.Encode()

	response, err := f.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var payload T
	err = json.Unmarshal(responseBody, &payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
