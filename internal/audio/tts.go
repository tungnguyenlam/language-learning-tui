package audio

import (
	"context"
	"strings"
)

type Synthesizer interface {
	Synthesize(ctx context.Context, text string) (string, error)
	ProviderName() string
	VoiceName() string
}

func normalizeSpeechText(text string) string {
	// Strip HTML tags
	text = stripHTML(text)

	// Strip Cloze markers like {{c1::text::hint}}
	text = stripCloze(text)

	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	return strings.Join(strings.Fields(text), " ")
}

func stripHTML(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func stripCloze(s string) string {
	res := s
	for {
		start := strings.Index(res, "{{c")
		if start == -1 {
			break
		}
		end := strings.Index(res[start:], "}}")
		if end == -1 {
			break
		}
		end += start
		content := res[start+2 : end]
		parts := strings.Split(content, "::")
		if len(parts) >= 2 {
			res = res[:start] + parts[1] + res[end+2:]
		} else {
			res = res[:start] + res[end+2:]
		}
	}
	return res
}
