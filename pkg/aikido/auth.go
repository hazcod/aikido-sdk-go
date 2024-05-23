package aikido

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	API_BASE_URL = "https://app.aikido.dev/api/"
)

// TokenResponse represents the JSON structure of the response from the token endpoint
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type tokenGenerator struct {
	value   string
	expires time.Time

	clientSecret string
	clientID     string
}

func newTokenGenerator(clientID, clientSecret string) (*tokenGenerator, error) {
	tokenResponse, err := getNewAccessToken(clientID, clientSecret)
	if err != nil {
		return nil, fmt.Errorf("access token failed: %w", err)
	}

	gen := tokenGenerator{
		clientID:     clientID,
		clientSecret: clientSecret,
		value:        tokenResponse.AccessToken,
		expires:      time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second),
	}

	if gen.isExpired() {
		return nil, errors.New("access token is already expired")
	}

	return &gen, nil
}

func (c *tokenGenerator) isExpired() bool {
	return c.expires.Before(time.Now())
}

func (c *tokenGenerator) GetToken() (string, error) {
	if c.isExpired() {
		newToken, err := getNewAccessToken(c.clientID, c.clientSecret)
		if err != nil {
			return "", fmt.Errorf("could not fetch new token: %w", err)
		}

		c.value = newToken.AccessToken
		c.expires = time.Now().Add(time.Duration(newToken.ExpiresIn) * time.Second)
	}

	return c.value, nil
}

// GetAccessToken fetches an OAuth2 token using client credentials
func getNewAccessToken(clientID, clientSecret string) (*tokenResponse, error) {
	url := API_BASE_URL + "oauth/token"

	// Create the Basic auth header
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))

	// Create the request body
	body := []byte("grant_type=client_credentials")
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, respBody)
	}

	// Parse the JSON response
	var tokenResponse tokenResponse
	if err := json.Unmarshal(respBody, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	if tokenResponse.AccessToken == "" {
		return nil, fmt.Errorf("no access token returned")
	}

	return &tokenResponse, nil
}
