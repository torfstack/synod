package logging

import (
	"context"
	"fmt"
	"log/slog"
)

func Debugf(ctx context.Context, msg string, args ...interface{}) {
	fields := logAttributeFields(ctx)
	m := fmt.Sprintf(msg, args...)
	slog.LogAttrs(ctx, slog.LevelDebug, m, fields...)
}

func Infof(ctx context.Context, msg string, args ...interface{}) {
	fields := logAttributeFields(ctx)
	m := fmt.Sprintf(msg, args...)
	slog.LogAttrs(ctx, slog.LevelInfo, m, fields...)
}

func Warnf(ctx context.Context, msg string, args ...interface{}) {
	fields := logAttributeFields(ctx)
	m := fmt.Sprintf(msg, args...)
	slog.LogAttrs(ctx, slog.LevelWarn, m, fields...)
}

func Errorf(ctx context.Context, msg string, args ...interface{}) {
	fields := logAttributeFields(ctx)
	m := fmt.Sprintf(msg, args...)
	slog.LogAttrs(ctx, slog.LevelError, m, fields...)
}

func Fatalf(ctx context.Context, msg string, args ...interface{}) {
	fields := logAttributeFields(ctx)
	m := fmt.Sprintf(msg, args...)
	panic(fmt.Sprint(m, fields))
}

func logAttributeFields(ctx context.Context) []slog.Attr {
	var fields []slog.Attr
	for _, attr := range allStaticLogAttributes {
		if value, ok := ctx.Value(attr).(int); ok {
			fields = append(fields, slog.Int(string(attr), value))
		}
	}
	return fields
}

type StaticLogAttribute string

const (
	LogAttributeUserId StaticLogAttribute = "user_id"
)

var allStaticLogAttributes = []StaticLogAttribute{
	LogAttributeUserId,
}

func WithStaticLogAttribute(ctx context.Context, key StaticLogAttribute, value any) context.Context {
	return context.WithValue(ctx, key, value)
}

func WithLogAttributeUserId(ctx context.Context, userId int) context.Context {
	return WithStaticLogAttribute(ctx, LogAttributeUserId, userId)
}

func SetLogLevel(level slog.Level) {
	slog.SetLogLoggerLevel(level)
}
