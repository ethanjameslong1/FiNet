package handler

import (
	"log"
	"net/http"
)

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SessionCookie")
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "SessionCookie",
			Value:    "",
			Path:     "/finet/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	} else {
		sessionID := cookie.Value
		_, err = h.UserSessionDBService.DeleteSessionByID(r.Context(), sessionID)
		if err != nil {
			log.Printf("Error deleting session from database: %v", err)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "SessionCookie",
		Value:    "",
		Path:     "/finet/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/Logout", http.StatusSeeOther)
}
