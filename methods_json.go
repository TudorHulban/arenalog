package arenalog

import (
	"strconv"

	"github.com/tudorhulban/arenalog/helpers"
)

// appendJSON builds a JSON log entry.
// `msg` must already be a fully formatted string.
// Caller info is included only when file != "" and line > 0.
func (l *Logger) appendJSON(buffer []byte, level, file string, line int, msg []byte) []byte {
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
	buffer = append(buffer, `"}`...)
	buffer = append(buffer, '\n')

	return buffer
}

// TODO: Formats to test.
// {"ts":"2026-03-18T14:27:09.123Z","level":"info","msg":"user login"}
// {"ts":"2026-03-18T14:27:09.123Z","level":"info","msg":"user login","user_id":3847291,"ip":"10.44.12.189","device":"mobile","session_id":"sess_abc123xyz"}
// {"ts":"…","level":"warn","msg":"slow database query","duration_ms":342,"query":"SELECT …","rows":124}
// {"ts":"…","level":"error","msg":"payment failed","error":"card_declined","code":"AUTH_402","attempt":3}
