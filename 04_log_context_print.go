package arenalog

import (
	"runtime"
	"strconv"

	"github.com/tudorhulban/arenalog/helpers"
)

func (ctx *LogContext) Print(args ...any) {
	cfg := ctx.cfg.Load()

	region, errWrite := ctx.logger.ingestor.TryWrite(ctx.logger.estimatedMessageSizeOverall)
	if errWrite != nil {
		return
	}

	buffer := region.Buf()[:0]

	if ctx.logger.withJSON {
		var (
			file string
			line int
		)

		if ctx.logger.withCaller {
			_, fileCaller, lineCaller, _ := runtime.Caller(int(ctx.logger.callerLevel))
			file = fileCaller
			line = lineCaller
		}

		if cfg.root != nil {
			buffer = ctx.logger.appendJSONRoot(
				buffer,
				helpers.AppendArgs(nil, args...),
				cfg,
				file,
				line,
			)
		}

		copy(region.Buf(), buffer)
		ctx.logger.ingestor.EndWrite(region)

		return
	}

	// Non‑JSON path
	if ctx.logger.fnTimestamp != nil {
		buffer = ctx.logger.fnTimestamp(buffer)
		buffer = append(buffer, ' ')
	}

	if ctx.logger.withCaller {
		_, file, line, _ := runtime.Caller(int(ctx.logger.callerLevel))

		buffer = append(buffer, file...)
		buffer = append(buffer, ' ')
		buffer = append(buffer, 'L', 'i', 'n', 'e')
		buffer = strconv.AppendInt(buffer, int64(line), 10)
		buffer = append(buffer, ' ')
	}

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

	buffer = append(buffer, "msg="...)
	buffer = helpers.AppendArgs(buffer, args...)
	buffer = append(buffer, '\n')

	copy(region.Buf(), buffer)

	ctx.logger.ingestor.EndWrite(region)
}
