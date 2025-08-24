package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"service-api/handlers"
	"service-api/middleware"
	"service-api/repository"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize repository (in-memory by default)
	var repo repository.ServiceRepository
	if os.Getenv("USE_SQLITE") == "true" {
		db, err := sql.Open("sqlite3", "services.db")
		if err != nil {
			log.Fatalf("Failed to open SQLite DB: %v", err)
		}

		repo, err = repository.NewSQLiteServiceRepository(db)
		if err != nil {
			log.Fatalf("Database Migration failed: %v", err)
		}

		fmt.Println("Using SQLite repository")
	} else {
		repo = repository.NewMemoryServiceRepository()
	}

	// Set up router and middleware
	r := mux.NewRouter()
	r.Use(middleware.AuthMiddleware)

	// Public endpoint for login
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	// Protected service endpoints
	r.HandleFunc("/services", handlers.ListServicesHandler(repo)).Methods("GET")
	r.HandleFunc("/services/{id}", handlers.GetServiceHandler(repo)).Methods("GET")
	r.HandleFunc("/services/{id}/versions", handlers.GetServiceVersionsHandler(repo)).Methods("GET")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Server running on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
