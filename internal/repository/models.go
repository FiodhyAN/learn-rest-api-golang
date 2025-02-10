// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                              uuid.UUID
	Name                            string
	Username                        string
	Email                           string
	Password                        string
	Role                            string
	EmailVerified                   bool
	EmailVerificationToken          sql.NullString
	EmailVerificationTokenExpiresAt sql.NullTime
	DeletedAt                       sql.NullTime
	CreatedAt                       time.Time
	UpdatedAt                       sql.NullTime
}
