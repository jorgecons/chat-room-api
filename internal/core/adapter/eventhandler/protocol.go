package eventhandler

import (
	"errors"
	"time"

	"chat-room-api/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var InvalidRoomError = errors.New("invalid room")

type (
	Handler func(*gin.Context, *websocket.Conn, chan Message)
	Message struct {
		Room     string    `json:"room"`
		Username string    `json:"username"`
		Text     string    `json:"text"`
		Date     time.Time `json:"date"`
	}
)

func BuildMessage(message Message) domain.Message {
	t := message.Date
	if t.IsZero() {
		t = time.Now().UTC()
	}
	return domain.Message{
		Room:     message.Room,
		Username: message.Username,
		Text:     message.Text,
		Date:     t,
	}
}

func CreateErrorMessage(room, text string) Message {
	return Message{
		Room:     room,
		Username: "System",
		Text:     "Unknown command: " + text,
	}
}
