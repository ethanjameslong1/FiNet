package handler

import (
	"log"
	"net/http"
	"text/template"
)

func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) { // TODO: I think this whole file can go
	if r.URL.Path != "/finet/" && r.URL.Path != "/" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	// _, ok := r.Context().Value(userContextKey).(database.User) //TODO: If keeping file this feels like security vulnerability
	// if ok {
	// 	http.Redirect(w, r, "/finet/stock", http.StatusSeeOther)
	// }
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
