package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("SessionCookie")
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				log.Printf("AuthMiddleware: Error getting cookie: %v", err)
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		sessionID := cookie.Value
		session, dbErr := h.SessionDBService.GetSessionByID(r.Context(), sessionID)
		if dbErr != nil {
			log.Printf("AuthMiddleware: Session validation failed for session ID '%s': %v", sessionID, dbErr)
			http.SetCookie(w, &http.Cookie{
				Name:     "SessionCookie",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			})
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if time.Now().After(session.ExpiresAt) {
			log.Printf("AuthMiddleware: Session ID '%s' has expired for user ID '%d'", sessionID, session.UserID)
			go func() {
				deleteCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_, delErr := h.SessionDBService.DeleteSessionByID(deleteCtx, sessionID)
				if delErr != nil {
					log.Printf("AuthMiddleware: Failed to delete expired session '%s': %v", sessionID, delErr)
				}
			}()

			http.SetCookie(w, &http.Cookie{
				Name:     "SessionCookie",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			})
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		user, userErr := h.UserDBService.GetUserByID(r.Context(), session.UserID)
		if userErr != nil {
			log.Printf("AuthMiddleware: Failed to get user details for user ID '%d': %v", session.UserID, userErr)
			http.SetCookie(w, &http.Cookie{Name: "SessionCookie", Value: "", Path: "/", Expires: time.Unix(0, 0), HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode})
			http.Error(w, "Authentication error", http.StatusInternalServerError)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
