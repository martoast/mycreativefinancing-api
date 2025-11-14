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

	// Public property routes (read-only and create)
	router.HandleFunc("/properties", controllers.GetProperties).Methods("GET")
	router.HandleFunc("/properties/{PropertyId}", controllers.GetPropertyById).Methods("GET")
	router.HandleFunc("/properties", controllers.CreateProperty).Methods("POST") // Users can submit properties

	// Protected routes (authentication required) - Admin operations
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.AuthMiddleware)

	// Protected property operations - only authenticated users can edit/delete
	apiRouter.HandleFunc("/properties/{PropertyId}", controllers.UpdateProperty).Methods("PUT")
	apiRouter.HandleFunc("/properties/{PropertyId}", controllers.DeleteProperty).Methods("DELETE")
}
