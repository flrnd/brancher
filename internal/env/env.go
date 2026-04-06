// Package env
package env

import (
	"os"
	"sync"

	"github.com/flrnd/brancher/internal/provider"
	"github.com/joho/godotenv"
)

const (
	GitHubToken = "BRANCHER_GITHUB_TOKEN"
	JiraToken   = "BRANCHER_JIRA_TOKEN"
	GitLabToken = "BRANCHER_GITLAB_TOKEN"
)

var loadOnce sync.Once

func Load() {
	loadOnce.Do(func() {
		_ = godotenv.Load()
	})
}

func Get(key string) string {
	Load()
	return os.Getenv(key)
}

func ProviderToken(name provider.Name) string {
	switch name {
	case provider.GitHub:
		return GitHubToken
	default:
		return ""
	}
}
