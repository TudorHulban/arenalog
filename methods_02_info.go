package arenalog

import "github.com/tudorhulban/arenalog/helpers"

func (l *Logger) labelInfo() string {
	if l.withColor {
		return colorInfo(logLevels[LevelInfo])
	}

	return logLevels[LevelInfo]
}

func (l *Logger) Info(args ...any) {
	if Level(l.logLevel.Load()) > LevelInfo {
		return
	}

	l.logWithLabel(
		l.labelInfo(),
		helpers.GetEstimatedMessageSize("", args),
		args,
	)
}

func (l *Logger) Infof(format string, args ...any) {
	if Level(l.logLevel.Load()) > LevelInfo {
		return
	}

	l.logfWithLabel(
		l.labelInfo(),
		format,
		helpers.GetEstimatedMessageSize(format, args),
		args,
	)
}

func (l *Logger) Infow(msg string, keysAndValues ...any) {
	if Level(l.logLevel.Load()) > LevelInfo {
		return
	}

	if (len(keysAndValues) & 1) != 0 {
		keysAndValues = append(keysAndValues, "(MISSING)")

		msg = "LOG_ERR(odd_args): " + msg
	}

	l.logwWithLabel(
		l.labelInfo(),
		msg,
		uint32(len(msg))+helpers.GetEstimatedMessageSize("", keysAndValues),
		keysAndValues...,
	)
}
