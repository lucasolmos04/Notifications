package service

import (
	"Notifications/src/internal/domain"
	"uuid"
)

type TokenClaims struct {
	UserID uuid.UUID
}

//go:generate mockgen -source=token_provider.go -destination=mocks/token_provider_mock.go -package=mocks

type TokenProvider interface {
	GenerateToken(user domain.User) (string, error)
	VerifyToken(token string) (*TokenClaims, error)
}
