package repository

import (
	"context"

	"github.com/ArtemChadaev/GoCreateHistory/internal/domain"
	"github.com/jmoiron/sqlx"
)

type authRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) domain.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(ctx context.Context, email string, password string) (int, error) {
	query := `INSERT INTO user_history (email, password) VALUES ($1, $2)`
	res, err := r.db.ExecContext(ctx, query, email, password)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
func (r *authRepository) GetUserId(ctx context.Context, email string, password string) (int, error) {
	query := `SELECT id FROM user_history WHERE email = $1 AND password = $2`
	row := r.db.QueryRowContext(ctx, query, email, password)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
