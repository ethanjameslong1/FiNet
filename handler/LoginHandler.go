// LoginHandler.go
package handler

import (
	"github.com/ethanjameslong1/GoCloudProject.git/database"
	"log"
	"net/http"
	"text/template"
)

type User struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type PageData struct {
	Error string
}

type Handler struct {
	DBService *database.Service
}

func NewHandler(DBService *database.Service) *Handler {
	h, err := database.NewService(database.DriverName, database.DataSource)
	if err != nil {

	}

}

func ShowLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/login.html")
	if err != nil {
		log.Printf("Error parsing login template: %v", err)
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var person User
	Service := database.NewService(r.Context)
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	person.Name = r.FormValue("username")
	person.Password = r.FormValue("password")
	if person.Name == "" || person.Password == "" {
		http.Error(w, "name or password empty", http.StatusNoContent)
		return
	}
	tmpl, err := template.ParseFiles("static/showUser.html")
	if err != nil {
		log.Printf("Error parsing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, person)
	if err != nil {
		log.Printf("Error executing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
