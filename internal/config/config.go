// Package config
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/flrnd/brancher/internal/provider"
	"gopkg.in/yaml.v3"
)

const (
	ConfigDir  = ".brancher"
	ConfigFile = "config.yml"
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

func Path() string {
	return filepath.Join(ConfigDir, ConfigFile)
}

func Load() (*Config, error) {
	path := Path()

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	path := Path()

	if err := os.MkdirAll(ConfigDir, 0o755); err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
