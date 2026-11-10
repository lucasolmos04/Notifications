package controller

import (
	"Notifications/src/internal/service"
	"net/http"
)

func NewRouter(authService *service.AuthenticationService, authHandler *AuthHandler, notificationHandler *NotificationHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)

	auth := AuthMiddleware(authService)
	mux.Handle("POST /notifications", auth(http.HandlerFunc(notificationHandler.Create)))
	mux.Handle("PUT /notifications/{id}", auth(http.HandlerFunc(notificationHandler.Update)))
	mux.Handle("DELETE /notifications/{id}", auth(http.HandlerFunc(notificationHandler.Delete)))
	mux.Handle("GET /notifications", auth(http.HandlerFunc(notificationHandler.List)))

	return mux
}
