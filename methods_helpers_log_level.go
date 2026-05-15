package arenalog

func (l *Logger) GetLogLevel() Level {
	return Level(l.logLevel.Load())
}

func (l *Logger) SetLogLevel(level Level) {
	clamped := convertLevel(level)

	l.logLevel.Store(uint32(clamped))
}
