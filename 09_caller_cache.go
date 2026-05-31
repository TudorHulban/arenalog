package arenalog

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

type callerCacheEntry struct {
	file string
	pc   uintptr
	line int
}

const cacheSize = 1024
const cacheMask = cacheSize - 1

var callerTable [cacheSize]unsafe.Pointer

func (*Logger) slowPathCaller(pc uintptr, idx uintptr) (string, int) {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "", 0
	}

	file, line := fn.FileLine(pc)

	// Create new entry
	newEntry := &callerCacheEntry{
		pc:   pc,
		file: file,
		line: line,
	}

	// Atomic store (last-writer-wins, race-safe)
	atomic.StorePointer(&callerTable[idx], unsafe.Pointer(newEntry))

	return file, line
}

// getCallerData retrieves the file name and line number for a specific stack depth.
// It uses a lock-free cache to avoid the high overhead of runtime.FuncForPC and
// runtime.FileLine on every log call.
//
// Parameters:
//   - skip: The number of stack frames to skip, where 0 is the immediate caller of getCallerData.
//
// Implementation Details:
//   - It adds an internal offset of 2 to 'skip' to account for runtime.Callers and getCallerData itself.
//   - It subtracts 1 from the PC to convert the return address into the call site address,
//     ensuring the line number matches the actual invocation.
//   - The cache uses a simple hash-mapped array with atomic pointers; if a collision occurs,
//     the "last writer wins," which maintains thread safety without locking.
func (l *Logger) getCallerData(skip int) (string, int) {
	// 1. Get the PC at the specified skip level.
	// We use a small stack buffer.
	var pcs [1]uintptr

	// Increment skip by 2:
	// +1 to get out of runtime.Callers
	// +1 to get out of getCallerData
	n := runtime.Callers(skip+2, pcs[:])
	if n == 0 {
		return "", 0
	}

	// Use pc-1 to ensure we are inside the call site's range
	pc := pcs[0] - 1

	// 2. Hash the PC
	idx := (pc >> 4) & cacheMask

	ptr := atomic.LoadPointer(&callerTable[idx])
	if ptr != nil {
		entry := (*callerCacheEntry)(ptr)

		// Race-Safe: Even if another goroutine evicts or replaces this slot in
		// callerTable between our check and return, 'entry' points to a heap-allocated,
		// immutable struct that cannot be garbage collected or mutated.
		if entry.pc == pc {
			return entry.file,
				entry.line
		}
	}

	// 3. Slow Path: Resolve symbols only if PC not in cache
	return l.slowPathCaller(pc, idx)
}
