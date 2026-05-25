package helpers

import (
	"fmt"
	"strconv"
)

func GetEstimatedArgsSize(args []any) uint32 {
	var size uint32
	numArgs := len(args)
	if numArgs == 0 {
		return 0
	}

	for i := 0; i < numArgs; i++ {
		arg := args[i]
		if arg == nil {
			size = size + 4 // "<nil>"
			continue
		}

		switch v := arg.(type) {
		case string:
			size = size + uint32(len(v))
		case []byte:
			size = size + uint32(len(v))

		case int:
			size = size + DigitsInt(int64(v))
		case int64:
			size = size + digitsInt64(v)
		case int32:
			size = size + digitsInt64(int64(v))
		case int16:
			size = size + digitsInt64(int64(v))
		case int8:
			size = size + digitsInt64(int64(v))
		case uint:
			size = size + digitsUint64(uint64(v))
		case uint64:
			size = size + digitsUint64(v)
		case uint32:
			size = size + digitsUint64(uint64(v))

		case bool:
			if v {
				size = size + 4 // "true"
			} else {
				size = size + 5 // "false"
			}

		case float64:
			size = size + uint32(len(strconv.FormatFloat(v, 'g', -1, 64)))
		case float32:
			size = size + uint32(len(strconv.FormatFloat(float64(v), 'g', -1, 32)))

		default:
			// Fallback for complex types (structs, maps, arrays)
			size = size + uint32(len(fmt.Sprint(v)))
		}
	}

	// Account for separators between arguments (e.g., spaces or commas)
	// If mimicking fmt.Println, add 1 byte per space between arguments
	size = size + uint32(numArgs-1)

	return size
}

func GetEstimatedMessageSize(format string, args []any) uint32 {
	var size uint32

	ai := 0
	flen := len(format)

	for i := 0; i < flen; i++ {
		c := format[i]

		if c != '%' {
			size++

			continue
		}

		// "%%"
		if i+1 < flen && format[i+1] == '%' {
			size++
			i++

			continue
		}

		i++
		if i >= flen {
			size++ // lone '%'
			break
		}

		if ai >= len(args) {
			size = size + 2 // "%x"

			continue
		}

		arg := args[ai]
		ai++

		switch format[i] {
		case 's':
			switch v := arg.(type) {
			case string:
				size = size + uint32(len(v))
			case []byte:
				size = size + uint32(len(v))
			default:
				size = size + uint32(len(fmt.Sprint(v)))
			}

		case 'd':
			switch v := arg.(type) {
			case int:
				size = size + DigitsInt(int64(v))
			case int64:
				size = size + digitsInt64(v)
			case int32:
				size = size + digitsInt64(int64(v))
			case uint:
				size = size + digitsUint64(uint64(v))
			case uint64:
				size = size + digitsUint64(v)
			default:
				size = size + uint32(len(fmt.Sprint(v)))
			}

		case 't':
			if b, ok := arg.(bool); ok && !b {
				size = size + 5 // "false"
			} else {
				size = size + 4 // "true" (or fallback for invalid type)
			}

		case 'f':
			switch v := arg.(type) {
			case float64:
				size = size + float64Len(v)
			case float32:
				size = size + float64Len(float64(v))
			default:
				size = size + uint32(len(fmt.Sprint(v)))
			}

		case 'v':
			size = size + uint32(len(fmt.Sprint(arg)))

		default:
			size = size + 2 // "%x"
		}
	}

	return size
}

func DigitsInt(value int64) uint32 {
	if value == 0 {
		return 1
	}

	var result uint32

	if value < 0 {
		result = 1
		value = -value
	}

	for value > 0 {
		value = value / 10
		result++
	}

	return result
}

func digitsInt64(value int64) uint32 {
	if value == 0 {
		return 1
	}

	var result uint32

	if value < 0 {
		result = 1
		value = -value
	}

	for value > 0 {
		value = value / 10
		result++
	}

	return result
}

func digitsUint64(v uint64) uint32 {
	if v == 0 {
		return 1
	}

	var result uint32

	for v > 0 {
		v /= 10
		result++
	}

	return result
}

func float64Len(v float64) uint32 {
	// conservative upper bound (cheap + safe for sizing)
	// "-12345.6789"
	var result uint32

	if v < 0 {
		result = 1
		v = -v
	}

	// integer part
	i := uint64(v)
	result = result + digitsUint64(i)

	// decimal part (assume ".+" even if not always present)
	return result + 12
}
