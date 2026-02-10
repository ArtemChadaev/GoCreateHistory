package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/ArtemChadaev/GoCreateHistory/internal/domain"
	"github.com/ArtemChadaev/GoCreateHistory/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
)

func (h *Handler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Начинаем отсчет времени
		start := time.Now()

		//Лог о начале запроса
		ctx := logger.WithFields(r.Context(),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path))

		logger.Info(ctx, "start request")

		// Оборачиваем ResponseWriter, чтобы узнать статус-код в конце
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Передаем управление дальше
		next.ServeHTTP(ww, r)

		// Лог о завершении запроса
		logger.Info(ctx, "request completed",
			slog.Int("status", ww.Status()),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

func (h *Handler) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var t domain.Token

		// 1. Декодируем JSON из тела запроса в структуру
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		userID, err := h.service.ParseToken(t.Token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Оборачиваем контекст: добавляем userID для бизнес-логики
		ctx := context.WithValue(r.Context(), "user_id", userID)

		// Контекст для логирования
		ctx = logger.WithField(ctx, "user_id", userID)

		// Передаем обновленный контекст дальше
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
