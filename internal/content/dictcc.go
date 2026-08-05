package content

import (
	"archive/zip"
	"bufio"
	"fmt"
	"html"
	"io"
	"regexp"
	"strings"

	"deutsch-tui/internal/core"
)

var (
	genderRegex  = regexp.MustCompile(`\{([mfn]|pl|m/f|m/n|f/n)\}`)
	formsRegex   = regexp.MustCompile(`<([^>]+)>`)
	bracketRegex = regexp.MustCompile(`\[([^\]]+)\]`)
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
// Auto-detects DE-EN vs EN-DE layouts from header comments and extracts genders, forms, and tags.
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

		// Decode HTML entities (e.g. &#346; -> Ś, &amp; -> &)
		line = html.UnescapeString(line)

		parts := strings.Split(line, "\t")
		if len(parts) >= 2 {
			enVal := strings.TrimSpace(parts[enIndex])
			deVal := strings.TrimSpace(parts[deIndex])

			gender := ""
			var formsList []string
			var tagSet []string

			// Extract gender annotation {...} anywhere in German term
			if strings.IndexByte(deVal, '{') != -1 {
				if match := genderRegex.FindStringSubmatch(deVal); len(match) > 1 {
					gender = match[1]
					deVal = genderRegex.ReplaceAllString(deVal, "")
				}
			}

			// Extract forms <...> from German term or English term
			if strings.IndexByte(deVal, '<') != -1 {
				if match := formsRegex.FindAllStringSubmatch(deVal, -1); len(match) > 0 {
					for _, m := range match {
						if len(m) > 1 {
							formsList = append(formsList, strings.TrimSpace(m[1]))
						}
					}
					deVal = formsRegex.ReplaceAllString(deVal, "")
				}
			}
			if strings.IndexByte(enVal, '<') != -1 {
				if match := formsRegex.FindAllStringSubmatch(enVal, -1); len(match) > 0 {
					for _, m := range match {
						if len(m) > 1 {
							formsList = append(formsList, strings.TrimSpace(m[1]))
						}
					}
					enVal = formsRegex.ReplaceAllString(enVal, "")
				}
			}

			// Extract inline brackets [...] from German or English terms
			if strings.IndexByte(deVal, '[') != -1 {
				if match := bracketRegex.FindAllStringSubmatch(deVal, -1); len(match) > 0 {
					for _, m := range match {
						if len(m) > 1 {
							tagSet = append(tagSet, "["+strings.TrimSpace(m[1])+"]")
						}
					}
					deVal = bracketRegex.ReplaceAllString(deVal, "")
				}
			}
			if strings.IndexByte(enVal, '[') != -1 {
				if match := bracketRegex.FindAllStringSubmatch(enVal, -1); len(match) > 0 {
					for _, m := range match {
						if len(m) > 1 {
							tagSet = append(tagSet, "["+strings.TrimSpace(m[1])+"]")
						}
					}
					enVal = bracketRegex.ReplaceAllString(enVal, "")
				}
			}

			wordClean := cleanWhitespace(deVal)
			enClean := cleanWhitespace(enVal)

			if wordClean == "" || enClean == "" {
				continue
			}

			entry := core.DictionaryEntry{
				ID:          fmt.Sprintf("dict-cc-%d", lineNum),
				Word:        wordClean,
				Translation: enClean,
				Gender:      gender,
			}

			if len(formsList) > 0 {
				entry.Forms = strings.Join(formsList, "; ")
			}

			if len(parts) > 2 {
				entry.WordClass = strings.TrimSpace(parts[2])
			}

			if len(parts) > 3 {
				tagVal := strings.TrimSpace(parts[3])
				if tagVal != "" {
					tagSet = append(tagSet, tagVal)
				}
			}
			if len(tagSet) > 0 {
				entry.Tags = deduplicateStrings(tagSet)
			}

			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read dict.cc file: %w", err)
	}

	return core.ConsolidateDictionaryEntries(entries), nil
}

func cleanWhitespace(s string) string {
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

func deduplicateStrings(slice []string) []string {
	seen := make(map[string]bool)
	var res []string
	for _, item := range slice {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" && !seen[trimmed] {
			seen[trimmed] = true
			res = append(res, trimmed)
		}
	}
	return res
}

// ParseDictCC is kept for backward compatibility (using default EN-DE columns)
func ParseDictCC(r io.Reader) ([]core.DictionaryEntry, error) {
	return ParseDictCCStream(r)
}

// FormatDictionaryCardFront attaches the German definite article for nouns
// based on gender, leaving words that already include der/die/das unchanged.
func FormatDictionaryCardFront(word, gender string) string {
	w := strings.TrimSpace(word)
	lower := strings.ToLower(w)
	if strings.HasPrefix(lower, "der ") || strings.HasPrefix(lower, "die ") || strings.HasPrefix(lower, "das ") {
		return w
	}
	switch strings.ToLower(strings.TrimSpace(gender)) {
	case "m", "masc", "der":
		return "der " + w
	case "f", "fem", "die":
		return "die " + w
	case "n", "neut", "das":
		return "das " + w
	case "pl", "plural":
		return "die " + w + " (pl.)"
	}
	return w
}

// DictionaryEntryExtra builds the Extra metadata block used on dictionary-sourced flashcards.
func DictionaryEntryExtra(entry core.DictionaryEntry) string {
	var extraParts []string
	if entry.Forms != "" {
		extraParts = append(extraParts, "Forms: "+entry.Forms)
	}
	if entry.WordClass != "" {
		extraParts = append(extraParts, "Class: ["+strings.ToUpper(entry.WordClass)+"]")
	}
	if entry.Gender != "" {
		extraParts = append(extraParts, "Gender: {"+entry.Gender+"}")
	}
	if len(entry.Examples) > 0 {
		extraParts = append(extraParts, "Examples:\n• "+strings.Join(entry.Examples, "\n• "))
	}
	return strings.Join(extraParts, "\n")
}

// DictionaryEntriesToNotes converts dictionary entries into Anki-friendly notes for TSV export.
func DictionaryEntriesToNotes(entries []core.DictionaryEntry, deckID string) []core.Note {
	notes := make([]core.Note, 0, len(entries))
	for i, entry := range entries {
		noteID := "dict-export-" + entry.ID
		if entry.ID == "" {
			noteID = fmt.Sprintf("dict-export-%d", i+1)
		}
		tags := append([]string(nil), entry.Tags...)
		tags = append(tags, "dictionary")
		notes = append(notes, core.Note{
			ID:     noteID,
			DeckID: deckID,
			Type:   "flashcard",
			Front:  FormatDictionaryCardFront(entry.Word, entry.Gender),
			Back:   entry.Translation,
			Extra:  DictionaryEntryExtra(entry),
			Tags:   tags,
		})
	}
	return notes
}
