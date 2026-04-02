package notify

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/smtp"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/proxy"
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

func dialWithProxy(ctx context.Context, network, addr string) (net.Conn, error) {
	netDialer := &net.Dialer{Timeout: smtpDialTimeout}

	proxyURL := getProxyURL()
	if proxyURL == nil {
		return netDialer.DialContext(ctx, network, addr)
	}

	if shouldBypassProxy(proxyURL, addr) {
		return netDialer.DialContext(ctx, network, addr)
	}

	switch proxyURL.Scheme {
	case "socks5", "socks5h":
		dialer, err := proxy.FromURL(proxyURL, netDialer)
		if err != nil {
			return nil, err
		}
		if ctxDialer, ok := dialer.(proxy.ContextDialer); ok {
			return ctxDialer.DialContext(ctx, network, addr)
		}
		return dialer.Dial(network, addr)
	case "http", "https":
		return dialHTTPProxy(ctx, netDialer, proxyURL, addr)
	default:
		// unsupported scheme; fallback direct
		return netDialer.DialContext(ctx, network, addr)
	}
}

func getProxyURL() *url.URL {
	keys := []string{"ALL_PROXY", "all_proxy", "HTTPS_PROXY", "https_proxy", "HTTP_PROXY", "http_proxy"}
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			if !strings.Contains(value, "://") {
				value = "http://" + value
			}
			u, err := url.Parse(value)
			if err == nil {
				return u
			}
		}
	}
	return nil
}

func shouldBypassProxy(proxyURL *url.URL, targetAddr string) bool {
	noProxy := os.Getenv("NO_PROXY")
	if noProxy == "" {
		noProxy = os.Getenv("no_proxy")
	}
	if noProxy == "" {
		return false
	}

	host, _, err := net.SplitHostPort(targetAddr)
	if err != nil {
		host = targetAddr
	}

	for _, h := range strings.Split(noProxy, ",") {
		h = strings.TrimSpace(h)
		if h == "" {
			continue
		}
		if strings.EqualFold(h, host) || strings.HasPrefix(host, strings.TrimPrefix(h, "*")) {
			return true
		}
	}
	return false
}

func dialHTTPProxy(ctx context.Context, dialer *net.Dialer, proxyURL *url.URL, targetAddr string) (net.Conn, error) {
	proxyAddr := proxyURL.Host
	if proxyAddr == "" {
		return nil, fmt.Errorf("invalid proxy address")
	}

	conn, err := dialer.DialContext(ctx, "tcp", proxyAddr)
	if err != nil {
		return nil, err
	}

	req := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n", targetAddr, targetAddr)
	if user := proxyURL.User; user != nil {
		password, _ := user.Password()
		cred := fmt.Sprintf("%s:%s", user.Username(), password)
		req += "Proxy-Authorization: Basic " + base64.StdEncoding.EncodeToString([]byte(cred)) + "\r\n"
	}
	req += "Connection: close\r\n\r\n"

	if _, err := conn.Write([]byte(req)); err != nil {
		_ = conn.Close()
		return nil, err
	}

	br := bufio.NewReader(conn)
	tp := textproto.NewReader(br)
	statusLine, err := tp.ReadLine()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	var proto string
	var statusCode int
	_, err = fmt.Sscanf(statusLine, "HTTP/%s %d", &proto, &statusCode)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	_, err = tp.ReadMIMEHeader()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	if statusCode != http.StatusOK {
		_ = conn.Close()
		return nil, fmt.Errorf("proxy CONNECT failed: %d", statusCode)
	}

	return conn, nil
}

// Send sends an email notification.
func (e *EmailClient) Send(ctx context.Context, title, message string) error {
	addr := e.host + ":" + e.port

	conn, err := dialWithProxy(ctx, "tcp", addr)
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

	body := buildEmailBody(e.from, e.to, title, message)

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

func buildEmailBody(from string, to []string, subject, htmlBody string) string {
	return fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		strings.Join(to, ", "),
		from,
		subject,
		htmlBody,
	)
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
