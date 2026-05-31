package query

import (
	"errors"
	"fmt"
	"strings"
)

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

		for i := 0; i < len(kv); i = i + 2 {
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

func (e LogSet) HasKeyWithValueInDelta(name string, value any, delta float64, noTimes uint) error {
	var count uint

	for _, item := range e {
		if exists, val := item.HasKey(name); exists {
			if valuesMatchDelta(val, value, delta) {
				count++
			}
		}
	}

	if count != noTimes {
		return fmt.Errorf(
			"expected key %q with value %v (%T) within delta %f to appear %d times, but found %d",
			name,
			value,
			value,
			delta,
			noTimes,
			count,
		)
	}

	return nil
}

// Helper to safely convert interface{} numbers and check delta
func valuesMatchDelta(actual, expected any, delta float64) bool {
	actFloat, okAct := convertToFloat64(actual)
	expFloat, okExp := convertToFloat64(expected)

	if !okAct || !okExp {
		return false
	}

	diff := actFloat - expFloat
	if diff < 0 {
		diff = -diff
	}

	return diff <= delta
}

// Helper to handle any incoming numeric type from the test or the JSON
func convertToFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case int32:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true

	default:
		return 0, false
	}
}
