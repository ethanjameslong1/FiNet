package handler

import (
	"bytes"
	"encoding/json"
	"github.com/ethanjameslong1/FiNet/cmd/stock/analysis"
	"github.com/ethanjameslong1/FiNet/database"
	"log"
	"net/http"
	"text/template"
	"time"
)

type PageData struct {
	UserData     UserLoginData
	Error        error
	StockWeights analysis.StockWeights //remove
	Interval     string
	Predictions  []database.Prediction //remove
}
type rawAnalysisData struct {
	Symbol string
	Period string
	Data   []*analysis.StockDataWeekly
}

func (h *Handler) StockRequestPageHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(database.User)
	if !ok {
		log.Printf("Error: User not found in context for StockHandler. Redirecting to login.")
		http.Redirect(w, r, "/finet/", http.StatusFound)
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

// i need to decouple this bad.
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

	tmpl, err := template.ParseFiles("static/stockAnalysisRequestComplete.html")
	if err != nil {
		log.Printf("Error parsing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// START API TO BACKEND STOCK ANALYSIS
	jsonData, err := json.Marshal(symbolList)
	if err != nil {
		http.Error(w, "Failed to prepare data for stock API call", http.StatusInternalServerError)
		return
	}
	stockAPIURL := "http://analysis:8001/item" //TODO put this somewhere
	req, err := http.NewRequest("POST", stockAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to create request for stock API", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "failed to reach the stock API", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		log.Printf("Stock API returned an error: %s", resp.Status)
		http.Error(w, "Stock API failed to process the item", resp.StatusCode)
		return
	}
	log.Println("Successfully forwarded item to stock API.")
	//END API BACKEND CALL

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
	// user, ok := r.Context().Value(userContextKey).(database.User)
	// if !ok {
	// 	log.Printf("Error: User not found in context for StockHandler. Redirecting to login.")
	// 	http.Redirect(w, r, "/finet/", http.StatusFound)
	// 	return
	// }
	// data := PageData{
	// 	UserData: UserLoginData{
	// 		Name: user.Username,
	// 	},
	// 	Error:    nil,
	// 	Interval: "Weekly",
	// }

	tmpl, err := template.ParseFiles("static/rawDataRequest.html")
	if err != nil {
		log.Printf("Error parsing stock template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, PageData{})
	if err != nil {
		log.Printf("Error executing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RawDataHandler(w http.ResponseWriter, r *http.Request) {
	// _, ok := r.Context().Value(userContextKey).(database.User)
	// if !ok {
	// 	log.Printf("Error: User not found in context for StockHandler. Redirecting to login.")
	// 	http.Redirect(w, r, "/finet/login", http.StatusNotFound)
	// 	return
	// }

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
	//TODO make this an API call to analysis container instead and remove uneeded dependencies
	//TODO or somehow make specifically this easy call to alpha vantage decoupled enough for the front end server to handle it.
	dataSlice, err := analysis.MakeWeeklyDataSlice(r.Context(), symbol)
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
