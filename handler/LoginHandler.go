package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	var person User
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if person.Name == "" || person.Password == "" {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	fmt.Fprintf(w, "Username is: %v. Password is: %v", person.Name, person.Password)

}
