package timestamp

const nanosPerDay = int64(24 * 60 * 60 * 1e9)

// RFC3339 with millisecond precision, UTC timezone
const rfc3339MilliLayout = "2006-01-02T15:04:05.000Z"
