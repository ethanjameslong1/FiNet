package main

import (
	"fmt"
	"https://github.com/ethanjameslong1/GoCloudProject.git/handlers/handler"
	"net/http"
)

func main() {
	server := http.FileServer(http.Dir("../static"))

	http.Handle("/", server)
	http.HandleFunc("/hello", handler.HelloHandler)
	http.HandleFunc("/form", handler.FormHandler)
	//we need a handler package, which will be imported through the github link
	//Hello handler I think will work with the index.html and the formHandler will deal with the form.html

}

func handleRoot(
	w http.ResponseWriter,
	r *http.Request,
) {
	fmt.Fprintf(w, "Hello World")
}
