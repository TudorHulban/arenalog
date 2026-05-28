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

// Msgs allocates but does not estimate.
func (e *Entry) Msgs(msg string) {
	if Level(e.formatter.logger.logLevel.Load()) > e.level {
		return
	}

	cfg := e.formatter.cfg.Load()

	// -------------------------------------------------------------------------
	// JSON MODE
	// -------------------------------------------------------------------------
	if e.formatter.logger.withJSON {
		// Calculate precise capacity to avoid runtime array growth allocations
		capacityJSON := _EstimatedBaselineCapacityJSON

		if e.formatter.logger.fnTimestamp != nil {
			capacityJSON = capacityJSON + 35 // "ts":"..." baseline
		}

		capacityJSON = capacityJSON + 15 + len(e.level.String()) // "level":"..."

		if cfg.root != nil {
			capacityJSON = capacityJSON + len(cfg.root.key) + len(cfg.root.valueString) + 40
		}

		// Account for context fields
		for ix := range cfg.fields {
			capacityJSON = capacityJSON + len(cfg.fields[ix].key) + len(cfg.fields[ix].valueString) + 40
		}

		// Account for dynamic entry fields
		for ix := range e.fieldCount {
			capacityJSON = capacityJSON + len(e.fields[ix].key) + len(e.fields[ix].valueString) + 40
		}

		// Account for message (Worst case JSON escaping doubles string size roughly, plus keys)
		capacityJSON = capacityJSON + len(msg)*2 + 16

		// Allocate local buffer from stack or small heap allotment bound
		buffer := make([]byte, 0, capacityJSON)
		buffer = append(buffer, '{')

		// timestamp
		if e.formatter.logger.fnTimestamp != nil {
			buffer = append(buffer, `"ts":`...)
			buffer = helpers.AppendJSON_Quoted(
				buffer,
				e.formatter.logger.fnTimestamp(buffer),
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
				buffer = helpers.AppendJSON_Quoted(buffer, []byte(fld.valueString))
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
				buffer = helpers.AppendJSON_Quoted(buffer, []byte(fld.valueString))
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
				buffer = helpers.AppendJSON_Quoted(buffer, []byte(fld.valueString))
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

		if len(buffer) == 0 {
			entryPool.Put(e)
			return
		}

		// Safely write exact length to ring buffer
		region, errWrite := e.formatter.logger.ingestor.TryWrite(uint32(len(buffer)))
		if errWrite == nil {
			copy(region.Buf(), buffer)
			e.formatter.logger.ingestor.EndWrite(region)
		}

		entryPool.Put(e)

		return
	}

	// -------------------------------------------------------------------------
	// NON-JSON PATH
	// -------------------------------------------------------------------------
	capacityRaw := len(msg) + _EstimatedBaselineCapacityRaw // msg="string"\n baseline

	if e.formatter.logger.fnTimestamp != nil {
		capacityRaw = capacityRaw + 21
	}

	capacityRaw = capacityRaw + 16 + len(e.level.StringColored()) // level=...

	if cfg.root != nil {
		capacityRaw = capacityRaw + len(cfg.root.key) + _EstimatedBaselineRoot
	}

	for ix := range cfg.fields {
		capacityRaw = capacityRaw + len(cfg.fields[ix].key) + _EstimatedBaselineFields
	}

	for ix := range e.fieldCount {
		capacityRaw = capacityRaw + len(e.fields[ix].key) + _EstimatedBaselineFieldCount
	}

	buffer := make([]byte, 0, capacityRaw)

	// timestamp
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

	if len(buffer) == 0 {
		entryPool.Put(e)
		return
	}

	// Safely write exact length to ring buffer
	region, errWrite := e.formatter.logger.ingestor.TryWrite(uint32(len(buffer)))
	if errWrite == nil {
		copy(region.Buf(), buffer)
		e.formatter.logger.ingestor.EndWrite(region)
	}

	entryPool.Put(e)
}
