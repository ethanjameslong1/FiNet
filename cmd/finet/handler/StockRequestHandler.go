package handler

import (
	"log"
	"net/http"
	"text/template"

	"github.com/ethanjameslong1/FiNet/cmd/stock/analysis"
	"github.com/ethanjameslong1/FiNet/database"
)

type PageData struct {
	UserData     UserLoginData
	Error        error
	StockWeights analysis.StockWeights
	Interval     string
	Predictions  []database.Prediction
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

// i need to decouple this bad.
func (h *Handler) StockRequestHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(database.User)
	if !ok {
		log.Printf("Error: User not found in context for StockHandler. Redirecting to login.")
		http.Redirect(w, r, "/login", http.StatusNotFound)
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

	//TODO add new place for this to route to, somewhere to look at raw data, maybe just add a user main page that can route to all of these individually would be smart.

	tmpl, err := template.ParseFiles("static/stockAnalysisRequestComplete.html")
	if err != nil {
		log.Printf("Error parsing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//TODO from here I need to move this analysis logic elsewhere.

	dataSlice, err := analysis.MakeWeeklyDataSlice(r.Context(), symbolList)
	if err != nil {
		log.Printf("Error creating data slice for analysis: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	DataMap, err := analysis.StoreWeeklyDataV1(dataSlice, "", 1)
	if err != nil {
		log.Printf("Error colelcting weekly stock data: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	Pred, err := analysis.AnalyzeStoredDataV1(DataMap)
	if err != nil {
		log.Printf("Error analyzing stored stock data: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	for _, prediction := range Pred {
		log.Printf("AddPrediction begin called with %s (predictable), %s (predictor) and %f (correlation)", prediction.PredictableSym, prediction.PredictorSym, prediction.Correlation)
		h.StockDBService.AddPrediction(r.Context(), prediction.PredictableSym, prediction.PredictorSym, prediction.Correlation, "First Draft", user.ID)
	}

	//TODO here the analysis logic is complete.

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
		http.Redirect(w, r, "/login", http.StatusNotFound)
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
	_, ok := r.Context().Value(userContextKey).(database.User)
	if !ok {
		log.Printf("Error: User not found in context for StockHandler. Redirecting to login.")
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}
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
