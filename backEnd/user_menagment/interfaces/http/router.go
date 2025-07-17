package http

import (
	apierrorwrapper "github.com/413ksz/BlueFox/backEnd/pkg/api_error_wrapper"
	"github.com/gorilla/mux"
)

// RegisterUserRoutes adds routes from the user handler to the given gorilla/mux router.
// It takes the main router and the initialized UserHandler.
func RegisterUserRoutes(r *mux.Router, userHandler UserHandler) {

	// ------ Register User Routes ------

	// Create a sub-router for all routes under /api/user
	userRouter := r.PathPrefix("/api/user").Subrouter()

	// Register the routes for the userRouter

	// POST /api/user - Create a new user
	userRouter.HandleFunc("", apierrorwrapper.ErrorWrapper(userHandler.UserCreateHandler, "user_create_handler")).Methods("POST")

	// GET /api/user/{id} - Get a user by ID
	userRouter.HandleFunc("/{id}", apierrorwrapper.ErrorWrapper(userHandler.UserCreateHandler, "get_user_handler")).Methods("GET")

	// DELETE /api/user/{id} - Delete a user
	userRouter.HandleFunc("/{id}", apierrorwrapper.ErrorWrapper(userHandler.UserCreateHandler, "delete_user_handler")).Methods("DELETE")

	// POST /api/user/login - User login
	userRouter.HandleFunc("/login", apierrorwrapper.ErrorWrapper(userHandler.UserCreateHandler, "login_user_handler")).Methods("POST")

	// PATCH /api/user/{id} - Update an existing user
	userRouter.HandleFunc("/{id}", apierrorwrapper.ErrorWrapper(userHandler.UserCreateHandler, "update_user_handler")).Methods("PATCH")

}
