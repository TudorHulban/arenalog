package main

import (
	"context"
	"io"

	log "github.com/tudorhulban/arenalog"

	fiberlog "github.com/gofiber/fiber/v3/log"
)

// see https://docs.gofiber.io/api/log#global-log

var _ fiberlog.AllLogger[*log.Logger] = (*FiberLogger)(nil)

type FiberLogger struct {
	L *log.Logger
}

// TODO: inline code for more performance instead of function call.

// --- Trace ---
func (f *FiberLogger) Trace(args ...any)                 { f.L.Trace(args...) }
func (f *FiberLogger) Tracef(format string, args ...any) { f.L.Tracef(format, args...) }
func (f *FiberLogger) Tracew(msg string, keysAndValues ...any) {
	f.L.Tracew(msg, keysAndValues...)
}

// --- Debug ---
func (f *FiberLogger) Debug(args ...any)                 { f.L.Debug(args...) }
func (f *FiberLogger) Debugf(format string, args ...any) { f.L.Debugf(format, args...) }
func (f *FiberLogger) Debugw(msg string, keysAndValues ...any) {
	f.L.Debugw(msg, keysAndValues...)
}

// --- Info ---
func (f *FiberLogger) Info(args ...any)                 { f.L.Info(args...) }
func (f *FiberLogger) Infof(format string, args ...any) { f.L.Infof(format, args...) }
func (f *FiberLogger) Infow(msg string, keysAndValues ...any) {
	f.L.Infow(msg, keysAndValues...)
}

// --- Warn ---
func (f *FiberLogger) Warn(args ...any)                 { f.L.Warn(args...) }
func (f *FiberLogger) Warnf(format string, args ...any) { f.L.Warnf(format, args...) }
func (f *FiberLogger) Warnw(msg string, keysAndValues ...any) {
	f.L.Warnw(msg, keysAndValues...)
}

// --- Error ---
func (f *FiberLogger) Error(args ...any)                 { f.L.Error(args...) }
func (f *FiberLogger) Errorf(format string, args ...any) { f.L.Errorf(format, args...) }
func (f *FiberLogger) Errorw(msg string, keysAndValues ...any) {
	f.L.Errorw(msg, keysAndValues...)
}

// --- Fatal ---
func (f *FiberLogger) Fatal(args ...any)                 { f.L.Fatal(args...) }
func (f *FiberLogger) Fatalf(format string, args ...any) { f.L.Fatalf(format, args...) }
func (f *FiberLogger) Fatalw(msg string, keysAndValues ...any) {
	f.L.Fatalw(msg, keysAndValues...)
}

// --- Panic ---
func (f *FiberLogger) Panic(args ...any)                 { f.L.Panic(args...) }
func (f *FiberLogger) Panicf(format string, args ...any) { f.L.Panicf(format, args...) }
func (f *FiberLogger) Panicw(msg string, keysAndValues ...any) {
	f.L.Panicw(msg, keysAndValues...)
}

// --- Print ---
func (f *FiberLogger) Print(args ...any)                 { f.L.Print(args...) }
func (f *FiberLogger) Printf(format string, args ...any) { f.L.Printf(format, args...) }
func (f *FiberLogger) Printw(msg string, keysAndValues ...any) {
	f.L.Printw(msg, keysAndValues...)
}

// --- Fiber-required structured logging ---
func (f *FiberLogger) With(args ...any) fiberlog.AllLogger[*log.Logger] {
	// Your logger does not support structured fields natively.
	// No-op is acceptable.
	return f
}

func (f *FiberLogger) WithGroup(name string) fiberlog.AllLogger[*log.Logger] {
	// Same: no-op unless you want grouping.
	return f
}

func (f *FiberLogger) Logger() *log.Logger {
	return f.L
}

func (f *FiberLogger) SetLevel(level fiberlog.Level) {
	switch level { //nolint:revive
	case fiberlog.LevelTrace:
		f.L.SetLogLevel(log.LevelTrace)
	case fiberlog.LevelDebug:
		f.L.SetLogLevel(log.LevelDebug)
	case fiberlog.LevelInfo:
		f.L.SetLogLevel(log.LevelInfo)
	case fiberlog.LevelWarn:
		f.L.SetLogLevel(log.LevelWarn)
	case fiberlog.LevelError:
		f.L.SetLogLevel(log.LevelError)
	case fiberlog.LevelFatal:
		f.L.SetLogLevel(log.LevelFatal)
	case fiberlog.LevelPanic:
		f.L.SetLogLevel(log.LevelPanic)

	default:
		// fallback to Info
		f.L.SetLogLevel(log.LevelInfo)
	}
}

func (*FiberLogger) SetOutput(w io.Writer) {
	// Your logger does not expose SetOutput directly,
	// but the standard pattern is to redirect PrintRaw / PrintMessage
	// through a writer. If your logger has no writer concept,
	// you must store it and ignore it.
	//
	// Minimal no‑op implementation:
	_ = w
}

// should be below, requires changes to the ingestor.
// func (f *FiberLogger) SetOutput(w io.Writer) {
//     if f.L == nil || f.L.ingestor == nil {
//         return
//     }
//     f.L.ingestor.SetWriter(w)
// }

func (f *FiberLogger) WithContext(ctx context.Context) fiberlog.CommonLogger {
	// Your logger does not use context, so this is a no‑op.
	return f
}
