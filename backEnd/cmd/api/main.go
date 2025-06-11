package main

import (
	"log"
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/database"
	"github.com/413ksz/BlueFox/backEnd/pkg/router"
	"github.com/gorilla/mux"
)

var appRouter *mux.Router

func init() {
	log.Println("Serverless function initializing...")

	if err := database.InitAppDB(); err != nil {
		log.Fatalf("Failed to initialize global database: %v", err)
	}
	log.Println("Global database connection successfully initialized.")

	appRouter = mux.NewRouter()

	router.RegisterRoutes(appRouter)
	log.Println("API routes initialized.")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: Method=%s, Path=%s", r.Method, r.URL.Path)
	appRouter.ServeHTTP(w, r)
}
