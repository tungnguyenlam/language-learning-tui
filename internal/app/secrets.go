package app

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// SecretsFileName holds API keys and provider-specific knobs (model,
// custom base URL). It lives next to config.json but is written at 0600
// so other users on the host can't read the credentials.
const SecretsFileName = "secrets.json"

// ProviderCreds is the per-provider credential block. BaseURL is optional;
// when empty the provider falls back to its canonical endpoint.
type ProviderCreds struct {
	APIKey  string `json:"api_key,omitempty"`
	Model   string `json:"model,omitempty"`
	BaseURL string `json:"base_url,omitempty"`
}

// Secrets bundles credentials for every key-based provider. Add a new
// field per provider — keep field names stable so existing secrets.json
// files keep deserialising.
//
// Ollama is keyless (it runs on the user's own machine), but its model and
// base URL live here too so all provider configuration stays in one place.
type Secrets struct {
	OpenAI    ProviderCreds `json:"openai"`
	Anthropic ProviderCreds `json:"anthropic"`
	Ollama    ProviderCreds `json:"ollama"`
}

// LoadSecrets reads secrets.json, creating an empty one on first run.
// Missing keys are returned as empty strings; callers decide whether
// that's a fatal condition.
func LoadSecrets(dataDir string) (Secrets, error) {
	path := filepath.Join(dataDir, SecretsFileName)
	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		empty := Secrets{}
		return empty, SaveSecrets(dataDir, empty)
	}
	if err != nil {
		return Secrets{}, err
	}
	var s Secrets
	if err := json.Unmarshal(raw, &s); err != nil {
		return Secrets{}, err
	}
	return s, nil
}

// SaveSecrets writes secrets.json at 0600. We deliberately re-set the
// mode on every write to recover from a file that was accidentally
// created world-readable (e.g. via an editor or `cp`).
func SaveSecrets(dataDir string, s Secrets) error {
	raw, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	path := filepath.Join(dataDir, SecretsFileName)
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		return err
	}
	// WriteFile preserves the file's existing mode if it already exists,
	// so chmod afterwards to be explicit.
	return os.Chmod(path, 0o600)
}

// MaskAPIKey returns a fixed-shape redaction for display in the UI.
// We keep the last 4 chars so the user can tell at a glance whether
// the value matches what they pasted, without leaking the full key.
func MaskAPIKey(key string) string {
	if key == "" {
		return "(not set)"
	}
	if len(key) <= 4 {
		return "****"
	}
	return "****" + key[len(key)-4:]
}
