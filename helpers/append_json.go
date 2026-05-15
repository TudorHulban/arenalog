package helpers

// AppendJSON writes s into buf with JSON string escaping but without
// the surrounding quotes. Used by appendArgsQuotedJSON.
func AppendJSON(buf, word []byte) []byte {
	const hex = "0123456789abcdef"

	for ix := range word {
		char := word[ix]

		switch char {
		case '\\', '"':
			buf = append(buf, '\\', char)
		case '\n':
			buf = append(buf, '\\', 'n')
		case '\r':
			buf = append(buf, '\\', 'r')
		case '\t':
			buf = append(buf, '\\', 't')
		case '\b':
			buf = append(buf, '\\', 'b')
		case '\f':
			buf = append(buf, '\\', 'f')

		default:
			if char < 0x20 {
				buf = append(buf, '\\', 'u', '0', '0', hex[char>>4], hex[char&0xF])
			} else {
				buf = append(buf, char)
			}
		}
	}

	return buf
}

func AppendJSON_Quoted(buf []byte, word []byte) []byte {
	buf = append(buf, '"')

	const hex = "0123456789abcdef"

	for ix := range word {
		char := word[ix]

		switch char {
		case '\\', '"':
			buf = append(buf, '\\', char)
		case '\n':
			buf = append(buf, '\\', 'n')
		case '\r':
			buf = append(buf, '\\', 'r')
		case '\t':
			buf = append(buf, '\\', 't')
		case '\b':
			buf = append(buf, '\\', 'b')
		case '\f':
			buf = append(buf, '\\', 'f')

		default:
			if char < 0x20 {
				buf = append(buf, '\\', 'u', '0', '0', hex[char>>4], hex[char&0xF])
			} else {
				buf = append(buf, char)
			}
		}
	}

	return append(buf, '"')
}
