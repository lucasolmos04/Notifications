package service

import (
	"Notifications/src/internal/domain"
	"context"
	"uuid"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
}
