package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"pwt_go/playwrightsetup"
	"testing"

	"github.com/stretchr/testify/require"
)

// Função utilitária para fazer requisições para a API Alpha Vantage
func makeRequest(t *testing.T, endpoint string) map[string]interface{} {
	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	require.NotEmpty(t, apiKey, "ALPHA_VANTAGE_API_KEY não definida")

	pm, err := playwrightsetup.NewPlaywrightManager()
	require.NoError(t, err, "Failed to start Playwright")
	defer pm.Close()

	requestContext, err := pm.Playwright.Request.NewContext()
	require.NoError(t, err, "Failed to create request context")

	response, err := requestContext.Get(fmt.Sprintf("https://www.alphavantage.co/query?%s&apikey=%s", endpoint, apiKey))
	require.NoError(t, err, "Failed to fetch API")
	require.Equal(t, http.StatusOK, response.Status(), "Expected status 200")

	responseText, err := response.Text()
	require.NoError(t, err, "Failed to read response text")

	fmt.Println("Response Body:", responseText)

	var body map[string]interface{}
	err = json.Unmarshal([]byte(responseText), &body)
	require.NoError(t, err, "Failed to parse JSON response")

	return body
}

// Testa dados intraday de uma ação
func TestIntradayStockData(t *testing.T) {
	symbol := "IBM"
	interval := "5min"
	endpoint := fmt.Sprintf("function=TIME_SERIES_INTRADAY&symbol=%s&interval=%s", symbol, interval)
	body := makeRequest(t, endpoint)

	metaData, ok := body["Meta Data"].(map[string]interface{})
	require.True(t, ok, "Meta Data not found in the response")
	require.Equal(t, symbol, metaData["2. Symbol"], "Expected correct symbol in Meta Data")
	require.Equal(t, interval, metaData["4. Interval"], "Expected correct interval in Meta Data")

	timeSeries, ok := body["Time Series (5min)"].(map[string]interface{})
	require.True(t, ok, "Time Series (5min) not found in the response")

	for time, data := range timeSeries {
		dataMap, ok := data.(map[string]interface{})
		require.True(t, ok, fmt.Sprintf("Expected map for time series data at %s", time))
		require.Contains(t, dataMap, "1. open", fmt.Sprintf("Expected open price at %s", time))
		require.Contains(t, dataMap, "4. close", fmt.Sprintf("Expected close price at %s", time))
		break // Verifica apenas o primeiro registro
	}
}

// Testa a taxa de câmbio entre duas moedas
func TestExchangeRate(t *testing.T) {
	fromCurrency := "USD"
	toCurrency := "EUR"
	endpoint := fmt.Sprintf("function=CURRENCY_EXCHANGE_RATE&from_currency=%s&to_currency=%s", fromCurrency, toCurrency)
	body := makeRequest(t, endpoint)

	exchangeRate, ok := body["Realtime Currency Exchange Rate"].(map[string]interface{})
	require.True(t, ok, "Realtime Currency Exchange Rate not found in the response")

	require.Equal(t, fromCurrency, exchangeRate["1. From_Currency Code"], "Expected correct from currency")
	require.Equal(t, toCurrency, exchangeRate["3. To_Currency Code"], "Expected correct to currency")
}

// Testa um endpoint inválido
func TestInvalidEndpoint(t *testing.T) {
	endpoint := "function=INVALID_ENDPOINT"
	body := makeRequest(t, endpoint)

	errorMessage, ok := body["Error Message"].(string)
	require.True(t, ok, "Expected error message in the response")
	require.Equal(t, "This API function (INVALID_ENDPOINT) does not exist.", errorMessage, "Expected error message for invalid endpoint")
}

// Testa a cotação global de uma ação
func TestGlobalQuote(t *testing.T) {
	symbol := "IBM"
	endpoint := fmt.Sprintf("function=GLOBAL_QUOTE&symbol=%s", symbol)
	body := makeRequest(t, endpoint)

	globalQuote, ok := body["Global Quote"].(map[string]interface{})
	require.True(t, ok, "Global Quote not found in the response")

	require.Contains(t, globalQuote, "01. symbol", "Expected symbol in Global Quote")
	require.Equal(t, symbol, globalQuote["01. symbol"], "Expected correct symbol")
}
