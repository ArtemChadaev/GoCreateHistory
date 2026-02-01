package logger

import (
	"context"
	"log/slog"
)

type ctxKey string

const loggerKey ctxKey = "logger"

// WithValue добавляет логгер с предустановленными полями в контекст
func WithValue(ctx context.Context, key string, value any) context.Context {
	l := slog.Default().With(slog.Any(key, value))
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext достает логгер из контекста или возвращает дефолтный
func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
