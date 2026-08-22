package audio

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"testing"
)

func TestNativeTTS(t *testing.T) {
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		t.Skip("skipping native tts test on unsupported os")
	}

	cacheDir := t.TempDir()
	synth := NewNativeTTS(cacheDir)

	ctx := context.Background()
	path, err := synth.Synthesize(ctx, "Hallo Welt")
	if err != nil {
		t.Logf("native tts might not be available in CI environment: %v", err)
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected audio file %s to exist", path)
	}
}

func TestMultiTTS(t *testing.T) {
	cacheDir := t.TempDir()

	// Mock primary that always fails
	primary := &mockSynth{err: fmt.Errorf("primary failed")}
	secondary := NewNativeTTS(cacheDir)

	multi := NewMultiTTS(primary, secondary)

	ctx := context.Background()
	path, err := multi.Synthesize(ctx, "Test fallback")

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		if err != nil {
			t.Logf("fallback might fail if native tools missing: %v", err)
		} else if path == "" {
			t.Error("expected path from secondary")
		}
	}
}

type mockSynth struct {
	err error
}

func (m *mockSynth) Synthesize(ctx context.Context, text string) (string, error) {
	return "", m.err
}
func (m *mockSynth) VoiceName() string { return "mock" }
