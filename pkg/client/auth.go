// Package client provides the API client for the sign-in service.
package client

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Login authenticates the user and returns the access token.
func (c *Client) Login(ctx context.Context, studentID, password string) (string, error) {
	pubKey, err := parsePublicKey(PublicKeyBase64)
	if err != nil {
		return "", fmt.Errorf("parse public key: %w", err)
	}

	encPass, err := encryptPassword(password, pubKey)
	if err != nil {
		return "", fmt.Errorf("encrypt password: %w", err)
	}

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", studentID)
	data.Set("password", encPass)

	encodedBody := []byte(data.Encode())
	resp, err := c.doRequestWithRetry(ctx, true, requestSpec{
		method:   http.MethodPost,
		path:     pathLogin,
		rawQuery: "OpenId=",
		bodyBytes: func() []byte {
			return encodedBody
		},
		apply: func(req *http.Request) {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Origin", "https://"+c.config.Host)
		},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("%w: status %d, body: %s", ErrLoginFailed, resp.StatusCode, string(body))
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return "Bearer " + loginResp.AccessToken, nil
}

// parsePublicKey parses a Base64 encoded public key.
func parsePublicKey(publicKeyBase64 string) (*rsa.PublicKey, error) {
	der, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed: %w", err)
	}

	pub, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		block, _ := pem.Decode(der)
		if block == nil {
			return nil, fmt.Errorf("parse DER public key failed: %w", err)
		}
		pub, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse PEM public key failed: %w", err)
		}
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not RSA")
	}

	return rsaPub, nil
}

// encryptPassword encrypts the password using the public key.
func encryptPassword(password string, publicKey *rsa.PublicKey) (string, error) {
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(password))
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
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
