// FileSearchServer.go
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
	mux.HandleFunc("GET /login", handler.ShowLogin)
	mux.HandleFunc("POST /login", handler.LoginHandler)

	fmt.Printf("port running on localhost:8080/\n")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}
