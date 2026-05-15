package helpers

import (
	"fmt"
	"strconv"
)

// Precondition: len(kv) is even.
func AppendKeyValues(dst []byte, kv ...any) []byte {
	for i := 0; i < len(kv); i = i + 2 {
		key := kv[i]
		val := kv[i+1]

		// --- KEY ---
		switch k := key.(type) {
		case string:
			dst = append(dst, k...)
		case []byte:
			dst = append(dst, k...)
		default:
			// Keys should never be exotic. Keep this branch cold.
			dst = append(dst, fmt.Sprint(k)...)
		}

		dst = append(dst, '=') // separator

		// --- VALUE ---
		switch v := val.(type) {
		case string:
			dst = append(dst, v...)
		case []byte:
			dst = append(dst, v...)
		case int:
			dst = strconv.AppendInt(dst, int64(v), 10)
		case int64:
			dst = strconv.AppendInt(dst, v, 10)
		case int32:
			dst = strconv.AppendInt(dst, int64(v), 10)
		case uint:
			dst = strconv.AppendUint(dst, uint64(v), 10)
		case uint64:
			dst = strconv.AppendUint(dst, v, 10)

		case float64:
			dst = AppendFloat(dst, v, 12)
		case float32:
			dst = AppendFloat(dst, float64(v), 6)

		case bool:
			dst = strconv.AppendBool(dst, v)
		case error:
			dst = append(dst, v.Error()...)
		case nil:
			dst = append(dst, "null"...)

		default:
			// Exotic types only — keeps hot path clean.
			// faster would be: dst = append(dst, v.String()...)
			dst = append(dst, fmt.Sprint(v)...)
		}

		if i+2 < len(kv) {
			dst = append(dst, ' ')
		}
	}

	return dst
}

// Precondition: len(kv) is even.
func AppendJSONKeyValuesIntoObject(dst []byte, kv ...any) []byte {
	for i := 0; i < len(kv); i += 2 {
		key := kv[i]
		val := kv[i+1]

		// --- KEY ---
		dst = append(dst, '"')

		switch k := key.(type) {
		case string:
			dst = append(dst, k...)
		case []byte:
			dst = append(dst, k...)
		default:
			dst = append(dst, fmt.Sprint(k)...)
		}

		dst = append(dst, '"', ':')

		// --- VALUE ---
		switch v := val.(type) {
		case string:
			dst = append(dst, '"')
			dst = append(dst, v...)
			dst = append(dst, '"')

		case []byte:
			dst = append(dst, '"')
			dst = append(dst, v...)
			dst = append(dst, '"')

		case int:
			dst = strconv.AppendInt(dst, int64(v), 10)
		case int64:
			dst = strconv.AppendInt(dst, v, 10)
		case int32:
			dst = strconv.AppendInt(dst, int64(v), 10)
		case uint:
			dst = strconv.AppendUint(dst, uint64(v), 10)
		case uint64:
			dst = strconv.AppendUint(dst, v, 10)

		case float64:
			dst = AppendFloat(dst, v, 12)
		case float32:
			dst = AppendFloat(dst, float64(v), 6)

		case bool:
			dst = strconv.AppendBool(dst, v)

		case error:
			dst = append(dst, '"')
			dst = append(dst, v.Error()...)
			dst = append(dst, '"')

		case nil:
			dst = append(dst, "null"...)

		default:
			dst = append(dst, '"')
			dst = append(dst, fmt.Sprint(v)...)
			dst = append(dst, '"')
		}

		dst = append(dst, ',')
	}

	return dst
}
