package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("TestHandler received a request.")
	w.WriteHeader(http.StatusOK) // Set HTTP status code 200 OK
	fmt.Fprintf(w, "Hello from the TestHandler! You hit route: %s\n", r.URL.Path)
}
