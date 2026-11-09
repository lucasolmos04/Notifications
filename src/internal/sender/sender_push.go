package sender

import (
	"Notifications/src/internal/domain"
	"Notifications/src/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type PushSender struct{}

func (e *PushSender) Validate(ctx context.Context, notification domain.Notification, user domain.User) error {
	if user.DeviceToken == "" {
		return fmt.Errorf("user has no device token registered")
	}
	return nil
}

func (p *PushSender) Send(ctx context.Context, n domain.Notification, user domain.User) (*service.SenderResult, error) {
	payload := struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}{
		Title: n.Title,
		Body:  n.Content,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return &service.SenderResult{}, err
	}

	sendResult := &service.SenderResult{
		Channel: domain.ChannelPush,
		SendAt:  time.Now(),
		Detail:  string(payloadJSON),
	}

	return sendResult, nil
}
