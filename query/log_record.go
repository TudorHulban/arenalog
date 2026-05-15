package query

type LogRecord struct {
	timestamp string
	keyValues map[string]any
	raw       string // The original, untouched line
}

func (e LogRecord) String() string {
	return e.raw
}

func (e LogRecord) IsRAW() bool {
	return len(e.raw) != 0 && (len(e.timestamp) == 0 || len(e.keyValues) == 0)
}

func (e LogRecord) HasTimestamp() bool {
	return len(e.timestamp) != 0
}

func (e LogRecord) HasKey(key string) (bool, any) {
	val, exists := e.keyValues[key]

	return exists, val
}
