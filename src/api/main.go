package main

import (
	"Notifications/src/internal/infrastructure"
	"Notifications/src/internal/repository"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Welcome to main() function")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	db, err := infrastructure.NewPostgresConnection(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatalf("Can not connection database: %v", err)
	}

	defer db.Close()

	fmt.Println("\033[32mConnection database successful\033[0m")

	// Inject Dependencies
	userRepo := repository.NewSQLUserRepository(db)

	_ = userRepo
}
