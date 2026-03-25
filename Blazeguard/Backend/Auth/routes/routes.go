package routes

import (
	"auth/handler"

	"github.com/gorilla/mux"
)

func SetupRoutes(authHandler *handler.AuthHandler) *mux.Router {
	router := mux.NewRouter()

	router.Use(authHandler.CORSMiddleware)

	router.HandleFunc("/health", authHandler.HealthCheck).Methods("GET")

	router.HandleFunc("/auth/login", authHandler.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST", "GET", "OPTIONS")

	api := router.PathPrefix("/api").Subrouter()
	api.Use(authHandler.AuthMiddleware)
	api.HandleFunc("/profile", authHandler.GetProfile).Methods("GET")

	return router
}
