package sockethandler

import (
	"chat-room-api/internal/core/domain"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type (
	ChatFeature interface {
		Chat(context.Context, domain.Message) error
	}

	chat struct {
		chatFeature ChatFeature
		jwtSecret   []byte
	}
)

func NewChat(cf ChatFeature, jwtSecret []byte) HandlerFunc {
	handler := chat{
		chatFeature: cf,
		jwtSecret:   jwtSecret,
	}
	handlerFunc := func(c *gin.Context, conn *websocket.Conn, broadcast chan Message) error {
		return handler.handle(c, conn, broadcast)
	}
	return handlerFunc
}

func (c *chat) handle(ctx *gin.Context, conn *websocket.Conn, broadcast chan Message) error {
	for {
		room := ctx.Param("room")
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			logrus.WithContext(ctx).WithError(err).WithField("room", room).Errorln("Error reading message")
			continue
		}
		username := ctx.Value(UserContextKey).(string)
		if err = ValidateUsername(username, msg); err != nil {
			logrus.WithContext(ctx).WithError(err).WithField("username", msg.Username).Errorln("Error validating user")
			_ = socketResponseAsJson(conn, BuildErrorMessage(room, err.Error()).Message)
			continue
		}
		if err = ValidateRoom(room, msg); err != nil {
			logrus.WithContext(ctx).WithError(err).WithField("room", room).Errorln("Error validating message")
			_ = socketResponseAsJson(conn, BuildErrorMessage(room, err.Error()).Message)
			continue
		}

		message := BuildMessage(msg)
		if err = c.chatFeature.Chat(ctx, message); err != nil {
			logrus.WithContext(ctx).WithError(err).WithField("message", msg).Errorln("Error chatting")
			_ = socketResponseAsJson(conn, BuildErrorMessage(room, err.Error()).Message)
			continue
		}
		msg = BuildSocketMessage(message)

		broadcast <- msg
	}
}
