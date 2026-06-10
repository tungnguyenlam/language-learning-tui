package content

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"strings"

	"deutsch-tui/internal/core"
)

// ParseDictCCZip opens a zip file, locates the first txt file, and parses it.
func ParseDictCCZip(zipPath string) ([]core.DictionaryEntry, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("open zip file: %w", err)
	}
	defer r.Close()

	var txtFile *zip.File
	for _, f := range r.File {
		if strings.HasSuffix(strings.ToLower(f.Name), ".txt") {
			txtFile = f
			break
		}
	}
	if txtFile == nil {
		return nil, fmt.Errorf("no txt file found in zip archive")
	}

	rc, err := txtFile.Open()
	if err != nil {
		return nil, fmt.Errorf("open txt file from zip: %w", err)
	}
	defer rc.Close()

	return ParseDictCCStream(rc)
}

// ParseDictCCStream parses a dict.cc vocabulary stream (tab-separated).
// Auto-detects DE-EN vs EN-DE layouts from header comments and extracts genders.
func ParseDictCCStream(r io.Reader) ([]core.DictionaryEntry, error) {
	var entries []core.DictionaryEntry
	scanner := bufio.NewScanner(r)

	// Increase buffer size for extremely long lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	enIndex := 0
	deIndex := 1
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// Check for comment headers to detect column order
		if strings.HasPrefix(line, "#") {
			if strings.Contains(line, "DE-EN") {
				deIndex = 0
				enIndex = 1
			} else if strings.Contains(line, "EN-DE") {
				enIndex = 0
				deIndex = 1
			}
			continue
		}

		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) >= 2 {
			enVal := strings.TrimSpace(parts[enIndex])
			deVal := strings.TrimSpace(parts[deIndex])

			// Extract gender if present at the end of the German word
			gender := ""
			wordClean := deVal

			if strings.HasSuffix(wordClean, " {m}") {
				gender = "m"
				wordClean = strings.TrimSuffix(wordClean, " {m}")
			} else if strings.HasSuffix(wordClean, " {f}") {
				gender = "f"
				wordClean = strings.TrimSuffix(wordClean, " {f}")
			} else if strings.HasSuffix(wordClean, " {n}") {
				gender = "n"
				wordClean = strings.TrimSuffix(wordClean, " {n}")
			} else if strings.HasSuffix(wordClean, " {pl}") {
				gender = "pl"
				wordClean = strings.TrimSuffix(wordClean, " {pl}")
			}

			entry := core.DictionaryEntry{
				ID:          fmt.Sprintf("dict-cc-%d", lineNum),
				Word:        wordClean,
				Translation: enVal,
				Gender:      gender,
			}

			if len(parts) > 2 {
				entry.WordClass = strings.TrimSpace(parts[2])
			}
			if len(parts) > 3 {
				tagVal := strings.TrimSpace(parts[3])
				if tagVal != "" {
					entry.Tags = []string{tagVal}
				}
			}

			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read dict.cc file: %w", err)
	}

	return entries, nil
}

// ParseDictCC is kept for backward compatibility (using default EN-DE columns)
func ParseDictCC(r io.Reader) ([]core.DictionaryEntry, error) {
	return ParseDictCCStream(r)
}
