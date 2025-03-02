package message

import (
	"context"

	"chat-room-api/internal/core/domain"

	"github.com/jackc/pgx/v5"
)

const (
	sqlSaveMsg = "INSERT INTO messages (room, username, text, date) VALUES ($1, $2, $3, $4)"
	sqlSearch  = "SELECT room, username, text, date FROM messages WHERE room = $1 ORDER BY date DESC LIMIT 50"
)

type (
	Storage struct {
		client *pgx.Conn
	}
)

func NewMessageStorage(dbClient *pgx.Conn) *Storage {
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

func (r *Storage) GetLastMessages(ctx context.Context, room string) ([]domain.Message, error) {
	rows, err := r.client.Query(ctx, sqlSearch, room)
	if err != nil {
		return nil, domain.WrapError(domain.ErrGettingMessages, err)
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var msg message
		if err = rows.Scan(&msg.Room, &msg.Username, &msg.Text, &msg.Date); err != nil {
			return nil, domain.WrapError(domain.ErrGettingMessages, err)
		}

		messages = append(messages, domain.NewMessage(msg.Room, msg.Username, msg.Text, msg.Date))
	}

	return messages, nil
}
