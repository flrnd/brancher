// Package branch
package branch

import (
	"github.com/flrnd/brancher/internal/task"
)

type Generator struct {
	strategy GeneratorStrategy
}

func NewGenerator(strategy GeneratorStrategy) Generator {
	return Generator{strategy: strategy}
}

func (g Generator) Generate(task task.Task) string {
	return g.strategy.Generate(task)
}
