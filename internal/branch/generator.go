// Package branch
package branch

import "github.com/flrnd/brancher/internal/task"

type Generator struct {
	Strategy Strategy
}

func (g Generator) Generate(t task.Task) string {
	return ""
}
