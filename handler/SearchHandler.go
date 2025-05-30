package handler

import (
	"fmt"
	"net/http"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Fprintf(w, "Username is: %v. Password is: %v", username, password)

}
