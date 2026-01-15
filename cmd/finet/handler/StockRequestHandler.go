package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/ethanjameslong1/FiNet/database"
)

type StockDataWeekly struct { // NOTE: Needs to be kept up to date with version in AlphaVantageAPI.go
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

type PageData struct {
	UserData    UserLoginData
	Error       error
	Interval    string
	Predictions []database.Prediction // remove
}

type rawAnalysisData struct {
	Symbol string
	Period string
	Data   []*StockDataWeekly
}

var (
	stockRawDataURL  = "http://analysis:8001/rawDataRequest"
	stockAnalysisURL = "http://analysis:8001/analysisRequest"
)

func (h *Handler) StockRequestPageHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(database.User)
	if !ok {
		log.Printf("Error: User not found in context for StockHandler. Redirecting to login.")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	data := PageData{
		UserData: UserLoginData{
			Name: user.Username,
		},
		Error:    nil,
		Interval: "Weekly",
	}

	tmpl, err := template.ParseFiles("static/stockAnalysisRequest.html")
	if err != nil {
		log.Printf("Error parsing stock template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// StockRequestHandler need to decouple this bad.
func (h *Handler) StockRequestHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(database.User)
	if !ok {
		log.Printf("Error: User not found in context for StockHandler. Redirecting to login.")
		http.Redirect(w, r, "/finet/login", http.StatusNotFound)
		return
	}
	uData := UserLoginData{Name: user.Username}
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Printf("error parsing form: %v", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	symbolList := r.PostForm["stocks"]

	err = h.callForAnalysis(w, r, symbolList, user.ID)
	if err != nil {
		log.Printf("error in request call: %v", err)
		http.Error(w, "Failed Request", http.StatusInternalServerError)
	}
	tmpl, err := template.ParseFiles("static/stockAnalysisRequestComplete.html")
	if err != nil {
		log.Printf("Error parsing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// stockData, err := h.callForWeeklyDataSlice(w, r, symbolList, "1", user.ID)
	_, err = h.callForWeeklyDataSlice(w, r, symbolList, "1", user.ID)
	if err != nil {
		http.Error(w, "failed to reach the stock API", http.StatusServiceUnavailable)
		return
	}
	// TODO: store the data in db

	err = tmpl.Execute(w, PageData{UserData: uData, Error: nil})
	if err != nil {
		log.Printf("Error executing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ShowPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(database.User)
	if !ok {
		log.Printf("Error: User not found in context for StockHandler. Redirecting to login.")
		http.Redirect(w, r, "/finet/login", http.StatusNotFound)
		return
	}

	predictions, err := h.StockDBService.GetAllPredictionsForUser(r.Context(), user.ID)
	if err != nil {
		log.Printf("Error fetching predictions: %v", err)
		http.Error(w, "Could not retrieve predictions.", http.StatusInternalServerError)
		return
	}
	data := PageData{
		Predictions: predictions,
	}
	tmpl, err := template.ParseFiles("static/show_predictions.html")
	if err != nil {
		log.Printf("Error parsing show_predictions template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing show_predictions template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RawDataRequest(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(database.User)
	if !ok {
		log.Printf("Error: User not found in context for StockHandler. Redirecting to login.")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	data := PageData{
		UserData: UserLoginData{
			Name: user.Username,
		},
		Error:    nil,
		Interval: "Weekly",
	}

	tmpl, err := template.ParseFiles("static/rawDataRequest.html")
	if err != nil {
		log.Printf("Error parsing stock template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RawDataHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(database.User)
	if !ok {
		log.Printf("Error: User not found in context for StockHandler. Redirecting to login.")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("error parsing form: %v", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	symbol := r.PostForm["stockSymbol"]
	period := r.PostForm["period"]
	pageData := rawAnalysisData{
		Symbol: symbol[0],
		Period: period[0],
	}

	tmpl, err := template.ParseFiles("static/showRawAnalysis.html")
	if err != nil {
		log.Printf("Error parsing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	dataSlice, err := h.callForWeeklyDataSlice(w, r, []string{pageData.Symbol}, pageData.Period, user.ID)
	if err != nil {
		log.Printf("Error creating data slice for analysis: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	pageData.Data = dataSlice
	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Printf("Error executing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

type apiCall struct {
	SymbolList []string `json:"symbolList"`
	TimePeriod string   `json:"time"`
	UserID     int      `json:"Id"`
}

func (h *Handler) callForWeeklyDataSlice(w http.ResponseWriter, r *http.Request, sym []string, timePeriod string, userID int) ([]*StockDataWeekly, error) { // NOTE: used for rawDataRequests
	defer r.Body.Close()
	requestData := apiCall{
		SymbolList: sym,
		TimePeriod: timePeriod,
		UserID:     userID,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		http.Error(w, "Failed to prepare data for stock API call", http.StatusInternalServerError)
		return nil, err

	}
	req, err := http.NewRequest("POST", stockRawDataURL, bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to create request for stock API", http.StatusInternalServerError)
		return nil, err

	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to reach the stock API", http.StatusServiceUnavailable)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		log.Printf("Stock API returned an error: %s", resp.Status)
		http.Error(w, "Stock API failed to process the symbol", resp.StatusCode)
		return nil, fmt.Errorf("stock API Returned Error: %v", err)
	}

	var results []*StockDataWeekly
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		log.Printf("Failed to decode response: %v", err)
		http.Error(w, "Failed to decode stock data", http.StatusInternalServerError)
		return nil, err
	}
	return results, nil
}

type anlysisCall struct { // any changes made here need to be reflected in stock api server
	SymbolList []string `json:"symbolList"`
	UserID     int      `json:"Id"`
}

func (h *Handler) callForAnalysis(w http.ResponseWriter, r *http.Request, symbols []string, userID int) error {
	// TODO: add panic function to return or something, right now errors just keep rolling
	defer r.Body.Close()
	requestData := anlysisCall{
		SymbolList: symbols,
		UserID:     userID,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		http.Error(w, "Failed to prepare data for stock API call", http.StatusInternalServerError)
		return err

	}
	req, err := http.NewRequest("POST", stockAnalysisURL, bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to create request for stock API", http.StatusInternalServerError)
		return err

	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to reach the stock API", http.StatusServiceUnavailable)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		log.Printf("Stock API returned an error: %s", resp.Status)
		http.Error(w, "Stock API failed to process request", resp.StatusCode)
		return fmt.Errorf("stock API Returned Error: %v", resp.StatusCode)
	}

	return nil
}

func (h *Handler) ClearPredictions(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(database.User)
	if !ok {
		log.Printf("Error: User not found in context for StockHandler. Redirecting to login.")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	err := h.StockDBService.RemoveAllPredictionsForUser(r.Context(), user.ID)
	if err != nil {
		log.Printf("Error removing all predictions: %v", err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/finet/predictions", http.StatusFound)
}
