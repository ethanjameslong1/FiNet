package analysis

// imports... obvi
// checking github correctness
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// global constants
const (
	APIKeyEnvValue          = "API_KEY"
	APIURL                  = "https://www.alphavantage.co/query?function=%v&symbol=%v&apikey=%v"
	WeeklyAdjustedFunction  = "TIME_SERIES_WEEKLY_ADJUSTED"
	MonthlyAdjustedFunction = "TIME_SERIES_MONTHLY_ADJUSTED"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second, // Global timeout for all requests using this client
}

// structs

// StockDataMonthly ...
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
	Information  string `json:"Information"`
}

// StockDataWeekly ...
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
	Information  string `json:"Information"`
}

type AlphaVantageParam struct {
	Function  string
	Symbol    string
	Datatype  string // defaults to JSON, i don't think csv will ever come into play
	APIKey    string
	StartDate time.Time
	EndDate   time.Time
}

func RetrieveStockDataWeekly(ctx context.Context, params AlphaVantageParam) (*StockDataWeekly, error) {
	if params.Function != "TIME_SERIES_WEEKLY_ADJUSTED" || params.Symbol == "" || params.APIKey == "" {
		return nil, fmt.Errorf("required params are missing or wrong")
	}
	if params.Datatype == "" {
		params.Datatype = "json"
	}
	apiRequestURL := fmt.Sprintf(APIURL, params.Function, params.Symbol, params.APIKey)
	time.Sleep(1 * time.Second)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiRequestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	typeAccepting := fmt.Sprintf("application/%v", params.Datatype)
	req.Header.Set("Accept", typeAccepting)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(errorBody))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	var stockData StockDataWeekly
	err = json.Unmarshal(body, &stockData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	if len(stockData.MetaData.Symbol) == 0 {
		fmt.Printf("No data returned, response: %+v", stockData)
	}
	if params.StartDate.IsZero() && params.EndDate.IsZero() {
		return &stockData, nil
	}
	if params.StartDate.IsZero() {
		params.StartDate, err = time.Parse(dateFormat, "2001-01-01")
		if err != nil {
			return nil, fmt.Errorf("internal Server Error: DateFormat: %v", err)
		}

	}
	if params.EndDate.IsZero() {
		params.EndDate, err = time.Parse(dateFormat, findMostRecentFriday())
		if err != nil {
			return nil, fmt.Errorf("internal Server Error: DateFormat: %v", err)
		}
	}
	for dateStr := range stockData.TimeSeriesWeekly {
		recordDate, err := time.Parse(dateFormat, dateStr)
		if err != nil {
			continue
		}
		if recordDate.Before(params.StartDate) || recordDate.After(params.EndDate) {
			delete(stockData.TimeSeriesWeekly, dateStr)
		}
	}
	return &stockData, nil
}

func MakeWeeklyDataSlice(ctx context.Context, symbols []string, timePer string) ([]*StockDataWeekly, error) {
	var startDate time.Time
	endDate := time.Now()

	switch timePer {
	case "1mo":
		startDate = endDate.AddDate(0, -1, 0)
	case "3mo":
		startDate = endDate.AddDate(0, -3, 0)
	case "6mo":
		startDate = endDate.AddDate(0, -6, 0)
	case "1y":
		startDate = endDate.AddDate(-1, 0, 0)
	case "5y":
		startDate = endDate.AddDate(-5, 0, 0)
	case "max":
	default:
		startDate = endDate.AddDate(0, -6, 0)
	}

	paramTemplate := AlphaVantageParam{Function: WeeklyAdjustedFunction, APIKey: os.Getenv(APIKeyEnvValue), StartDate: startDate, EndDate: endDate}
	dataSlice := make([]*StockDataWeekly, len(symbols))
	var allErrors []error
	var err error
	for i, s := range symbols {
		paramTemplate.Symbol = s
		dataSlice[i], err = RetrieveStockDataWeekly(ctx, paramTemplate)
		if err != nil {
			log.Printf("Error retrieving stock data for symbol %q: %v", s, err)
			allErrors = append(allErrors, fmt.Errorf("symbol %q: %w", s, err))
			if len(allErrors) > 3 {
				return nil, fmt.Errorf("too many failed API calls, check symbols list and API. atleast 3 failed API calls")
			}
			continue
		}
	}
	return dataSlice, nil
}

func AnalysisStoreWeeklyDataSlice(ctx context.Context, symbols []string) ([]*StockDataWeekly, error) {
	paramTemplate := AlphaVantageParam{Function: WeeklyAdjustedFunction, APIKey: os.Getenv(APIKeyEnvValue)}
	dataSlice := make([]*StockDataWeekly, len(symbols))
	var allErrors []error
	var err error
	for i, s := range symbols {
		paramTemplate.Symbol = s
		dataSlice[i], err = RetrieveStockDataWeekly(ctx, paramTemplate)
		if err != nil {
			log.Printf("Error retrieving stock data for symbol %q: %v", s, err)
			allErrors = append(allErrors, fmt.Errorf("symbol %q: %w", s, err))
			if len(allErrors) > 3 {
				return nil, fmt.Errorf("too many failed API calls, check symbols list and API. atleast 3 failed API calls")
			}
			continue
		}
	}
	return dataSlice, nil
}
