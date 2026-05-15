package helpers

import (
	"fmt"
	"strconv"
)

// AppendArgs appends each arg to destination without reflection or fmt.
// Covers the types that appear in practice. Falls back to fmt only for
// exotic types — still no alloc on the hot path.
func AppendArgs(destination []byte, args ...any) []byte {
	for ix, arg := range args {
		if ix > 0 {
			destination = append(destination, ' ')
		}

		switch value := arg.(type) {
		case string:
			destination = append(destination, value...)
		case []byte:
			destination = append(destination, value...)
		case int:
			destination = strconv.AppendInt(destination, int64(value), 10)
		case int64:
			destination = strconv.AppendInt(destination, value, 10)
		case int32:
			destination = strconv.AppendInt(destination, int64(value), 10)
		case uint:
			destination = strconv.AppendUint(destination, uint64(value), 10)
		case uint64:
			destination = strconv.AppendUint(destination, value, 10)

		case float64:
			destination = AppendFloat(destination, value, 12)
		case float32:
			destination = AppendFloat(destination, float64(value), 6)

		case bool:
			destination = strconv.AppendBool(destination, value)
		case error:
			destination = append(destination, []byte(value.Error())...)
		case nil:
			destination = append(destination, "null"...)

		default:
			// Exotic types only — keeps hot path clean.
			destination = append(destination, fmt.Sprint(value)...)
		}
	}

	return destination
}

// Appendf produces a []byte message without allocations.
func Appendf(destination []byte, format string, args []any) []byte {
	ai := 0
	flen := len(format)

	for i := 0; i < flen; i++ {
		c := format[i]

		if c != '%' {
			destination = append(destination, c)

			continue
		}

		// Handle "%%"
		if i+1 < flen && format[i+1] == '%' {
			destination = append(destination, '%')
			i++

			continue
		}

		// Move to verb
		i++
		if i >= flen {
			// malformed trailing '%'
			destination = append(destination, '%')

			break
		}

		if ai >= len(args) {
			// no argument provided
			destination = append(destination, '%', format[i])

			continue
		}

		arg := args[ai]
		ai++

		switch format[i] {
		case 's':
			switch word := arg.(type) {
			case string:
				destination = append(destination, word...)
			case []byte:
				destination = append(destination, word...)
			default:
				destination = append(destination, fmt.Sprint(word)...) // exotic fallback
			}

		case 'd':
			switch v := arg.(type) {
			case int:
				destination = strconv.AppendInt(destination, int64(v), 10)
			case int64:
				destination = strconv.AppendInt(destination, v, 10)
			case int32:
				destination = strconv.AppendInt(destination, int64(v), 10)
			case uint:
				destination = strconv.AppendUint(destination, uint64(v), 10)
			case uint64:
				destination = strconv.AppendUint(destination, v, 10)
			default:
				destination = append(destination, fmt.Sprint(v)...)
			}

		case 'v':
			switch v := arg.(type) {
			case string:
				destination = append(destination, v...)
			case []byte:
				destination = append(destination, v...)
			case int:
				destination = strconv.AppendInt(destination, int64(v), 10)
			case int64:
				destination = strconv.AppendInt(destination, v, 10)
			case int32:
				destination = strconv.AppendInt(destination, int64(v), 10)
			case uint:
				destination = strconv.AppendUint(destination, uint64(v), 10)
			case uint64:
				destination = strconv.AppendUint(destination, v, 10)
			case float64:
				destination = strconv.AppendFloat(destination, v, 'f', -1, 64)
			case float32:
				destination = strconv.AppendFloat(destination, float64(v), 'f', -1, 32)
			case bool:
				destination = strconv.AppendBool(destination, v)
			case error:
				destination = append(destination, []byte(v.Error())...)
			case nil:
				destination = append(destination, "null"...)
			default:
				destination = append(destination, fmt.Sprint(v)...)
			}

		case 't':
			switch v := arg.(type) {
			case bool:
				destination = strconv.AppendBool(destination, v)
			default:
				destination = append(destination, fmt.Sprint(v)...)
			}

		case 'f':
			switch v := arg.(type) {
			case float64:
				destination = strconv.AppendFloat(destination, v, 'f', -1, 64)
			case float32:
				destination = strconv.AppendFloat(destination, float64(v), 'f', -1, 32)
			default:
				destination = append(destination, fmt.Sprint(v)...)
			}

		default:
			// unsupported verb → literal fallback
			destination = append(destination, '%', format[i])
		}
	}

	return destination
}
