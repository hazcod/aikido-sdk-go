package aikido

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	ISSUES_URL = API_BASE_URL + "public/v1/issues/export"
)

func (a *Aikido) GetIssues(onlyOpen bool) ([]Issue, error) {
	issues := make([]Issue, 0)

	token, err := a.authToken.GetToken()
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	fullURL := ISSUES_URL
	if onlyOpen {
		fullURL = fullURL + "?filter_status=open"
	}

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues: %w", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read issues: %w", err)
	}

	if err := json.Unmarshal(body, &issues); err != nil {
		return nil, fmt.Errorf("failed to decode issues: %w", err)
	}

	return issues, nil
}
