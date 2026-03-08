package client

import (
	"errors"
	"net/http"
	"time"
)

var (
	// ErrLoginFailed is returned when the login request fails.
	ErrLoginFailed = errors.New("login failed")
	// ErrGetTaskListFailed is returned when the task list request fails.
	ErrGetTaskListFailed = errors.New("get task list failed")
	// ErrGetLeaveListFailed is returned when the leave list request fails.
	ErrGetLeaveListFailed = errors.New("get leave list failed")
	// ErrSignFailed is returned when the sign request fails.
	ErrSignFailed = errors.New("sign failed")
	// ErrTaskNotFound is returned when a task is not found.
	ErrTaskNotFound = errors.New("task not found")
)

// Config holds the configuration for the Client.
type Config struct {
	Host      string `toml:"Host"`
	UserAgent string `toml:"UserAgent"`
	Latitude  string `toml:"Latitude"`
	Longitude string `toml:"Longitude"`
	SignSite  string `toml:"SignSite"`
}

// Client is the API client for the sign-in service.
type Client struct {
	httpClient *http.Client
	config     Config
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		Host:      DefaultHost,
		UserAgent: DefaultUserAgent,
		Latitude:  DefaultLatitude,
		Longitude: DefaultLongitude,
		SignSite:  DefaultSignSite,
	}
}

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithHTTPClient sets the HTTP client for the Client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithConfig sets the configuration for the Client.
// Note: This will override the entire configuration.
func WithConfig(cfg Config) Option {
	return func(c *Client) {
		c.config = cfg
	}
}

// New creates a new Client with the given configuration.
func New(opts ...Option) *Client {
	const defaultTimeout = 30 * time.Second
	c := &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		config:     DefaultConfig(),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) setCommonHeaders(req *http.Request) {
	req.Header.Set("User-Agent", c.config.UserAgent)
	req.Header.Set("Host", c.config.Host)
	req.Header.Set("Referer", "http://"+c.config.Host+"/Phone/index.html")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh-Hans;q=0.9")
	req.Header.Set("Connection", "keep-alive")
}
