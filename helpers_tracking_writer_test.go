package arenalog

import (
	"io"
	"strings"
)

type trackingWriter struct {
	w     io.Writer
	state *state
}

func newTrackingWriter(w io.Writer, s *state) *trackingWriter {
	return &trackingWriter{w: w, state: s}
}

func (tw *trackingWriter) Write(p []byte) (int, error) {
	bytesWritten, err := tw.w.Write(p)
	if err != nil {
		return bytesWritten, err
	}

	data := string(p)

	// Use the SAME mutex as emitData
	tw.state.muDictionary.Lock()

	for key, msg := range tw.state.dictionary {
		if !msg.wasReceived && strings.Contains(data, key) {
			msg.wasReceived = true
			tw.state.dictionary[key] = msg // reassign to persist struct update
		}
	}

	tw.state.muDictionary.Unlock()

	return bytesWritten, nil
}

func (tw *trackingWriter) UnreceivedCount() int {
	tw.state.muDictionary.Lock()
	defer tw.state.muDictionary.Unlock()

	count := 0

	for _, msg := range tw.state.dictionary {
		if !msg.wasReceived {
			count++
		}
	}

	return count
}

func (tw *trackingWriter) UnreceivedMessages() []string {
	tw.state.muDictionary.Lock()
	defer tw.state.muDictionary.Unlock()

	var missing []string

	for key, msg := range tw.state.dictionary {
		if !msg.wasReceived {
			missing = append(missing, key)
		}
	}

	return missing
}
