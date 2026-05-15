package arenalog

import (
	"fmt"
)

/*
Colour Styles
	//Text Colour
	Black   = 30
	Red     = 31
	Green   = 32
	Yellow  = 33
	Blue    = 34
	Magenta = 35
	Cyan    = 36
	White   = 37
	Grey    = 90

	//Style
	Bold      = 1
	Dim       = 2
	Italic    = 3
	Underline = 4
	Blinkslow = 5
	Blinkfast = 6
	Inverse   = 7
	Hidden    = 8
	Strikeout = 9
*/

func colorTrace(word string) string {
	return fmt.Sprintf(
		"%s%v\x1b[0m",

		"\x1b[38;5;135m",
		word,
	)
}

func colorDebug(word string) string {
	return fmt.Sprintf(
		"%s%v\x1b[0m",

		"\x1b[1;34m",
		word,
	)
}

func colorInfo(word string) string {
	return fmt.Sprintf(
		"%s%v\x1b[0m",

		"\x1b[1;32m",
		word,
	)
}

func colorWarn(word string) string {
	return fmt.Sprintf(
		"%s%v\x1b[0m",

		"\x1b[1;33m",
		word,
	)
}

func colorError(word string) string {
	return fmt.Sprintf(
		"%s%v\x1b[0m",

		"\x1b[1;31m",
		word,
	)
}
