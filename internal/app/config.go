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
	AIProvider          string                       `json:"ai_provider"`
	DictionaryProvider  string                       `json:"dictionary_provider"`
	TTSProvider         string                       `json:"tts_provider"`
	TTSVoice            string                       `json:"tts_voice"`
	LogLevel            string                       `json:"log_level"`
	AutoPlayAudio       bool                         `json:"autoplay_audio"`
	StrictNormalization bool                         `json:"strict_normalization"`
	RevealSpeed         int                          `json:"reveal_speed"` // 0: instant, 1-10: slow to fast
	AITemplates         map[string]map[string]string `json:"ai_templates,omitempty"`
}

func DefaultConfig() Config {
	return Config{
		AIProvider:         "disabled",
		DictionaryProvider: "Local TUI",
		TTSProvider:        "edge",
		TTSVoice:           "de-DE-KatjaNeural",
		LogLevel:           "info",
		AutoPlayAudio:      false,
		RevealSpeed:        5,
		AITemplates: map[string]map[string]string{
			"vocabulary": {
				"front":   "{{.Topic}}",
				"back":    "Translation: German prompt for {{.Topic}}.\nPlural: die {{.Topic}}e (example)\nGender: der/die/das",
				"example": "Ich lerne {{.Topic}}.",
			},
			"grammar": {
				"front":   "Ich {{c1::...}} {{.Topic}}.",
				"back":    "Grammar: {{.Topic}}\nRule: Explanation of the rule for {{.Topic}}.",
				"example": "Ich {{c1::bin}} {{.Topic}}.",
			},
			"explanation": {
				"front":   "Explain the grammar rule for: {{.Topic}}",
				"back":    "Grammar Point: {{.Topic}}\n\nCore Rule: [Detailed explanation here]\nUsage: [When and how to use it]\nCommon Pitfalls: [Frequent mistakes]",
				"example": "Practice Sentence: [Example 1]\nPractice Sentence: [Example 2]",
			},
			"articles": {
				"front":   "{{c1::...}} {{.Topic}}",
				"back":    "Article for {{.Topic}}: der/die/das",
				"example": "MCQ:der,die,das",
			},
			"conjugation": {
				"front":   "{{.Topic}} (ich/du/er/sie/es) -> {{c1::...}}",
				"back":    "Conjugate {{.Topic}}\nIch {{c1::lerne}}, du {{c2::lernst}}, er {{c3::lernt}}.",
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
	if c.AIProvider == "" {
		c.AIProvider = defaults.AIProvider
	}
	if c.DictionaryProvider == "" {
		c.DictionaryProvider = defaults.DictionaryProvider
	}
	if c.TTSProvider == "" {
		c.TTSProvider = defaults.TTSProvider
	}
	if c.TTSVoice == "" {
		c.TTSVoice = defaults.TTSVoice
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
		AIProvider          string          `json:"ai_provider"`
		DictionaryProvider  string          `json:"dictionary_provider"`
		TTSProvider         string          `json:"tts_provider"`
		TTSVoice            string          `json:"tts_voice"`
		LogLevel            string          `json:"log_level"`
		AutoPlayAudio       bool            `json:"autoplay_audio"`
		StrictNormalization bool            `json:"strict_normalization"`
		RevealSpeed         int             `json:"reveal_speed"`
		AITemplates         json.RawMessage `json:"ai_templates,omitempty"`
	}

	var file configFile
	if err := json.Unmarshal(raw, &file); err != nil {
		return err
	}
	cfg.AIProvider = file.AIProvider
	cfg.DictionaryProvider = file.DictionaryProvider
	cfg.TTSProvider = file.TTSProvider
	cfg.TTSVoice = file.TTSVoice
	cfg.LogLevel = file.LogLevel
	cfg.AutoPlayAudio = file.AutoPlayAudio
	cfg.StrictNormalization = file.StrictNormalization
	cfg.RevealSpeed = file.RevealSpeed

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
