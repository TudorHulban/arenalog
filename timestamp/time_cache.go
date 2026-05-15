package timestamp

import (
	"strconv"
	"sync/atomic"
	"time"
)

type timeBuf struct {
	valueMillisecond int64 // millisecond-epoch at which this buf was built
	valueDay         int64 // Unix-day at which the date prefix was built

	length  int // total byte length of the formatted timestamp
	dateLen int // byte length of the date prefix inside data

	output [32]byte // the formatted timestamp bytes
}

type timeCache struct {
	active atomic.Pointer[timeBuf]
}

var (
	timeCacheStandard  timeCache
	timeCacheYYYYMonth timeCache

	timeCacheRFC3339 timeCache
)

func buildRFC3339Cache(now time.Time) {
	nowNano := now.UnixNano()
	nowMs := nowNano / 1e6

	prev := timeCacheRFC3339.active.Load()
	if prev != nil && prev.valueMillisecond == nowMs {
		return // already current
	}

	next := new(timeBuf)
	next.valueMillisecond = nowMs

	// Get the correct day representation using the local date,
	// flexible and not tied to UTC math.
	year, month, day := now.Date()
	nowDay := int64(year)*10000 + int64(month)*100 + int64(day) // Unique daily integer YYYYMMDD

	if prev != nil {
		*next = *prev
	}

	if prev == nil || nowDay != prev.valueDay {
		next.valueDay = nowDay
		buffer := next.output[:0]
		buffer = strconv.AppendInt(buffer, int64(year), 10)

		buffer = append(buffer, '-')
		if month < 10 {
			buffer = append(buffer, '0')
		}

		buffer = strconv.AppendInt(buffer, int64(month), 10)

		buffer = append(buffer, '-')
		if day < 10 {
			buffer = append(buffer, '0')
		}

		buffer = strconv.AppendInt(buffer, int64(day), 10)
		buffer = append(buffer, 'T')

		next.dateLen = len(buffer)
	}

	hour, minute, sec := now.Clock()
	milli := now.Nanosecond() / 1e6
	buffer := next.output[next.dateLen:next.dateLen]

	if hour < 10 {
		buffer = append(buffer, '0')
	}

	buffer = strconv.AppendInt(buffer, int64(hour), 10)

	buffer = append(buffer, ':')
	if minute < 10 {
		buffer = append(buffer, '0')
	}

	buffer = strconv.AppendInt(buffer, int64(minute), 10)

	buffer = append(buffer, ':')
	if sec < 10 {
		buffer = append(buffer, '0')
	}

	buffer = strconv.AppendInt(buffer, int64(sec), 10)

	buffer = append(buffer, '.')
	if milli < 100 {
		buffer = append(buffer, '0')
	}

	if milli < 10 {
		buffer = append(buffer, '0')
	}

	buffer = strconv.AppendInt(buffer, int64(milli), 10)

	// Dynamically calculate and append the passed offset.
	_, offsetSecs := now.Zone()
	if offsetSecs == 0 {
		buffer = append(buffer, 'Z')
	} else {
		if offsetSecs > 0 {
			buffer = append(buffer, '+')
		} else {
			buffer = append(buffer, '-')
			offsetSecs = -offsetSecs
		}

		offsetMins := offsetSecs / 60
		offsetHours := offsetMins / 60
		offsetMins = offsetMins % 60

		if offsetHours < 10 {
			buffer = append(buffer, '0')
		}

		buffer = strconv.AppendInt(buffer, int64(offsetHours), 10)

		buffer = append(buffer, ':')
		if offsetMins < 10 {
			buffer = append(buffer, '0')
		}

		buffer = strconv.AppendInt(buffer, int64(offsetMins), 10)
	}

	next.length = next.dateLen + len(buffer)
	timeCacheRFC3339.active.Store(next)
}

func buildStandardCache(now time.Time) {
	nowNano := now.UnixNano()
	nowMs := nowNano / 1e6

	prev := timeCacheStandard.active.Load()
	if prev != nil && prev.valueMillisecond == nowMs {
		return
	}

	next := new(timeBuf)
	next.valueMillisecond = nowMs
	nowDay := nowNano / nanosPerDay

	if prev != nil {
		*next = *prev
	}

	if prev == nil || nowDay != prev.valueDay {
		next.valueDay = nowDay
		year, month, day := now.Date()
		buffer := next.output[:0]
		buffer = strconv.AppendInt(buffer, int64(year), 10)

		buffer = append(buffer, '/')
		if month < 10 {
			buffer = append(buffer, '0')
		}

		buffer = strconv.AppendInt(buffer, int64(month), 10)

		buffer = append(buffer, '/')
		if day < 10 {
			buffer = append(buffer, '0')
		}

		buffer = strconv.AppendInt(buffer, int64(day), 10)
		buffer = append(buffer, ' ')

		next.dateLen = len(buffer)
	}

	hour, minute, sec := now.Clock()
	milli := now.Nanosecond() / 1e6
	buffer := next.output[next.dateLen:next.dateLen]

	if hour < 10 {
		buffer = append(buffer, '0')
	}

	buffer = strconv.AppendInt(buffer, int64(hour), 10)

	buffer = append(buffer, ':')
	if minute < 10 {
		buffer = append(buffer, '0')
	}

	buffer = strconv.AppendInt(buffer, int64(minute), 10)

	buffer = append(buffer, ':')
	if sec < 10 {
		buffer = append(buffer, '0')
	}

	buffer = strconv.AppendInt(buffer, int64(sec), 10)

	buffer = append(buffer, '.')
	if milli < 100 {
		buffer = append(buffer, '0')
	}

	if milli < 10 {
		buffer = append(buffer, '0')
	}

	buffer = strconv.AppendInt(buffer, int64(milli), 10)

	next.length = next.dateLen + len(buffer)
	timeCacheStandard.active.Store(next)
}

func buildYYYYMonthCache(now time.Time) {
	nowNano := now.UnixNano()
	nowMs := nowNano / 1e6

	prev := timeCacheYYYYMonth.active.Load()
	if prev != nil && prev.valueMillisecond == nowMs {
		return
	}

	next := new(timeBuf)
	next.valueMillisecond = nowMs
	nowDay := nowNano / nanosPerDay

	if prev != nil {
		*next = *prev
	}

	if prev == nil || nowDay != prev.valueDay {
		next.valueDay = nowDay
		year, month, day := now.Date()
		buffer := next.output[:0]

		buffer = strconv.AppendInt(buffer, int64(year), 10)
		if month < 10 {
			buffer = append(buffer, '0')
		}

		buffer = strconv.AppendInt(buffer, int64(month), 10)

		buffer = append(buffer, ' ')
		if day < 10 {
			buffer = append(buffer, '0')
		}

		buffer = strconv.AppendInt(buffer, int64(day), 10)
		buffer = append(buffer, ' ')
		next.dateLen = len(buffer)
	}

	hour, minute, sec := now.Clock()
	milli := now.Nanosecond() / 1e6
	buffer := next.output[next.dateLen:next.dateLen]

	if hour < 10 {
		buffer = append(buffer, '0')
	}

	buffer = strconv.AppendInt(buffer, int64(hour), 10)

	buffer = append(buffer, ':')
	if minute < 10 {
		buffer = append(buffer, '0')
	}

	buffer = strconv.AppendInt(buffer, int64(minute), 10)

	buffer = append(buffer, ':')
	if sec < 10 {
		buffer = append(buffer, '0')
	}

	buffer = strconv.AppendInt(buffer, int64(sec), 10)

	buffer = append(buffer, '.')
	if milli < 100 {
		buffer = append(buffer, '0')
	}

	if milli < 10 {
		buffer = append(buffer, '0')
	}

	buffer = strconv.AppendInt(buffer, int64(milli), 10)

	next.length = next.dateLen + len(buffer)
	timeCacheYYYYMonth.active.Store(next)
}
