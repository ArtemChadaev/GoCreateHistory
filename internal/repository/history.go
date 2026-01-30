package repository

import (
	"context"

	"github.com/ArtemChadaev/GoCreateHistory/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type historyRepository struct {
	db *sqlx.DB
}

func NewHistoryRepository(db *sqlx.DB) domain.HistoryRepository {
	return &historyRepository{db: db}
}

func (r *historyRepository) Create(ctx context.Context, h *domain.History) error {
	query := `
	INSERT INTO history (uuid, book_title, user_request, status, error, created_at, update_at, chapters) 
	VALUES (:uuid, :book_title, :user_request, :status, :error, :created_at, :update_at, :chapters)
	`
	_, err := r.db.NamedExecContext(ctx, query, h)
	if err != nil {
		return err
	}
	return nil
}

func (r *historyRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.History, error) {
	var history domain.History
	query := `SELECT * FROM history WHERE uuid = $1`

	err := r.db.GetContext(ctx, &history, query, id)
	if err != nil {
		return nil, err
	}

	return &history, nil
}

func (r *historyRepository) GetByUserID(ctx context.Context, id int) (*[]domain.History, error) {
	var history []domain.History
	query := `SELECT * FROM history WHERE user_request->user_id = $1`
	err := r.db.SelectContext(ctx, &history, query, id)
	if err != nil {
		return nil, err
	}
	return &history, nil
}

func (r *historyRepository) Update(ctx context.Context, h *domain.History) error {
	query := `
	UPDATE history
	SET book_title = :book_title, status = :status, error = :error, chapters = :chapters
	WHERE uuid = :uuid
	`
	_, err := r.db.NamedExecContext(ctx, query, h)
	if err != nil {
		return err
	}
	return nil
}

func (r *historyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	UPDATE history 
	SET status = 'deleted', book_title = null, status = null, error = null, chapters = null
	WHERE uuid = $1
`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *historyRepository) Freeze(ctx context.Context, id uuid.UUID, frozen bool) error {
	var err error
	query := "UPDATE history SET status = $1 WHERE uuid = $2"

	if frozen {
		_, err = r.db.ExecContext(ctx, query, domain.StatusFrozen, id)
	} else {
		_, err = r.db.ExecContext(ctx, query, domain.StatusPending, id)
	}

	if err != nil {
		return err
	}
	return nil
}

func (r *historyRepository) CountActiveTasks(ctx context.Context, userID string) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM history WHERE user_request->user_id = $1"
	err := r.db.GetContext(ctx, &count, query, userID)
	if err != nil {
		return 0, err
	}
	return count, nil
}
