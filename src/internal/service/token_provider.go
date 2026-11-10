package service

import (
	"Notifications/src/internal/domain"
	"uuid"
)

type TokenClaims struct {
	UserID uuid.UUID
}

type TokenProvider interface {
	GenerateToken(user domain.User) (string, error)
	VerifyToken(token string) (*TokenClaims, error)
}
