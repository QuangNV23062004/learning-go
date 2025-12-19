package utils

import (
	"bytes"
	"fmt"
	"learning-go/internal/config"
	"text/template"

	"embed"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config *config.MailConfig
}

func NewEmailService(mailConfig *config.MailConfig) *EmailService {
	return &EmailService{
		config: mailConfig,
	}
}

func (e *EmailService) SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.config.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	port := e.config.Port
	d := gomail.NewDialer(
		e.config.Host,
		port,
		e.config.Username,
		e.config.Password,
	)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// RenderEmailTemplate renders an HTML email template with the provided data
func (e *EmailService) RenderEmailTemplate(name string, data any) (string, error) {
	// Parse the template file
	tpl, err := template.ParseFS(templatesFS, "internal/pkg/users/templates/"+name)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template with data
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

var templatesFS embed.FS
