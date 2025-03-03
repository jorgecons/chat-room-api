package application

import (
	"context"

	"chat-room-api/internal/core/adapter/sockethandler"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

const (
	logo = `─────── ⚡ is up and running ⚡ ───────`
)

type (
	App struct {
		router        *gin.Engine
		queueConsumer consumer
		socket        socket
		configuration configuration
		context       context.Context
		handlers      Handlers
	}

	consumer struct {
		queue      amqp.Queue
		consumer   *amqp.Channel
		connection *amqp.Connection
		consumerFn func(context.Context, chan bool, <-chan amqp.Delivery)
	}

	socket struct {
		socketClients map[*websocket.Conn]bool
		broadcast     chan sockethandler.Message
	}
)

// NewApp constructor for fury web application.
func NewApp() *App {
	a := new(App)
	a.configureContext()
	return a
}

func (a *App) configureContext() {
	a.context = context.Background()
}

func (a *App) BuildConsumer() *App {
	// Step 1: Connect to RabbitMQ
	conn, err := amqp.Dial(a.configuration.RabbitURL)
	if err != nil {
		logrus.WithContext(a.context).WithError(err).Fatalln("Failed to connect to RabbitMQ")
	}

	// Step 2: Create a channel
	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("Failed to open a channel: %v", err)
	}

	// Step 3: Declare a queue
	q, err := ch.QueueDeclare(
		a.configuration.RabbitQueue, // name
		true,                        // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
	)
	if err != nil {
		logrus.Fatalf("Failed to declare a queue: %v", err)
	}

	a.queueConsumer = consumer{
		queue:      q,
		consumer:   ch,
		connection: conn,
	}
	return a
}

func (a *App) BuildRouter() *App {
	a.router = gin.Default()
	return a
}

func (a *App) BuildSocket() *App {
	a.socket = socket{
		socketClients: make(map[*websocket.Conn]bool),
	}
	return a
}

func (a *App) Consume() *App {
	msgs, err := a.queueConsumer.consumer.Consume(
		a.queueConsumer.queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.WithContext(a.context).WithError(err).Fatalln("Failed to register a consumer")
	}

	forever := make(chan bool)
	go a.queueConsumer.consumerFn(a.context, forever, msgs)
	logrus.WithContext(a.context).Println("Waiting for messages.")
	<-forever
	defer a.queueConsumer.consumer.Close()
	defer a.queueConsumer.connection.Close()
	return a
}

func (a *App) Run() *App {
	go func() {
		if err := a.router.Run(a.configuration.RouterURL); err != nil {
			logrus.Panic("couldn't run app", logrus.WithError(err))
		}
		logrus.Println(logo)
	}()
	return a
}
