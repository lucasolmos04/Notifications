package sender

import (
	"Notifications/src/internal/domain"
	"Notifications/src/internal/service"
	"bytes"
	"context"
	"embed"
	"html/template"
	"net/mail"
	"time"
)

//go:embed templates/email.html
var emailTemplateFS embed.FS

var emailTmpl = template.Must(template.ParseFS(emailTemplateFS, "templates/email.html"))

type EmailSender struct{}

func (e *EmailSender) Validate(ctx context.Context, notification domain.Notification, user domain.User) error {
	email := user.Email
	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	}

	return nil
}

func (e *EmailSender) Send(ctx context.Context, notification domain.Notification, user domain.User) (*service.SenderResult, error) {

	type emailTemplateData struct {
		Title   string
		Content string
	}

	var buf bytes.Buffer
	data := emailTemplateData{
		Title:   notification.Title,
		Content: notification.Content,
	}

	if err := emailTmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	htmlContent := buf.String()

	sendResult := service.SenderResult{
		Channel: domain.ChannelEmail,
		SendAt:  time.Now(),
		Detail:  htmlContent,
	}

	return &sendResult, nil
}
