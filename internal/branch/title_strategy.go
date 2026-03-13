package branch

import (
	"github.com/flrnd/brancher/internal/task"
	"github.com/flrnd/brancher/pkg/slug"
)

type titleStrategy struct{}

func (titleStrategy) Generate(t task.Task) string {
	return slug.Generate(t.Title)
}

func init() {
	RegisterStrategy(StrategyTitle, titleStrategy{})
}
