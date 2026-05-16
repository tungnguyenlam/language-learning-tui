package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadOrCreateConfig(t *testing.T) {
	dir := t.TempDir()
	if err := EnsureDataDir(dir); err != nil {
		t.Fatalf("ensure data dir: %v", err)
	}

	cfg, err := LoadOrCreateConfig(dir)
	if err != nil {
		t.Fatalf("load or create config: %v", err)
	}
	if cfg.AIProvider != "disabled" {
		t.Fatalf("AIProvider = %q, want disabled", cfg.AIProvider)
	}
	if cfg.TTSProvider != "edge" {
		t.Fatalf("TTSProvider = %q, want edge", cfg.TTSProvider)
	}
	if cfg.TTSVoice != "de-DE-KatjaNeural" {
		t.Fatalf("TTSVoice = %q, want de-DE-KatjaNeural", cfg.TTSVoice)
	}
	if _, err := os.Stat(filepath.Join(dir, ConfigFileName)); err != nil {
		t.Fatalf("config file was not created: %v", err)
	}
}

func TestLoadOrCreateConfigMigratesLegacyAITemplates(t *testing.T) {
	dir := t.TempDir()
	if err := EnsureDataDir(dir); err != nil {
		t.Fatalf("ensure data dir: %v", err)
	}
	raw := []byte(`{
  "theme": "system",
  "keymap": "default",
  "ai_provider": "template",
  "log_level": "info",
  "ai_templates": {
    "front": "{{.Topic}}",
    "back": "German prompt for {{.Topic}}.",
    "example": "Practice sentence using {{.Topic}}."
  }
}`)
	if err := os.WriteFile(filepath.Join(dir, ConfigFileName), raw, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadOrCreateConfig(dir)
	if err != nil {
		t.Fatalf("load legacy config: %v", err)
	}
	if got := cfg.AITemplates["vocabulary"]["front"]; got != "{{.Topic}}" {
		t.Fatalf("vocabulary front = %q, want legacy template", got)
	}
	if _, ok := cfg.AITemplates["grammar"]; !ok {
		t.Fatal("grammar default template was not restored")
	}
}

func TestLoadOrCreateConfigKeepsNestedAITemplates(t *testing.T) {
	dir := t.TempDir()
	if err := EnsureDataDir(dir); err != nil {
		t.Fatalf("ensure data dir: %v", err)
	}
	raw := []byte(`{
  "ai_templates": {
    "custom": {
      "front": "Front {{.Topic}}",
      "back": "Back {{.Topic}}",
      "example": "Example {{.Topic}}"
    }
  }
}`)
	if err := os.WriteFile(filepath.Join(dir, ConfigFileName), raw, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadOrCreateConfig(dir)
	if err != nil {
		t.Fatalf("load nested config: %v", err)
	}
	if got := cfg.AITemplates["custom"]["front"]; got != "Front {{.Topic}}" {
		t.Fatalf("custom front = %q, want nested template", got)
	}
	if _, ok := cfg.AITemplates["vocabulary"]; !ok {
		t.Fatal("vocabulary default template was not restored")
	}
}

func TestOpenLog(t *testing.T) {
	dir := t.TempDir()
	file, logger, err := OpenLog(dir)
	if err != nil {
		t.Fatalf("open log: %v", err)
	}
	logger.Print("hello")
	if err := file.Close(); err != nil {
		t.Fatalf("close log: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, LogFileName)); err != nil {
		t.Fatalf("log file was not created: %v", err)
	}
}

func TestLoadOrCreateConfigDisabledProviderPreserved(t *testing.T) {
	dir := t.TempDir()
	if err := EnsureDataDir(dir); err != nil {
		t.Fatalf("ensure data dir: %v", err)
	}
	raw := []byte(`{
  "ai_provider": "disabled"
}`)
	if err := os.WriteFile(filepath.Join(dir, ConfigFileName), raw, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadOrCreateConfig(dir)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.AIProvider != "disabled" {
		t.Fatalf("AIProvider = %q, want disabled", cfg.AIProvider)
	}
}
