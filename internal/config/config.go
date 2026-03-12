// Package config
package config

type Config struct {
	Provider string        `yaml:"provider"`
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
