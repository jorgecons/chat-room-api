package application

import (
	"os"
)

type (
	configuration struct {
		RouterURL      string `json:"router_url"`
		RabbitURL      string `json:"rabbit_url"`
		ConsumerQueue  string `json:"consumer_queue"`
		PublisherQueue string `json:"publisher_queue"`
		DatabaseURL    string `json:"database_url"`
		StockPriceURL  string `json:"stock_price_url"`
		JWTSecret      string `json:"-"`
	}
)

const (
	routerURLEnv      = "ROUTER_URL"
	rabbitURLEnv      = "RABBITMQ_URL"
	consumerQueueEnv  = "CONSUMER_QUEUE"
	publisherQueueEnv = "PUBLISHER_QUEUE"
	databaseURLEnv    = "DB_URL"
	stockPriceURLEnv  = "STOCK_PRICE_URL"
	JWTSecretEnv      = "JWT_SECRET"
)

func (a *App) BuildConfiguration() *App {
	c := configuration{
		RouterURL:      os.Getenv(routerURLEnv),
		RabbitURL:      os.Getenv(rabbitURLEnv),
		ConsumerQueue:  os.Getenv(consumerQueueEnv),
		PublisherQueue: os.Getenv(publisherQueueEnv),
		DatabaseURL:    os.Getenv(databaseURLEnv),
		StockPriceURL:  os.Getenv(stockPriceURLEnv),
		JWTSecret:      os.Getenv(JWTSecretEnv),
	}
	a.configuration = c
	return a
}
