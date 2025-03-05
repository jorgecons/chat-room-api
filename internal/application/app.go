package application

import (
	"context"
	"time"

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
	var (
		maxRetries = 10
		conn       *amqp.Connection
		err        error
	)

	for i := 0; i < maxRetries && conn == nil; i++ {
		conn, err = amqp.Dial(a.configuration.RabbitURL)
		if err == nil {
			break
		}
		t := time.Duration(i*2) * time.Second
		logrus.Printf("RabbitMQ not available, retrying... (%d/%d) wait %s", i+1, maxRetries, t)
		time.Sleep(t) // Exponential backoff
	}
	if conn == nil {
		logrus.WithError(err).Panic("Failed to open a channel")
	}
	// Step 2: Create a channel
	ch, err := conn.Channel()
	if err != nil {
		logrus.WithError(err).Panic("Failed to open a channel")
	}
	// Step 3: Declare a queue
	q, err := ch.QueueDeclare(
		a.configuration.ConsumerQueue, // name
		true,                          // durable
		false,                         // delete when unused
		false,                         // exclusive
		false,                         // no-wait
		nil,                           // arguments
	)
	if err != nil {
		logrus.WithError(err).Panic("Failed to declare a queue")
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
			logrus.WithError(err).Panic("couldn't run app")
		}
		logrus.Println(logo)
	}()
	return a
}
