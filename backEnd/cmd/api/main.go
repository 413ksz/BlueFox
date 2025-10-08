// Handler is the primary entry point for the serverless function hosted on Vercel
package handler

import (
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/app"
	"github.com/rs/zerolog/log"
)

// init is a function that is called when the application is started.
// It sets up the global variables for warm starts.
func init() {
	app.Init()
}

// Handler is the primary entry point for the serverless function hosted on serverless environment.
// It handles incoming HTTP requests and forwards them to the router for processing.
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Str("component", "main_app_handler").
		Str("event", "http_request_received").
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("remote_addr", r.RemoteAddr).
		Msg("Incoming HTTP request")

	app.AppRouter.ServeHTTP(w, r)
}
