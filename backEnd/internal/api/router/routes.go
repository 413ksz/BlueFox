package router

import (
	"github.com/413ksz/BlueFox/backEnd/internal/api/handler"
	"github.com/413ksz/BlueFox/backEnd/internal/api/handler/user"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// RegisterRoutes adds routes from the handlers package to the given gorilla/mux router.
func RegisterRoutes(r *mux.Router) {
	log.Info().
		Str("component", "router").
		Str("event", "routes_register_start").
		Msg("Started registering routes...")

	r.HandleFunc("/api/test", handler.TestHandler).Methods("GET")
	r.HandleFunc("/api/user/{id}", user.UserGetHandler).Methods("GET")
	r.HandleFunc("/api/user", user.UserCreateHandler).Methods("PUT")
	r.HandleFunc("/api/user/{id}", handler.TestHandler).Methods("DELETE")
	r.HandleFunc("/api/user/login", user.UserLoginHandler).Methods("POST")
	r.HandleFunc("/api/user/{id}", user.UserUpdateHandler).Methods("PATCH")

	log.Info().
		Str("component", "router").
		Str("event", "routes_register_finished").
		Msg("Finished registering routes...")

}
