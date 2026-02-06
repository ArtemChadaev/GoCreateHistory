package logger

import (
	"context"
	"log/slog"
)

type ctxKey string

const fieldsKey ctxKey = "log_fields"

// WithValue добавляет поле в список полей контекста
func WithValue(ctx context.Context, key string, value any) context.Context {
	fields, _ := ctx.Value(fieldsKey).([]slog.Attr)
	newFields := append(fields, slog.Any(key, value))
	return context.WithValue(ctx, fieldsKey, newFields)
}

// ErrorWithContext создает лог с ошибкой, вытягивая всё из контекста
// TODO: Переделать чтобы бли возвращал логгер а не это
func LogError(ctx context.Context, msg string, err error) {
	fields, _ := ctx.Value(fieldsKey).([]slog.Attr)

	attrs := make([]any, len(fields)+1)
	for i, v := range fields {
		attrs[i] = v
	}
	attrs[len(fields)] = slog.Any("error", err.Error())

	slog.ErrorContext(ctx, msg, attrs...)
}

// TODO: Сделать лог чтобы сохранялся в кастомную ошибку
type SlogError struct {
	Msg    string
	Err    error
	Fields map[string]any
}

//func (e *SlogError) Error() string {
//	return e.Err.Error()
//}
//
//func (e *SlogError) Unwrap() error {
//
//}
//
//func (e *SlogError) WrapError(ctx context.Context, msg string, err error) SlogError {
//
//}
