package router

import (
	"github.com/413ksz/BlueFox/backEnd/domain/checks"
	"github.com/413ksz/BlueFox/backEnd/pkg/aggregator"
	apierrorwrapper "github.com/413ksz/BlueFox/backEnd/pkg/api_error_wrapper"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// RegisterRoutes adds routes from the handlers package to the given gorilla/mux router.
func RegisterRoutes(r *mux.Router, aggregator aggregator.HandlerAggregator) {
	log.Info().
		Str("component", "router").
		Str("event", "routes_register_start").
		Msg("Started registering routes...")

	r.HandleFunc("/api/test", apierrorwrapper.ErrorWrapper(checks.TestHandler, "test_handler")).Methods("GET")
	r.HandleFunc("/api/user/{id}", apierrorwrapper.ErrorWrapper(checks.TestHandler, "test_handler")).Methods("GET")
	r.HandleFunc("/api/user", apierrorwrapper.ErrorWrapper(aggregator.UserHandler.UserCreateHandler, "user_create_handler")).Methods("POST")
	r.HandleFunc("/api/user/{id}", apierrorwrapper.ErrorWrapper(checks.TestHandler, "test_handler")).Methods("DELETE")
	r.HandleFunc("/api/user/login", apierrorwrapper.ErrorWrapper(checks.TestHandler, "test_handler")).Methods("POST")
	r.HandleFunc("/api/user/{id}", apierrorwrapper.ErrorWrapper(checks.TestHandler, "test_handler")).Methods("PATCH")

	log.Info().
		Str("component", "router").
		Str("event", "routes_register_finished").
		Msg("Finished registering routes...")

}
