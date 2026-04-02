package notify

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

// EmailClient is a client for sending email notifications.
var (
	// ErrEmailSendFailed is returned when sending email notification fails.
	ErrEmailSendFailed = fmt.Errorf("email send failed")
)

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

const smtpDialTimeout = 15 * time.Second

// Send sends an email notification.
func (e *EmailClient) Send(ctx context.Context, title, message string) error {
	addr := e.host + ":" + e.port

	dialer := &net.Dialer{Timeout: smtpDialTimeout}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrEmailSendFailed, err)
	}

	client, err := smtp.NewClient(conn, e.host)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("%w: %v", ErrEmailSendFailed, err)
	}
	defer client.Close()

	if err := client.StartTLS(&tls.Config{ServerName: e.host}); err != nil {
		return fmt.Errorf("%w: %v", ErrEmailSendFailed, err)
	}

	auth := smtp.PlainAuth("", e.username, e.password, e.host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("%w: %v", ErrEmailSendFailed, err)
	}

	if err := client.Mail(e.from); err != nil {
		return fmt.Errorf("%w: %v", ErrEmailSendFailed, err)
	}

	for _, rcpt := range e.to {
		if err := client.Rcpt(rcpt); err != nil {
			return fmt.Errorf("%w: %v", ErrEmailSendFailed, err)
		}
	}

	body := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n\r\n%s",
		strings.Join(e.to, ", "),
		e.from,
		title,
		message,
	)

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrEmailSendFailed, err)
	}

	if _, err := fmt.Fprint(w, body); err != nil {
		_ = w.Close()
		return fmt.Errorf("%w: %v", ErrEmailSendFailed, err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("%w: %v", ErrEmailSendFailed, err)
	}

	return client.Quit()
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
