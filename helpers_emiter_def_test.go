package arenalog

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type message struct {
	level       Level
	wasReceived bool
}

type state struct {
	dictionary map[string]message // all requests
	invariants map[Level]*uint32

	prints *uint32
	totals uint32

	muDictionary sync.Mutex
}

func (s *state) String() string {
	return s.stringVerbose(false)
}

func (s *state) StringVerbose() string {
	return s.stringVerbose(true)
}

func (s *state) stringVerbose(verbose bool) string {
	var sb strings.Builder

	sb.WriteString("state{\n")

	// totals
	fmt.Fprintf(&sb, "  totals: %d,\n", s.totals)

	// prints
	if s.prints != nil {
		fmt.Fprintf(&sb, "  prints: %d,\n", *s.prints)
	} else {
		fmt.Fprint(&sb, "  prints: <nil>,\n")
	}

	// invariants
	sb.WriteString("  invariants: [")

	first := true

	for lvl := LevelTrace; lvl <= LevelPanic; lvl++ {
		if !first {
			sb.WriteString(" ")
		}

		first = false

		if cnt, ok := s.invariants[lvl]; ok && cnt != nil {
			fmt.Fprintf(&sb, "%s:%d", lvl.String(), *cnt)
		} else {
			fmt.Fprintf(&sb, "%s:*", lvl.String())
		}
	}

	sb.WriteString("],\n")

	// dictionary summary
	received := 0
	byLevel := make(map[Level]uint32)

	s.muDictionary.Lock()
	for _, msg := range s.dictionary {
		if msg.wasReceived {
			received++
		}

		byLevel[msg.level]++
	}
	s.muDictionary.Unlock()

	sb.WriteString("  dict: {\n")
	fmt.Fprintf(&sb, "    total: %d,\n", len(s.dictionary))
	fmt.Fprintf(&sb, "    received: %d,\n", received)
	sb.WriteString("    byLevel: [")

	vFirst := true

	for lvl := LevelTrace; lvl <= LevelPanic; lvl++ {
		if !vFirst {
			sb.WriteString(" ")
		}

		vFirst = false

		fmt.Fprintf(&sb, "%s:%d", lvl.String(), byLevel[lvl])
	}

	sb.WriteString("]\n")
	sb.WriteString("  }")

	// entries (always show in verbose, optionally in concise)
	if verbose || len(s.dictionary) <= 5 {
		sb.WriteString(",\n  entries: [")

		ids := make([]string, 0, len(s.dictionary))

		for id := range s.dictionary {
			ids = append(ids, id)
		}

		sort.Strings(ids)

		for _, id := range ids {
			msg := s.dictionary[id]
			rcv := ""

			if msg.wasReceived {
				rcv = " ✓"
			}

			fmt.Fprintf(
				&sb,
				"\n    %s: %s%s",
				id,
				msg.level.String(),
				rcv,
			)
		}

		sb.WriteString("\n  ]")
	}

	sb.WriteString("\n}")

	return sb.String()
}

func createEmitData(total uint32, invariants map[Level]*uint32, prints *int) (*state, error) {
	// Minimal validation: reserved must not exceed total
	var reserved uint32

	for _, cnt := range invariants {
		if cnt != nil {
			reserved = reserved + *cnt
		}
	}

	if prints != nil {
		reserved = reserved + uint32(*prints) //nolint:gosec
	}

	if reserved > total {
		return nil,
			fmt.Errorf("reserved (%d) > total (%d)", reserved, total)
	}

	// Collect flexible levels: nil invariants, skip LevelNONE if prints is fixed
	var flexible []Level

	for lvl := LevelTrace; lvl <= LevelPanic; lvl++ {
		if lvl == LevelPanic && prints != nil {
			continue
		}

		if invariants[lvl] == nil {
			flexible = append(flexible, lvl)
		}
	}

	// Build dictionary with exactly `total` entries
	dictionary := make(map[string]message, total)

	var reqID uint32

	// Assign fixed invariant levels
	for lvl, cnt := range invariants {
		if cnt != nil {
			for i := uint32(0); i < *cnt; i++ {
				dictionary[fmt.Sprintf("req-%d", reqID)] =
					message{
						level:       lvl,
						wasReceived: false,
					}

				reqID++
			}
		}
	}

	// Assign PRINT (LevelNONE) messages
	if prints != nil {
		for i := uint32(0); i < uint32(*prints); i++ { //nolint:gosec
			dictionary[fmt.Sprintf("req-%d", reqID)] =
				message{
					level:       LevelPanic,
					wasReceived: false,
				}

			reqID++
		}
	}

	// Distribute remaining via round-robin to flexible levels
	remaining := total - reqID

	if len(flexible) > 0 && remaining > 0 {
		for i := range remaining {
			lvl := flexible[i%uint32(len(flexible))] //nolint:gosec

			dictionary[fmt.Sprintf("req-%d", reqID)] =
				message{
					level:       lvl,
					wasReceived: false,
				}

			reqID++
		}
	}

	// Initialize state prints pointer
	var statePrints *uint32

	if prints != nil {
		v := uint32(*prints) //nolint:gosec
		statePrints = &v
	}

	return &state{
			invariants: invariants,
			dictionary: dictionary,
			prints:     statePrints,
			totals:     total,
		},
		nil
}

func emitData(data *state, logger *Logger) []error { //nolint:revive
	var (
		errs []error
		idx  uint32
	)

	consecFails := make(map[string]int)

	// Deterministic iteration: map order is randomized in Go
	data.muDictionary.Lock()

	ids := make([]string, 0, len(data.dictionary))

	for id := range data.dictionary {
		ids = append(ids, id)
	}

	data.muDictionary.Unlock()

	sort.Strings(ids)

	for _, id := range ids {
		idx++

		data.muDictionary.Lock()
		lvl := data.dictionary[id].level
		data.muDictionary.Unlock()

		var (
			call   func() error
			method string
		)

		switch lvl {
		case LevelDebug:
			switch idx % 3 {
			case 0:
				call, method = func() error { logger.Debug(id); return nil }, "Debug"
			case 1:
				call, method = func() error { logger.Debugf("%s", id); return nil }, "Debugf"
			case 2:
				call, method = func() error { logger.Debugw(id); return nil }, "Debugw"
			}

		case LevelInfo:
			switch idx % 3 {
			case 0:
				call, method = func() error { logger.Info(id); return nil }, "Info"
			case 1:
				call, method = func() error { logger.Infof("%s", id); return nil }, "Infof"
			case 2:
				call, method = func() error { logger.Infow(id); return nil }, "Infow"
			}

		case LevelWarn:
			switch idx % 3 {
			case 0:
				call, method = func() error { logger.Warn(id); return nil }, "Warn"
			case 1:
				call, method = func() error { logger.Warnf("%s", id); return nil }, "Warnf"
			case 2:
				call, method = func() error { logger.Warnw(id); return nil }, "Warnw"
			}

		case LevelError:
			switch idx % 3 {
			case 0:
				call, method = func() error { logger.Error(id); return nil }, "Error"
			case 1:
				call, method = func() error { logger.Errorf("%s", id); return nil }, "Errorf"
			case 2:
				call, method = func() error { logger.Errorw(id); return nil }, "Errorw"
			}

		case LevelPanic:
			switch idx % 5 {
			case 0:
				call, method = func() error { logger.Print(id); return nil }, "Print"
			case 1:
				call, method = func() error { logger.Printf("%s", id); return nil }, "Printf"
			case 2:
				call, method = func() error { logger.Printw(id); return nil }, "Printw"
			case 3:
				call, method = func() error { logger.Msg(id); return nil }, "PrintMessage"
			case 4:
				call, method = func() error { logger.PrintRaw([]byte(id)); return nil }, "PrintRaw"
			}

		// Fallback for undefined levels (TRACE, FATAL, PANIC, etc.)
		default:
			call, method = func() error { logger.Debug(id); return nil }, "DebugFallback"
		}

		// Execute and track errors
		if errCall := call(); errCall != nil {
			errs = append(errs, fmt.Errorf("%s(%q): %w", method, id, errCall))
			consecFails[method]++

			if consecFails[method] >= 2 {
				return errs // Abort on 2 consecutive failures for this method
			}
		} else {
			consecFails[method] = 0 // Reset on success
		}
	}

	return errs
}
