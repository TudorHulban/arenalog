package arenalog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendJSONRoot(t *testing.T) {
	type testCase struct { //nolint:govet
		description string
		timestamp   string
		rootField   *field
		fields      []field
		file        string
		line        int
		msg         string
		expected    string
	}

	tests := []testCase{
		{
			description: "1. Full output: TS, Root, Fields, Caller, and Msg",
			timestamp:   "2023-10-27T10:00:00Z",
			rootField:   &field{key: "service", kind: kindString, valueString: "auth-api"},
			fields: []field{
				{key: "env", kind: kindString, valueString: "prod"},
				{key: "version", kind: kindInt, valueInt: 2},
			},
			file:     "main.go",
			line:     10,
			msg:      "started",
			expected: `{"service":"auth-api","env":"prod","version":2,"caller":"main.go","line":10,"msg":"started"}`,
		},
		{
			description: "2. No timestamp or caller info",
			timestamp:   "",
			rootField:   &field{key: "app", kind: kindString, valueString: "tester"},
			fields:      nil,
			file:        "",
			line:        0,
			msg:         "no-extras",
			expected:    `{"app":"tester","msg":"no-extras"}`,
		},
		{
			description: "3. Diverse field types (bool and float)",
			timestamp:   "",
			rootField:   nil,
			fields: []field{
				{key: "active", kind: kindBool, valueBool: true},
				{key: "ratio", kind: kindFloat, valueFloat: 0.75},
			},
			file:     "logic.go",
			line:     55,
			msg:      "data",
			expected: `{"active":true,"ratio":0.750000000000,"caller":"logic.go","line":55,"msg":"data"}`,
		},
		{
			description: "4. Minimal log (only message)",
			timestamp:   "",
			rootField:   nil,
			fields:      nil,
			file:        "",
			line:        0,
			msg:         "minimal",
			expected:    `{"msg":"minimal"}`,
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				var l Logger

				logContext := NewLogContext(&l)
				require.NotNil(t, logContext)

				fields := make([]field, len(tc.fields))
				copy(fields, tc.fields)

				logContext.cfg.Store(
					&formatterConfig{
						root:   tc.rootField,
						fields: fields,
					},
				)

				buf := make([]byte, 0, 512)
				out := l.appendJSONRoot(
					buf,
					[]byte(tc.msg),
					logContext.cfg.Load(),
					tc.file,
					tc.line,
				)

				require.Equal(t, tc.expected+"\n", string(out))
			},
		)
	}
}
