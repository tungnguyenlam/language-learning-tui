package audio

import (
	"context"
	"fmt"
)

type MultiTTS struct {
	primary   Synthesizer
	secondary Synthesizer
}

func NewMultiTTS(primary, secondary Synthesizer) *MultiTTS {
	return &MultiTTS{
		primary:   primary,
		secondary: secondary,
	}
}

func (m *MultiTTS) ProviderName() string {
	return "multi"
}

func (m *MultiTTS) VoiceName() string {
	if m.primary != nil {
		return m.primary.VoiceName()
	}
	if m.secondary != nil {
		return m.secondary.VoiceName()
	}
	return "none"
}

func (m *MultiTTS) Synthesize(ctx context.Context, text string) (string, error) {
	if m.primary != nil {
		path, err := m.primary.Synthesize(ctx, text)
		if err == nil {
			return path, nil
		}
		// Fallback to secondary if primary fails
		if m.secondary != nil {
			return m.secondary.Synthesize(ctx, text)
		}
		return "", err
	}
	if m.secondary != nil {
		return m.secondary.Synthesize(ctx, text)
	}
	return "", fmt.Errorf("no synthesizer available")
}
