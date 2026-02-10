package repository

import (
	"github.com/ArtemChadaev/GoCreateHistory/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	History domain.HistoryRepository
	Auth    domain.AuthRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		History: NewHistoryRepository(db),
		Auth:    NewAuthRepository(db),
	}
}
