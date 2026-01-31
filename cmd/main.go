package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ArtemChadaev/GoCreateHistory/internal/config"
	"github.com/ArtemChadaev/GoCreateHistory/internal/domain"
	"github.com/ArtemChadaev/GoCreateHistory/internal/handler"
	"github.com/ArtemChadaev/GoCreateHistory/internal/repository"
	"github.com/ArtemChadaev/GoCreateHistory/internal/service"
)

func main() {
	// Сначала создаю slog для всех ошибок и комментариев

	// Потом подтягиваю конфигурацию
	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	// Создаю контекст для выхода из программы
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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
		log.Fatalf("failed to initialize db: %s", err.Error())
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
			log.Fatalf("Ошибка при запуске сервера: %s\n", err)
		}
	}()

	// Выход из программы
	<-ctx.Done()
}
