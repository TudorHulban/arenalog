package arenalog

import (
	"fmt"
	"strconv"
	"strings"
)

func (l *Logger) format(level Level, args ...any) string {
	var buffer strings.Builder

	if l.withJSON {
		buffer.WriteByte('{')

		if l.fnTimestamp != nil {
			buffer.WriteString(`"ts":"`)
			buffer.Write(l.fnTimestamp(nil))
			buffer.WriteString(`",`)
		}

		buffer.WriteString(`"level":`)
		buffer.WriteString(strconv.Quote(logLevels[level]))
		buffer.WriteString(`,"msg":`)
		buffer.WriteString(strconv.Quote(fmt.Sprint(args...)))

		buffer.WriteByte('}')
	} else {
		if l.fnTimestamp != nil {
			buffer.Write(l.fnTimestamp(nil))
			buffer.WriteByte(' ')
		}

		buffer.WriteByte('[')
		buffer.WriteString(logLevels[level])
		buffer.WriteString("] ")
		buffer.WriteString(fmt.Sprint(args...))
	}

	return buffer.String()
}

func (l *Logger) formatf(level Level, format string, args ...any) string {
	var buffer strings.Builder

	if l.withJSON {
		buffer.WriteByte('{')

		if l.fnTimestamp != nil {
			buffer.WriteString(`"ts":"`)
			buffer.Write(l.fnTimestamp(nil))
			buffer.WriteString(`",`)
		}

		buffer.WriteString(`"level":`)
		buffer.WriteString(strconv.Quote(logLevels[level]))
		buffer.WriteString(`,"msg":`)
		buffer.WriteString(strconv.Quote(fmt.Sprintf(format, args...)))

		buffer.WriteByte('}')
	} else {
		// Non-JSON / Raw section
		if l.fnTimestamp != nil {
			buffer.Write(l.fnTimestamp(nil))
			buffer.WriteByte(' ')
		}

		buffer.WriteByte('[')
		buffer.WriteString(logLevels[level])
		buffer.WriteString("] ")

		buffer.WriteString(fmt.Sprintf(format, args...))
	}

	return buffer.String()
}

func (l *Logger) formatw(level Level, msg string, keysAndValues ...any) string {
	if (len(keysAndValues) & 1) != 0 {
		keysAndValues = append(keysAndValues, "(MISSING)")

		msg = "LOG_ERR(odd_args): " + msg
	}

	var buffer strings.Builder

	if l.withJSON {
		// 1. Start the JSON object
		buffer.WriteByte('{')

		// 2. Handle timestamp inside the object
		if l.fnTimestamp != nil {
			buffer.WriteString(`"ts":"`)
			buffer.Write(l.fnTimestamp(nil))
			buffer.WriteString(`",`)
		}

		// 3. Add level and msg
		buffer.WriteString(`"level":`)
		buffer.WriteString(strconv.Quote(logLevels[level]))
		buffer.WriteString(`,"msg":`)
		buffer.WriteString(strconv.Quote(msg))

		for i := 0; i < len(keysAndValues); i = i + 2 {
			key := keysAndValues[i]
			val := keysAndValues[i+1]

			ks, couldCast := key.(string)
			if !couldCast {
				return `{"level":"fatal","error":"fatalw: key must be string"}`
			}

			buffer.WriteByte(',')
			buffer.WriteString(strconv.Quote(ks))
			buffer.WriteByte(':')

			switch value := val.(type) {
			case string:
				buffer.WriteString(strconv.Quote(value))
			case int:
				buffer.WriteString(strconv.Itoa(value))
			case bool:
				if value {
					buffer.WriteString("true")
				} else {
					buffer.WriteString("false")
				}
			default:
				buffer.WriteString(strconv.Quote(fmt.Sprint(value)))
			}
		}

		// 4. Close the object
		buffer.WriteByte('}')
	} else {
		// Non-JSON path
		if l.fnTimestamp != nil {
			tsBuf := l.fnTimestamp(nil)
			buffer.Write(tsBuf)
			buffer.WriteByte(' ')
		}

		// log level
		buffer.WriteByte('[')
		buffer.WriteString(logLevels[level])
		buffer.WriteString("] ")

		buffer.Grow(len(msg) + len(keysAndValues)*8)
		buffer.WriteString(msg)

		for i := 0; i < len(keysAndValues); i = i + 2 {
			buffer.WriteByte(' ')
			fmt.Fprint(&buffer, keysAndValues[i])
			buffer.WriteByte('=')
			fmt.Fprint(&buffer, keysAndValues[i+1])
		}
	}

	return buffer.String()
}
