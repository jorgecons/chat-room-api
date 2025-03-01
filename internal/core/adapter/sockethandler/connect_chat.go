package sockethandler

import (
	"context"
	"log"
	"sync"

	"chat-room-api/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type (
	GetLastMessages interface {
		GetLastMessages(context.Context, string) ([]domain.Message, error)
	}
	connectChat struct {
		getLastMessages GetLastMessages
		mutex           sync.Mutex
		jwtSecret       []byte
	}
)

func NewConnectChat(lm GetLastMessages, jwtSecret []byte) HandlerFunc {
	handler := connectChat{
		getLastMessages: lm,
		mutex:           sync.Mutex{},
		jwtSecret:       jwtSecret,
	}
	handlerFunc := func(c *gin.Context, conn *websocket.Conn, broadcast chan Message) error {
		return handler.handle(c, conn, broadcast)
	}
	return handlerFunc
}

func (c *connectChat) handle(ctx *gin.Context, conn *websocket.Conn, broadcast chan Message) error {
	var (
		room        = ctx.Param("room")
		tokenString = ctx.GetHeader(TokenHeader)
	)
	if tokenString == "" {
		logrus.WithContext(ctx).WithField("room", room).Error("Missing Token")
		err := BuildErrorMessage(room, InvalidTokenMsg)
		_ = conn.WriteJSON(err)
		return err
	}

	claims, err := validateJWT(c.jwtSecret, tokenString)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithField("room", room).Error(InvalidTokenMsg)
		err = BuildErrorMessage(room, InvalidTokenMsg)
		_ = conn.WriteJSON(err)
		return err
	}
	context.WithValue(ctx, UserContextKey, claims["username"].(string))

	if room == "" {
		logrus.WithContext(ctx).Error("Error missing room")
		return InvalidRoomError
	}
	c.mutex.Lock()
	if Chatrooms[room] == nil {
		Chatrooms[room] = make(map[*websocket.Conn]bool)
	}
	Chatrooms[room][conn] = true
	c.mutex.Unlock()
	// TODO: refactor
	lm, _ := c.getLastMessages.GetLastMessages(ctx, room)
	for i := len(lm) - 1; i >= 0; i = i - 1 {
		broadcast <- BuildSocketMessage(lm[i])
	}
	return nil
}

func HandleMessages(broadcast chan Message) {
	for {
		msg := <-broadcast
		for client := range Chatrooms[msg.Room] {
			err := client.WriteJSON(msg)
			if err != nil {
				// TODO improve
				log.Println("Write error:", err)
				_ = client.Close()
				delete(Chatrooms[msg.Room], client)
			}
		}
	}
}
