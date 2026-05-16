package audio

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lib-x/edgetts"
)

const (
	ProviderDisabled = "disabled"
	ProviderEdgeTTS  = "edge"

	DefaultEdgeVoice = "de-DE-KatjaNeural"
)

type Synthesizer interface {
	Synthesize(ctx context.Context, text string) (string, error)
	ProviderName() string
	VoiceName() string
}

type EdgeTTS struct {
	client  edgeTTSClient
	voice   string
	cache   string
	timeout time.Duration
}

type edgeTTSClient interface {
	Save(ctx context.Context, text, path string, opts ...edgetts.Option) error
}

func NewEdgeTTS(cacheDir, voice string) *EdgeTTS {
	if strings.TrimSpace(voice) == "" {
		voice = DefaultEdgeVoice
	}
	return &EdgeTTS{
		client:  edgetts.New(edgetts.WithVoice(voice)),
		voice:   voice,
		cache:   cacheDir,
		timeout: 30 * time.Second,
	}
}

func (e *EdgeTTS) ProviderName() string {
	return ProviderEdgeTTS
}

func (e *EdgeTTS) VoiceName() string {
	return e.voice
}

func (e *EdgeTTS) Synthesize(ctx context.Context, text string) (string, error) {
	text = normalizeSpeechText(text)
	if text == "" {
		return "", errors.New("no text available for TTS")
	}
	if strings.TrimSpace(e.cache) == "" {
		return "", errors.New("TTS cache directory is not configured")
	}
	if e.client == nil {
		return "", errors.New("Edge TTS client is not configured")
	}
	if err := os.MkdirAll(e.cache, 0o755); err != nil {
		return "", err
	}

	path := filepath.Join(e.cache, e.cacheFilename(text))
	if info, err := os.Stat(path); err == nil && info.Size() > 0 {
		return path, nil
	}

	runCtx := ctx
	cancel := func() {}
	if e.timeout > 0 {
		runCtx, cancel = context.WithTimeout(ctx, e.timeout)
	}
	defer cancel()

	tmp := path + ".tmp"
	_ = os.Remove(tmp)
	if err := e.client.Save(runCtx, text, tmp); err != nil {
		_ = os.Remove(tmp)
		return "", fmt.Errorf("edge-tts failed: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return "", err
	}
	return path, nil
}

func (e *EdgeTTS) cacheFilename(text string) string {
	sum := sha256.Sum256([]byte(e.voice + "\x00" + text))
	return hex.EncodeToString(sum[:]) + ".mp3"
}

func normalizeSpeechText(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	return strings.Join(strings.Fields(text), " ")
}

func singleLine(text string) string {
	return strings.Join(strings.Fields(text), " ")
}
