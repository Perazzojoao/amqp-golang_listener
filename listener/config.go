package listener

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var QUEUE_NAME string
var RABBITHOST string
var MONGODB_URI string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file - %v", err)
	}

	QUEUE_NAME = os.Getenv("QUEUE_NAME")
	RABBITHOST = os.Getenv("RABBITHOST")
	MONGODB_URI = os.Getenv("MONGODB_URI")
}

type Config struct {
	QueueName  string
	RabbitHost string
	MongoUri   string
}

func NewConfig() *Config {
	return &Config{
		QueueName:  QUEUE_NAME,
		RabbitHost: RABBITHOST,
		MongoUri:   MONGODB_URI,
	}
}
