package application

import (
	"context"

	"bot/internal/core/adapter/eventhandler"
	"bot/internal/core/repository/message"
	"bot/internal/core/repository/publisher"
	"bot/internal/core/repository/stock"
	"bot/internal/core/usecase/botusecase"

	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type (
	Handlers struct {
		ConsumeBotMessageHandler func(context.Context, []byte) error
	}

	querier interface {
		Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
		Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
		QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	}
)

func (a *App) BuildHandlers() *App {
	p := NewPublisher(a.configuration)
	stockPriceClient := resty.New().SetBaseURL(a.configuration.StockPriceURL)
	dbClient := newDBClient(a.configuration.DatabaseURL)

	publisherRepo := publisher.NewPublisher(p, a.configuration.PublisherQueue)
	stockRepo := stock.NewStockPriceRepo(stockPriceClient)
	messageStorage := message.NewMessageStorage(dbClient)

	botUseCase := botusecase.NewUseCase(stockRepo, messageStorage, publisherRepo)

	a.handlers = Handlers{
		ConsumeBotMessageHandler: eventhandler.NewBotHandler(botUseCase),
	}

	return a
}

func newDBClient(url string) querier {
	db, err := connectDB(url)
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to connect to database")
	}
	return db
}

func connectDB(url string) (querier, error) {
	connConfig, err := pgx.ParseConfig(url)
	if err != nil {
		logrus.WithError(err).Fatalln("Unable to parse connection string")
	}

	// Connect to the database
	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewPublisher(config configuration) *amqp.Channel {
	conn, err := amqp.Dial(config.RabbitURL)
	if err != nil {
		logrus.WithError(err).Panic("Failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		logrus.WithError(err).Panic("Failed to open a channel")
	}
	_, err = ch.QueueDeclare(
		config.PublisherQueue,
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		logrus.WithError(err).Panic("Failed to declare queue")
	}
	return ch
}
