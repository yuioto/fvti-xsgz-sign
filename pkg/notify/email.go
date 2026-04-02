package notify

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

// EmailClient is a client for sending email notifications.
type EmailClient struct {
	host     string
	port     string
	username string
	password string
	from     string
	to       []string
}

// EmailConfig holds the configuration for the email client.
type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	To       string
}

// NewEmail creates a new email notification client.
func NewEmail(cfg EmailConfig) *EmailClient {
	to := splitAndTrim(cfg.To, ",")
	return &EmailClient{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
		from:     cfg.From,
		to:       to,
	}
}

// Send sends an email notification.
func (e *EmailClient) Send(_ context.Context, title, message string) error {
	addr := e.host + ":" + e.port

	auth := smtp.PlainAuth("", e.username, e.password, e.host)

	body := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n\r\n%s",
		strings.Join(e.to, ", "),
		e.from,
		title,
		message,
	)

	if err := smtp.SendMail(addr, auth, e.from, e.to, []byte(body)); err != nil {
		return fmt.Errorf("send email failed: %w", err)
	}
	return nil
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
