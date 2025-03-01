package eventhandler

import (
	"errors"
	"strings"

	"chat-room-api/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var InvalidRoomError = errors.New("invalid room")

type (
	Handler func(*gin.Context, *websocket.Conn, chan Message)
	Message struct {
		Room     string `json:"room"`     // Chatroom name
		Username string `json:"username"` // Sender
		Text     string `json:"text"`     // Message content
	}
)

func BuildMessage(message Message) domain.Message {
	return domain.Message{
		Room:     message.Room,
		Username: message.Username,
		Text:     message.Text,
	}
}

func GetStockName(text string) string {
	return strings.TrimPrefix(text, "/stock=")
}

func ValidateMessage(room string, message Message) error {
	if room == message.Room {
		return nil
	}
	return InvalidRoomError
}

func CreateErrorMessage(room, text string) Message {
	return Message{
		Room:     room,
		Username: "System",
		Text:     "Unknown command: " + text,
	}
}
