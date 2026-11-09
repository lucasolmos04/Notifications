package service

import (
	"Notifications/src/internal/domain"
	"context"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Save(ctx context.Context, user *domain.User) error
}
