package branch

import (
	"fmt"

	"github.com/flrnd/brancher/internal/task"
)

type Strategy string

const (
	StrategyTitle Strategy = "title"
)

type GeneratorStrategy interface {
	Generate(task task.Task) string
}

var strategies = map[Strategy]GeneratorStrategy{}

func AvailableStrategies() []Strategy {
	var list []Strategy

	for s := range strategies {
		list = append(list, s)
	}

	return list
}

func RegisterStrategy(name Strategy, strategy GeneratorStrategy) {
	if _, exists := strategies[name]; exists {
		panic("branch strategy already registered: " + string(name))
	}

	strategies[name] = strategy
}

func ResolveStrategy(name Strategy) (GeneratorStrategy, error) {
	s, ok := strategies[name]
	if !ok {
		return nil, fmt.Errorf("unknown branch strategy: %s", name)
	}
	return s, nil
}
