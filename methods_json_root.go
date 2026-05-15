package arenalog

import (
	"strconv"

	"github.com/tudorhulban/arenalog/helpers"
)

func (l *Logger) appendJSONRoot(buffer, msg []byte, cfg *formatterConfig, file string, line int) []byte {
	buffer = append(buffer, '{')

	// timestamp
	if l.fnTimestamp != nil {
		buffer = append(buffer, `"ts":"`...)
		buffer = l.fnTimestamp(buffer)
		buffer = append(buffer, `",`...)
	}

	// root
	if cfg.root != nil {
		fld := cfg.root

		buffer = append(buffer, '"')
		buffer = append(buffer, fld.key...)
		buffer = append(buffer, '"', ':')

		switch fld.kind {
		case kindString:
			buffer = helpers.AppendJSON_Quoted(
				buffer,
				[]byte(fld.valueString),
			)
		case kindInt:
			buffer = strconv.AppendInt(buffer, fld.valueInt, 10)
		case kindBool:
			buffer = strconv.AppendBool(buffer, fld.valueBool)
		case kindFloat:
			buffer = helpers.AppendFloat(buffer, fld.valueFloat, _PrecisionFloat)
		}

		buffer = append(buffer, ',')
	}

	// context fields
	for ix := range cfg.fields {
		fld := &cfg.fields[ix]

		buffer = append(buffer, '"')
		buffer = append(buffer, fld.key...)
		buffer = append(buffer, '"', ':')

		switch fld.kind {
		case kindString:
			buffer = helpers.AppendJSON_Quoted(
				buffer,
				[]byte(fld.valueString),
			)
		case kindInt:
			buffer = strconv.AppendInt(buffer, fld.valueInt, 10)
		case kindBool:
			buffer = strconv.AppendBool(buffer, fld.valueBool)
		case kindFloat:
			buffer = helpers.AppendFloat(buffer, fld.valueFloat, _PrecisionFloat)
		}

		buffer = append(buffer, ',')
	}

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
