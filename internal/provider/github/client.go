package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	baseURL   = "https://api.github.com"
	userAgent = "brancher"
)

type client struct {
	httpClient *http.Client
	token      string
	owner      string
	repo       string
}

type issue struct {
	Number      int    `json:"number"`
	Title       string `json:"title"`
	State       string `json:"state"`
	HTMLURL     string `json:"html_url"`
	Labels      []label
	PullRequest *struct{} `json:"pull_request,omitempty"`
}

type label struct {
	Name string `json:"name"`
}

func newClient(token, owner, repo string) *client {
	return &client{
		httpClient: http.DefaultClient,
		token:      token,
		owner:      owner,
		repo:       repo,
	}
}

func (c *client) listIssues(ctx context.Context) ([]issue, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/repos/%s/%s/issues?state=open&per_page=100", baseURL, c.owner, c.repo),
		nil,
	)
	if err != nil {
		return nil, err
	}

	var issues []issue
	if err := c.doJSON(req, &issues); err != nil {
		return nil, err
	}

	filtered := issues[:0]
	for _, it := range issues {
		if it.PullRequest == nil {
			filtered = append(filtered, it)
		}
	}

	return filtered, nil
}

func (c *client) getIssue(ctx context.Context, id string) (issue, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/repos/%s/%s/issues/%s", baseURL, c.owner, c.repo, id),
		nil,
	)
	if err != nil {
		return issue{}, err
	}

	var out issue
	if err := c.doJSON(req, &out); err != nil {
		return issue{}, err
	}

	if out.PullRequest != nil {
		return issue{}, fmt.Errorf("task %s is a pull request, not an issue", id)
	}

	return out, nil
}

func (c *client) doJSON(req *http.Request, dst any) error {
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		msg := strings.TrimSpace(string(body))
		if msg == "" {
			msg = resp.Status
		}
		return fmt.Errorf("github api request failed: %s", msg)
	}

	if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
		return fmt.Errorf("decode github response: %w", err)
	}

	return nil
}
