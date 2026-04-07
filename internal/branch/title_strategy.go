package branch

import (
	"strings"

	"github.com/flrnd/brancher/internal/task"
	"github.com/flrnd/brancher/pkg/slug"
)

type titleStrategy struct{}

func (titleStrategy) Generate(t task.Task) string {
	parts := []string{
		slug.Generate(t.ID),
		slug.Generate(t.Title),
	}

	filtered := parts[:0]
	for _, part := range parts {
		if part != "" {
			filtered = append(filtered, part)
		}
	}

	return strings.Join(filtered, "-")
}

func init() {
	RegisterStrategy(StrategyTitle, titleStrategy{})
}
