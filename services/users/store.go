package users

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/FiodhyAN/learn-rest-api-golang/types"
)

type Store struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUser(username string) (*types.User, error) {
	user := new(types.User)

	if err := s.db.QueryRow("SELECT * FROM users WHERE email = $1 OR username = $1", username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("query database error")
	}

	return user, nil
}

func (s *Store) CreateUser(user types.User) (*types.User, error) {
	query := `INSERT INTO users (name, username, email, password) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING _id, name, username, email, password, created_at`

	created_user := new(types.User)

	err := s.db.QueryRow(query, user.Name, user.Username, user.Email, user.Password).
		Scan(&created_user.ID, &created_user.Name, &created_user.Username, &created_user.Email, &created_user.Password, &created_user.CreatedAt)

	if err != nil {
		return created_user, err
	}

	return created_user, nil
}

func (s *Store) UpdateUserVerificationExpired(user *types.User, expired time.Time, token string) error {
	query := `UPDATE users SET email_verification_token = $1, email_verification_token_expires_at = $2 WHERE _id = $3`

	_, err := s.db.Exec(query, token, expired, user.ID)
	if err != nil {
		return err
	}

	return nil
}
