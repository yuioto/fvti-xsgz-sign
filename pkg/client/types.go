package client

// TaskList represents the response from the task list API.
type TaskList struct {
	List List `json:"List"`
}

// List contains the items.
type List struct {
	Items []Item `json:"Items"`
}

// Item represents a single task item.
type Item struct {
	ID         string `json:"Id"`
	Name       string `json:"Name"`
	QD         string `json:"QD"`
	SignID     string `json:"SignID"`
	QDTimeText string `json:"QDTimeText"`
}

// LoginResponse represents the response from the login API.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	UserType     string `json:"UserType"`
	IsActive     bool   `json:"IsActive"`
	Msg          string `json:"Msg"`
}
