package content

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// CompoundPart represents a single component of a German compound word.
type CompoundPart struct {
	Word        string // The component word
	Translation string // English translation (filled in by caller via dictionary lookup)
}

// WordValidator checks if a word exists in the dictionary (case-insensitive).
type WordValidator func(word string) bool

// DecomposeCompound attempts to split a German compound noun into its components.
// German compounds are formed by concatenating nouns (right-to-left head),
// sometimes with linking elements (Fugenelemente) like -s-, -n-, -en-, -er-, -e-, -es-.
//
// If validate is non-nil, only splits where both parts are valid dictionary words are returned.
// If validate is nil, returns the best heuristic split based on component length balance.
//
// Returns nil if the word is too short or no valid split is found.
func DecomposeCompound(word string, validate WordValidator) []CompoundPart {
	// Strip common articles
	w := strings.TrimSpace(word)
	lower := strings.ToLower(w)
	for _, art := range []string{"der ", "die ", "das ", "den ", "dem ", "des "} {
		if strings.HasPrefix(lower, art) {
			w = w[len(art):]
			break
		}
	}

	w = strings.TrimSpace(w)
	if utf8.RuneCountInString(w) < 6 {
		return nil
	}

	runes := []rune(w)
	n := len(runes)

	type candidate struct {
		left  string
		right string
		score int
	}

	// Common linking elements (Fugenelemente) in German compounds.
	// Ordered by frequency: direct join is most common.
	fugen := []string{"", "s", "n", "en", "er", "e", "es", "ns"}

	var candidates []candidate
	for _, fuge := range fugen {
		fugeRunes := []rune(fuge)
		fugeLen := len(fugeRunes)
		// Left part min 3 chars, right part min 3 chars
		for splitPos := 3; splitPos <= n-3-fugeLen; splitPos++ {
			// Check if the linking element matches at the split position
			if fugeLen > 0 {
				segment := strings.ToLower(string(runes[splitPos : splitPos+fugeLen]))
				if segment != fuge {
					continue
				}
			}

			left := string(runes[:splitPos])
			right := string(runes[splitPos+fugeLen:])

			// Both parts must start with a letter
			leftFirst, _ := utf8.DecodeRuneInString(left)
			rightFirst, _ := utf8.DecodeRuneInString(right)
			if !unicode.IsLetter(leftFirst) || !unicode.IsLetter(rightFirst) {
				continue
			}

			rightCap := capitalizeFirst(strings.ToLower(right))
			leftCap := capitalizeFirst(strings.ToLower(left))

			// If validator provided, both parts must be real words
			if validate != nil {
				if !validate(leftCap) || !validate(rightCap) {
					continue
				}
			}

			// Score: prefer balanced splits (both parts are reasonably long)
			// and shorter/empty Fugen
			leftLen := utf8.RuneCountInString(left)
			rightLen := utf8.RuneCountInString(right)
			score := leftLen + rightLen // Total coverage
			// Penalize very unbalanced splits
			if leftLen < rightLen {
				score -= (rightLen - leftLen)
			} else {
				score -= (leftLen - rightLen)
			}
			// Bonus for empty Fuge (most natural)
			if fugeLen == 0 {
				score += 3
			}
			// Bonus if right part starts uppercase in original (real noun)
			if unicode.IsUpper(rightFirst) {
				score += 2
			}

			candidates = append(candidates, candidate{
				left:  leftCap,
				right: rightCap,
				score: score,
			})
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	// Return the best candidate
	best := candidates[0]
	for _, c := range candidates[1:] {
		if c.score > best.score {
			best = c
		}
	}

	return []CompoundPart{
		{Word: best.left},
		{Word: best.right},
	}
}

// capitalizeFirst uppercases the first rune of a string (German nouns are capitalized).
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}
