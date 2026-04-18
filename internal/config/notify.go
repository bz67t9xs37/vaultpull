package config

import "github.com/vaultpull/internal/notify"

// NotifyConfig holds notification configuration.
type NotifyConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Channel  string `yaml:"channel"`
	FilePath string `yaml:"file_path"`
	OnlyDiff bool   `yaml:"only_diff"`
}

// DefaultNotifyConfig returns sensible defaults.
func DefaultNotifyConfig() *NotifyConfig {
	return &NotifyConfig{
		Enabled:  false,
		Channel:  string(notify.ChannelStdout),
		OnlyDiff: true,
	}
}

// ApplyNotifyDefaults fills zero-value fields with defaults.
func ApplyNotifyDefaults(cfg *NotifyConfig) {
	if cfg == nil {
		return
	}
	if cfg.Channel == "" {
		cfg.Channel = string(notify.ChannelStdout)
	}
}

// ToNotifyConfig converts to the notify package Config.
func (c *NotifyConfig) ToNotifyConfig() notify.Config {
	if c == nil {
		return notify.Config{}
	}
	return notify.Config{
		Enabled:  c.Enabled,
		Channel:  notify.Channel(c.Channel),
		FilePath: c.FilePath,
		OnlyDiff: c.OnlyDiff,
	}
}
