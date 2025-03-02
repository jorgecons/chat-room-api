package sockethandler

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"

	"chat-room-api/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gorilla/websocket"
)

var (
	InvalidRoomError = errors.New("invalid room")
	InvalidUserError = errors.New("invalid user")
	ChatroomError    = errors.New("error connecting to chat room")
	Chatrooms        = make(map[string]map[*websocket.Conn]bool)
)

const (
	UserContextKey  = "user"
	SystemUser      = "System"
	TokenHeader     = "X-Access-Token"
	InvalidTokenMsg = "Invalid Token"
	textFormat      = "%s@room_%s: %s at %s"
)

type (
	HandlerFunc func(*gin.Context, *websocket.Conn, chan Message) error
	Message     struct {
		Room     string    `json:"room"`
		Username string    `json:"username"`
		Text     string    `json:"text"`
		Date     time.Time `json:"date"`
	}

	ErrorMessage struct {
		Message
	}
)

func (e ErrorMessage) Error() string {
	return e.Text
}

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

func BuildSocketMessage(message domain.Message) Message {
	return Message{
		Room:     message.Room,
		Username: message.Username,
		Text:     message.Text,
		Date:     message.Date,
	}
}

func BuildErrorMessage(room, text string) ErrorMessage {
	return ErrorMessage{
		Message{
			Room:     room,
			Username: SystemUser,
			Text:     text,
			Date:     time.Now().UTC(),
		},
	}
}

func ValidateRoom(room string, message Message) error {
	if room == message.Room {
		return nil
	}
	return InvalidRoomError
}

func ValidateUsername(username string, message Message) error {
	if username == message.Username {
		return nil
	}
	return InvalidUserError
}

func validateJWT(secret []byte, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func socketResponseAsJson(client *websocket.Conn, msg Message) error {
	err := client.WriteJSON(msg)
	if err != nil {
		logrus.
			WithError(err).
			WithField("username", msg.Username).
			Error("error writing message")
		_ = client.Close()
	}
	return err
}

func socketResponseAsText(client *websocket.Conn, msg Message) error {
	t := fmt.Sprintf(textFormat, msg.Username, msg.Room, msg.Text, msg.Date)
	err := client.WriteMessage(websocket.TextMessage, []byte(t))
	if err != nil {
		logrus.
			WithError(err).
			WithField("username", msg.Username).
			Error("error writing message")
		_ = client.Close()
	}
	return err
}
