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
			if match := genderRegex.FindStringSubmatch(deVal); len(match) > 1 {
				gender = match[1]
				deVal = genderRegex.ReplaceAllString(deVal, "")
			}

			// Extract forms <...> from German term or English term
			if match := formsRegex.FindAllStringSubmatch(deVal, -1); len(match) > 0 {
				for _, m := range match {
					if len(m) > 1 {
						formsList = append(formsList, strings.TrimSpace(m[1]))
					}
				}
				deVal = formsRegex.ReplaceAllString(deVal, "")
			}
			if match := formsRegex.FindAllStringSubmatch(enVal, -1); len(match) > 0 {
				for _, m := range match {
					if len(m) > 1 {
						formsList = append(formsList, strings.TrimSpace(m[1]))
					}
				}
				enVal = formsRegex.ReplaceAllString(enVal, "")
			}

			// Extract inline brackets [...] from German or English terms
			if match := bracketRegex.FindAllStringSubmatch(deVal, -1); len(match) > 0 {
				for _, m := range match {
					if len(m) > 1 {
						tagSet = append(tagSet, "["+strings.TrimSpace(m[1])+"]")
					}
				}
				deVal = bracketRegex.ReplaceAllString(deVal, "")
			}
			if match := bracketRegex.FindAllStringSubmatch(enVal, -1); len(match) > 0 {
				for _, m := range match {
					if len(m) > 1 {
						tagSet = append(tagSet, "["+strings.TrimSpace(m[1])+"]")
					}
				}
				enVal = bracketRegex.ReplaceAllString(enVal, "")
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

	return entries, nil
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
