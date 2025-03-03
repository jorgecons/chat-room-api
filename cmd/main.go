package main

import (
	"chat-room-api/internal/application"
)

func main() {
	application.NewApp().
		BuildConfiguration().
		BuildConsumer().
		BuildRouter().
		BuildSocket().
		BuildHandlers().
		MapSocket().
		MapEventRoutes().
		MapWebRoutes().
		Run().
		Consume()
}
