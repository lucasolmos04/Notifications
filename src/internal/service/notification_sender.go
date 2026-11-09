package service

import (
	"Notifications/src/internal/domain"
	"context"
	"time"
)

type SenderResult struct {
	Channel domain.Channel
	SendAt  time.Time
	Detail  string
}

type NotificationSender interface {
	Validate(ctx context.Context, notification domain.Notification, user domain.User) error
	Send(ctx context.Context, notification domain.Notification, user domain.User) (*SenderResult, error)
}
