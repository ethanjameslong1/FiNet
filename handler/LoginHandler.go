// LoginHandler.go
package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/ethanjameslong1/GoCloudProject.git/database"
)

type contextKey string

const userContextKey contextKey = "authenticatedUser"

type Guy struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type PageData struct {
	Error error
	Guy   Guy
}

type Handler struct {
	DBService *database.Service
}

func NewHandler(DBService *database.Service) (*Handler, error) {
	h, err := database.NewService(database.DriverName, database.DataSource)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}
	return &Handler{DBService: h}, nil

}

func (h *Handler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/login.html")
	if err != nil {
		log.Printf("Error parsing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, PageData{}) //don't see a reason for adding Request context into this execution
	if err != nil {
		log.Printf("Error executing login template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Form Logic
	var user Guy
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if r.FormValue("username") == "" || r.FormValue("password") == "" {
		http.Error(w, "name or password empty", http.StatusNoContent)
		return
	}
	//Referencing DataBase
	person, err := h.DBService.LoginQuery(r.Context(), r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		log.Printf("Login failed for user '%s': %v", user.Name, err) // Log the actual error for debugging
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "Invalid Username or Password") {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Login failed due to a server error", http.StatusInternalServerError) // 500 Internal Server Error
		return
	}
	//moving on with successful login, going to stock home page
	ctx := context.WithValue(r.Context(), userContextKey, person)
	r = r.WithContext(ctx)
	http.Redirect(w, r, "/stock", http.StatusSeeOther)
}
