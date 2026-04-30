package app

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	ConfigFileName = "config.json"
	LogFileName    = "deutsch-tui.log"
)

type Config struct {
	Theme         string            `json:"theme"`
	Keymap        string            `json:"keymap"`
	AIProvider    string            `json:"ai_provider"`
	LogLevel      string            `json:"log_level"`
	AutoPlayAudio bool              `json:"autoplay_audio"`
	AITemplates   map[string]string `json:"ai_templates,omitempty"`
}

func DefaultConfig() Config {
	return Config{
		Theme:         "system",
		Keymap:        "default",
		AIProvider:    "disabled",
		LogLevel:      "info",
		AutoPlayAudio: false,
		AITemplates: map[string]string{
			"front":   "{{.Topic}}",
			"back":    "German prompt for {{.Topic}}.",
			"example": "Practice sentence using {{.Topic}}.",
		},
	}
}

func ResolveDataDir(override string) (string, error) {
	if override != "" {
		return override, nil
	}
	root, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "deutsch-tui"), nil
}

func EnsureDataDir(path string) error {
	if path == "" {
		return errors.New("data dir is required")
	}
	return os.MkdirAll(path, 0o755)
}

func LoadOrCreateConfig(dataDir string) (Config, error) {
	path := filepath.Join(dataDir, ConfigFileName)
	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		cfg := DefaultConfig()
		return cfg, SaveConfig(dataDir, cfg)
	}
	if err != nil {
		return Config{}, err
	}
	cfg := DefaultConfig()
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return Config{}, err
	}
	return cfg.withDefaults(), nil
}

func SaveConfig(dataDir string, cfg Config) error {
	raw, err := json.MarshalIndent(cfg.withDefaults(), "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	return os.WriteFile(filepath.Join(dataDir, ConfigFileName), raw, 0o644)
}

func (c Config) withDefaults() Config {
	defaults := DefaultConfig()
	if c.Theme == "" {
		c.Theme = defaults.Theme
	}
	if c.Keymap == "" {
		c.Keymap = defaults.Keymap
	}
	if c.AIProvider == "" {
		c.AIProvider = defaults.AIProvider
	}
	if c.LogLevel == "" {
		c.LogLevel = defaults.LogLevel
	}
	if c.AITemplates == nil {
		c.AITemplates = defaults.AITemplates
	}
	return c
}
