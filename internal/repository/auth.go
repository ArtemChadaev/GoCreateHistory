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

func (r *authRepository) CreateUser(ctx context.Context, email string, password string) error {
	query := `INSERT INTO user_history (email, password) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, email, password)
	return err
}
func (r *authRepository) GetUserId(ctx context.Context, email string, password string) (int, error) {
	query := `SELECT user_id FROM user_history WHERE email = $1 AND password = $2`

	var id int
	err := r.db.QueryRowContext(ctx, query, email, password).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}
