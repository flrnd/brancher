package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	ConfigDir  = ".brancher"
	ConfigFile = "config.yml"
)

func Load() (*Config, error) {
	path := filepath.Join(ConfigDir, ConfigFile)

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
