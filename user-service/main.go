package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	// "github.com/jackc/pgx/v5/pgxpool"
	// "github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"user-service/events"
	"user-service/handlers"
	"user-service/repositries"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User Service is running")
}

func main() {

	// Load .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env")
		return
	}

	// Connect to NATS
	nc, err := events.ConnectNats()
	if err != nil {
		fmt.Println("NATS connection failed:", err)
		return
	}
	defer nc.Close()

	// jetstream

	js, err := events.CreateJetStream(nc)
	if err != nil {
		fmt.Println("JetStream setup failed:", err)
		return
	}
	fmt.Println("JetStream ready")
	_ = js

	// Get database URL
	databaseURL := os.Getenv("DATABASE_URL")

	// Connect to Neon PostgreSQL
	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		fmt.Println("Database connection failed:", err)
		return
	}

	defer db.Close()

	// Check database connection
	err = db.Ping(context.Background())
	if err != nil {
		fmt.Println("Database ping failed:", err)
		return
	}

	fmt.Println("Connected to Neon PostgreSQL")

	// Create repository
	repo := repositries.NewUserRepository(db)

	// Create handler
	userHandler := handlers.NewUserHandler(repo, js)

	// Routes
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/users", userHandler.CreateUserHandler)
	http.HandleFunc("/login", userHandler.LoginHandler)

	fmt.Println("User Service running on port 8081")

	err = http.ListenAndServe(":8081", nil)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}
