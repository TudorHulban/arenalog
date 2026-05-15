package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type LogSet []LogRecord

// TODO: review error conditions
func NewLogset(from string) (LogSet, error) {
	if len(from) == 0 {
		return LogSet{}, errors.New("empty input")
	}

	lines := strings.Split(strings.TrimSpace(from), "\n")
	entries := make(LogSet, 0, len(lines))

	for _, line := range lines {
		cleanLine := ansiRegex.ReplaceAllString(line, "")
		cleanLine = nonPrintableRegex.ReplaceAllString(cleanLine, "")

		if len(strings.TrimSpace(cleanLine)) == 0 {
			continue
		}

		entry := LogRecord{
			raw:       line, // Store the original line immediately
			keyValues: make(map[string]any),
		}

		idx := strings.IndexByte(cleanLine, '{')

		// If it is not JSON, we still keep it as a 'raw' entry
		// but it will not have keyValues or a timestamp.
		if idx != -1 {
			jsonPart := cleanLine[idx:]

			var rawMap map[string]any

			if err := json.Unmarshal([]byte(jsonPart), &rawMap); err == nil {
				for k, v := range rawMap {
					if k == "ts" {
						if tsStr, ok := v.(string); ok {
							entry.timestamp = tsStr
						}
					} else {
						entry.keyValues[k] = v
					}
				}
			}
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (e LogSet) String() string {
	entries := make([]string, len(e))

	for ix, entry := range e {
		entries[ix] = entry.String()
	}

	return strings.Join(entries, "\n")
}

func (e LogSet) WithTimestamp() LogSet {
	var filtered LogSet

	for _, item := range e {
		if item.HasTimestamp() {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

func (e LogSet) WithNoTimestamp() LogSet {
	var filtered LogSet

	for _, item := range e {
		if !item.HasTimestamp() {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

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

// HasKey checks if a key exists a specific number of times across all entries.
func (e LogSet) HasKey(name string, noTimes uint) error {
	var count uint

	for _, item := range e {
		if exists, _ := item.HasKey(name); exists {
			count++
		}
	}

	if count != noTimes {
		return fmt.Errorf(
			"expected key %q to appear %d times, but found it %d times",
			name,
			noTimes,
			count,
		)
	}

	return nil
}

// HasKeyWithValue checks if a key with a specific value exists a specific number of times.
func (e LogSet) HasKeyWithValue(name string, value any, noTimes uint) error {
	var count uint

	for _, item := range e {
		if exists, val := item.HasKey(name); exists {
			if valuesMatch(val, value) {
				count++
			}
		}
	}

	if count != noTimes {
		return fmt.Errorf(
			"expected key %q with value %v (%T) to appear %d times, but found %d",
			name,
			value,
			value,
			noTimes,
			count,
		)
	}

	return nil
}

// HasKeyWithValueLike matches also numbers of bool.
func (e LogSet) HasKeyWithValueLike(name, value string, noTimes uint) error {
	var count uint

	for _, item := range e {
		if exists, val := item.HasKey(name); exists {
			if strings.Contains(fmt.Sprint(val), value) {
				count++
			}
		}
	}

	if count != noTimes {
		return fmt.Errorf(
			"expected key %q with value like %q to appear %d times, but found %d",
			name,
			value,
			noTimes,
			count,
		)
	}

	return nil
}

func (e LogSet) HasKeysWithValues(noTimes uint, kv ...any) error {
	if len(kv)%2 != 0 {
		return errors.New(
			"hasKeysWithValues requires an even number of kv arguments",
		)
	}

	var count uint

	for _, item := range e {
		matchAll := true

		for i := 0; i < len(kv); i += 2 {
			key, ok := kv[i].(string)
			if !ok {
				return fmt.Errorf(
					"key at index %d must be a string",
					i,
				)
			}

			if exists, actual := item.HasKey(key); !exists || !valuesMatch(actual, kv[i+1]) {
				matchAll = false

				break
			}
		}

		if matchAll {
			count++
		}
	}

	if count != noTimes {
		return fmt.Errorf(
			"expected %d entries matching %v, but found %d",
			noTimes,
			kv,
			count,
		)
	}

	return nil
}

// FilterBy returns a new subset of entries where the key matches the expected value.
func (e LogSet) FilterBy(key string, value any) LogSet {
	var filtered LogSet

	for _, item := range e {
		if exists, val := item.HasKey(key); exists {
			if valuesMatch(val, value) {
				filtered = append(filtered, item)
			}
		}
	}

	return filtered
}

func (e LogSet) First() LogSet {
	if len(e) == 0 {
		return LogSet{}
	}

	return []LogRecord{
		e[0],
	}
}

func (e LogSet) Last() LogSet {
	if len(e) == 0 {
		return LogSet{}
	}

	return []LogRecord{
		e[len(e)-1],
	}
}

// At returns a LogSet containing only the record at the specific index.
// Returns an empty LogSet if the index is out of bounds.
func (e LogSet) At(index int) LogSet {
	if index < 0 || index >= len(e) {
		return LogSet{}
	}

	return []LogRecord{e[index]}
}

// Skip returns a new LogSet starting after the first n records.
// Useful for shifting the "window" before calling First() or At().
func (e LogSet) Skip(n int) LogSet {
	if n <= 0 {
		return e
	}

	if n >= len(e) {
		return LogSet{}
	}

	return e[n:]
}

// SortByTimestamp reorders the entries based on the ts field.
// If desc is true, it sorts newest to oldest.
func (e LogSet) SortByTimestamp(desc bool) LogSet {
	// Create a copy to avoid mutating the original slice during a test
	sorted := make(LogSet, len(e))
	copy(sorted, e)

	sort.SliceStable(
		sorted,
		func(i, j int) bool {
			if desc {
				return sorted[i].timestamp > sorted[j].timestamp
			}

			return sorted[i].timestamp < sorted[j].timestamp
		},
	)

	return sorted
}
