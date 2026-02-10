package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

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
		// 1. Получаем значение заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// 2. Заголовок должен быть в формате "Bearer <token>"
		// Проверяем префикс и отрезаем его
		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			http.Error(w, "Invalid authorization format. Use 'Bearer <token>'", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, prefix)

		// 3. Передаем чистый токен в сервис для парсинга
		userID, err := h.service.ParseToken(tokenString)
		if err != nil {
			logger.Error(r.Context(), "Invalid token", err)
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
