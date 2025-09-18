package handler

import (
	"log"
	"net/http"
)

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SessionCookie")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sessionID := cookie.Value
	_, err = h.UserSessionDBService.DeleteSessionByID(r.Context(), sessionID)
	if err != nil {
		log.Printf("Error deleting session from database: %v", err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "SessionCookie",
		Value:    "", // Clear the value
		Path:     "/",
		MaxAge:   -1, // Expire immediately
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
