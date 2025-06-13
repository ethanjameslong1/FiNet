package handler

import (
	"github.com/ethanjameslong1/GoCloudProject.git/database"
	"log"
	"net/http"
	"text/template"
)

type StockWeights struct {
	OpenPriceWeight  string
	HighPriceWeight  string
	ClosePriceWeight string
	LowPriceWeight   string
	VolumeWeight     string
	PercChangeWeight string
	PercRangeWeight  string
}

type PageData struct {
	Guy          Guy
	Error        error
	StockWeights StockWeights
	Interval     string
}

func (h *Handler) StockPageHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(database.Person)
	if !ok {
		log.Printf("Error: User not found in context for StockHandler. Redirecting to login.")
		http.Redirect(w, r, "/login", http.StatusFound) // StatusFound (302) is common for redirection
		return
	}
	data := PageData{
		Guy: Guy{
			Name: user.Username,
		},
		Error: nil,
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

func (h *Handler) StockRequestHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(database.Person)
	if !ok {
		log.Printf("Error: User not found in context for StockHandler. Redirecting to login.")
		http.Redirect(w, r, "/login", http.StatusFound) // StatusFound (302) is common for redirection
		return
	}
	uData := Guy{Name: user.Username}
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	sData := StockWeights{OpenPriceWeight: r.FormValue("weightCurrentOpen"), HighPriceWeight: r.FormValue("weightCurrentHigh"), ClosePriceWeight: r.FormValue("weightCurrentClose"), LowPriceWeight: r.FormValue("weightCurrentLow"), VolumeWeight: r.FormValue("weightCurrentVolume"), PercChangeWeight: r.FormValue("weightCloseOpenPctChange"), PercRangeWeight: r.FormValue("weightHighLowPctRange")}
	tmpl, err := template.ParseFiles("static/stockAnalysisRequestComplete.html")
	if err != nil {
		log.Printf("Error parsing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, PageData{Guy: uData, StockWeights: sData, Error: nil})
	if err != nil {
		log.Printf("Error executing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

//API CALL AND ANALYSIS LOGIC
