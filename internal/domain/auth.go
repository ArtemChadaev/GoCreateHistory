package domain

import "context"

type Auth struct {
	UserID   int    `json:"-"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

type AuthService interface {
	CreateUser(ctx context.Context, email string, password string) error
	ParseToken(accessToken string) (int, error)
	GenerateToken(ctx context.Context, email, password string) (string, error)
}

type AuthRepository interface {
	CreateUser(ctx context.Context, email string, password string) error
	GetUserId(ctx context.Context, email string, password string) (int, error)
}
