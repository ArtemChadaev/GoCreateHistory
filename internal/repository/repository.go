package repository

import (
	"github.com/ArtemChadaev/GoCreateHistory/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	History domain.HistoryRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		History: NewHistoryRepository(db),
	}
}
