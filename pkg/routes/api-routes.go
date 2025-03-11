package routes

import (
	"api/pkg/controllers"
	"api/pkg/middleware"

	"github.com/gorilla/mux"
)

var RegisterRoutes = func(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/auth/register", controllers.Register).Methods("POST")
	router.HandleFunc("/auth/login", controllers.Login).Methods("POST")

	// Public property routes (no auth required)
	router.HandleFunc("/properties", controllers.GetProperties).Methods("GET")
	router.HandleFunc("/properties/{PropertyId}", controllers.GetPropertyById).Methods("GET")

	// Protected routes (authentication required)
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.AuthMiddleware)

	// Protected property routes (auth required)
	apiRouter.HandleFunc("/properties", controllers.CreateProperty).Methods("POST")
	apiRouter.HandleFunc("/properties/{PropertyId}", controllers.DeleteProperty).Methods("DELETE")
	apiRouter.HandleFunc("/properties/{PropertyId}", controllers.UpdateProperty).Methods("PUT")
}
