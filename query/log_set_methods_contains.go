package query

import (
	"fmt"
	"regexp"
	"strings"
)

// Contains checks if a substring exists in the raw log line a specific number of times
// across all entries.
func (e LogSet) Contains(noTimes uint, text string) error {
	var count uint

	for _, item := range e {
		// We check the 'raw' field of the LogRecord
		if strings.Contains(item.raw, text) {
			count++
		}
	}

	if count != noTimes {
		return fmt.Errorf(
			"expected text %q to be contained in %d entries, but found it in %d entries",
			text,
			noTimes,
			count,
		)
	}

	return nil
}

// ContainsEach ensures that EACH provided word exists as a distinct,
// whole word (case-insensitive) in exactly noTimes entries.
func (e LogSet) ContainsEach(noTimes uint, words ...string) error {
	if len(words) == 0 {
		return nil
	}

	// We'll map the compiled regex to the original word for error reporting
	type wordMatcher struct {
		re       *regexp.Regexp
		original string
		count    uint
	}

	matchers := make([]*wordMatcher, len(words))

	for ix, word := range words {
		// \b ensures the match starts and ends at a word boundary
		// (?i) makes the entire expression case-insensitive
		pattern := `(?i)\b` + regexp.QuoteMeta(word) + `\b`

		re, errCompile := regexp.Compile(pattern)
		if errCompile != nil {
			return fmt.Errorf("invalid word pattern for %q: %w", word, errCompile)
		}

		matchers[ix] = &wordMatcher{original: word, re: re}
	}

	// Iterate through logs
	for _, item := range e {
		for _, m := range matchers {
			if m.re.MatchString(item.raw) {
				m.count++
			}
		}
	}

	// Validation
	for _, matcher := range matchers {
		if matcher.count != noTimes {
			return fmt.Errorf(
				"expected exact word %q to be found in %d entries, but found it in %d",
				matcher.original, noTimes, matcher.count,
			)
		}
	}

	return nil
}

// ContainsLike ensures that EACH provided word is found in exactly noTimesEach entries.
// The check is case-insensitive.
func (e LogSet) ContainsLike(noTimesEach uint, words ...string) error {
	// map to track hits for each search term
	// key: lowercase search term, value: count of matching lines
	counts := make(map[string]uint)

	// mapping for error reporting (original casing)
	originalCasing := make(map[string]string)

	for _, t := range words {
		lower := strings.ToLower(t)
		counts[lower] = 0
		originalCasing[lower] = t
	}

	for _, item := range e {
		line := strings.ToLower(item.raw)

		for lowerTerm := range counts {
			if strings.Contains(line, lowerTerm) {
				counts[lowerTerm]++
			}
		}
	}

	// Check results for each term
	for lowerTerm, foundCount := range counts {
		if foundCount != noTimesEach {
			return fmt.Errorf(
				"expected text like %q to be found in %d entries, but found it in %d entries",
				originalCasing[lowerTerm],
				noTimesEach,
				foundCount,
			)
		}
	}

	return nil
}
