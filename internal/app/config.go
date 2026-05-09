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
	Theme               string                       `json:"theme"`
	Keymap              string                       `json:"keymap"`
	AIProvider          string                       `json:"ai_provider"`
	LogLevel            string                       `json:"log_level"`
	AutoPlayAudio       bool                         `json:"autoplay_audio"`
	StrictNormalization bool                         `json:"strict_normalization"`
	AITemplates         map[string]map[string]string `json:"ai_templates,omitempty"`
}

func DefaultConfig() Config {
	return Config{
		Theme:         "system",
		Keymap:        "default",
		AIProvider:    "disabled",
		LogLevel:      "info",
		AutoPlayAudio: false,
		AITemplates: map[string]map[string]string{
			"vocabulary": {
				"front":   "{{.Topic}}",
				"back":    "German prompt for {{.Topic}}.",
				"example": "Practice sentence using {{.Topic}}.",
			},
			"grammar": {
				"front":   "Ich {{c1::...}} {{.Topic}}.",
				"back":    "Grammar: {{.Topic}}",
				"example": "Ich {{c1::bin}} {{.Topic}}.",
			},
			"articles": {
				"front":   "{{c1::...}} {{.Topic}}",
				"back":    "Article for {{.Topic}}",
				"example": "MCQ:der,die,das",
			},
			"conjugation": {
				"front":   "{{.Topic}} (ich) -> {{c1::...}}",
				"back":    "Conjugate {{.Topic}}",
				"example": "Ich {{c1::lerne}} Deutsch.",
			},
			"passive_voice": {
				"front":   "Das Haus {{c1::...}} gebaut.",
				"back":    "Passive voice with 'werden'.",
				"example": "Das Haus {{c1::wird}} gebaut. MCQ:wird,wurde,worden,werden",
			},
			"subjunctive": {
				"front":   "Wenn ich Zeit {{c1::...}}, ...",
				"back":    "Konjunktiv II (hätte/wäre/würde).",
				"example": "Wenn ich Zeit {{c1::hätte}}, würde ich kommen. MCQ:hätte,wäre,würde,habe",
			},
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
	if err := unmarshalConfig(raw, &cfg); err != nil {
		return Config{}, err
	}
	return cfg.withDefaults(), nil
}

func SaveConfig(dataDir string, cfg Config) error {
	raw, err := json.MarshalIndent(cfg, "", "  ")
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
	} else {
		for set, template := range defaults.AITemplates {
			if _, ok := c.AITemplates[set]; !ok {
				c.AITemplates[set] = template
			}
		}
	}
	return c
}

func unmarshalConfig(raw []byte, cfg *Config) error {
	type configFile struct {
		Theme               string          `json:"theme"`
		Keymap              string          `json:"keymap"`
		AIProvider          string          `json:"ai_provider"`
		LogLevel            string          `json:"log_level"`
		AutoPlayAudio       bool            `json:"autoplay_audio"`
		StrictNormalization bool            `json:"strict_normalization"`
		AITemplates         json.RawMessage `json:"ai_templates,omitempty"`
	}

	var file configFile
	if err := json.Unmarshal(raw, &file); err != nil {
		return err
	}
	cfg.Theme = file.Theme
	cfg.Keymap = file.Keymap
	cfg.AIProvider = file.AIProvider
	cfg.LogLevel = file.LogLevel
	cfg.AutoPlayAudio = file.AutoPlayAudio
	cfg.StrictNormalization = file.StrictNormalization

	if len(file.AITemplates) == 0 || string(file.AITemplates) == "null" {
		return nil
	}
	var nested map[string]map[string]string
	if err := json.Unmarshal(file.AITemplates, &nested); err == nil {
		cfg.AITemplates = nested
		return nil
	}
	var legacy map[string]string
	if err := json.Unmarshal(file.AITemplates, &legacy); err != nil {
		return err
	}
	cfg.AITemplates = map[string]map[string]string{"vocabulary": legacy}
	return nil
}
