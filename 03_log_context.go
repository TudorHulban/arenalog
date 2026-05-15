package arenalog

import (
	"sync/atomic"
)

type formatterConfig struct {
	root   *field  // nil if no root
	fields []field // ephemeral fields
}

// LogContext acts as the root.
type LogContext struct {
	logger *Logger
	cfg    atomic.Pointer[formatterConfig]
}

func NewLogContext(logger *Logger) *LogContext {
	f := LogContext{
		logger: logger,
	}

	f.cfg.Store(
		&formatterConfig{
			fields: nil,
		},
	)

	return &f
}

func (ctx *LogContext) WithRoot(key string, value any) *LogContext {
	old := ctx.cfg.Load()

	// copy ephemeral fields
	newFields := make([]field, len(old.fields))
	copy(newFields, old.fields)

	// replace root
	ctx.cfg.Store(
		&formatterConfig{
			root:   makeFieldPtr(key, value),
			fields: newFields,
		},
	)

	return ctx
}

func (ctx *LogContext) SetString(key, value string) *LogContext {
	old := ctx.cfg.Load()

	newFields := make([]field, len(old.fields)+1)
	copy(newFields, old.fields)

	newFields[len(old.fields)] = field{
		key:         key,
		kind:        kindString,
		valueString: value,
	}

	ctx.cfg.Store(
		&formatterConfig{
			root:   old.root,  // keep root
			fields: newFields, // updated ephemeral fields
		},
	)

	return ctx
}

func (ctx *LogContext) SetInt(key string, value int64) *LogContext {
	old := ctx.cfg.Load()

	newFields := make([]field, len(old.fields)+1)
	copy(newFields, old.fields)

	newFields[len(old.fields)] = field{
		key:      key,
		kind:     kindInt,
		valueInt: value,
	}

	ctx.cfg.Store(
		&formatterConfig{
			root:   old.root,
			fields: newFields,
		},
	)

	return ctx
}

func (ctx *LogContext) SetBool(key string, value bool) *LogContext {
	old := ctx.cfg.Load()

	newFields := make([]field, len(old.fields)+1)
	copy(newFields, old.fields)

	newFields[len(old.fields)] = field{
		key:       key,
		kind:      kindBool,
		valueBool: value,
	}

	ctx.cfg.Store(
		&formatterConfig{
			root:   old.root,
			fields: newFields,
		},
	)

	return ctx
}

func (ctx *LogContext) Clear() {
	old := ctx.cfg.Load()

	ctx.cfg.Store(
		&formatterConfig{
			root:   old.root, // keep root
			fields: nil,      // clear ephemeral fields
		},
	)
}

func (ctx *LogContext) Reset() {
	ctx.cfg.Store(&formatterConfig{})
}

func (ctx *LogContext) With(key string, value any) *Entry {
	e, _ := entryPool.Get().(*Entry) //nolint:revive

	e.formatter = ctx
	e.fieldCount = 0

	e.fields[e.fieldCount] = makeField(key, value)
	e.fieldCount++

	return e
}

func (ctx *LogContext) WithString(key, value string) *Entry {
	e, _ := entryPool.Get().(*Entry) //nolint:revive

	e.formatter = ctx
	e.fieldCount = 0

	e.fields[e.fieldCount] = field{
		key:         key,
		valueString: value,
	}
	e.fieldCount++

	return e
}

func (ctx *LogContext) WithInt(key string, value int64) *Entry {
	e := entryPool.Get().(*Entry) //nolint:revive

	e.formatter = ctx
	e.fieldCount = 0

	e.fields[e.fieldCount] = field{
		key:      key,
		valueInt: value,
	}
	e.fieldCount++

	return e
}

func (ctx *LogContext) WithFloat(key string, value float64) *Entry {
	e := entryPool.Get().(*Entry) //nolint:revive

	e.formatter = ctx
	e.fieldCount = 0

	e.fields[e.fieldCount] = field{
		key:        key,
		valueFloat: value,
	}
	e.fieldCount++

	return e
}

func (ctx *LogContext) WithBool(key string, value bool) *Entry {
	e := entryPool.Get().(*Entry) //nolint:revive

	e.formatter = ctx
	e.fieldCount = 0

	e.fields[e.fieldCount] = field{
		key:       key,
		valueBool: value,
	}
	e.fieldCount++

	return e
}
