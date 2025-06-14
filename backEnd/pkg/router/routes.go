package router

import (
	"log"

	"github.com/413ksz/BlueFox/backEnd/pkg/handlers"
	"github.com/gorilla/mux"
)

// RegisterRoutes adds routes from the handlers package to the given gorilla/mux router.
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/test", handlers.TestHandler).Methods("GET")

	r.HandleFunc("/api/user/{id}", handlers.UserGetHandler).Methods("GET")
	r.HandleFunc("/api/user", handlers.UserCreateHandler).Methods("PUT")
	r.HandleFunc("/api/user/{id}", handlers.TestHandler).Methods("DELETE")
	r.HandleFunc("/api/user/login/{email}", handlers.TestHandler).Methods("POST")
	r.HandleFunc("/api/user/{id}", handlers.UserUpdateHandler).Methods("PATCH")

	log.Println("Handlers package registered routes.")
}
