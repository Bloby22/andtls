package config

import (
	"time"
)

// Config holds runtime configuration options for the andtls application
type Config struct {
	// ADBPath is the binary path or name for adb (default: "adb")
	ADBPath string

	// CommandTimeout is the maximum duration to wait for ADB commands
	CommandTimeout time.Duration

	// PollInterval is the duration between automatic device list refreshes
	PollInterval time.Duration
}

// DefaultConfig returns the standard default configuration
func DefaultConfig() Config {
	return Config{
		ADBPath:        "adb",
		CommandTimeout: 5 * time.Second,
		PollInterval:   1 * time.Second,
	}
}
