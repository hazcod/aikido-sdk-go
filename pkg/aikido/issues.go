package aikido

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	ISSUE_URL = API_BASE_URL + "public/v1/issues/export"
)

type Issue struct {
	ID                int    `json:"id"`
	GroupID           int    `json:"group_id"`
	AttackSurface     string `json:"attack_surface"`
	Status            string `json:"status"`
	Severity          string `json:"severity"`
	SeverityScore     int    `json:"severity_score"`
	Type              string `json:"type"`
	Rule              string `json:"rule"`
	AffectedPackage   string `json:"affected_package"`
	AffectedFile      string `json:"affected_file"`
	FirstDetectedAt   int    `json:"first_detected_at"`
	CodeRepoName      string `json:"code_repo_name"`
	CodeRepoID        int    `json:"code_repo_id"`
	ContainerRepoID   int    `json:"container_repo_id"`
	ContainerRepoName string `json:"container_repo_name"`
	SLADays           int    `json:"sla_days"`
	SLARemediateBy    int    `json:"sla_remediate_by"`
	IgnoredAt         any    `json:"ignored_at"`
	ClosedAt          any    `json:"closed_at"`
}

func toTitle(s string) string {
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func (i *Issue) GetName() string {
	return toTitle(i.Severity) + " " + i.Rule + " in " + i.AffectedPackage + " in repo " + i.CodeRepoName + " "
}

func (i *Issue) IsIgnored() bool {
	return strings.EqualFold(i.Status, "ignored")
}

func (a *Aikido) GetIssues(onlyOpen bool) ([]Issue, error) {
	issues := make([]Issue, 0)

	token, err := a.authToken.GetToken()
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	fullURL := ISSUE_URL
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
