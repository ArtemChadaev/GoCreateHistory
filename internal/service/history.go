package service

import (
	"context"
	"log/slog"

	"github.com/ArtemChadaev/GoCreateHistory/internal/domain"
	"github.com/ArtemChadaev/GoCreateHistory/pkg/logger"
	"github.com/google/uuid"
)

type historyService struct {
	repo domain.HistoryRepository
}

func NewHistoryService(repo domain.HistoryRepository) domain.HistoryService {
	return &historyService{
		repo: repo,
	}
}

func (s *historyService) Create(ctx context.Context, req domain.UserRequest) (hID uuid.UUID, err error) {
	hID, err = uuid.NewUUID()
	if err != nil {
		return
	}

	history := &domain.History{
		UUID:        hID,
		UserRequest: req,
		Status:      domain.StatusPending,
	}
	err = s.repo.Create(ctx, history)
	if err != nil {
		return hID, logger.WithLog(err,
			slog.String("op", "history.Create"),
			slog.Any("history", history),
		)
	}

	logger.Info(ctx, "create history success", slog.String("history_ID", hID.String()))
	// TODO: Сделать в канал запись про ид чтобы потихоньку начинал делать, скорее всего в отдельной горутине
	return hID, nil
}

func (s *historyService) Get(ctx context.Context, id uuid.UUID) (*domain.History, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *historyService) Freeze(ctx context.Context, id uuid.UUID, frozen bool) error {
	// TODO: Сделать чтобы убирался из очереди для создания
	if err := s.repo.Freeze(ctx, id, frozen); err != nil {
		return logger.WithLog(err,
			slog.String("op", "history.Freeze"),
			slog.String("history_ID", id.String()),
			slog.Bool("frozen", frozen),
		)
	}
	logger.Info(ctx, "freeze or unfreeze history success",
		slog.String("history_ID", id.String()),
		slog.Bool("frozen", frozen),
	)
	return nil
}

func (s *historyService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return logger.WithLog(err,
			slog.String("op", "history.Delete"),
			slog.String("history_ID", id.String()),
		)
	}
	logger.Info(ctx, "delete history success", slog.String("history_ID", id.String()))
	return nil
}
