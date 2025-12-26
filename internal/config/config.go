package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/faisal/yo/internal/state"
)

// Config represents user configuration
type Config struct {
	WatchDirs     []string `json:"watch_dirs"`
	Notifications bool     `json:"notifications"`
	Editor        string   `json:"editor"`
	MaxBypassDay  int      `json:"max_bypass_day"`
	MaxBypassWeek int      `json:"max_bypass_week"`
}

// Default returns the default configuration
func Default() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		WatchDirs:     []string{filepath.Join(homeDir, "Dev")},
		Notifications: true,
		Editor:        os.Getenv("EDITOR"),
		MaxBypassDay:  1,
		MaxBypassWeek: 5,
	}
}

// getConfigPath returns the path to config.json
func getConfigPath() (string, error) {
	yoDir, err := state.GetYoDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(yoDir, "config.json"), nil
}

// Load loads configuration from disk
func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg := Default()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

// Save saves configuration to disk
func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// Set sets a configuration value
func (c *Config) Set(key, value string) error {
	switch key {
	case "notifications":
		c.Notifications = value == "on" || value == "true"
	case "editor":
		c.Editor = value
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
	return nil
}

// Get gets a configuration value
func (c *Config) Get(key string) (string, error) {
	switch key {
	case "notifications":
		if c.Notifications {
			return "on", nil
		}
		return "off", nil
	case "editor":
		return c.Editor, nil
	case "watch_dirs":
		data, _ := json.Marshal(c.WatchDirs)
		return string(data), nil
	default:
		return "", fmt.Errorf("unknown config key: %s", key)
	}
}
