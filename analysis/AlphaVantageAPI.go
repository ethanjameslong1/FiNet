package analysis

//imports... obvi
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// global constants
const (
	ApiKey                  = "FTMAJAOUIWD2Z7L9"
	ApiURL                  = "https://www.alphavantage.co/query?function=%v&symbol=%v&apikey=%v"
	WeeklyAdjustedFunction  = "TIME_SERIES_WEEKLY_ADJUSTED"
	MonthlyAdjustedFunction = "TIME_SERIES_MONTHLY_ADJUSTED"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second, // Global timeout for all requests using this client
}

// structs
type StockDataWeekly struct {
	MetaData struct {
		Information   string `json:"1. Information"`
		Symbol        string `json:"2. Symbol"`
		LastRefreshed string `json:"3. Last Refreshed"`
		TimeZone      string `json:"4. Time Zone"`
	} `json:"Meta Data"`

	TimeSeriesWeekly map[string]struct {
		Open      string `json:"1. open"`
		High      string `json:"2. high"`
		Low       string `json:"3. low"`
		Close     string `json:"4. close"`
		AdjClose  string `json:"5. adjusted close"`
		Volume    string `json:"6. volume"`
		DivAmount string `json:"7. dividend amount"`
	} `json:"Weekly Adjusted Time Series"`
	ErrorMessage string `json:"Error Message"`
	Note         string `json:"Note"`
}

type StockDataMonthly struct {
	MetaData struct {
		Information   string `json:"1. Information"`
		Symbol        string `json:"2. Symbol"`
		LastRefreshed string `json:"3. Last Refreshed"`
		TimeZone      string `json:"4. Time Zone"`
	} `json:"Meta Data"`

	TimeSeriesMonthly map[string]struct {
		Open      string `json:"1. open"`
		High      string `json:"2. high"`
		Low       string `json:"3. low"`
		Close     string `json:"4. close"`
		AdjClose  string `json:"5. adjusted close"`
		Volume    string `json:"6. volume"`
		DivAmount string `json:"7. dividend amount"`
	} `json:"Monthly Adjusted Time Series"`
	ErrorMessage string `json:"Error Message"`
	Note         string `json:"Note"`
}

type StockWeights struct {
	OpenPriceWeight  string
	HighPriceWeight  string
	ClosePriceWeight string
	LowPriceWeight   string
	VolumeWeight     string
	PercChangeWeight string
	PercRangeWeight  string
}

type AlphaVantageParam struct {
	Function     string
	Symbol       string
	Datatype     string //defaults to JSON, i don't think csv will ever come into play
	APIKey       string
	StartDate    string //for internal use, not a parameter
	EndDate      string //for internal use, not a parameter
	StockWeights StockWeights
}

func RetrieveStockDataWeekly(ctx context.Context, params AlphaVantageParam) (*StockDataWeekly, error) {
	if params.Function != "TIME_SERIES_WEEKLY_ADJUSTED" || params.Symbol == "" || params.APIKey == "" {
		return nil, fmt.Errorf("Required params are missing or wrong")
	}
	if params.Datatype == "" {
		params.Datatype = "json"
	}
	apiRequestUrl := fmt.Sprintf(ApiURL, params.Function, params.Symbol, params.APIKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiRequestUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %w", err)
	}
	typeAccepting := fmt.Sprintf("application/%v", params.Datatype)
	req.Header.Set("Accept", typeAccepting)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(errorBody))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %w", err)
	}
	var stockData StockDataWeekly
	err = json.Unmarshal(body, &stockData)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling response: %w", err)
	}
	return &stockData, nil
}

func MakeWeeklyDataSlice(ctx context.Context, symbols []string) ([]*StockDataWeekly, error) {
	paramTemplate := AlphaVantageParam{Function: WeeklyAdjustedFunction, APIKey: ApiKey}
	dataSlice := make([]*StockDataWeekly, len(symbols))
	var allErrors []error
	var err error
	for i, s := range symbols {
		paramTemplate.Symbol = s
		dataSlice[i], err = RetrieveStockDataWeekly(ctx, paramTemplate)
		if err != nil {
			log.Printf("Error retrieving stock data for symbol %q: %v", s, err)
			allErrors = append(allErrors, fmt.Errorf("symbol %q: %w", s, err))
			if len(allErrors) > 8 {
				return nil, fmt.Errorf("Too many failed API calls, check symbols list and API. atleast 9 failed API calls")
			}
			continue
		}
	}
	return dataSlice, nil
}
