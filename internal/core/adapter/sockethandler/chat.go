package sockethandler

import (
	"chat-room-api/internal/core/domain"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type (
	ChatFeature interface {
		Chat(context.Context, domain.Message) error
	}
	StockFeature interface {
		Stock(context.Context, string, string) error
	}

	chat struct {
		chatFeature  ChatFeature
		stockFeature StockFeature
		jwtSecret    []byte
	}
)

func NewChat(cf ChatFeature, sf StockFeature, jwtSecret []byte) HandlerFunc {
	handler := chat{
		chatFeature:  cf,
		stockFeature: sf,
		jwtSecret:    jwtSecret,
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
			log.Println("Read error:", err)
			continue
		}
		username := ctx.Value(UserContextKey).(string)
		if err = ValidateUsername(username, msg); err != nil {
			logrus.WithContext(ctx).WithField("username", msg.Username).WithError(err).Errorln("Error validating user")
			_ = conn.WriteJSON(BuildErrorMessage(room, err.Error()))
			continue
		}
		if err = ValidateRoom(room, msg); err != nil {
			logrus.WithContext(ctx).WithField("room", room).WithError(err).Errorln("Error validating message")
			_ = conn.WriteJSON(BuildErrorMessage(room, err.Error()))
			continue
		}

		// Validate and handle unrecognized messages
		if strings.HasPrefix(msg.Text, "/") {
			if !strings.HasPrefix(msg.Text, "/stock=") {
				logrus.WithContext(ctx).WithField("command", msg.Text).Errorln("Error validating command")
				m := BuildErrorMessage(room, fmt.Sprintf(UnknownCommand, msg.Text))
				_ = conn.WriteJSON(m)
				continue
			}
			_ = c.stockFeature.Stock(ctx, msg.Room, GetStockName(msg.Text))
			// continue //if command should not be a message for all participants
		}
		if err = c.chatFeature.Chat(ctx, BuildMessage(msg)); err != nil {
			logrus.WithContext(ctx).WithField("message", msg).WithError(err).Errorln("Error chatting")
			_ = conn.WriteJSON(BuildErrorMessage(room, fmt.Sprintf(UnknownCommand, msg.Text)))
			continue
		}

		broadcast <- msg
	}
}
