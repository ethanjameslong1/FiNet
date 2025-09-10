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
	UserSessionDBService *database.DBService
	StockDBService       *database.DBService
	SessionDuration      time.Duration
}

func NewHandler(UserSessionDB *database.DBService, StockDB *database.DBService, sessionDuration time.Duration) (*Handler, error) {
	if UserSessionDB == nil {
		return nil, errors.New("database services cannot be nil")
	}
	if StockDB == nil {
		return nil, errors.New("database services cannot be nil")
	}

	return &Handler{
		UserSessionDBService: UserSessionDB,
		StockDBService:       StockDB,
		SessionDuration:      sessionDuration,
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

	user, err := h.UserSessionDBService.AuthenticateUser(r.Context(), username, password)
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
	_, err = h.UserSessionDBService.AddSession(r.Context(), sessionID, user.ID, h.SessionDuration)
	if err != nil {
		log.Printf("Error adding session for user '%s': %v", user.Username, err)
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:  "SessionCookie",
		Value: sessionID.String(), Path: "/",
		Expires:  time.Now().Add(h.SessionDuration),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/homepage", http.StatusSeeOther)

}

func (h *Handler) ShowRegistration(w http.ResponseWriter, r *http.Request) {

	_, ok := r.Context().Value(userContextKey).(database.User)
	if ok {
		http.Redirect(w, r, "/stock", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("static/registration.html")
	if err != nil {
		log.Printf("Error parsing registration template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, PageData{})
	if err != nil {
		log.Printf("Error executing registration template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
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

	succuess, err := h.UserSessionDBService.AddUser(r.Context(), username, password)
	if err != nil || !succuess {
		log.Printf("Register attempt failed for user '%s': %v", username, err)
		http.Error(w, "Register failed due to a server error", http.StatusInternalServerError)
		return
	}
	_, err = h.UserSessionDBService.GetUserByName(r.Context(), username)
	if err != nil {
		log.Printf("Error finding recently added user %s: %v", username, err)
		http.Error(w, "Failed to find recently added user", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
