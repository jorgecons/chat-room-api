package application

import (
	"context"
	"log"

	"chat-room-api/internal/core/usecase/connectchatusecase"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"chat-room-api/internal/core/adapter/eventhandler"
	"chat-room-api/internal/core/adapter/sockethandler"
	"chat-room-api/internal/core/adapter/webhandler"
	"chat-room-api/internal/core/repository/message"
	"chat-room-api/internal/core/repository/publisher"
	"chat-room-api/internal/core/repository/stock"
	"chat-room-api/internal/core/repository/user"
	"chat-room-api/internal/core/usecase/botusecase"
	"chat-room-api/internal/core/usecase/chatusecase"
	"chat-room-api/internal/core/usecase/createaccountusecase"
	"chat-room-api/internal/core/usecase/loginusecase"
)

type (
	Handlers struct {
		ConnectChatHandler       sockethandler.HandlerFunc
		ChatHandler              sockethandler.HandlerFunc
		PostMessageHandler       func(chan sockethandler.Message)
		ConsumeBotMessageHandler func(context.Context, []byte) error
		CreateAccountHandler     gin.HandlerFunc
		LoginHandler             gin.HandlerFunc
	}
)

func (a *App) BuildHandlers() *App {
	secret := []byte(a.configuration.JWTSecret)
	p := NewPublisher(a.configuration)
	stockPriceClient := resty.New().SetBaseURL(a.configuration.StockPriceURL)
	dbClient := newDBClient(a.configuration.DatabaseURL)

	publisherRepo := publisher.NewPublisher(p, a.configuration.RabbitQueue)
	stockRepo := stock.NewStockPriceRepo(stockPriceClient)
	messageStorage := message.NewMessageStorage(dbClient)
	userStorage := user.NewUserStorage(dbClient)

	chatUseCase := chatusecase.NewUseCase(messageStorage, publisherRepo)
	botUseCase := botusecase.NewUseCase(stockRepo, publisherRepo)
	createAccountUseCase := createaccountusecase.NewUseCase(userStorage)
	loginUseCase := loginusecase.NewUseCase(userStorage, []byte(a.configuration.JWTSecret))
	connectChat := connectchatusecase.NewUseCase(messageStorage)

	a.handlers = Handlers{
		ConnectChatHandler:       sockethandler.NewConnectChat(connectChat, secret),
		ChatHandler:              sockethandler.NewChat(chatUseCase, secret),
		PostMessageHandler:       sockethandler.HandleMessages,
		ConsumeBotMessageHandler: eventhandler.NewBotHandler(botUseCase),
		CreateAccountHandler:     webhandler.NewCreateAccount(createAccountUseCase),
		LoginHandler:             webhandler.NewLogin(loginUseCase),
	}

	return a
}

func newDBClient(url string) *pgx.Conn {
	db, err := connectDB(url)
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to connect to database")
	}
	return db
}

func connectDB(url string) (*pgx.Conn, error) {
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
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	_, err = ch.QueueDeclare(
		config.RabbitQueue,
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}
	return ch
}
