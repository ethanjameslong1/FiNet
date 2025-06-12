// FileSearchServer.go
package main

import (
	"fmt"
	"github.com/ethanjameslong1/GoCloudProject.git/database"
	"github.com/ethanjameslong1/GoCloudProject.git/handler"

	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	service, err := database.NewService(database.DriverName, database.DataSource)
	if err != nil {
		log.Fatal(err)
	}
	handle, err := handler.NewHandler(service)
	if err != nil {
		log.Fatal(err)
	}

	// server := http.FileServer(http.Dir("../static"))

	mux.HandleFunc("/", handle.RootHandler)
	mux.HandleFunc("GET /login", handle.ShowLogin)
	mux.HandleFunc("POST /login", handle.LoginHandler)
	mux.HandleFunc("GET /stock", handle.StockHandler)

	fmt.Printf("port running on localhost:8080/\n")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}
