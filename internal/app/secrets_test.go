package app

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSecretsRoundTripPreservesValues(t *testing.T) {
	dir := t.TempDir()
	s := Secrets{
		OpenAI: ProviderCreds{
			APIKey:  "sk-test-1",
			Model:   "gpt-4o-mini",
			BaseURL: "https://api.openai.com/v1",
		},
		Anthropic: ProviderCreds{
			APIKey: "sk-ant-test",
			Model:  "claude-haiku-4-5-20251001",
		},
	}
	if err := SaveSecrets(dir, s); err != nil {
		t.Fatalf("SaveSecrets: %v", err)
	}
	got, err := LoadSecrets(dir)
	if err != nil {
		t.Fatalf("LoadSecrets: %v", err)
	}
	if got.OpenAI != s.OpenAI {
		t.Errorf("OpenAI mismatch: got %+v want %+v", got.OpenAI, s.OpenAI)
	}
	if got.Anthropic != s.Anthropic {
		t.Errorf("Anthropic mismatch: got %+v want %+v", got.Anthropic, s.Anthropic)
	}
}

func TestSecretsLoadCreatesEmptyFile(t *testing.T) {
	dir := t.TempDir()
	got, err := LoadSecrets(dir)
	if err != nil {
		t.Fatalf("LoadSecrets: %v", err)
	}
	if got.OpenAI.APIKey != "" || got.Anthropic.APIKey != "" {
		t.Errorf("expected empty secrets on first load, got %+v", got)
	}
	// secrets.json should exist after first load.
	if _, err := os.Stat(filepath.Join(dir, SecretsFileName)); err != nil {
		t.Errorf("secrets.json should exist after LoadSecrets: %v", err)
	}
}

func TestSecretsFileMode0600(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX file modes not enforced on Windows")
	}
	dir := t.TempDir()
	if err := SaveSecrets(dir, Secrets{OpenAI: ProviderCreds{APIKey: "x"}}); err != nil {
		t.Fatalf("SaveSecrets: %v", err)
	}
	info, err := os.Stat(filepath.Join(dir, SecretsFileName))
	if err != nil {
		t.Fatalf("stat secrets.json: %v", err)
	}
	mode := info.Mode().Perm()
	if mode != 0o600 {
		t.Errorf("secrets.json mode = %o, want 0600 (group/world readable would leak API keys)", mode)
	}
}

func TestMaskAPIKey(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", "(not set)"},
		{"abc", "****"},
		{"abcd", "****"},
		{"sk-test1234", "****1234"},
		{"verylongapikey0001", "****0001"},
	}
	for _, c := range cases {
		got := MaskAPIKey(c.in)
		if got != c.want {
			t.Errorf("MaskAPIKey(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
