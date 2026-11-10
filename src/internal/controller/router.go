package controller

import "net/http"

func NewRouter(authHandler *AuthHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)

	return mux
}
