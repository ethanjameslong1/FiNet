package handler

import (
	"github.com/ethanjameslong1/GoCloudProject.git/database"
	"log"
	"net/http"
	"text/template"
)

func (h *Handler) StockHandler(w http.ResponseWriter, r *http.Request) {
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

	tmpl, err := template.ParseFiles("static/stock.html")
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
