package main

import (
	"Notifications/src/internal/auth"
	"Notifications/src/internal/controller"
	"Notifications/src/internal/domain"
	"Notifications/src/internal/infrastructure"
	"Notifications/src/internal/repository"
	"Notifications/src/internal/sender"
	"Notifications/src/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const tokenTTL = 24 * time.Hour

func main() {
	fmt.Println("Welcome to main() Notifications")

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading config from the environment")
	}

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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	// Repositories
	userRepo := repository.NewSQLUserRepository(db)
	notificationRepo := repository.NewSQLNotificationRepository(db)

	// Auth
	passwordEncoder := auth.NewBcryptPasswordEncoder()
	tokenProvider := auth.NewJWTTokenProvider(jwtSecret, tokenTTL)
	authenticationService := service.NewAuthenticationService(userRepo, passwordEncoder, tokenProvider)

	// Notification senders, by channel
	senders := service.SenderRegistry{
		domain.ChannelEmail: &sender.EmailSender{},
		domain.ChannelSMS:   &sender.SmsSender{},
		domain.ChannelPush:  &sender.PushSender{},
	}
	notificationService := service.NewNotificationService(notificationRepo, userRepo, senders)
	_ = notificationService // TODO: wire up notification handlers once the auth middleware exists

	// HTTP
	authHandler := controller.NewAuthHandler(authenticationService)
	router := controller.NewRouter(authHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
