package analysis

//imports... obvi
//checking github correctness
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
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
			if len(allErrors) > 3 {
				return nil, fmt.Errorf("Too many failed API calls, check symbols list and API. atleast 9 failed API calls")
			}
			continue
		}
	}
	return dataSlice, nil
}

func RetrieveStockDataMonthly(ctx context.Context, params AlphaVantageParam) (*StockDataMonthly, error) {
	if params.Function != MonthlyAdjustedFunction || params.Symbol == "" || params.APIKey == "" {
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
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var stockData StockDataMonthly
	if err := json.NewDecoder(resp.Body).Decode(&stockData); err != nil {
		return nil, fmt.Errorf("Error unmarshalling response: %w", err)
	}
	return &stockData, nil
}

func MakeMonthlyDataSlice(ctx context.Context, symbols []string) ([]*StockDataMonthly, error) {
	const (
		rateLimitDelay = 12 * time.Second // 5 calls per minute → 12s between calls
		maxRetries     = 3
		maxErrors      = 5
	)

	paramTemplate := AlphaVantageParam{Function: MonthlyAdjustedFunction, APIKey: ApiKey}
	dataSlice := make([]*StockDataMonthly, 0, len(symbols))
	var allErrors []error

	for _, s := range symbols {
		paramTemplate.Symbol = s

		var stock *StockDataMonthly
		var err error
		for attempt := 1; attempt <= maxRetries; attempt++ {
			stock, err = RetrieveStockDataMonthly(ctx, paramTemplate)
			if err != nil {
				log.Printf("Attempt %d: Error retrieving %q: %v", attempt, s, err)
				time.Sleep(5 * time.Second) // short retry wait
				continue
			}

			// Check for API limit Note or missing data
			if stock.Note != "" {
				log.Printf("AlphaVantage limit reached for %q: %s", s, stock.Note)
				time.Sleep(rateLimitDelay) // wait and retry
				err = fmt.Errorf("rate limited, retrying")
				continue
			}
			if stock.ErrorMessage != "" {
				err = fmt.Errorf("API error for %q: %s", s, stock.ErrorMessage)
				break
			}

			// Success
			dataSlice = append(dataSlice, stock)
			break
		}

		

		if err != nil {
			log.Printf("Failed to retrieve data for %q after %d attempts", s, maxRetries)
			allErrors = append(allErrors, fmt.Errorf("symbol %q: %w", s, err))
			if len(allErrors) > maxErrors {
				return dataSlice, fmt.Errorf("too many failed API calls: %v", allErrors)
			}
		}

		// Respect free-tier rate limit
		time.Sleep(rateLimitDelay)
	}

	for stock, r := range dataSlice {
		fmt.Println(stock, len(r.TimeSeriesMonthly))
	}

	//check return length of each stock, truncate to shortest length
	minLength := -1
	for _, r := range dataSlice {
		length := len(r.TimeSeriesMonthly)
		if minLength == -1 || length < minLength {
			minLength = length
		}
	}
	if minLength == -1 || minLength < 2 {
		return dataSlice, fmt.Errorf("not enough valid data retrieved")
	}
	
	// Truncate all to minLength
	for i, r := range dataSlice {
		if len(r.TimeSeriesMonthly) > minLength {
			// Create a new map with only the most recent minLength entries
			newTimeSeries := make(map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}, minLength)
			
			// Extract and sort the dates to ensure chronological order
			var dates []string
			for date := range r.TimeSeriesMonthly {
				dates = append(dates, date)
			}
			sort.Strings(dates)
			
			// Keep only the most recent minLength dates
			for _, date := range dates[len(dates)-minLength:] {
				newTimeSeries[date] = r.TimeSeriesMonthly[date]
			}
			
			// Replace the old map with the new truncated map
			dataSlice[i].TimeSeriesMonthly = newTimeSeries
		}
	}

	return dataSlice, nil
}
