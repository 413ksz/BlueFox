package router

import (
	"log"

	"github.com/413ksz/BlueFox/backEnd/pkg/handlers"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/test", handlers.TestHandler).Methods("GET")

	log.Println("Handlers package registered routes.")
}
