package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/ArtemChadaev/GoCreateHistory/internal/domain"
	"github.com/ArtemChadaev/GoCreateHistory/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

const (
	signingKey     = "awsg8s#@4Sf86DS#$2dF"
	accessTokenTTL = time.Hour * 24 * 365
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

type authService struct {
	repo domain.AuthRepository // Используем интерфейс из domain
}

func NewAuthService(repo domain.AuthRepository) domain.AuthService { return &authService{repo: repo} }

func (s *authService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, logger.WithLog(err, slog.String("op", "auth.ParseToken"), slog.String("token", accessToken))
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserId, nil
}

func (s *authService) GenerateToken(ctx context.Context, email, password string) (string, error) {
	userId, err := s.repo.GetUserId(ctx, email, password)

	if err != nil {
		return "", logger.WithLog(err, slog.String("op", "repo.GetUserId"))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId: userId,
	})

	accessToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", logger.WithLog(err, slog.String("op", "token.SignedString"))
	}
	return accessToken, nil
}

func (s *authService) CreateUser(ctx context.Context, email string, password string) (int, error) {
	userId, err := s.repo.GetUserId(ctx, email, password)
	if err != nil {
		return 0, logger.WithLog(err, slog.String("op", "repo.GetUserId"))
	}
	return userId, nil
}
