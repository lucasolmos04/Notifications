package service

import (
	"Notifications/src/internal/domain"
	"context"
	"uuid"
)

//go:generate mockgen -source=notification_repository.go -destination=mocks/notification_repository_mock.go -package=mocks

type NotificationRepository interface {
	Create(ctx context.Context, notification domain.Notification) error
	Update(ctx context.Context, notification domain.Notification) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error)
	FindByUserId(ctx context.Context, userId uuid.UUID) ([]domain.Notification, error)
}
