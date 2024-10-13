package fmp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	apiURL string
	apiKey string

	client *http.Client
}

func NewClient(apiURL string, apiKey string) *Client {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	return &Client{
		apiURL: apiURL,
		apiKey: apiKey,
		client: client,
	}
}

func (s *Client) GetAvailableSymbolsByExchange(exchange string) ([]StockSymbol, error) {
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

	var sList []StockSymbol
	err = json.Unmarshal(responseBody, &sList)
	if err != nil {
		return nil, err
	}

	return sList, nil
}
