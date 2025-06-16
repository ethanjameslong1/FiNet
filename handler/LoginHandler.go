package handler

import (
	"database/sql"
	"errors"
	"github.com/ethanjameslong1/GoCloudProject.git/database"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

type contextKey string

const userContextKey contextKey = "authenticatedUser"

type UserLoginData struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type Handler struct {
	UserDBService    *database.DBService
	SessionDBService *database.DBService
	SessionDuration  time.Duration
}

func NewHandler(userDB *database.DBService, sessionDB *database.DBService, sessionDuration time.Duration) (*Handler, error) {
	if userDB == nil || sessionDB == nil {
		return nil, errors.New("database services cannot be nil")
	}
	return &Handler{
		UserDBService:    userDB,
		SessionDBService: sessionDB,
		SessionDuration:  sessionDuration,
	}, nil
}

func (h *Handler) ShowLogin(w http.ResponseWriter, r *http.Request) {

	_, ok := r.Context().Value(userContextKey).(database.User)
	if ok {
		http.Redirect(w, r, "/stock", http.StatusSeeOther)
		return
	}

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

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Username or password cannot be empty", http.StatusBadRequest)
		return
	}

	user, err := h.UserDBService.AuthenticateUser(r.Context(), username, password)
	if err != nil {
		log.Printf("Login attempt failed for user '%s': %v", username, err)
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "user not found") || strings.Contains(err.Error(), "invalid password") {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Login failed due to a server error", http.StatusInternalServerError)
		return
	}

	sessionID := uuid.New()
	_, err = h.SessionDBService.AddSession(r.Context(), sessionID, user.ID, h.SessionDuration)
	if err != nil {
		log.Printf("Error adding session for user '%s': %v", user.Username, err)
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "SessionCookie",
		Value:    sessionID.String(),
		Path:     "/",
		Expires:  time.Now().Add(h.SessionDuration),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	tmpl, err := template.ParseFiles("static/stockAnalysisRequest.html")
	if err != nil {
		log.Printf("Error parsing stock template after login: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, PageData{UserData: UserLoginData{Name: user.Username}, Error: nil}) // Use user.Username
	if err != nil {
		log.Printf("Error executing stock template after login: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
