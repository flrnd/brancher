package task

import (
	"github.com/flrnd/brancher/internal/config"
	"github.com/flrnd/brancher/internal/provider"
)

type Definition struct {
	Name     provider.Name
	Required []string
	New      func(*config.Config) (Provider, error)
}

var registry = map[provider.Name]Definition{}

func RegisterProvider(def Definition) {
	registry[def.Name] = def
}

func GetProviderDefinition(name provider.Name) (Definition, bool) {
	def, ok := registry[name]
	return def, ok
}
