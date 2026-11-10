package controller

import (
	"Notifications/src/internal/service"
	"context"
	"net/http"
	"strings"
)

func AuthMiddleware(authService *service.AuthenticationService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			const prefix = "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				writeError(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}
			token := strings.TrimPrefix(authHeader, prefix)

			claims, err := authService.ValidateToken(token)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
