package handler

import (
	"time"

	"github.com/ArtemChadaev/GoCreateHistory/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	service service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: *service,
	}
}

func (h *Handler) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(h.loggingMiddleware)

	r.Route("/history", func(r chi.Router) {
		r.Use(h.auth)
		r.Post("/", h.createHistory)
		//r.Get("/", getAll)

		r.Route("/{id}", func(r chi.Router) {
			//r.Post("/")
		})
	})

	return r
}

//TODO: slog для ошибок и информации
//TODO: доки сделать хоть какие то, а мб и не надо, самодокументирующий код типа
