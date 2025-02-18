package listener

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

var QUEUE_NAME string
var RABBITHOST string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	QUEUE_NAME = os.Getenv("QUEUE_NAME")
	RABBITHOST = os.Getenv("RABBITHOST")
}

type Config struct {
	QueueName  string
	RabbitHost string
}

func NewConfig() *Config {
	return &Config{
		QueueName:  QUEUE_NAME,
		RabbitHost: RABBITHOST,
	}
}

type Listener struct {
	Config *Config
}

func NewListener(config *Config) *Listener {
	return &Listener{Config: config}
}

func (l *Listener) Listen() {
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
			if ok {
				log.Printf("Received a message: %s", d.Body)
			} else {
				log.Printf("Channel closed")
			}

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
