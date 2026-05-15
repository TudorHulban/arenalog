package helpers

var digitPairs = [...]byte{
	'0', '0', '0', '1', '0', '2', '0', '3', '0', '4', '0', '5', '0', '6', '0', '7', '0', '8', '0', '9',
	'1', '0', '1', '1', '1', '2', '1', '3', '1', '4', '1', '5', '1', '6', '1', '7', '1', '8', '1', '9',
	'2', '0', '2', '1', '2', '2', '2', '3', '2', '4', '2', '5', '2', '6', '2', '7', '2', '8', '2', '9',
	'3', '0', '3', '1', '3', '2', '3', '3', '3', '4', '3', '5', '3', '6', '3', '7', '3', '8', '3', '9',
	'4', '0', '4', '1', '4', '2', '4', '3', '4', '4', '4', '5', '4', '6', '4', '7', '4', '8', '4', '9',
	'5', '0', '5', '1', '5', '2', '5', '3', '5', '4', '5', '5', '5', '6', '5', '7', '5', '8', '5', '9',
	'6', '0', '6', '1', '6', '2', '6', '3', '6', '4', '6', '5', '6', '6', '6', '7', '6', '8', '6', '9',
	'7', '0', '7', '1', '7', '2', '7', '3', '7', '4', '7', '5', '7', '6', '7', '7', '7', '8', '7', '9',
	'8', '0', '8', '1', '8', '2', '8', '3', '8', '4', '8', '5', '8', '6', '8', '7', '8', '8', '8', '9',
	'9', '0', '9', '1', '9', '2', '9', '3', '9', '4', '9', '5', '9', '6', '9', '7', '9', '8', '9', '9',
}

var pow10 = [...]uint64{
	1,
	10,
	100,
	1_000,
	10_000,
	100_000,
	1_000_000,
	10_000_000,
	100_000_000,
	1_000_000_000,
	10_000_000_000,
	100_000_000_000,
	1_000_000_000_000,
	10_000_000_000_000,
	100_000_000_000_000,
	1_000_000_000_000_000,
	10_000_000_000_000_000,
	100_000_000_000_000_000,
	1_000_000_000_000_000_000,
}

func appendUint64(destination []byte, value uint64) []byte {
	var buf [20]byte // Write into a fixed-size scratch buffer, then copy the tail.

	i := len(buf)

	for value >= 100 {
		q := value / 100
		r := value - q*100
		p := int(r) * 2

		i = i - 2
		buf[i+0] = digitPairs[p+0]
		buf[i+1] = digitPairs[p+1]

		value = q
	}

	// Final 1–2 digits
	if value < 10 {
		i--
		buf[i] = byte('0' + value)
	} else {
		p := int(value) * 2
		i = i - 2
		buf[i+0] = digitPairs[p+0]
		buf[i+1] = digitPairs[p+1]
	}

	return append(destination, buf[i:]...)
}

func appendUintZeroPadded(destination []byte, value uint64, width int) []byte {
	var buf [20]byte // Write into a fixed-size scratch, then copy tail.

	i := len(buf)

	for value >= 100 {
		q := value / 100
		r := value - q*100
		p := int(r) * 2

		i -= 2
		buf[i+0] = digitPairs[p+0]
		buf[i+1] = digitPairs[p+1]

		value = q
	}

	if value < 10 {
		i--
		buf[i] = byte('0' + value)
	} else {
		p := int(value) * 2
		i -= 2
		buf[i+0] = digitPairs[p+0]
		buf[i+1] = digitPairs[p+1]
	}

	// Ensure at least width digits (zero‑padded on the left).
	for len(buf)-i < width {
		i--
		buf[i] = '0'
	}

	return append(destination, buf[i:]...)
}

func AppendFloat(destination []byte, value float64, precision int) []byte {
	// Handle NaN and Inf explicitly
	if value != value {
		return append(destination, 'n', 'a', 'n')
	}

	if value > 1e308 {
		return append(destination, 'i', 'n', 'f')
	}

	if value < -1e308 {
		return append(destination, '-', 'i', 'n', 'f')
	}

	// Sign
	if value < 0 {
		destination = append(destination, '-')
		value = -value
	}

	// Integer part (no zero padding)
	intPart := uint64(value)
	destination = appendUint64(destination, intPart)

	// Fractional part
	if precision > 0 {
		destination = append(destination, '.')

		fractional := value - float64(intPart)

		if precision >= len(pow10) {
			precision = len(pow10) - 1
		}

		scale := float64(pow10[precision])

		// Truncate, do not round.
		fv := fractional * scale
		if fv < 0 {
			fv = 0
		}

		fracInt := uint64(fv)

		destination = appendUintZeroPadded(destination, fracInt, precision)
	}

	return destination
}
