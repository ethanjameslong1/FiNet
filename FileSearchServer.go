package main

import (
	"fmt"
	"github.com/ethanjameslong1/GoCloudProject.git/handler"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	// server := http.FileServer(http.Dir("../static"))

	mux.HandleFunc("/", handler.RootHandler)
	mux.HandleFunc("POST /search", handler.SearchHandler)

	fmt.Printf("port running on localhost:8080/\n")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}
