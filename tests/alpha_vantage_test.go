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
	// Busca a chave da API a partir da variável de ambiente
	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if apiKey == "" {
		t.Fatal("ALPHA_VANTAGE_API_KEY não definida")
	}

	pm, err := playwrightsetup.NewPlaywrightManager()
	require.NoError(t, err, "Failed to start Playwright")
	defer pm.Close()

	requestContext, err := pm.Playwright.Request.NewContext()
	require.NoError(t, err, "Failed to create request context")

	// Faz a requisição
	response, err := requestContext.Get(fmt.Sprintf("https://www.alphavantage.co/query?%s&apikey=%s", endpoint, apiKey))
	require.NoError(t, err, "Failed to fetch API")
	require.Equal(t, http.StatusOK, response.Status(), "Expected status 200")

	// Faz o parse da resposta
	responseText, err := response.Text()
	require.NoError(t, err, "Failed to read response text")

	// Imprime a resposta para depuração
	fmt.Println("Response Body:", responseText)

	var body map[string]interface{}
	err = json.Unmarshal([]byte(responseText), &body)
	require.NoError(t, err, "Failed to parse JSON response")

	// Verifica se existe uma mensagem de erro
	if errorMessage, exists := body["Error Message"]; exists {
		t.Fatalf("API Error: %s", errorMessage)
	}

	return body
}

// Cenário 1: Verificar dados de uma ação intraday
func TestIntradayStockData(t *testing.T) {
	symbol := "IBM"
	interval := "5min"
	body := makeRequest(t, fmt.Sprintf("function=TIME_SERIES_INTRADAY&symbol=%s&interval=%s", symbol, interval))

	// Valida se "Meta Data" está correto
	if metaData, ok := body["Meta Data"].(map[string]interface{}); ok {
		require.Equal(t, symbol, metaData["2. Symbol"], "Expected correct symbol in Meta Data")
		require.Equal(t, interval, metaData["4. Interval"], "Expected correct interval in Meta Data")
	} else {
		t.Fatalf("Meta Data not found in the response")
	}

	// Verifica se o Time Series contém preços de abertura, fechamento, etc.
	if timeSeries, ok := body["Time Series (5min)"].(map[string]interface{}); ok {
		for time, data := range timeSeries {
			require.Contains(t, data.(map[string]interface{}), "1. open", fmt.Sprintf("Expected open price at %s", time))
			require.Contains(t, data.(map[string]interface{}), "4. close", fmt.Sprintf("Expected close price at %s", time))
			break // Verifica apenas o primeiro registro
		}
	} else {
		t.Fatalf("Time Series (5min) not found in the response")
	}
}

// Cenário 2: Verificar taxa de câmbio entre duas moedas
func TestExchangeRate(t *testing.T) {
	fromCurrency := "USD"
	toCurrency := "EUR"
	body := makeRequest(t, fmt.Sprintf("function=CURRENCY_EXCHANGE_RATE&from_currency=%s&to_currency=%s", fromCurrency, toCurrency))

	if exchangeRate, ok := body["Realtime Currency Exchange Rate"].(map[string]interface{}); ok {
		require.Equal(t, fromCurrency, exchangeRate["1. From_Currency Code"], "Expected correct from currency")
		require.Equal(t, toCurrency, exchangeRate["3. To_Currency Code"], "Expected correct to currency")
	} else {
		t.Fatalf("Realtime Currency Exchange Rate not found in the response")
	}
}

// Cenário 3: Testar endpoint inválido
func TestInvalidEndpoint(t *testing.T) {
	pm, err := playwrightsetup.NewPlaywrightManager()
	require.NoError(t, err, "Failed to start Playwright")
	defer pm.Close()

	requestContext, err := pm.Playwright.Request.NewContext()
	require.NoError(t, err, "Failed to create request context")

	// Faz uma requisição para um endpoint inválido
	response, err := requestContext.Get("https://www.alphavantage.co/query?function=INVALID_ENDPOINT&apikey=" + os.Getenv("ALPHA_VANTAGE_API_KEY"))
	require.NoError(t, err, "Failed to fetch API")

	// Verifica se o status da resposta é 200, pois a API retorna status 200 com mensagem de erro
	require.Equal(t, http.StatusOK, response.Status(), "Expected status 200")

	// Lê o corpo da resposta e verifica a mensagem de erro
	var body map[string]interface{}
	responseText, err := response.Text()
	require.NoError(t, err, "Failed to read response text")

	err = json.Unmarshal([]byte(responseText), &body)
	require.NoError(t, err, "Failed to parse JSON response")

	// Verifica se a mensagem de erro está presente
	require.Equal(t, "This API function (INVALID_ENDPOINT) does not exist.", body["Error Message"], "Expected error message for invalid endpoint")
}

// Cenário 4: Verificar o total de casas (ações ou outras informações)
func TestGlobalQuote(t *testing.T) {
	symbol := "IBM"
	body := makeRequest(t, fmt.Sprintf("function=GLOBAL_QUOTE&symbol=%s", symbol))

	// Verifica se o Global Quote está presente
	if globalQuote, ok := body["Global Quote"].(map[string]interface{}); ok {
		require.Contains(t, globalQuote, "01. symbol", "Expected symbol in Global Quote")
		require.Equal(t, symbol, globalQuote["01. symbol"], "Expected correct symbol")
	} else {
		t.Fatalf("Global Quote not found in the response")
	}
}
