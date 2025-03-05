package application

import (
	"os"
)

type (
	configuration struct {
		RabbitURL      string `json:"rabbit_url"`
		ConsumerQueue  string `json:"consumer_queue"`
		PublisherQueue string `json:"publisher_queue"`
		DatabaseURL    string `json:"database_url"`
		StockPriceURL  string `json:"stock_price_url"`
	}
)

const (
	rabbitURLEnv      = "RABBITMQ_URL"
	consumerQueueEnv  = "CONSUMER_QUEUE"
	publisherQueueEnv = "PUBLISHER_QUEUE"
	databaseURLEnv    = "DB_URL"
	stockPriceURLEnv  = "STOCK_PRICE_URL"
)

func (a *App) BuildConfiguration() *App {
	c := configuration{
		RabbitURL:      os.Getenv(rabbitURLEnv),
		ConsumerQueue:  os.Getenv(consumerQueueEnv),
		PublisherQueue: os.Getenv(publisherQueueEnv),
		DatabaseURL:    os.Getenv(databaseURLEnv),
		StockPriceURL:  os.Getenv(stockPriceURLEnv),
	}
	a.configuration = c
	return a
}
