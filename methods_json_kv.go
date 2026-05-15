package arenalog

import (
	"strconv"

	"github.com/tudorhulban/arenalog/helpers"
)

func (l *Logger) appendJSONKV(buffer []byte, level, file string, line int, msg []byte, kv ...any) []byte {
	buffer = append(buffer, '{')

	// timestamp
	if l.fnTimestamp != nil {
		buffer = append(buffer, `"ts":"`...)
		buffer = l.fnTimestamp(buffer)
		buffer = append(buffer, `",`...)
	}

	// level
	buffer = append(buffer, `"level":"`...)
	buffer = append(buffer, level...)
	buffer = append(buffer, `",`...)

	// caller info
	if len(file) > 0 && line > 0 {
		buffer = append(buffer, `"caller":"`...)
		buffer = append(buffer, file...)
		buffer = append(buffer, `","line":`...)
		buffer = strconv.AppendInt(buffer, int64(line), 10)
		buffer = append(buffer, ',')
	}

	// message
	buffer = append(buffer, `"msg":"`...)
	buffer = helpers.AppendJSON(buffer, msg)
	buffer = append(buffer, `",`...)

	// key/value pairs
	buffer = helpers.AppendJSONKeyValuesIntoObject(buffer, kv...)

	// remove trailing comma if present
	if buffer[len(buffer)-1] == ',' {
		buffer = buffer[:len(buffer)-1]
	}

	buffer = append(buffer, '}')
	buffer = append(buffer, '\n')

	return buffer
}
