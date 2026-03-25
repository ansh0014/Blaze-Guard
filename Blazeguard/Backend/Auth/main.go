package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"auth/config"
	"auth/handler"
	"auth/models"
	"auth/routes"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Println("Configuration loaded successfully")

	db, err := sql.Open("postgres", cfg.GetDatabaseDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connection established")

	store := sessions.NewCookieStore([]byte(cfg.Session.Secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}



	userRepo := models.NewUserRepository(db)

	authHandler, err := handler.NewAuthHandler(userRepo, cfg, store)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}



	router := routes.SetupRoutes(authHandler)
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on http://localhost%s", addr)


	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
