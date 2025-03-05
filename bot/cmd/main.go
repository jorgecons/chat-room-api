package main

import (
	"bot/internal/application"
)

func main() {
	application.NewApp().
		BuildConfiguration().
		BuildConsumer().
		BuildHandlers().
		MapEventRoutes().
		Consume()
}
