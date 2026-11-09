package repository

import (
	"Notifications/src/internal/domain"
	"context"
	"database/sql"
	"time"
	"uuid"
)

type SQLNotificationRepository struct {
	db *sql.DB
}

func NewSQLNotificationRepository(db *sql.DB) domain.NotificationRepository {
	return &SQLNotificationRepository{db: db}
}

func (S SQLNotificationRepository) Create(ctx context.Context, notification domain.Notification) {
	//TODO implement me
	panic("implement me")
}

func (S SQLNotificationRepository) Update(ctx context.Context, notification domain.Notification) {
	//TODO implement me
	panic("implement me")
}

func (S SQLNotificationRepository) Delete(ctx context.Context, id uuid.UUID) {
	//TODO implement me
	panic("implement me")
}

func (S SQLNotificationRepository) FindByUserId(ctx context.Context, userId uuid.UUID) {
	//TODO implement me
	panic("implement me")
}

type NotificationDB struct {
	ID       string
	UserID   string
	Title    string
	Content  string
	Channel  string
	Status   string
	SentAt   time.Time
	CreateAt time.Time
	UpdateAt time.Time
}

func (n *NotificationDB) toDomain() *domain.Notification {
	return &domain.Notification{
		ID:       uuid.MustParse(n.ID),
		UserID:   uuid.MustParse(n.UserID),
		Title:    n.Title,
		Content:  n.Content,
		Channel:  domain.Channel(n.Channel),
		Status:   domain.StatusNotification(n.Status),
		SentAt:   n.SentAt,
		CreateAt: n.CreateAt,
		UpdateAt: n.UpdateAt,
	}
}

func fromDomainNotification(n *domain.Notification) *NotificationDB {
	return &NotificationDB{
		ID:       n.ID.String(),
		UserID:   n.UserID.String(),
		Title:    n.Title,
		Content:  n.Content,
		Channel:  string(n.Channel),
		Status:   string(n.Status),
		SentAt:   n.SentAt,
		CreateAt: n.CreateAt,
		UpdateAt: n.UpdateAt,
	}
}
