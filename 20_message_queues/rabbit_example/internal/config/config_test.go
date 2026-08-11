package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	for _, key := range []string{
		"AMQP_URL", "PREFETCH", "RETRY_DELAY", "MAX_RETRIES",
		"PUBLISH_TIMEOUT", "SHUTDOWN_TIMEOUT",
	} {
		t.Setenv(key, "")
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Prefetch != 8 || cfg.MaxRetries != 3 || cfg.RetryDelay != 5*time.Second {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestLoadRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{name: "zero prefetch", key: "PREFETCH", value: "0"},
		{name: "negative retries", key: "MAX_RETRIES", value: "-1"},
		{name: "invalid retry delay", key: "RETRY_DELAY", value: "later"},
		{name: "zero publish timeout", key: "PUBLISH_TIMEOUT", value: "0s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, key := range []string{
				"PREFETCH", "RETRY_DELAY", "MAX_RETRIES",
				"PUBLISH_TIMEOUT", "SHUTDOWN_TIMEOUT",
			} {
				t.Setenv(key, "")
			}
			t.Setenv(tt.key, tt.value)
			if _, err := Load(); err == nil {
				t.Fatal("Load() error = nil, want validation error")
			}
		})
	}
}
