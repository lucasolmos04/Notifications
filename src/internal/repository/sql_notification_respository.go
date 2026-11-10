package repository

import (
	"Notifications/src/internal/domain"
	"Notifications/src/internal/service"
	"context"
	"database/sql"
	"errors"
	"time"
	"uuid"
)

type SQLNotificationRepository struct {
	db *sql.DB
}

func NewSQLNotificationRepository(db *sql.DB) service.NotificationRepository {
	return &SQLNotificationRepository{db: db}
}

func (s SQLNotificationRepository) Create(ctx context.Context, notification domain.Notification) error {
	n := fromDomainNotification(&notification)
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO notification (id, user_id, title, content, channel, status, sent_at, create_at, update_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		n.ID, n.UserID, n.Title, n.Content, n.Channel, n.Status, n.SentAt, n.CreateAt, n.UpdateAt)
	return err
}

func (s SQLNotificationRepository) Update(ctx context.Context, notification domain.Notification) error {
	n := fromDomainNotification(&notification)
	result, err := s.db.ExecContext(ctx,
		`UPDATE notification SET title=$1, content=$2, channel=$3, status=$4, sent_at=$5, update_at=$6 WHERE id=$7`,
		n.Title, n.Content, n.Channel, n.Status, n.SentAt, n.UpdateAt, n.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return service.ErrNotFound
	}
	return nil
}

func (s SQLNotificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM notification WHERE id = $1`, id.String())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return service.ErrNotFound
	}
	return nil
}

func (s SQLNotificationRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, user_id, title, content, channel, status, sent_at, create_at, update_at FROM notification WHERE id = $1`, id.String())

	var n NotificationDB
	if err := row.Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.Channel, &n.Status, &n.SentAt, &n.CreateAt, &n.UpdateAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrNotFound
		}
		return nil, err
	}
	return n.toDomain(), nil
}

func (s SQLNotificationRepository) FindByUserId(ctx context.Context, userId uuid.UUID) ([]domain.Notification, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, user_id, title, content, channel, status, sent_at, create_at, update_at FROM notification WHERE user_id = $1`, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var notifications []domain.Notification
	for rows.Next() {
		var n NotificationDB
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.Channel, &n.Status, &n.SentAt, &n.CreateAt, &n.UpdateAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, *n.toDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notifications, nil
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
