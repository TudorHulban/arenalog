package query

import (
	"regexp"
	"strconv"
	"strings"
)

func NewRawset(from string) (LogSet, error) {
	lines := strings.Split(strings.TrimSpace(from), "\n")
	entries := make(LogSet, 0, len(lines))

	// Regex for standard key=value pairs
	kvRegex := regexp.MustCompile(`([a-zA-Z0-9_-]+)=("[^"]*"|[^\s]+)`)

	// Regex for the "Linexx" pattern
	lineRegex := regexp.MustCompile(`Line(\d+)`)

	for _, line := range lines {
		cleanLine := ansiRegex.ReplaceAllString(line, "")
		cleanLine = nonPrintableRegex.ReplaceAllString(cleanLine, "")

		if len(strings.TrimSpace(cleanLine)) == 0 {
			continue
		}

		entry := LogRecord{
			raw:       line,
			keyValues: make(map[string]any),
		}

		words := strings.Fields(cleanLine)
		if len(words) > 0 {
			// 1. Extract Timestamp (First word)
			entry.timestamp = words[0]

			// 2. Identify Caller and Line Number from remaining words
			for _, word := range words[1:] {
				// a. If it starts with / it's the caller (Linux path)
				if strings.HasPrefix(word, "/") {
					entry.keyValues["caller"] = word

					continue
				}

				// b. If it matches LineXX, extract the number
				if match := lineRegex.FindStringSubmatch(word); len(match) > 1 {
					if val, err := strconv.ParseFloat(match[1], 64); err == nil {
						entry.keyValues["line"] = val
					}

					continue
				}
			}
		}

		// 3. Extract all other standard key=value pairs
		matches := kvRegex.FindAllStringSubmatch(cleanLine, -1)
		for _, match := range matches {
			key := match[1]
			val := match[2]

			if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
				val = strings.Trim(val, "\"")
			}

			if f, err := strconv.ParseFloat(val, 64); err == nil {
				entry.keyValues[key] = f
			} else if b, err := strconv.ParseBool(val); err == nil {
				entry.keyValues[key] = b
			} else {
				entry.keyValues[key] = val
			}
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
