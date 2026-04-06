package task

import (
	"fmt"

	"github.com/flrnd/brancher/internal/config"
	"github.com/flrnd/brancher/internal/env"
)

func NewProvider(cfg *config.Config) (Provider, error) {
	def, ok := GetProviderDefinition(cfg.Provider)
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", cfg.Provider)
	}

	for _, envar := range def.Required {
		if env.Get(envar) == "" {
			return nil, fmt.Errorf("missing required environment variable: %s", envar)
		}
	}

	return def.New(cfg)
}
