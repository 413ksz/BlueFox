package app

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/413ksz/BlueFox/backEnd/pkg/database"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	passwordHashing "github.com/413ksz/BlueFox/backEnd/pkg/password_hashing"
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
	AppRouter           *mux.Router
	dB                  *database.DB
	userHandler         *userRouter.UserHandler
	userRepo            *userRepository.UserRepository
	userSvc             *userService.UserServiceImpl
	pwnedPassword       *validation.PwnedPassword
	pwnedClient         *http.Client
	argon2ID            *passwordHashing.KriptoArgon2ID
	saltLength          uint8
	iterations          uint32
	memoryCostKiloBytes uint32
	threads             uint8
	keyLengthBytes      uint32
	pepperSecret        []byte
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

	argonSetupErr := SetupArgon2IDVars()
	if argonSetupErr != nil {
		log.Fatal().
			Err(argonSetupErr).
			Str("component", "main_app_initializer").
			Str("status", "failed").
			Str("event", "app_domain_init_user_failure").
			Str("errorcode", argonSetupErr.Code.String()).
			Str("message", argonSetupErr.Message).
			Interface("details", argonSetupErr.Details).
			Msg("Failed to initialize user domain")
	}

	argon2IDTemp, argonErr := passwordHashing.NewKriptoArgon2Id(
		saltLength,
		iterations,
		memoryCostKiloBytes,
		threads,
		keyLengthBytes,
		pepperSecret)
	if argonErr != nil {
		log.Fatal().
			Err(argonSetupErr).
			Str("component", "main_app_initializer").
			Str("status", "failed").
			Str("event", "app_domain_init_user_failure").
			Str("errorcode", argonSetupErr.Code.String()).
			Str("message", argonSetupErr.Message).
			Interface("details", argonSetupErr.Details).
			Msg("Failed to initialize user domain")

	}
	argon2ID = argon2IDTemp

	userSvc = userService.NewUserService(userRepo, pwnedPassword, argon2ID)

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

func SetupArgon2IDVars() *models.CustomError {
	pepperTemp := os.Getenv("ARGON2ID_PEPPER")
	saltLengthTemp := os.Getenv("ARGON2ID_SALT_LENGTH")
	iterationsTemp := os.Getenv("ARGON2ID_ITERATIONS")
	memoryCostKiloBytesTemp := os.Getenv("ARGON2ID_MEMORY_COST_KILOBYTES")
	threadsTemp := os.Getenv("ARGON2ID_THREADS")
	keyLengthBytesTemp := os.Getenv("ARGON2ID_KEY_LENGTH_BYTES")

	log.Info().
		Str("papper", pepperTemp).
		Str("saltLength", saltLengthTemp).
		Str("iterations", iterationsTemp).
		Str("memoryCostKiloBytes", memoryCostKiloBytesTemp).
		Str("threads", threadsTemp).
		Str("keyLengthBytes", keyLengthBytesTemp).
		Msg("ARGON2ID environment variables found")

	if pepperTemp == "" || saltLengthTemp == "" || iterationsTemp == "" || memoryCostKiloBytesTemp == "" || threadsTemp == "" || keyLengthBytesTemp == "" {
		return models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "ARGON2ID environment variables not set", nil, nil)
	}

	parsedSaltLength, err := strconv.ParseUint(saltLengthTemp, 10, 8)
	if err != nil {
		return models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "ARGON2ID envirenment variable ARGON2ID_SALT_LENGTH is not valid for uint8", nil, nil)
	}
	saltLength = uint8(parsedSaltLength)

	parsedMemoryCostKiloBytes, err := strconv.ParseUint(memoryCostKiloBytesTemp, 10, 32)
	if err != nil {
		return models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "ARGON2ID envirenment variable ARGON2ID_MEMORY_COST_KILOBYTES is not valid for uint32", nil, nil)
	}
	memoryCostKiloBytes = uint32(parsedMemoryCostKiloBytes)

	parsedIterations, err := strconv.ParseUint(iterationsTemp, 10, 32)
	if err != nil {
		return models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "ARGON2ID envirenment variable ARGON2ID_ITERATIONS is not valid for uint32", nil, nil)
	}
	iterations = uint32(parsedIterations)

	parsedThreads, err := strconv.ParseUint(threadsTemp, 10, 8)
	if err != nil {
		return models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "ARGON2ID envirenment variable ARGON2ID_THREADS is not valid for uint8", nil, nil)
	}
	threads = uint8(parsedThreads)

	parsedKeyLengthBytes, err := strconv.ParseUint(keyLengthBytesTemp, 10, 32)
	if err != nil {
		return models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "ARGON2ID envirenment variable ARGON2ID_KEY_LENGTH_BYTES is not valid for uint32", nil, nil)
	}
	keyLengthBytes = uint32(parsedKeyLengthBytes)

	pepperSecret = []byte(pepperTemp)

	return nil
}
