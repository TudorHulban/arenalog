package arenalog

import (
	"runtime"
	"strconv"

	"github.com/tudorhulban/arenalog/helpers"
)

func (l *Logger) logWithLabel(label string, estimatedMessageSize uint32, args []any) {
	var (
		callingFromFile string
		callingFromLine int
	)

	if l.withCaller {
		fileCaller, lineCaller := l.getCallerData(int(l.callerLevel))
		callingFromFile = fileCaller
		callingFromLine = lineCaller
	}

	if l.withJSON {
		region, errWrite := l.ingestor.TryWrite(
			estimatedMessageSize +
				l.estimateJSONOverhead(0, callingFromFile, callingFromLine, args),
		)
		if errWrite != nil {
			return
		}

		buf := region.Buf()[:0]

		buf = l.appendJSON(
			buf,
			label,
			callingFromFile,
			callingFromLine,
			helpers.AppendArgs(nil, args...),
		)

		copy(region.Buf(), buf)
		l.ingestor.EndWrite(region)

		return
	}

	// Non‑JSON path
	region, errWrite := l.ingestor.TryWrite(
		estimatedMessageSize +
			l.estimateTextOverhead(0, callingFromFile, callingFromLine, args),
	)
	if errWrite != nil {
		return
	}

	buf := region.Buf()[:0]

	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	if l.withCaller {
		buf = append(buf, callingFromFile...)
		buf = append(buf, ' ')
		buf = append(buf, 'L', 'i', 'n', 'e')
		buf = strconv.AppendInt(buf, int64(callingFromLine), 10)
		buf = append(buf, ' ')
	}

	buf = append(buf, label...)
	buf = append(buf, delim...)
	buf = helpers.AppendArgs(buf, args...)
	buf = append(buf, '\n')

	copy(region.Buf(), buf)

	l.ingestor.EndWrite(region)
}

func (l *Logger) logfWithLabel(label, format string, estimatedMessageSize uint32, args []any) {
	var (
		callingFromFile string
		callingFromLine int
	)

	if l.withCaller {
		fileCaller, lineCaller := l.getCallerData(int(l.callerLevel))
		callingFromFile = fileCaller
		callingFromLine = lineCaller
	}

	if l.withJSON {
		region, errWrite := l.ingestor.TryWrite(
			estimatedMessageSize +
				l.estimateJSONOverhead(0, callingFromFile, callingFromLine, args),
		)
		if errWrite != nil {
			return
		}

		buf := region.Buf()[:0]

		buf = l.appendJSON(
			buf,
			label,
			callingFromFile,
			callingFromLine,
			helpers.Appendf(nil, format, args),
		)

		copy(region.Buf(), buf)
		l.ingestor.EndWrite(region)

		return
	}

	// Non‑JSON path
	region, errWrite := l.ingestor.TryWrite(
		estimatedMessageSize +
			l.estimateTextOverhead(0, callingFromFile, callingFromLine, args),
	)
	if errWrite != nil {
		return
	}

	buffer := region.Buf()[:0]

	if l.fnTimestamp != nil {
		buffer = l.fnTimestamp(buffer)
		buffer = append(buffer, ' ')
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(int(l.callerLevel))

		buffer = append(buffer, file...)
		buffer = append(buffer, ' ')
		buffer = append(buffer, 'L', 'i', 'n', 'e')
		buffer = strconv.AppendInt(buffer, int64(line), 10)
		buffer = append(buffer, ' ')
	}

	buffer = append(buffer, label...)
	buffer = append(buffer, delim...)
	buffer = helpers.Appendf(buffer, format, args)
	buffer = append(buffer, '\n')

	copy(region.Buf(), buffer)

	l.ingestor.EndWrite(region)
}

func (l *Logger) logwWithLabel(label, msg string, estimatedMessageSize uint32, keysAndValues ...any) {
	var (
		callingFromFile string
		callingFromLine int
	)

	if l.withCaller {
		// using the cache is faster than
		// _, fileCaller, lineCaller, _ := runtime.Caller(int(l.callerLevel))
		fileCaller, lineCaller := l.getCallerData(int(l.callerLevel))

		callingFromFile = fileCaller
		callingFromLine = lineCaller
	}

	if l.withJSON {
		region, errWrite := l.ingestor.TryWrite(
			estimatedMessageSize +
				l.estimateJSONOverhead(len(msg), callingFromFile, callingFromLine, keysAndValues),
		)
		if errWrite != nil {
			return
		}

		buf := region.Buf()[:0]

		buf = l.appendJSONKV(
			buf,
			label,
			callingFromFile,
			callingFromLine,
			[]byte(msg),
			keysAndValues...,
		)

		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)

		return
	}

	// Non‑JSON path
	region, errWrite := l.ingestor.TryWrite(
		estimatedMessageSize +
			l.estimateTextOverhead(len(msg), callingFromFile, callingFromLine, keysAndValues),
	)
	if errWrite != nil {
		return
	}

	buf := region.Buf()[:0]

	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	buf = append(buf, msg...)
	buf = append(buf, ' ')

	if l.withCaller {
		buf = append(buf, callingFromFile...)
		buf = append(buf, ' ')
		buf = append(buf, 'L', 'i', 'n', 'e')
		buf = strconv.AppendInt(buf, int64(callingFromLine), 10)
		buf = append(buf, ' ')
	}

	buf = append(buf, label...)
	buf = append(buf, delim...)

	buf = helpers.AppendKeyValues(buf, keysAndValues...)
	buf = append(buf, '\n')

	copy(region.Buf(), buf)

	l.ingestor.EndWrite(region)
}

// estimateJSONOverhead includes log level label.
func (l *Logger) estimateJSONOverhead(msgLen int, file string, line int, args []any) uint32 {
	var size uint32 = 64 // base JSON overhead

	// timestamp
	if l.fnTimestamp != nil {
		size = size + 32 // worst case timestamp length
	}

	// level
	size = size + 5 + 10 // 5 - level label, 10 - for JSON

	// caller info
	if len(file) > 0 && line > 0 {
		size = size + uint32(len(file)) + 20 // "caller":"...","line":123,
	}

	// message field: "msg":"<escaped>"
	// worst case: every char becomes \u00XX (6 bytes)
	size = size + uint32(msgLen)*2
	size = size + 10 // field name + quotes + comma

	// key/value pairs
	// no assumption about even length
	for i := 0; i < len(args); i = i + 2 {
		key := args[i]

		// key
		switch k := key.(type) {
		case string:
			size = size + uint32(len(k))*2 + 4
		case []byte:
			size = size + uint32(len(k))*2 + 4
		default:
			size = size + 16
		}

		// value
		if i+1 < len(args) {
			val := args[i+1]

			switch v := val.(type) {
			case string:
				size = size + uint32(len(v))*2 + 4
			case []byte:
				size = size + uint32(len(v))*2 + 4

			case int, int32, int64, uint, uint64:
				size = size + 20

			case float32, float64:
				size = size + 32

			case bool:
				size = size + 5

			case nil:
				size = size + 4

			default:
				size = size + 32
			}
		} else {
			// dangling key without value
			// worst case: treat as string with unknown length
			size = size + 32
		}
	}

	// newline
	size = size + 1

	return size
}

// estimateTextOverhead includes log level label.
func (l *Logger) estimateTextOverhead(msgLen int, file string, line int, args []any) uint32 {
	var size uint32

	// 1. Timestamp (e.g., "2026-05-25T14:02:52.975Z ")
	if l.fnTimestamp != nil {
		size = size + 24 + 1 // "YYYY-MM-DDTHH:MM:SS.mmmZ" + space
	}

	// 2. Caller Info (e.g., "/path/to/file.go Line68 ")
	if len(file) > 0 && line > 0 {
		size = size + uint32(len(file)) + 4 + 10 + 1 // " Line" + max int digits + space
	}

	// 3. Level (e.g., "TRACE: " or "level=INFO ")
	// Looking at your logs, some use "TRACE:" and some use "level=INFO ".
	// We'll budget for the larger one ("level=ERROR " = 12 bytes).
	size = size + 12

	// 4. Key/Value Pairs (e.g., "service=auth req_id=12345 ")
	for i := 0; i < len(args); i = i + 2 {
		key := args[i]

		// Key size + '=' sign
		switch k := key.(type) {
		case string:
			size = size + uint32(len(k)) + 1
		case []byte:
			size = size + uint32(len(k)) + 1
		default:
			size = size + 16 + 1 // Fallback for weird key types
		}

		// Value size + trailing space
		if i+1 < len(args) {
			val := args[i+1]

			switch v := val.(type) {
			case string:
				// If it has spaces it gets quoted, add 2 bytes buffer for safety
				size = size + uint32(len(v)) + 2 + 1
			case []byte:
				size = size + uint32(len(v)) + 2 + 1

			case int, int32, int64, uint, uint64:
				size = size + 20 + 1 // Fast estimation without Digits logic

			case float32, float64:
				// Non-json text logs preserve floating-point ghost digits (e.g., 4.299999999999)
				size = size + 24 + 1

			case bool:
				size = size + 5 + 1 // "false" + space

			case nil:
				size = size + 4 + 1 // "nil" + space

			default:
				size = size + 32 + 1 // Fallback for structural strings
			}
		} else {
			// Dangling key without a value (e.g., " dangled= ")
			size = size + 4 + 1
		}
	}

	// 5. Message Field (e.g., "msg=hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh")
	if msgLen > 0 {
		size = size + 4 + uint32(msgLen) // "msg=" + message length
	}

	// 6. Trailing Newline
	size = size + 1

	return size
}
