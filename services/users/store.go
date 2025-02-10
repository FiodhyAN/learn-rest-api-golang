package users

import (
	"context"
	"database/sql"
	"time"

	"github.com/FiodhyAN/learn-rest-api-golang/internal/repository"
	"github.com/google/uuid"
)

type Store struct {
	queries *repository.Queries
	db      *sql.DB
}

func NewUserStore(db *sql.DB) *Store {
	return &Store{
		queries: repository.New(db),
		db:      db,
	}
}

func (s *Store) GetUser(ctx context.Context, username string) (*repository.User, error) {
	userRow, err := s.queries.GetUser(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user := repository.User{
		ID:                              userRow.ID,
		Username:                        userRow.Username,
		Email:                           userRow.Email,
		Password:                        userRow.Password,
		EmailVerified:                   userRow.EmailVerified,
		EmailVerificationTokenExpiresAt: userRow.EmailVerificationTokenExpiresAt,
	}

	return &user, nil
}

func (s *Store) CreateUser(ctx context.Context, user repository.CreateUserParams) (*repository.User, error) {
	createdUserRow, err := s.queries.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	createdUser := repository.User{
		ID:        createdUserRow.ID,
		Name:      createdUserRow.Name,
		Username:  createdUserRow.Username,
		Email:     createdUserRow.Email,
		Password:  createdUserRow.Password,
		Role:      createdUserRow.Role,
		CreatedAt: createdUserRow.CreatedAt,
	}

	return &createdUser, nil
}

func (s *Store) UpdateUserVerificationExpired(ctx context.Context, userID uuid.UUID, expired time.Time, token string) error {
	err := s.queries.UpdateUserVerificationExpired(ctx, repository.UpdateUserVerificationExpiredParams{
		EmailVerificationToken:          sql.NullString{String: token, Valid: true},
		EmailVerificationTokenExpiresAt: sql.NullTime{Time: expired, Valid: true},
		ID:                              userID,
	})
	return err
}

func (s *Store) GetUserById(ctx context.Context, userId uuid.UUID) (*repository.User, error) {
	userRow, err := s.queries.GetUserById(ctx, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user := repository.User{
		ID:                              userRow.ID,
		EmailVerified:                   userRow.EmailVerified,
		EmailVerificationToken:          userRow.EmailVerificationToken,
		EmailVerificationTokenExpiresAt: userRow.EmailVerificationTokenExpiresAt,
	}

	return &user, nil
}

func (s *Store) VerifyEmail(ctx context.Context, userID uuid.UUID) error {
	return s.queries.VerifyEmail(ctx, userID)
}
