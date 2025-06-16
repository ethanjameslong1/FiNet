package handler

import (
	"log"
	"net/http"
	"text/template"

	"github.com/ethanjameslong1/GoCloudProject.git/database"
)

func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	_, ok := r.Context().Value(userContextKey).(database.User)
	if ok {
		http.Redirect(w, r, "/stock", http.StatusSeeOther)
	}
	tmpl, err := template.ParseFiles("static/root.html")
	if err != nil {
		log.Printf("Error parsing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}
