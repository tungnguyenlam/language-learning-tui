package content

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"deutsch-tui/internal/core"

	"github.com/google/uuid"
)

// ParseDictCC parses a dict.cc vocabulary file (tab-separated).
// A typical dict.cc export is: English \t German \t [Word Class] \t ...
func ParseDictCC(r io.Reader) ([]core.DictionaryEntry, error) {
	var entries []core.DictionaryEntry
	scanner := bufio.NewScanner(r)

	// Increase buffer size in case of very long lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) >= 2 {
			en := strings.TrimSpace(parts[0])
			de := strings.TrimSpace(parts[1])

			// Sometimes words contain metadata like {m} or [noun], we just leave it for now
			// or could strip it, but it's useful context.

			entry := core.DictionaryEntry{
				ID:          uuid.New().String(),
				Word:        de,
				Translation: en,
			}

			// If there are extra tab columns, collect them
			if len(parts) > 2 {
				entry.WordClass = strings.TrimSpace(parts[2])
			}
			if len(parts) > 3 {
				entry.Tags = append(entry.Tags, strings.TrimSpace(parts[3]))
			}

			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read dict.cc file: %w", err)
	}

	return entries, nil
}
