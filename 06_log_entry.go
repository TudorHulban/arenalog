package arenalog

import (
	"github.com/tudorhulban/arenalog/helpers"
)

// Entry is not safe for concurrent use.
// Each goroutine should obtain its own Entry via Formatter.With.
type Entry struct {
	formatter  *LogContext
	fields     [8]field // per-request, owned by this Entry
	level      Level
	fieldCount int
}

// With allocates.
//
// Additional fields are silently
// ignored if the maximum field capacity is reached.
//
// Passing as any, the Go compiler must box the concrete value into an interface.
// The boxing causes a heap allocation when the value does not fit in a pointer word
// or the compiler cannot prove it escapes to the stack only.
func (e *Entry) With(key string, value any) *Entry {
	if e.fieldCount < len(e.fields) {
		e.fields[e.fieldCount] = makeField(key, value)

		e.fieldCount++
	}

	return e
}

// WithString appends a string field without boxing value into any.
//
// Additional fields are silently
// ignored if the maximum field capacity is reached.
func (e *Entry) WithString(key, value string) *Entry {
	if e.fieldCount < len(e.fields) {
		e.fields[e.fieldCount] = field{
			key:         key,
			kind:        kindString,
			valueString: value,
		}

		e.fieldCount++
	}

	return e
}

// Additional fields are silently
// ignored if the maximum field capacity is reached.
func (e *Entry) WithFloat(key string, value float64) *Entry {
	if e.fieldCount < len(e.fields) {
		e.fields[e.fieldCount] = field{
			key:        key,
			kind:       kindFloat,
			valueFloat: value,
		}

		e.fieldCount++
	}

	return e
}

// WithInt appends an int field without boxing value into any.
//
// Additional fields are silently
// ignored if the maximum field capacity is reached.
func (e *Entry) WithInt(key string, value int64) *Entry {
	if e.fieldCount < len(e.fields) {
		e.fields[e.fieldCount] = field{
			key:      key,
			kind:     kindInt,
			valueInt: value,
		}

		e.fieldCount++
	}

	return e
}

// WithBool appends a bool field without boxing value into any.
//
// Additional fields are silently
// ignored if the maximum field capacity is reached.
func (e *Entry) WithBool(key string, value bool) *Entry {
	if e.fieldCount < len(e.fields) {
		e.fields[e.fieldCount] = field{
			key:       key,
			kind:      kindBool,
			valueBool: value,
		}

		e.fieldCount++
	}

	return e
}

// Additional fields are silently
// ignored if the maximum field capacity is reached.
func (e *Entry) WithGoroutineID() *Entry {
	if e.fieldCount < len(e.fields) {
		e.fields[e.fieldCount] = field{
			key:      "g",
			kind:     kindInt,
			valueInt: helpers.GoroutineID(),
		}

		e.fieldCount++
	}

	return e
}

func (e *Entry) estimateFieldsSize() uint32 {
	var result uint32

	for _, f := range e.fields {
		// key
		result = result + uint32(len(`{"key":"`)+len(f.key)+len(`","value":""}`))

		switch f.kind {
		case kindString:
			result = result + uint32(len(f.valueString))

		case kindBool:
			if f.valueBool {
				result = result + 4 // true
			} else {
				result = result + 5 // false
			}

		case kindInt:
			result = result + helpers.DigitsInt(f.valueInt)
		}
	}

	return result
}
