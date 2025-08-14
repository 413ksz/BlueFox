package app

import (
	"net/http"
	"os"
	"time"

	"github.com/413ksz/BlueFox/backEnd/pkg/database"
	userService "github.com/413ksz/BlueFox/backEnd/user_menagment/application/service"
	userRepository "github.com/413ksz/BlueFox/backEnd/user_menagment/infrastructure/persistence/repository"
	userRouter "github.com/413ksz/BlueFox/backEnd/user_menagment/interfaces/http"
	"github.com/413ksz/BlueFox/backEnd/user_menagment/shared/validation"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Global variables for warm starts
var (
	AppRouter     *mux.Router
	dB            *database.DB
	userHandler   *userRouter.UserHandler
	userRepo      *userRepository.UserRepository
	userSvc       *userService.UserServiceImpl
	pwnedPassword *validation.PwnedPassword
	pwnedClient   *http.Client
)

// Init is the initialization function for the main application on cold starts
// it initializes every resource that the application needs to function correctly
//
// Functinality:
// 1. Set up zerolog logging: set the global log level, configure the logger
// 2. Set up database connection: initialize the global database
// 3. Set up domains: initialize the service, handler, repository layers for the application
// 4. Set up router: initialize the main router for handling API requests
//
// Error Conditions:
// - if the database connection fails
func Init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

	if os.Getenv("DEBUG") == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			With().Timestamp().Caller().Logger()
		log.Debug().
			Str("component", "main_app").
			Str("event", "debug_mode_enabled").
			Msg("Debug mode enabled, logging to console and setting debug level")
	}
	log.Info().
		Str("component", "main_app").
		Str("event", "app_logger_initialized").
		Msg("Logger successfully initialized")

	log.Info().
		Str("component", "main_app").
		Str("event", "app_init_start").
		Msg("Serverless function initializing...")

	log.Info().
		Str("component", "main_app").
		Str("event", "app_db_init_start").
		Msg("Initializing global database connection")
	db, dbErr := database.NewGormDb()
	if dbErr != nil {
		log.Fatal().
			Err(dbErr).
			Str("component", "main_app").
			Str("event", "app_db_init_failure").
			Msg("Failed to initialize global database")
	}

	dB = database.NewDB(db)
	log.Info().
		Str("component", "main_app").
		Str("event", "app_db_init_success").
		Msg("Global database connection initialized successfully")

	log.Info().
		Str("component", "main_app").
		Str("event", "app_domain_init_start").
		Msg("Initializing application domains")

	log.Info().
		Str("component", "main_app").
		Str("event", "app_domain_init_user").
		Msg("Initializing user domain")

	userRepo = userRepository.NewUserRepository(dB)

	pwnedClient = &http.Client{}
	pwnedPassword = validation.NewPwnedPassword(pwnedClient)

	userSvc = userService.NewUserService(userRepo, pwnedPassword)

	userHandler = userRouter.NewUserHandler(userSvc)

	AppRouter = mux.NewRouter()
	userRouter.RegisterUserRoutes(AppRouter, userHandler)

	log.Info().
		Str("component", "main_app").
		Str("event", "app_domain_init_user_success").
		Msg("User domain initialized successfully")

	log.Info().
		Str("component", "main_app").
		Str("event", "app_domain_init_success").
		Msg("Application domains initialized successfully")

	log.Info().
		Str("component", "main_app").
		Str("event", "app_init_success").
		Msg("Serverless function initialized successfully")
}
