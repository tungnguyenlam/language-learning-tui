package content

import (
	"strings"
)

func stripArticles(s string) string {
	s = strings.TrimSpace(s)
	lower := strings.ToLower(s)
	if strings.HasPrefix(lower, "der ") {
		return s[4:]
	}
	if strings.HasPrefix(lower, "die ") {
		return s[4:]
	}
	if strings.HasPrefix(lower, "das ") {
		return s[4:]
	}
	return s
}

func findWordInSentence(sentence, word string) (int, int) {
	s := strings.ToLower(sentence)
	w := strings.ToLower(word)

	idx := strings.Index(s, w)
	if idx == -1 {
		return -1, -1
	}

	return idx, idx + len(word)
}
