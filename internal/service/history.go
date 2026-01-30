package service

import (
	"context"

	"github.com/ArtemChadaev/GoCreateHistory/internal/domain"
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
		return
	}
	// TODO: Сделать в канал запись про ид чтобы потихоньку начинал делать, скорее всего в отдельной горутине
	return hID, nil
}

func (s *historyService) Get(ctx context.Context, id uuid.UUID) (*domain.History, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *historyService) Freeze(ctx context.Context, id uuid.UUID, frozen bool) error {
	// TODO: Сделать чтобы убирался из очереди для создания
	return s.repo.Freeze(ctx, id, frozen)
}

func (s *historyService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
