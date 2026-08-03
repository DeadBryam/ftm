package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	AppName              = "foundry-tunnel"
	ConfigFile           = "config.yaml"
	NotificationPending  = "pending"
	NotificationGranted  = "granted"
	NotificationRejected = "rejected"

	// ConfigVersion is the current on-disk schema version. Bump it whenever a
	// stored config needs rewriting on load, and add a step to migrate.
	ConfigVersion = 2

	// LanguageAuto follows the operating system language. It is the default so
	// that a config which was never touched does not look like a deliberate
	// choice of English.
	LanguageAuto = "auto"
)

type Config struct {
	Version  int            `yaml:"version"`
	Tunnels  []TunnelConfig `yaml:"tunnels"`
	WebPort  int            `yaml:"web_port,omitempty"`
	Language string         `yaml:"language"`

	NotificationsStatus string `yaml:"notifications_status"`
	NotificationSound   bool   `yaml:"notification_sound"`

	ExpirationThresholds      []int          `yaml:"expiration_thresholds"`
	ProviderExpirationMinutes map[string]int `yaml:"provider_expiration_minutes"`
}

func DefaultConfig() *Config {
	return &Config{
		Version:              ConfigVersion,
		Language:             LanguageAuto,
		Tunnels:              []TunnelConfig{},
		NotificationsStatus:  NotificationPending,
		NotificationSound:    true,
		ExpirationThresholds: []int{30, 15, 10, 5, 1},
		ProviderExpirationMinutes: map[string]int{
			"pinggy":       60,
			"serveo":       0,
			"cloudflared":  0,
			"tunnelmole":   0,
			"localhostrun": 0,
		},
	}
}

// migrate brings a stored config up to ConfigVersion, reporting whether
// anything changed and therefore needs writing back.
func (c *Config) migrate() bool {
	changed := false

	// v1 wrote `language: en` as its default. Because InitFromConfig only
	// detected the system language when the field was empty, that default made
	// every install look like the user had explicitly chosen English and left
	// detection permanently dead. A v1 config never offered that choice, so
	// hand these back to auto-detection.
	if c.Version < 2 {
		if c.Language == "" || c.Language == "en" {
			c.Language = LanguageAuto
		}
		changed = true
	}

	if c.Version != ConfigVersion {
		c.Version = ConfigVersion
		changed = true
	}

	return changed
}

func (c *Config) NormalizeNotificationsStatus() {
	if c.NotificationsStatus != NotificationGranted && c.NotificationsStatus != NotificationRejected {
		c.NotificationsStatus = NotificationPending
	}
}

func ConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".config", AppName)
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), ConfigFile)
}

func (c *Config) Save() error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	path := ConfigPath()
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func Load() (*Config, error) {
	path := ConfigPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		cfg := DefaultConfig()
		if err := cfg.Save(); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	cfg.NormalizeNotificationsStatus()

	if cfg.migrate() {
		if err := cfg.Save(); err != nil {
			return nil, fmt.Errorf("failed to persist migrated config: %w", err)
		}
	}

	return &cfg, nil
}

func (c *Config) GetTunnel(id string) *TunnelConfig {
	for i := range c.Tunnels {
		if c.Tunnels[i].ID == id {
			return &c.Tunnels[i]
		}
	}
	return nil
}

func (c *Config) AddTunnel(t TunnelConfig) {
	c.Tunnels = append(c.Tunnels, t)
}

func (c *Config) RemoveTunnel(id string) bool {
	for i, t := range c.Tunnels {
		if t.ID == id {
			c.Tunnels = append(c.Tunnels[:i], c.Tunnels[i+1:]...)
			return true
		}
	}
	return false
}
