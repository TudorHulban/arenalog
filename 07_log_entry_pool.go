package arenalog

import "sync"

var entryPool = sync.Pool{
	New: func() any {
		return &Entry{
			fields: [8]field{}, // small reusable buffer
		}
	},
}
