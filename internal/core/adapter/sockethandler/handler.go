package sockethandler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"chat-room-api/internal/core/domain"

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

	GetLastMessages interface {
		GetLastMessages(context.Context, string) ([]domain.Message, error)
	}
	chat struct {
		chatFeature     ChatFeature
		stockFeature    StockFeature
		getLastMessages GetLastMessages
		mutex           sync.Mutex
		jwtSecret       []byte
	}
)

func NewChat(cf ChatFeature, sf StockFeature, lm GetLastMessages, jwtSecret []byte) Handler {
	handler := chat{
		chatFeature:     cf,
		stockFeature:    sf,
		getLastMessages: lm,
		mutex:           sync.Mutex{},
		jwtSecret:       jwtSecret,
	}
	handlerFunc := func(c *gin.Context, conn *websocket.Conn, broadcast chan Message) {
		handler.handle(c, conn, broadcast)
	}
	return handlerFunc
}

func (c *chat) handle(ctx *gin.Context, conn *websocket.Conn, broadcast chan Message) {
	// Extract JWT token from request headers
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		_ = conn.WriteJSON(Message{Username: "system", Text: "❌ Missing JWT token"})
		return
	}

	// Validate the token
	claims, err := validateJWT(c.jwtSecret, tokenString)
	if err != nil {
		_ = conn.WriteJSON(Message{Username: "system", Text: "❌ Invalid token: " + err.Error()})
		return
	}
	username := claims["username"].(string)

	room := ctx.Param("room") // Get chatroom from URL param
	if room == "" {
		logrus.WithContext(ctx).Errorln("Error missing room")
		return
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

	for {
		var msg Message
		err = conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			continue
		}
		if err = ValidateUsername(username, msg); err != nil {
			logrus.WithContext(ctx).WithField("username", msg.Username).WithError(err).Errorln("Error validating user")
			_ = conn.WriteJSON(CreateErrorMessage(msg.Username, err.Error()))
		}
		if err = ValidateRoom(room, msg); err != nil {
			logrus.WithContext(ctx).WithField("room", room).WithError(err).Errorln("Error validating message")
			_ = conn.WriteJSON(CreateErrorMessage(room, err.Error()))
			continue
		}

		// Validate and handle unrecognized messages
		if strings.HasPrefix(msg.Text, "/") {
			if !strings.HasPrefix(msg.Text, "/stock=") {
				logrus.WithContext(ctx).WithField("command", msg.Text).Errorln("Error validating command")
				m := CreateErrorMessage(room, fmt.Sprintf(UnknownCommand, msg.Text))
				_ = conn.WriteJSON(m)
				continue
			}
			_ = c.stockFeature.Stock(ctx, msg.Room, GetStockName(msg.Text))
			// continue //if command should not be a message for all participants
		}
		if err = c.chatFeature.Chat(ctx, BuildMessage(msg)); err != nil {
			logrus.WithContext(ctx).WithField("message", msg).WithError(err).Errorln("Error chatting")
			_ = conn.WriteJSON(CreateErrorMessage(room, fmt.Sprintf(UnknownCommand, msg.Text)))
			continue
		}

		broadcast <- msg
	}
}

func HandleMessages(broadcast chan Message) {
	for {
		msg := <-broadcast
		for client := range Chatrooms[msg.Room] {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Write error:", err)
				_ = client.Close()
				delete(Chatrooms[msg.Room], client)
			}
		}
	}
}
