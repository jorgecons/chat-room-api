package application

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type (
	App struct {
		queueConsumer consumer
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
