package stocks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dmihai/stocks/pkg/data"
)

type polygon struct {
	apiURL string
	apiKey string

	client *http.Client
}

func NewPolygonClient(apiURL string, apiKey string) Client {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	return &fmp{
		apiURL: apiURL,
		apiKey: apiKey,
		client: client,
	}
}

func (p *polygon) GetAvailableSymbolsByExchange(exchange string) ([]data.Symbol, error) {
	return nil, errors.New("unimplemented")
}

func (p *polygon) GetAllRealtimePrices() ([]data.Intraday, error) {
	return nil, errors.New("unimplemented")
}

func (p *polygon) GetSymbolDetails(symbol string) (*SymbolDetails, error) {
	url := fmt.Sprintf("%s/v3/reference/tickers/%s", p.apiURL, symbol)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))

	response, err := p.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var details PolygonResponse[PolygonSymbolDetails]
	err = json.Unmarshal(responseBody, &details)
	if err != nil {
		return nil, err
	}

	return &SymbolDetails{
		Symbol:            details.Results.Symbol,
		Name:              details.Results.Name,
		IpoDate:           details.Results.IpoDate,
		SharesOutstanding: details.Results.SharesOutstanding,
	}, nil
}
