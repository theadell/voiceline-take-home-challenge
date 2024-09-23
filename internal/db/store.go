package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type SqlStore struct {
	Querier
	db *sql.DB
}

func NewSqlStore(db *sql.DB, q Querier) *SqlStore {
	return &SqlStore{db: db, Querier: q}
}

func (s *SqlStore) LoginWithProvider(ctx context.Context, provider string, providerUserID string, email string) (User, error) {

	arg := GetUserAndProviderInfoParams{
		Provider:       provider,
		ProviderUserID: providerUserID,
		Email:          email,
	}

	result, err := s.GetUserAndProviderInfo(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			var user User
			err = ExecTransaction(ctx, s.db, func(q *Queries) error {

				createUserParams := CreateUserParams{
					Email: email,
				}
				user, err = s.CreateUser(ctx, createUserParams)
				if err != nil {
					return fmt.Errorf("failed to create user: %w", err)
				}

				createOAuthParams := CreateUserProviderParams{
					UserID:         user.ID,
					Provider:       provider,
					ProviderUserID: providerUserID,
				}
				_, err = s.CreateUserProvider(ctx, createOAuthParams)
				return err
			})

			if err != nil {
				return User{}, fmt.Errorf("failed to create user and link provider: %w", err)
			}

			return user, nil
		}

		return User{}, fmt.Errorf("error querying user and provider info: %w", err)
	}

	if result.UserID != 0 && result.OauthID.Valid {
		return User{
			ID:    result.UserID,
			Email: result.UserEmail,
		}, nil
	}

	if result.UserID != 0 && !result.OauthID.Valid {

		err = ExecTransaction(ctx, s.db, func(q *Queries) error {
			createOAuthParams := CreateUserProviderParams{
				UserID:         result.UserID,
				Provider:       provider,
				ProviderUserID: providerUserID,
			}
			_, err := s.CreateUserProvider(ctx, createOAuthParams)
			return err
		})

		if err != nil {
			return User{}, fmt.Errorf("failed to link provider: %w", err)
		}

		return User{
			ID:    result.UserID,
			Email: result.UserEmail,
		}, nil
	}

	return User{}, fmt.Errorf("unexpected error")
}
