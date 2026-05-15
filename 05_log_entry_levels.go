package arenalog

func (e *Entry) Trace() *Entry {
	e.level = LevelTrace

	return e
}

func (e *Entry) Debug() *Entry {
	e.level = LevelDebug

	return e
}

func (e *Entry) Info() *Entry {
	e.level = LevelInfo

	return e
}

func (e *Entry) Warn() *Entry {
	e.level = LevelWarn

	return e
}

func (e *Entry) Error() *Entry {
	e.level = LevelError

	return e
}

func (e *Entry) Fatal() *Entry {
	e.level = LevelFatal

	return e
}

func (e *Entry) Panic() *Entry {
	e.level = LevelPanic

	return e
}
