package sender

import (
	"Notifications/src/internal/domain"
	"Notifications/src/internal/service"
	"context"
	"fmt"
	"time"
	"unicode/utf8"
)

type SmsSender struct{}

func (s *SmsSender) Validate(ctx context.Context, notification domain.Notification, user domain.User) error {
	text := notification.Content

	// 2. Con utf8.RuneCountInString() (cuenta caracteres reales)
	runesCount := utf8.RuneCountInString(text)

	maxCharacters := 160
	if runesCount > maxCharacters {
		return fmt.Errorf("Max lenght characters")
	}

	return nil
}

func (s *SmsSender) Send(ctx context.Context, notification domain.Notification, user domain.User) (*service.SenderResult, error) {

	sendResult := &service.SenderResult{
		Channel: domain.ChannelSMS,
		SendAt:  time.Now(),
		Detail:  notification.Content,
	}
	return sendResult, nil
}
