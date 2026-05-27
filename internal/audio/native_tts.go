package audio

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type NativeTTS struct {
	cache string
	goos  string
}

func NewNativeTTS(cacheDir string) *NativeTTS {
	return &NativeTTS{
		cache: cacheDir,
		goos:  runtime.GOOS,
	}
}

func (n *NativeTTS) ProviderName() string {
	return "native"
}

func (n *NativeTTS) VoiceName() string {
	return "default"
}

func (n *NativeTTS) Synthesize(ctx context.Context, text string) (string, error) {
	text = normalizeSpeechText(text)
	if text == "" {
		return "", fmt.Errorf("no text available for TTS")
	}

	if err := os.MkdirAll(n.cache, 0o755); err != nil {
		return "", err
	}

	path := filepath.Join(n.cache, n.cacheFilename(text))
	if info, err := os.Stat(path); err == nil && info.Size() > 0 {
		return path, nil
	}

	switch n.goos {
	case "darwin":
		// say -o output.aiff "text"
		cmd := exec.CommandContext(ctx, "say", "-o", path, text)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("native tts (say) failed: %w", err)
		}
		return path, nil
	case "linux":
		// espeak -w output.wav "text"
		cmd := exec.CommandContext(ctx, "espeak", "-v", "de", "-w", path, text)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("native tts (espeak) failed: %w", err)
		}
		return path, nil
	default:
		return "", fmt.Errorf("native tts not supported on %s", n.goos)
	}
}

func (n *NativeTTS) cacheFilename(text string) string {
	sum := sha256.Sum256([]byte("native\x00" + text))
	ext := ".wav"
	if n.goos == "darwin" {
		ext = ".aiff"
	}
	return hex.EncodeToString(sum[:]) + ext
}
