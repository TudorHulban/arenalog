package arenalog

import (
	"github.com/tudorhulban/arenalog/helpers"
)

// - Print / PrintMessage / PrintWithNoTimestamp / Printw / Printf
//   Build into their own buffer first, then reserve exactly what is needed.
//   Guaranteed no truncation.
//   One allocation per call.
//   Use when message size is unbounded or correctness is required.

// PrintRaw is always safe.
// The caller owns the buffer and the reservation is sized to match.

func (l *Logger) labelPrint() string {
	if l.withColor {
		return colorDebug(logLevels[l.GetLogLevel()])
	}

	return logLevels[l.GetLogLevel()]
}

func (l *Logger) Msg(msg string) {
	region, errWrite := l.ingestor.TryWrite(uint32(len(msg) + _DeltaEstimation))
	if errWrite != nil {
		return
	}

	buffer := region.Buf()[:0]

	if l.withJSON {
		buffer = l.appendJSON(
			buffer,
			l.labelPrint(),
			"",
			0,
			[]byte(msg),
		)

		copy(region.Buf(), buffer)
		l.ingestor.EndWrite(region)

		return
	}

	if l.fnTimestamp != nil {
		buffer = l.fnTimestamp(buffer)
		buffer = append(buffer, ' ')
	}

	buffer = append(buffer, msg...)
	buffer = append(buffer, '\n')

	copy(region.Buf(), buffer)
	l.ingestor.EndWrite(region)
}

func (l *Logger) Print(args ...any) {
	l.logWithLabel(
		l.labelPrint(),
		helpers.GetEstimatedMessageSize("", args),
		args,
	)
}

func (l *Logger) Printf(format string, args ...any) {
	l.logfWithLabel(
		l.labelPrint(),
		format,
		helpers.GetEstimatedMessageSize(format, args),
		args,
	)
}

func (l *Logger) Printw(msg string, keysAndValues ...any) {
	if (len(keysAndValues) & 1) != 0 {
		keysAndValues = append(keysAndValues, "(MISSING)")

		msg = "LOG_ERR(odd_args): " + msg
	}

	l.logwWithLabel(
		l.labelPrint(),
		msg,
		uint32(len(msg))+helpers.GetEstimatedMessageSize("", keysAndValues),
		keysAndValues...,
	)
}

func (l *Logger) PrintRaw(msg []byte) {
	region, errWrite := l.ingestor.TryWrite(uint32(len(msg))) //nolint:gosec
	if errWrite == nil {
		copy(region.Buf(), msg)

		l.ingestor.EndWrite(region)
	}
}
