package handler

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Router() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Heartbeat("/ping"))

	r.Route("/history", func(r chi.Router) {
		r.Get("/", getAll)

		r.Route("/{id}", func(r chi.Router) {
			r.Post("/")
		})
	})
}
