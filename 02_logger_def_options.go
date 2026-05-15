package arenalog

import (
	"context"

	"github.com/tudorhulban/arenalog/timestamp"
)

type Option func(*Logger)

func WithTimestampRFC3339UTC(ctx context.Context) Option {
	chReady := timestamp.StartRFC3339UTCCache(ctx)

	<-chReady

	return func(item *Logger) {
		item.fnTimestamp = timestamp.TimestampRFC3339UTC
	}
}

func WithTimestampRFC3339Bucharest(ctx context.Context) Option {
	chReady := timestamp.StartRFC3339BucharestCache(ctx)

	<-chReady

	return func(item *Logger) {
		item.fnTimestamp = timestamp.TimestampRFC3339Bucharest
	}
}

func WithTimestampRFC3339CustomLocation(ctx context.Context, location string) Option {
	chReady := timestamp.StartRFC3339CustomLocationCache(ctx, location)

	<-chReady

	return func(item *Logger) {
		item.fnTimestamp = timestamp.TimestampRFC3339CustomLocation
	}
}

func WithTimestampStandardLocal(ctx context.Context) Option {
	chReady := timestamp.StartStandardCache(ctx)

	<-chReady

	return func(item *Logger) {
		item.fnTimestamp = timestamp.TimestampStandard
	}
}

func WithTimestampYYYYMonthLocal(ctx context.Context) Option {
	chReady := timestamp.StartYYYYMonthCache(ctx)

	<-chReady

	return func(item *Logger) {
		item.fnTimestamp = timestamp.TimestampYYYYMonth
	}
}
