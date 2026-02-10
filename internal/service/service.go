package service

import (
	"github.com/ArtemChadaev/GoCreateHistory/internal/domain"
	"github.com/ArtemChadaev/GoCreateHistory/internal/repository"
)

type Service struct {
	domain.HistoryService
	domain.AuthService
}

func NewService(repos *repository.Repository) *Service {
	historyService := NewHistoryService(repos.History)
	authService := NewAuthService(repos.Auth)
	return &Service{
		HistoryService: historyService,
		AuthService:    authService,
	}
}
