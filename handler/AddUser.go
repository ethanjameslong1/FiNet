package handler

import (
	"log"
	"net/http"
	"text/template"
)

func (h *Handler) ShowAddUser(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/addUser.html")
	if err != nil {
		log.Printf("Error parsing addUser template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, PageData{})
	if err != nil {
		log.Printf("Error executing addUser template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	username, password := r.FormValue("username"), r.FormValue("password")
	_, err = h.UserDBService.AddUser(r.Context(), username, password)
	if err != nil {
		log.Printf("Error adding user: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
