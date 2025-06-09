package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello You")

}
