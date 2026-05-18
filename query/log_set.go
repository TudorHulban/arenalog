package query

import (
	"encoding/json"
	"errors"
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
