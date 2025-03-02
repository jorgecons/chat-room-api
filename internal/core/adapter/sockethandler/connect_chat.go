package sockethandler

import (
	"context"
	"sync"

	"chat-room-api/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type (
	ConnectChatFeature interface {
		ConnectChat(context.Context, string) ([]domain.Message, error)
	}

	connectChat struct {
		feature   ConnectChatFeature
		mutex     sync.Mutex
		jwtSecret []byte
	}
)

func NewConnectChat(f ConnectChatFeature, jwtSecret []byte) HandlerFunc {
	handler := connectChat{
		feature:   f,
		mutex:     sync.Mutex{},
		jwtSecret: jwtSecret,
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
		_ = socketResponseAsJson(conn, BuildErrorMessage(room, err.Error()).Message)
		return err
	}

	claims, err := validateJWT(c.jwtSecret, tokenString)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithField("room", room).Error(InvalidTokenMsg)
		err = BuildErrorMessage(room, InvalidTokenMsg)
		_ = socketResponseAsJson(conn, BuildErrorMessage(room, err.Error()).Message)
		return err
	}
	ctx.Set(UserContextKey, claims["username"].(string))

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

	msgs, err := c.feature.ConnectChat(ctx, room)
	if err != nil {
		_ = socketResponseAsJson(conn, BuildErrorMessage(room, err.Error()).Message)
		return err
	}

	for _, m := range msgs {
		_ = socketResponseAsJson(conn, BuildSocketMessage(m))
	}
	return nil
}

func HandleMessages(mutex *sync.Mutex, broadcast chan Message) {
	for {
		msg := <-broadcast
		for client := range Chatrooms[msg.Room] {
			if err := socketResponseAsJson(client, msg); err != nil {
				mutex.Lock()
				delete(Chatrooms[msg.Room], client)
				mutex.Unlock()
			}
		}
	}
}
