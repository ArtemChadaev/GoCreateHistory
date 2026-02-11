package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ArtemChadaev/GoCreateHistory/internal/config"
	"github.com/ArtemChadaev/GoCreateHistory/internal/domain"
	"github.com/ArtemChadaev/GoCreateHistory/internal/handler"
	"github.com/ArtemChadaev/GoCreateHistory/internal/repository"
	"github.com/ArtemChadaev/GoCreateHistory/internal/service"
	"github.com/ArtemChadaev/GoCreateHistory/pkg/logger"
)

func main() {
	// Подтягиваю конфигурацию
	cfg, err := config.Load()

	if err != nil {
		slog.Error("Config Die", "error", err)
		os.Exit(1)
	}

	// Создаю контекст для выхода из программы
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Логгер смотрю какой ставить по ситуации
	logger.InitLogger(cfg.AppEnv)

	// Инициализация ресурсов БД
	db, err := repository.NewPostgresDB(repository.PostgresConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		Username: cfg.DBUser,
		Database: cfg.DBName,
		Password: cfg.DBPassword,
		SSLMode:  "disable",
	})
	if err != nil {
		slog.Error("DB Die", "error", err)
		os.Exit(1)
	}

	// Вызов всех слоёв
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	// Запуск сервера http
	router := handlers.Router()

	srv := new(domain.Server)
	go func() {
		if err := srv.Run(cfg.Port, router); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server Die", "error", err)
			os.Exit(1)
		}
	}()

	// Выход из программы
	<-ctx.Done()

	//TODO выход из БД, и закрытие горутин всех и т.д.
}
