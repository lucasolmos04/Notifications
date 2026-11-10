package service

import (
	"Notifications/src/internal/domain"
	"context"
	"fmt"
	"time"
	"uuid"
)

type SenderRegistry map[domain.Channel]NotificationSender

type NotificationService struct {
	repo     NotificationRepository
	userRepo UserRepository
	senders  SenderRegistry
}

func NewNotificationService(repo NotificationRepository, userRepo UserRepository, senders SenderRegistry) *NotificationService {
	return &NotificationService{
		repo:     repo,
		userRepo: userRepo,
		senders:  senders,
	}
}

func (s *NotificationService) Create(ctx context.Context, req CreateNotificationRequest) (*domain.Notification, error) {
	sender, ok := s.senders[req.Channel]
	if !ok {
		return nil, fmt.Errorf("unsupported channel: %s", req.Channel)
	}

	user, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	notification := domain.Notification{
		ID:       uuid.New(),
		UserID:   req.UserID,
		Title:    req.Title,
		Content:  req.Content,
		Channel:  req.Channel,
		Status:   domain.Pending,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	if err := sender.Validate(ctx, notification, *user); err != nil {
		return nil, err
	}

	// Create notification
	if err := s.repo.Create(ctx, notification); err != nil {
		return nil, err
	}

	// Send notification
	result, sendErr := sender.Send(ctx, notification, *user)
	if sendErr != nil {
		notification.Status = domain.Failed
	} else {
		notification.Status = domain.Sent
		notification.SentAt = result.SendAt
	}
	notification.UpdateAt = time.Now()

	// Update notification after sent it
	if err := s.repo.Update(ctx, notification); err != nil {
		return nil, err
	}
	return &notification, nil
}

func (s *NotificationService) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Notification, error) {
	return s.repo.FindByUserId(ctx, userID)
}

func (s *NotificationService) Update(ctx context.Context, userID uuid.UUID, notificationID uuid.UUID, req UpdateNotificationRequest) (*domain.Notification, error) {
	if _, ok := s.senders[req.Channel]; !ok {
		return nil, fmt.Errorf("unsupported channel: %s", req.Channel)
	}

	notification, err := s.repo.FindByID(ctx, notificationID)
	if err != nil {
		return nil, err
	}
	if notification.UserID != userID {
		return nil, ErrNotFound
	}

	notification.Title = req.Title
	notification.Content = req.Content
	notification.Channel = req.Channel
	notification.UpdateAt = time.Now()

	if err := s.repo.Update(ctx, *notification); err != nil {
		return nil, err
	}
	return notification, nil
}

func (s *NotificationService) Delete(ctx context.Context, userID uuid.UUID, notificationID uuid.UUID) error {
	notification, err := s.repo.FindByID(ctx, notificationID)
	if err != nil {
		return err
	}
	if notification.UserID != userID {
		return ErrNotFound
	}

	return s.repo.Delete(ctx, notificationID)
}
