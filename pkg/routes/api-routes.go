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

	// Public property routes (no auth required) - ALL PROPERTY OPERATIONS
	router.HandleFunc("/properties", controllers.GetProperties).Methods("GET")
	router.HandleFunc("/properties/{PropertyId}", controllers.GetPropertyById).Methods("GET")
	router.HandleFunc("/properties", controllers.CreateProperty).Methods("POST")
	router.HandleFunc("/properties/{PropertyId}", controllers.UpdateProperty).Methods("PUT")
	router.HandleFunc("/properties/{PropertyId}", controllers.DeleteProperty).Methods("DELETE")

	// Protected routes (authentication required) - for future use
	// If you need protected endpoints later, uncomment and use this:
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.AuthMiddleware)

	// Example protected routes (currently none needed):
	// apiRouter.HandleFunc("/admin/stats", controllers.GetAdminStats).Methods("GET")
}
