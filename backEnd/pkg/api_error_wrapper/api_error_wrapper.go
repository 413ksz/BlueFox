package apierrorwrapper

import (
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/logging"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// AppHandler defines the signature for an application-specific HTTP handler with type parameter T.
// It returns a pointer to an ApiResponse for the response body and a pointer
// to a CustomError if an error occurred during processing.
type AppHandler[T any] func(w http.ResponseWriter, r *http.Request) (*models.ApiResponse[T], *models.CustomError)

// ErrorWrapper wraps an AppHandler to provide centralized error handling,
// request ID management, and standardized response logging.
//
// It ensures that exactly one HTTP response is sent per request, categorizing
// outcomes as:
//   - Scenario 1(error response): Handler returned a CustomError, which is logged and sent to the client.
//   - Scenario 2(success response): Handler processed successfully and returned a valid ApiResponse.
//   - Scenario 3(internal error): Handler returned (nil, nil), indicating an unexpected
//     internal issue, resulting in a 500 Internal Server Error being sent.
//
// Parameters:
// - handler: The AppHandler function to be wrapped.
// - componentName: A string identifying the component for logging purposes.
//
// Returns an http.HandlerFunc compatible with standard Go HTTP servers (e.g., net/http, gorilla/mux in this case).
func ErrorWrapper[T any](handler AppHandler[T], componentName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Execute the wrapped application handler
		apiResponse, customError := handler(w, r)

		// Retrieve or generate a unique request ID for logging and tracing.
		// X-Vercel-Id is typically provided in Vercel environments.
		requestId := r.Header.Get("X-Vercel-Id")
		if requestId == "" {
			requestId = uuid.New().String()
		}

		// --- Scenario 1: Custom Error Occurred ---
		// If the handler returned a custom error and an apiResponse, log it and send the appropriate response.
		// The apiResponse check is added to avoid a panic when the apiResponse is nil because we work with pointers
		if customError != nil && apiResponse != nil {

			logEvent := logging.LogLevelHelperForError(customError)

			logEvent.
				Str("request_id", requestId).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("component", componentName).
				Str("status", "failed").
				Str("event", "api_error_occurred").
				Str("errorcode", customError.Code).
				Str("message", customError.Message).
				Err(customError.Err).
				Interface("details", customError.Details).
				Msg("Error occurred in API handler")
			models.SendApiResponse(r, w, apiResponse, requestId, componentName)
			return
		}

		// --- Scenario 2: Successful Response ---
		// If no custom error occurred, we expect a valid apiResponse for a successful operation.
		// Log the success and send the response.
		if apiResponse != nil {
			log.Info().
				Str("request_id", requestId).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("component", componentName).
				Str("status", "success").
				Str("event", "api_handler_success").
				Msg("API handler completed successfully")
			models.SendApiResponse(r, w, apiResponse, requestId, componentName)
			return
		}

		// --- Scenario 3: Fallback - Unexpected Nil Responses ---
		// This block is reached if customError is nil AND apiResponse is also nil.
		// This indicates an unexpected state where the handler neither returned a
		// success response nor explicitly signaled an error. It's considered an internal bug.
		// Log the error and send a 500 Internal Server Error response.
		// if this happen we are in deep shit
		log.Error().
			Str("request_id", requestId).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("component", componentName).
			Str("status", "failed").
			Str("event", "api_handler_failed").
			Msg("API handler returned nil responses for both apiResponse and customError")

		apiResponse = models.NewApiResponse[T](nil, http.StatusInternalServerError, "Internal Server Error")
		models.SendApiResponse(r, w, apiResponse, requestId, componentName)
	}
}
