package types

import (
	"database/sql"
	"time"
)

type UserStore interface {
	GetUser(email string) (*User, error)
	CreateUser(user User) (*User, error)
	UpdateUserVerificationExpired(*User, time.Time, string) error
	GetUserById(id string) (*User, error)
	VerifyEmail(*User) error
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
