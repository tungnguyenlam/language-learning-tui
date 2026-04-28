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
	if _, err := os.Stat(filepath.Join(dir, ConfigFileName)); err != nil {
		t.Fatalf("config file was not created: %v", err)
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
