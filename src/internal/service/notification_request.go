package service

import (
	"Notifications/src/internal/domain"
	"uuid"
)

type CreateNotificationRequest struct {
	UserID  uuid.UUID
	Title   string
	Content string
	Channel domain.Channel
}
