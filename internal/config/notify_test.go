package config_test

import (
	"testing"

	"github.com/vaultpull/internal/config"
	"github.com/vaultpull/internal/notify"
)

func TestDefaultNotifyConfig_Values(t *testing.T) {
	cfg := config.DefaultNotifyConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.Channel != string(notify.ChannelStdout) {
		t.Errorf("expected stdout channel, got %s", cfg.Channel)
	}
	if !cfg.OnlyDiff {
		t.Error("expected OnlyDiff to be true by default")
	}
}

func TestApplyNotifyDefaults_FillsChannel(t *testing.T) {
	cfg := &config.NotifyConfig{}
	config.ApplyNotifyDefaults(cfg)
	if cfg.Channel != string(notify.ChannelStdout) {
		t.Errorf("expected channel filled, got %s", cfg.Channel)
	}
}

func TestApplyNotifyDefaults_PreservesExistingChannel(t *testing.T) {
	cfg := &config.NotifyConfig{Channel: string(notify.ChannelStderr)}
	config.ApplyNotifyDefaults(cfg)
	if cfg.Channel != string(notify.ChannelStderr) {
		t.Errorf("expected stderr preserved, got %s", cfg.Channel)
	}
}

func TestApplyNotifyDefaults_NilSafe(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panicked on nil config: %v", r)
		}
	}()
	config.ApplyNotifyDefaults(nil)
}

func TestToNotifyConfig_MapsFields(t *testing.T) {
	cfg := &config.NotifyConfig{
		Enabled:  true,
		Channel:  "file",
		FilePath: "/tmp/notify.log",
		OnlyDiff: false,
	}
	nc := cfg.ToNotifyConfig()
	if !nc.Enabled {
		t.Error("expected Enabled true")
	}
	if nc.Channel != notify.ChannelFile {
		t.Errorf("expected file channel, got %s", nc.Channel)
	}
	if nc.FilePath != "/tmp/notify.log" {
		t.Errorf("unexpected file path: %s", nc.FilePath)
	}
}

func TestToNotifyConfig_NilReturnsEmpty(t *testing.T) {
	var cfg *config.NotifyConfig
	nc := cfg.ToNotifyConfig()
	if nc.Enabled {
		t.Error("expected empty config")
	}
}
