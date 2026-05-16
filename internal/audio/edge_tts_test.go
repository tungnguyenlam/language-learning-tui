package audio

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/lib-x/edgetts"
)

type fakeEdgeClient struct {
	calls int
	err   error
}

func (f *fakeEdgeClient) Save(ctx context.Context, text, path string, opts ...edgetts.Option) error {
	f.calls++
	if f.err != nil {
		return f.err
	}
	return os.WriteFile(path, []byte("audio:"+text), 0o644)
}

func TestEdgeTTSSynthesizesAndCachesAudio(t *testing.T) {
	dir := t.TempDir()
	client := &fakeEdgeClient{}

	synth := NewEdgeTTS(filepath.Join(dir, "cache"), "de-DE-KatjaNeural")
	synth.client = client
	first, err := synth.Synthesize(context.Background(), "Guten Tag")
	if err != nil {
		t.Fatalf("synthesize first: %v", err)
	}
	second, err := synth.Synthesize(context.Background(), "Guten   Tag")
	if err != nil {
		t.Fatalf("synthesize second: %v", err)
	}
	if first != second {
		t.Fatalf("cache path changed: %q != %q", first, second)
	}
	if _, err := os.Stat(first); err != nil {
		t.Fatalf("cached audio missing: %v", err)
	}
	if client.calls != 1 {
		t.Fatalf("Save calls = %d, want 1 cached synthesis", client.calls)
	}
}

func TestEdgeTTSPropagatesClientError(t *testing.T) {
	want := errors.New("network unavailable")
	synth := NewEdgeTTS(t.TempDir(), DefaultEdgeVoice)
	synth.client = &fakeEdgeClient{err: want}
	if _, err := synth.Synthesize(context.Background(), "Hallo"); !errors.Is(err, want) {
		t.Fatalf("Synthesize error = %v, want %v", err, want)
	}
}
