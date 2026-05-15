package arenalog

import (
	"strconv"

	"github.com/tudorhulban/arenalog/helpers"
)

func (e *Entry) Msg(msg string) {
	if Level(e.formatter.logger.logLevel.Load()) > e.level {
		return
	}

	cfg := e.formatter.cfg.Load()

	region, errWrite := e.formatter.logger.ingestor.TryWrite(
		uint32(len(msg)) + e.estimateFieldsSize() + _DeltaEstimation,
	)
	if errWrite != nil {
		entryPool.Put(e)

		return
	}

	buffer := region.Buf()[:0]

	// JSON MODE
	if e.formatter.logger.withJSON {
		buffer = append(buffer, '{')

		// timestamp
		if e.formatter.logger.fnTimestamp != nil {
			buffer = append(buffer, `"ts":`...)
			buffer = helpers.AppendJSON_Quoted(
				buffer,
				e.formatter.logger.fnTimestamp(nil),
			)
			buffer = append(buffer, ',')
		}

		// level
		buffer = append(buffer, `"level":`...)
		buffer = helpers.AppendJSON_Quoted(
			buffer,
			[]byte(e.level.String()),
		)
		buffer = append(buffer, ',')

		// root field
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

		// entry fields
		for ix := range e.fieldCount {
			fld := &e.fields[ix]

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

		// message
		buffer = append(buffer, `"msg":"`...)
		buffer = helpers.AppendJSON(buffer, []byte(msg))
		buffer = append(buffer, '"', '}', '\n')

		copy(region.Buf(), buffer)
		e.formatter.logger.ingestor.EndWrite(region)
		entryPool.Put(e)

		return
	}

	// Non‑JSON path
	if e.formatter.logger.fnTimestamp != nil {
		buffer = e.formatter.logger.fnTimestamp(buffer)
		buffer = append(buffer, ' ')
	}

	// level
	buffer = append(buffer, `level=`...)
	if e.formatter.logger.withColor {
		buffer = append(buffer, []byte(e.level.StringColored())...)
	} else {
		buffer = append(buffer, []byte(e.level.String())...)
	}

	buffer = append(buffer, ' ')

	// root
	if cfg.root != nil {
		fld := cfg.root

		buffer = append(buffer, fld.key...)
		buffer = append(buffer, '=')

		switch fld.kind {
		case kindString:
			buffer = append(buffer, fld.valueString...)
		case kindInt:
			buffer = strconv.AppendInt(buffer, fld.valueInt, 10)
		case kindBool:
			buffer = strconv.AppendBool(buffer, fld.valueBool)
		case kindFloat:
			buffer = helpers.AppendFloat(buffer, fld.valueFloat, _PrecisionFloat)
		}

		buffer = append(buffer, ' ')
	}

	// context fields
	for ix := range cfg.fields {
		fld := &cfg.fields[ix]

		buffer = append(buffer, fld.key...)
		buffer = append(buffer, '=')

		switch fld.kind {
		case kindString:
			buffer = append(buffer, fld.valueString...)

		case kindInt:
			buffer = strconv.AppendInt(buffer, fld.valueInt, 10)

		case kindBool:
			buffer = strconv.AppendBool(buffer, fld.valueBool)
		}

		buffer = append(buffer, ' ')
	}

	// entry fields
	for ix := range e.fieldCount {
		fld := &e.fields[ix]

		buffer = append(buffer, fld.key...)
		buffer = append(buffer, '=')

		switch fld.kind {
		case kindString:
			buffer = append(buffer, fld.valueString...)

		case kindInt:
			buffer = strconv.AppendInt(buffer, fld.valueInt, 10)

		case kindBool:
			buffer = strconv.AppendBool(buffer, fld.valueBool)

		case kindFloat:
			buffer = helpers.AppendFloat(buffer, fld.valueFloat, _PrecisionFloat)
		}

		buffer = append(buffer, ' ')
	}

	buffer = append(buffer, "msg="...)
	buffer = helpers.AppendArgs(buffer, msg)
	buffer = append(buffer, '\n')

	copy(region.Buf(), buffer)
	e.formatter.logger.ingestor.EndWrite(region)

	entryPool.Put(e)
}
