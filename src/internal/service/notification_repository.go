package service

import (
	"Notifications/src/internal/domain"
	"context"
	"uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification domain.Notification) error
	Update(ctx context.Context, notification domain.Notification) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error)
	FindByUserId(ctx context.Context, userId uuid.UUID) ([]domain.Notification, error)
}
