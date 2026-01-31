package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func (h *Handler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Начинаем отсчет времени
		start := time.Now()

		// Достаем IP. Если ты за прокси (Nginx/Docker),
		// лучше использовать r.Header.Get("X-Forwarded-For")
		ip := r.RemoteAddr

		//TODO: Когда разберусь подумаю менять ли slog тут
		//Лог о начале запроса
		slog.Info("request started",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("ip", ip),
		)

		// Оборачиваем ResponseWriter, чтобы узнать статус-код в конце
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Передаем управление дальше
		next.ServeHTTP(ww, r)

		// Лог о завершении запроса
		slog.Info("request completed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", ww.Status()),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

func (h *Handler) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Сделать проверку авторизации и получения userID

		userID := 123

		// Оборачиваем контекст: добавляем userID для бизнес-логики и для логов
		ctx := context.WithValue(r.Context(), "user_id", userID)
		//TODO: Сделать контекст с user_id после создания slog
		//ctx = logger.WithValue(ctx, "user_id", userID) // Наш хелпер для slog

		// Передаем обновленный контекст дальше
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
