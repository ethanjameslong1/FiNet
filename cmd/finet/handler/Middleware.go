package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type AuthRequest struct {
	Token string `json:"authToken"`
}

func (h *Handler) Middleware(w http.ResponseWriter, r *http.Request) {
	var p AuthRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Print("DEBUG: middleware error decoding")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"authToken": "",
			"username":  "",
		})
	}

	log.Printf("DEBUG: authToken: %s", p.Token)
	sessionID := p.Token
	session, dbErr := h.UserSessionDBService.GetSessionByID(r.Context(), sessionID)
	if dbErr != nil {
		log.Printf("AuthMiddleware: Session validation failed for session ID '%s': %v", sessionID, dbErr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusFound)
		json.NewEncoder(w).Encode(map[string]string{
			"authToken": "",
			"username":  "",
		})
		return
	}
	if time.Now().After(session.ExpiresAt) {
		log.Printf("AuthMiddleware: Session ID '%s' has expired for user ID '%d'", sessionID, session.UserID)
		go func() {
			deleteCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, delErr := h.UserSessionDBService.DeleteSessionByID(deleteCtx, sessionID)
			if delErr != nil {
				log.Printf("AuthMiddleware: Failed to delete expired session '%s': %v", sessionID, delErr)
			}
		}()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusFound)
		json.NewEncoder(w).Encode(map[string]string{
			"authToken": "",
			"username":  "",
		})
		return
	}
	user, err := h.UserSessionDBService.GetUserByID(r.Context(), session.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"authToken": "",
			"username":  "",
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"authToken": sessionID,
		"username":  user.Username,
	})
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("SessionCookie")
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				log.Printf("AuthMiddleware: Error getting cookie: %v", err)
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			http.Redirect(w, r, "/finet/login", http.StatusFound)
			return
		}

		sessionID := cookie.Value
		session, dbErr := h.UserSessionDBService.GetSessionByID(r.Context(), sessionID)
		if dbErr != nil {
			log.Printf("AuthMiddleware: Session validation failed for session ID '%s': %v", sessionID, dbErr)
			http.SetCookie(w, &http.Cookie{
				Name:     "SessionCookie",
				Value:    "",
				Path:     "/finet/",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			})
			http.Redirect(w, r, "/finet/login", http.StatusFound)
			return
		}

		if time.Now().After(session.ExpiresAt) {
			log.Printf("AuthMiddleware: Session ID '%s' has expired for user ID '%d'", sessionID, session.UserID)
			go func() {
				deleteCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_, delErr := h.UserSessionDBService.DeleteSessionByID(deleteCtx, sessionID)
				if delErr != nil {
					log.Printf("AuthMiddleware: Failed to delete expired session '%s': %v", sessionID, delErr)
				}
			}()

			http.SetCookie(w, &http.Cookie{
				Name:     "SessionCookie",
				Value:    "",
				Path:     "/finet/",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			})
			http.Redirect(w, r, "/finet/login", http.StatusFound)
			return
		}

		user, userErr := h.UserSessionDBService.GetUserByID(r.Context(), session.UserID)
		if userErr != nil {
			log.Printf("AuthMiddleware: Failed to get user details for user ID '%d': %v", session.UserID, userErr)
			http.SetCookie(w, &http.Cookie{Name: "SessionCookie", Value: "", Path: "/finet/", Expires: time.Unix(0, 0), HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode})
			http.Error(w, "Authentication error", http.StatusInternalServerError)
			http.Redirect(w, r, "/finet/login", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
