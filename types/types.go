package types

import "time"

type UserStore interface {
	GetUser(email string) (*User, error)
	CreateUser(user User) (*User, error)
	UpdateUserVerificationExpired(*User, time.Time, string) error
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

type User struct {
	ID                         string
	Name                       string
	Username                   string
	Email                      string
	Password                   string
	EmailVerificationExpiresAt time.Time
	CreatedAt                  time.Time
}

// type EmailRequest struct {
// 	to      []string
// 	subject string
// 	body    string
// }
