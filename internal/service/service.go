package service

import (
	"github.com/ArtemChadaev/GoCreateHistory/internal/domain"
	"github.com/ArtemChadaev/GoCreateHistory/internal/repository"
)

type Service struct {
	domain.HistoryService
}

func NewService(repos *repository.Repository) *Service {
	historyService := NewHistoryService(repos.History)
	return &Service{
		HistoryService: historyService,
	}
}
