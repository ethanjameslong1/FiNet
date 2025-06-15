// LoginHandler.go
package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ethanjameslong1/GoCloudProject.git/database"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type contextKey string

const userContextKey contextKey = "authenticatedUser"

type UserLoginData struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type Handler struct {
	DBService *database.Service
}

func NewUserHandler(DBService *database.Service) (*Handler, error) {
	h, err := database.NewService(database.DriverName, database.UserDataSource)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}
	return &Handler{DBService: h}, nil

}

func (h *Handler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(userContextKey).(database.Person)
	if ok {
		http.Redirect(w, r, "/stock", http.StatusSeeOther)
	}
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
	var user UserLoginData
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

	// establishing UUID and Cookie
	newUUID := uuid.New()
	cookie := http.Cookie{
		Name:     "SessionCookie",
		Value:    newUUID.String(),
		Path:     "/",
		MaxAge:   24 * 60 * 60, //1 day
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	//HTML Stuff
	tmpl, err := template.ParseFiles("static/stockAnalysisRequest.html")
	if err != nil {
		log.Printf("Error parsing stock template after login: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, PageData{Guy: UserLoginData{Name: person.Username}, Error: nil})
	if err != nil {
		log.Printf("Error executing stock template after login: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}
