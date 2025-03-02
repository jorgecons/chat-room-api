package application

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"chat-room-api/internal/core/adapter/sockethandler"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (a *App) MapSocket() *App {
	a.socket.broadcast = make(chan sockethandler.Message)
	a.router.GET("/ws/:room", a.buildHandler(a.socket.broadcast))
	go sockethandler.HandleMessages(a.socket.broadcast)
	return a
}

// Handle new WebSocket connections
func (a *App) buildHandler(broadcast chan sockethandler.Message) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logrus.WithContext(c).WithError(err).Error("WebSocket upgrade error")
			return
		}
		defer func() { _ = conn.Close() }()

		if err = a.handlers.ConnectChatHandler(c, conn, broadcast); err != nil {
			return
		}
		_ = a.handlers.ChatHandler(c, conn, broadcast)
	}
}
