package main

import (
	"fmt"
	"github.com/ethanjameslong1/GoCloudProject.git/handler"
	"log"
	"net/http"
)

func main() {
	server := http.FileServer(http.Dir("../static"))

	http.Handle("/", server)
	http.HandleFunc("/hello", handler.HelloHandler)

	fmt.Printf("port running on localhost:8080/\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
