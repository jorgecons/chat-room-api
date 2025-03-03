package user

import (
	"context"
	"errors"

	"chat-room-api/internal/core/domain"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

const (
	sqlSaveUser = "INSERT INTO users (username, password_hash) VALUES ($1, $2)"
	sqlGetUser  = "SELECT username, password_hash FROM users WHERE username=$1"
)

var zeroUser = domain.User{}

type (
	Storage struct {
		client querier
	}

	querier interface {
		Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
		QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	}
)

func NewUserStorage(dbClient querier) *Storage {
	return &Storage{
		client: dbClient,
	}
}

// it can be improved by classifying the client errors in different types
func (s *Storage) Save(ctx context.Context, user domain.User) error {
	_, err := s.client.Exec(ctx, sqlSaveUser, user.Username, user.Password)
	if err != nil {
		return domain.WrapError(domain.ErrSavingUser, err)
	}
	return nil
}

func (s *Storage) Get(ctx context.Context, username string) (domain.User, error) {
	user := domain.User{}
	err := s.client.QueryRow(ctx, sqlGetUser, username).Scan(&user.Username, &user.Password)
	if errors.Is(err, pgx.ErrNoRows) {
		return zeroUser, domain.ErrInvalidUser
	}
	if err != nil {
		return zeroUser, domain.WrapError(domain.ErrGettingUser, err)
	}
	return user, nil
}
