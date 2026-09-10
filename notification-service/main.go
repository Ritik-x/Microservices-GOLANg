package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"notification-service/events"
	"notification-service/repositries"
)

func main() {

	// Load .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	// Get database URL
	databaseURL := os.Getenv("DATABASE_URL")

	// Connect to PostgreSQL
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

	fmt.Println("Connected to PostgreSQL")

	repo := repositries.NewNotificationRepository(db)
	// Connect to NATS
	nc, err := events.ConnectNats()
	if err != nil {
		fmt.Println("NATS connection failed:", err)
		return
	}

	defer nc.Close()

	// Start consumer
	err = events.StartConsumer(nc, repo)
	if err != nil {
		fmt.Println("Consumer failed:", err)
		return
	}

	fmt.Println("Notification Service running")

	select {}
}
