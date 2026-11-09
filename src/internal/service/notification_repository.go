package service

import (
	"Notifications/src/internal/domain"
	"context"
	"uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification domain.Notification)
	Update(ctx context.Context, notification domain.Notification)
	Delete(ctx context.Context, id uuid.UUID)
	FindByUserId(ctx context.Context, userId uuid.UUID)
}
