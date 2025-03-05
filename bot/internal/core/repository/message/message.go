package message

import (
	"context"

	"bot/internal/core/domain"

	"github.com/jackc/pgconn"
)

const (
	sqlSaveMsg = "INSERT INTO messages (room, username, text, date) VALUES ($1, $2, $3, $4)"
)

type (
	Storage struct {
		client querier
	}

	querier interface {
		Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	}
)

func NewMessageStorage(dbClient querier) *Storage {
	return &Storage{
		client: dbClient,
	}
}

func (r *Storage) Save(ctx context.Context, msg domain.Message) error {
	_, err := r.client.Exec(ctx,
		sqlSaveMsg,
		msg.Room,
		msg.Username,
		msg.Text,
		msg.Date,
	)
	if err != nil {
		return domain.WrapError(domain.ErrSavingMessage, err)
	}
	return nil
}
