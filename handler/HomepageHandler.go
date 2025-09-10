package handler

import (
	"github.com/ethanjameslong1/GoCloudProject.git/database"
	"log"
	"net/http"
	"text/template"
)

func (h *Handler) HomepageHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(database.User)
	if !ok {
		log.Printf("Error: User not found in context for StockHandler. Redirecting to login.")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	data := PageData{
		UserData: UserLoginData{
			Name: user.Username,
		},
		Error:    nil,
		Interval: "Weekly",
	}

	tmpl, err := template.ParseFiles("static/userHomepage.html")
	if err != nil {
		log.Printf("Error parsing homepage template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing homepage template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
