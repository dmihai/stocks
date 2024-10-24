package stocks

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	apiKey = "secret-api-key"
)

func TestFMPGetAvailableSymbolsByExchange(t *testing.T) {
	t.Run("successfully fetch symbols", func(t *testing.T) {
		exchange := "NASDAQ"
		endpoint := "/v3/symbol/"

		expected := []map[string]any{
			{
				"symbol":            "AAPL",
				"name":              "Apple Inc.",
				"price":             230.35,
				"marketCap":         3502264435000,
				"exchange":          "NASDAQ",
				"volume":            17202113,
				"avgVolume":         50834700,
				"open":              229.98,
				"previousClose":     230.76,
				"sharesOutstanding": 15204100000,
				"timestamp":         1729792519,
			},
			{
				"symbol":            "NFLX",
				"name":              "Netflix, Inc.",
				"price":             752.81,
				"marketCap":         321794656980,
				"exchange":          "NASDAQ",
				"volume":            1366645,
				"avgVolume":         3214039,
				"open":              751.97,
				"previousClose":     749.29,
				"sharesOutstanding": 427458000,
				"timestamp":         1729792507,
			},
			{
				"symbol":            "AMZN",
				"name":              "Amazon.com, Inc.",
				"price":             186.5474,
				"marketCap":         1957926891440,
				"exchange":          "NASDAQ",
				"volume":            12742507,
				"avgVolume":         37979045,
				"open":              185.25,
				"previousClose":     184.71,
				"sharesOutstanding": 10495600000,
				"timestamp":         1729792524,
			},
		}

		handler := func(w http.ResponseWriter, r *http.Request) {
			query := strings.ReplaceAll(r.URL.Path, endpoint, "")

			assert.Equal(t, exchange, query)
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, apiKey, r.URL.Query().Get("apikey"))

			w.WriteHeader(http.StatusOK)
			encoded, _ := json.Marshal(expected)
			_, _ = w.Write(encoded)
		}

		server := newTestServer(endpoint, handler)
		defer server.Close()

		client := NewFMPClient(server.URL, apiKey)

		actual, err := client.GetAvailableSymbolsByExchange(exchange)
		require.NoError(t, err)

		require.Len(t, actual, len(expected))
		for i := 0; i < len(actual); i++ {
			require.Equal(t, expected[0]["symbol"], actual[0].Name)
			require.Equal(t, expected[0]["previousClose"], actual[0].PrevDayClose)
			require.Equal(t, expected[0]["sharesOutstanding"], actual[0].Shares)
		}
	})
}

func TestFMPGetAllRealtimePrices(t *testing.T) {
	t.Run("successfully fetch realtime prices", func(t *testing.T) {
		timestamp := time.Now()
		endpoint := "/v3/stock/full/real-time-price"

		expected := []map[string]any{
			{
				"symbol":      "AAPL",
				"fmpLast":     12.345,
				"volume":      100.2,
				"lastUpdated": timestamp.UnixMilli(),
			},
			{
				"symbol":  "NFLX",
				"fmpLast": 1.45,
				"volume":  123456,
			},
		}

		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, apiKey, r.URL.Query().Get("apikey"))

			w.WriteHeader(http.StatusOK)
			encoded, _ := json.Marshal(expected)
			_, _ = w.Write(encoded)
		}

		server := newTestServer(endpoint, handler)
		defer server.Close()

		client := NewFMPClient(server.URL, apiKey)

		actual, err := client.GetAllRealtimePrices()
		require.NoError(t, err)

		require.Len(t, actual, 2)
		require.Equal(t, expected[0]["symbol"], actual[0].Symbol)
		require.Equal(t, expected[0]["fmpLast"], actual[0].Price.Price)
		require.Equal(t, int(expected[0]["volume"].(float64)), actual[0].Price.Volume)
		require.Equal(t, timestamp.Unix(), actual[0].Timestamp.Unix())
		require.Equal(t, expected[1]["symbol"], actual[1].Symbol)
		require.Equal(t, expected[1]["fmpLast"], actual[1].Price.Price)
		require.Equal(t, expected[1]["volume"], actual[1].Price.Volume)
	})
}

func TestFMPGetSymbolDetails(t *testing.T) {
	t.Run("successfully fetch symbol details", func(t *testing.T) {
		symbol := "AAPL"
		endpoint := "/v3/profile/"

		expected := []map[string]string{
			{
				"symbol":      symbol,
				"companyName": "Apple Inc.",
				"industry":    "Consumer Electronics",
				"sector":      "Technology",
				"ipoDate":     "1980-12-12",
			},
		}

		handler := func(w http.ResponseWriter, r *http.Request) {
			query := strings.ReplaceAll(r.URL.Path, endpoint, "")

			assert.Equal(t, symbol, query)
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, apiKey, r.URL.Query().Get("apikey"))

			w.WriteHeader(http.StatusOK)
			encoded, _ := json.Marshal(expected)
			_, _ = w.Write(encoded)
		}

		server := newTestServer(endpoint, handler)
		defer server.Close()

		client := NewFMPClient(server.URL, apiKey)

		actual, err := client.GetSymbolDetails(symbol)
		require.NoError(t, err)

		require.Equal(t, expected[0]["symbol"], actual.Symbol)
		require.Equal(t, expected[0]["companyName"], actual.Name)
		require.Equal(t, expected[0]["industry"], actual.Industry)
		require.Equal(t, expected[0]["sector"], actual.Sector)
		require.Equal(t, expected[0]["ipoDate"], actual.IpoDate)
	})
}

func newTestServer(endpoint string, handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(endpoint, handler)

	return httptest.NewServer(mux)
}
