package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client connections organized by chatroom
var chatrooms = make(map[string]map[*websocket.Conn]bool)

// Broadcast channel
var broadcast = make(chan Message)

// Message structure
type Message struct {
	Room     string `json:"room"`     // Chatroom name
	Username string `json:"username"` // Sender
	Text     string `json:"text"`     // Message content
}

func main() {
	r := gin.Default()

	// WebSocket endpoint
	r.GET("/ws/:room", handleConnections)

	// Start message handler
	go handleMessages()

	// Start server
	log.Println("Server started on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed:", err)
	}
}

// Handle new WebSocket connections
func handleConnections(c *gin.Context) {
	room := c.Param("room") // Get chatroom from URL param
	if room == "" {
		log.Println("Missing room name")
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// Ensure room exists
	if chatrooms[room] == nil {
		chatrooms[room] = make(map[*websocket.Conn]bool)
	}
	chatrooms[room][conn] = true

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			delete(chatrooms[room], conn)
			break
		}

		// Ensure message belongs to the correct room
		msg.Room = room

		// Validate and handle unrecognized messages
		if strings.HasPrefix(msg.Text, "/") && !strings.HasPrefix(msg.Text, "/stock=") {
			errMsg := Message{
				Room:     room,
				Username: "System",
				Text:     "⚠️ Unknown command: " + msg.Text,
			}
			conn.WriteJSON(errMsg)
			continue
		}

		// Send valid message to the broadcast channel
		broadcast <- msg
	}
}

// Handle incoming messages and broadcast to the correct chatroom
func handleMessages() {
	for {
		msg := <-broadcast
		for client := range chatrooms[msg.Room] {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Write error:", err)
				client.Close()
				delete(chatrooms[msg.Room], client)
			}
		}
	}
}
