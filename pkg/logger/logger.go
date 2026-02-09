package logger

import (
	"context"
	"errors"
	"log/slog"
	"os"
)

func InitLogger(app string) {
	var handler slog.Handler
	opts := &slog.HandlerOptions{AddSource: true}

	if app == "production" {
		// В продакшене — строгий JSON для машин (Loki/ELK)
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		// В разработке — красивый текст для людей
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}

// ctxKey — приватный тип для ключа контекста, чтобы избежать коллизий
type ctxKey struct{}

// LogCtx — контейнер для атрибутов лога
type LogCtx struct {
	Fields []slog.Attr
}

// AppError — обертка над ошибкой, которая позволяет прикрепить к ней метаданные
type AppError struct {
	Cause    error
	Metadata LogCtx
}

func (e *AppError) Error() string { return e.Cause.Error() }
func (e *AppError) Unwrap() error { return e.Cause }

// --- Работа с Контекстом ---

// FromContext извлекает LogCtx из контекста. Всегда возвращает объект, а не nil.
func FromContext(ctx context.Context) LogCtx {
	if v, ok := ctx.Value(ctxKey{}).(LogCtx); ok {
		return v
	}
	return LogCtx{Fields: make([]slog.Attr, 0)}
}

// WithField добавляет один атрибут в контекст. Создает копию данных (Thread-safe).
func WithField(ctx context.Context, key string, val any) context.Context {
	return WithFields(ctx, slog.Any(key, val))
}

// WithFields добавляет пачку атрибутов в контекст.
func WithFields(ctx context.Context, attrs ...slog.Attr) context.Context {
	prev := FromContext(ctx)

	newFields := make([]slog.Attr, len(prev.Fields), len(prev.Fields)+len(attrs))
	copy(newFields, prev.Fields)
	newFields = append(newFields, attrs...)

	return context.WithValue(ctx, ctxKey{}, LogCtx{Fields: newFields})
}

// --- Работа с Ошибками ---

// WithLog прикрепляет атрибуты к ошибке. Если ошибка уже AppError, дополняет её.
func WithLog(err error, attrs ...slog.Attr) error {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		appErr.Metadata.Fields = append(appErr.Metadata.Fields, attrs...)
		return appErr
	}

	return &AppError{
		Cause:    err,
		Metadata: LogCtx{Fields: attrs},
	}
}

// WrapWithContext извлекает все поля логов из контекста и прикрепляет их к ошибке.
func WrapWithContext(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	// 1. Достаем поля из контекста
	ctxData := FromContext(ctx)
	if len(ctxData.Fields) == 0 {
		return err // В контексте ничего нет, возвращаем ошибку как есть
	}

	// 2. Прикрепляем эти поля к ошибке через нашу существующую функцию
	return WithLog(err, ctxData.Fields...)
}

// --- Функции логирования ---

// Info логирует сообщение, автоматически подмешивая поля из контекста.
func Info(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelInfo, msg, nil, args...)
}

// Warn логирует предупреждение с полями из контекста.
func Warn(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelWarn, msg, nil, args...)
}

// Error логирует ошибку, объединяя поля из контекста и самой ошибки.
func Error(ctx context.Context, msg string, err error, args ...any) {
	log(ctx, slog.LevelError, msg, err, args...)
}

// log — внутренняя функция для сборки всех полей воедино.
func log(ctx context.Context, level slog.Level, msg string, err error, args ...any) {
	// 1. Поля из контекста
	ctxData := FromContext(ctx)

	// Оцениваем примерный объем памяти
	capacity := len(ctxData.Fields) + len(args)
	if err != nil {
		capacity += 1
		var appErr *AppError
		if errors.As(err, &appErr) {
			capacity += len(appErr.Metadata.Fields)
		}
	}

	finalArgs := make([]any, 0, capacity)
	finalArgs = append(finalArgs, args...)

	// 2. Поля из контекста
	for _, attr := range ctxData.Fields {
		finalArgs = append(finalArgs, attr)
	}

	// 3. Поля из ошибки
	if err != nil {
		finalArgs = append(finalArgs, slog.Any("error", err))
		var appErr *AppError
		if errors.As(err, &appErr) {
			for _, attr := range appErr.Metadata.Fields {
				finalArgs = append(finalArgs, attr)
			}
		}
	}

	slog.Log(ctx, level, msg, finalArgs...)
}
