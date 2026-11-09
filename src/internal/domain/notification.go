package domain

import (
	"time"
	"uuid"
)

type Notification struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Title    string
	Content  string
	Channel  Channel
	Status   StatusNotification
	SentAt   time.Time
	CreateAt time.Time
	UpdateAt time.Time
}
