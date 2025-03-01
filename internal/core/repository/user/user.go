package user

import (
	"context"
	"errors"

	"chat-room-api/internal/core/domain"

	"github.com/jackc/pgx/v5"
)

const (
	sqlSaveUser = "INSERT INTO users (username, password_hash) VALUES ($1, $2)"
	sqlGetUser  = "SELECT username, password_hash FROM users WHERE username=$1"
)

var zeroUser = domain.User{}

type (
	Storage struct {
		client *pgx.Conn
	}
)

func NewUserStorage(dbClient *pgx.Conn) *Storage {
	return &Storage{
		client: dbClient,
	}
}

// It can be improved by clasifficating the client error in different types
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
