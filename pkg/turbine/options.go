package turbine

import (
	"maps"
	"sync"
)

type OptionConfig struct {
	mu                   sync.Mutex
	Name                 string
	PluginConfig         map[string]string
	PlatformPluginConfig []string
}

func (o *OptionConfig) Apply(opts ...Option) *OptionConfig {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.PluginConfig == nil {
		o.PluginConfig = make(map[string]string)
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

type Option func(*OptionConfig)

func WithName(pluginName string) Option {
	return func(cfg *OptionConfig) {
		cfg.Name = pluginName
	}
}

func WithPluginConfig(pluginConfig map[string]string) Option {
	return func(cfg *OptionConfig) {
		maps.Copy(cfg.PluginConfig, pluginConfig)
	}
}

func WithPlatformConfig(ref string) Option {
	return func(cfg *OptionConfig) {
		cfg.PlatformPluginConfig = append(cfg.PlatformPluginConfig, ref)
	}
}
