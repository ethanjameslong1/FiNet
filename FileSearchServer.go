package main

import (
	"net/http";
	"fmt"
)

func main() {
	fmt.Println("hello, making server now")
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)

	fmt.Println("server listening to :8080")
	http.ListenAndServe(":8080", mux)


}

func handleRoot (
	w http.ResponseWriter,
	r *http.Request,
) {
	fmt.Fprintf(w, "Hello World")
}
