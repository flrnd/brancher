// Package config
package config

import (
	"fmt"

	"github.com/flrnd/brancher/internal/provider"
)

type Config struct {
	Provider provider.Name `yaml:"provider"`
	Project  ProjectConfig `yaml:"project"`
	Branch   BranchConfig  `yaml:"branch"`
}

type ProjectConfig struct {
	Owner string `yaml:"owner"`
	Repo  string `yaml:"repo"`
}

type BranchConfig struct {
	Strategy string `yaml:"strategy"`
}

type fieldRule struct {
	name  string
	field string
}

func (c *Config) Validate() error {
	rules := []fieldRule{
		{name: "provider", field: string(c.Provider)},
		{name: "project.owner", field: c.Project.Owner},
		{name: "project.repo", field: c.Project.Repo},
		{name: "branch.strategy", field: c.Branch.Strategy},
	}

	for _, rule := range rules {
		if rule.field == "" {
			return fmt.Errorf("%s is required", rule.name)
		}
	}

	return nil
}
