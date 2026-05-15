package timestamp

import (
	"strconv"
	"time"
)

// For no timestamp do not add a timestamp function.
type Timestamp func(appendTo []byte) []byte

// TimestampNano provides true nanosecond‑accurate timestamps.
// On Linux time.Now() costs ~40–70 ns by itself.
// UnixNano() + AppendInt adds ~10–15 ns.
//
// Due to above, cost is around 150 ns.
func TimestampNano(appendTo []byte) []byte {
	return strconv.AppendInt(appendTo, time.Now().UnixNano(), 10)
}

func TimestampRFC3339Nano(appendTo []byte) []byte {
	return time.Now().UTC().AppendFormat(appendTo, time.RFC3339Nano)
}

func TimestampRFC3339UTC(appendTo []byte) []byte {
	buf := timeCacheRFC3339.active.Load() // atomic pointer load, ~1 ns
	if buf != nil {
		return append(appendTo, buf.output[:buf.length]...)
	}

	return time.Now().UTC().AppendFormat(appendTo, rfc3339MilliLayout)
}

func TimestampRFC3339Bucharest(appendTo []byte) []byte {
	buf := timeCacheRFC3339.active.Load() // atomic pointer load, ~1 ns
	if buf != nil {
		return append(appendTo, buf.output[:buf.length]...)
	}

	loc, _ := time.LoadLocation("Europe/Bucharest")

	return time.Now().In(loc).AppendFormat(appendTo, time.RFC3339)
}

func TimestampRFC3339CustomLocation(appendTo []byte) []byte {
	buf := timeCacheRFC3339.active.Load() // atomic pointer load, ~1 ns
	if buf != nil {
		return append(appendTo, buf.output[:buf.length]...)
	}

	// fallback if cache not ready, allow time for the cache to start.
	return time.Now().UTC().AppendFormat(appendTo, rfc3339MilliLayout)
}

func TimestampYYYYMonth(appendTo []byte) []byte {
	buf := timeCacheYYYYMonth.active.Load()
	if buf != nil {
		return append(appendTo, buf.output[:buf.length]...)
	}

	// cold-start fallback: build manually, same layout as buildYYYYMonthCache
	now := time.Now()
	year, month, day := now.Date()
	hour, minute, sec := now.Clock()
	milli := now.Nanosecond() / 1e6

	appendTo = strconv.AppendInt(appendTo, int64(year), 10)
	if month < 10 {
		appendTo = append(appendTo, '0')
	}

	appendTo = strconv.AppendInt(appendTo, int64(month), 10)

	appendTo = append(appendTo, ' ')
	if day < 10 {
		appendTo = append(appendTo, '0')
	}

	appendTo = strconv.AppendInt(appendTo, int64(day), 10)

	appendTo = append(appendTo, ' ')
	if hour < 10 {
		appendTo = append(appendTo, '0')
	}

	appendTo = strconv.AppendInt(appendTo, int64(hour), 10)

	appendTo = append(appendTo, ':')
	if minute < 10 {
		appendTo = append(appendTo, '0')
	}

	appendTo = strconv.AppendInt(appendTo, int64(minute), 10)

	appendTo = append(appendTo, ':')
	if sec < 10 {
		appendTo = append(appendTo, '0')
	}

	appendTo = strconv.AppendInt(appendTo, int64(sec), 10)

	appendTo = append(appendTo, '.')
	if milli < 100 {
		appendTo = append(appendTo, '0')
	}

	if milli < 10 {
		appendTo = append(appendTo, '0')
	}

	return strconv.AppendInt(appendTo, int64(milli), 10)
}

func TimestampStandard(appendTo []byte) []byte {
	buf := timeCacheStandard.active.Load()
	if buf != nil {
		return append(appendTo, buf.output[:buf.length]...)
	}

	// cold-start fallback: YYYY/MM/DD HH:MM:SS.mmm
	now := time.Now()
	year, month, day := now.Date()
	hour, minute, sec := now.Clock()
	milli := now.Nanosecond() / 1e6

	appendTo = strconv.AppendInt(appendTo, int64(year), 10)

	appendTo = append(appendTo, '/')
	if month < 10 {
		appendTo = append(appendTo, '0')
	}

	appendTo = strconv.AppendInt(appendTo, int64(month), 10)

	appendTo = append(appendTo, '/')
	if day < 10 {
		appendTo = append(appendTo, '0')
	}

	appendTo = strconv.AppendInt(appendTo, int64(day), 10)

	appendTo = append(appendTo, ' ')
	if hour < 10 {
		appendTo = append(appendTo, '0')
	}

	appendTo = strconv.AppendInt(appendTo, int64(hour), 10)

	appendTo = append(appendTo, ':')
	if minute < 10 {
		appendTo = append(appendTo, '0')
	}

	appendTo = strconv.AppendInt(appendTo, int64(minute), 10)

	appendTo = append(appendTo, ':')
	if sec < 10 {
		appendTo = append(appendTo, '0')
	}

	appendTo = strconv.AppendInt(appendTo, int64(sec), 10)

	appendTo = append(appendTo, '.')
	if milli < 100 {
		appendTo = append(appendTo, '0')
	}

	if milli < 10 {
		appendTo = append(appendTo, '0')
	}

	return strconv.AppendInt(appendTo, int64(milli), 10)
}
