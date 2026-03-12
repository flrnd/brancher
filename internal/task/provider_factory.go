package task

import (
	"fmt"
	"os"

	"github.com/flrnd/brancher/internal/config"
)

func NewProvider(cfg *config.Config) (Provider, error) {
	def, ok := GetProviderDefinition(cfg.Provider)
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", cfg.Provider)
	}

	for _, env := range def.Required {
		if os.Getenv(env) == "" {
			return nil, fmt.Errorf("missing required environment variable: %s", env)
		}
	}

	return def.New(cfg)
}
