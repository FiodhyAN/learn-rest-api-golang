package types

import (
	"context"
	"database/sql"
	"time"

	"github.com/FiodhyAN/learn-rest-api-golang/internal/repository"
	"github.com/google/uuid"
)

type UserStore interface {
	GetUser(ctx context.Context, email string) (*repository.User, error)
	CreateUser(ctx context.Context, user repository.CreateUserParams) (*repository.User, error)
	UpdateUserVerificationExpired(ctx context.Context, userId uuid.UUID, expiresAt time.Time, token string) error
	GetUserById(ctx context.Context, id uuid.UUID) (*repository.User, error)
	VerifyEmail(ctx context.Context, userId uuid.UUID) error
}

type RegisterPayload struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type LoginUserPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type VerifyEmailPayload struct {
	UserId string `json:"userId" validate:"required"`
	Token  string `json:"token" validate:"required"`
}

type User struct {
	ID                         string
	Name                       string
	Username                   string
	Email                      string
	Password                   string
	Role                       string
	EmailVerified              bool
	EmailVerificationToken     sql.NullString
	EmailVerificationExpiresAt sql.NullTime
	CreatedAt                  time.Time
}
