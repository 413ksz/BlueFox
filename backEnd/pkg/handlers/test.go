package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing test handler.")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello from the /api/test endpoint! Routing works!\n")
}
