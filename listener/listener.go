package listener

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Listener struct {
	Config   *Config
	Database *Mongo
}

func NewListener(config *Config) *Listener {
	return &Listener{Config: config, Database: NewMongo(config.MongoUri, config)}
}

func (l *Listener) Listen() {
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Connect to MongoDB

	l.Database.Connect(ctx)
	defer l.Database.Disconnect(context.TODO())

	// Connect to RabbitMQ
	conn, err := amqp.Dial(l.Config.RabbitHost)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		l.Config.QueueName, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	log.Printf(" [*] Waiting for messages. Press CTRL+C to exit.")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case d, ok := <-msgs:
			if !ok {
				log.Printf("Channel closed")
			}
			payload := NewEventPayload(d.Body)
			if payload.UserId == "" {
				respondOnError(ch, d, "Missing user_id in the payload")
				continue
			}

			collection := l.Config.MongoCollection
			err := l.Database.InsertOne(context.Background(), collection, payload)
			if err != nil {
				log.Printf("Failed to insert a document: %v", err)
			}
			log.Printf("Received a message: %v", payload)
		case <-sigChan:
			log.Printf("Received signal to stop")
			os.Exit(0)
		}
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func respondOnError(ch *amqp.Channel, d amqp.Delivery, errMsg string) {
	log.Printf("Error: %s", errMsg)
	publishErr := ch.Publish(
		"",        // exchange
		d.ReplyTo, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: d.CorrelationId,
			Body:          []byte(fmt.Sprintf(`{"error": "%s"}`, errMsg)),
		})
	if publishErr != nil {
		log.Printf("Failed to publish error message: %v", publishErr)
	}
	d.Nack(false, false) // Reject the message whitout requeue
}
