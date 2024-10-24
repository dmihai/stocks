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

const (
	polygonStatusOK       = "OK"
	polygonStatusNotFound = "NOT_FOUND"
)

type polygon struct {
	apiURL string
	apiKey string

	client *http.Client
}

func NewPolygonClient(apiURL string, apiKey string) *polygon {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	return &polygon{
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
	endpoint := fmt.Sprintf("v3/reference/tickers/%s", symbol)

	details, err := getPolygonResponse[PolygonSymbolDetails](p, endpoint, nil)
	if err != nil {
		return nil, err
	}

	// TODO not found response
	if details == nil {
		return nil, nil
	}

	return &SymbolDetails{
		Symbol:            details.Symbol,
		Name:              details.Name,
		IpoDate:           details.IpoDate,
		SharesOutstanding: details.SharesOutstanding,
	}, nil
}

func getPolygonResponse[T any](p *polygon, endpoint string, params map[string]string) (*T, error) {
	url := fmt.Sprintf("%s/%s", p.apiURL, endpoint)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if params != nil {
		query := request.URL.Query()
		for param, value := range params {
			query.Add(param, value)
		}

		request.URL.RawQuery = query.Encode()
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

	var payload PolygonResponse[T]
	err = json.Unmarshal(responseBody, &payload)
	if err != nil {
		return nil, err
	}

	switch payload.Status {
	case polygonStatusOK:
		return &payload.Results, nil
	case polygonStatusNotFound:
		return nil, nil
	default:
		return nil, fmt.Errorf("invalid API status")
	}
}
