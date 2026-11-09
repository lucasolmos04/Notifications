package domain

import (
	"context"
	"uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification Notification)
	Update(ctx context.Context, notification Notification)
	Delete(ctx context.Context, id uuid.UUID)
	FindByUserId(ctx context.Context, userId uuid.UUID)
}
