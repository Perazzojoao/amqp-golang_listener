package listener

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var QUEUE_NAME string
var RABBITHOST string
var MONGODB_URI string
var MONGODB_COLLECTION string
var MONGODB_DATABASE string
var MONGODB_USERNAME string
var MONGODB_PASSWORD string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file - %v", err)
	}

	QUEUE_NAME = os.Getenv("QUEUE_NAME")
	RABBITHOST = os.Getenv("RABBITHOST")
	MONGODB_URI = os.Getenv("MONGODB_URI")
	MONGODB_COLLECTION = os.Getenv("MONGODB_COLLECTION")
	MONGODB_DATABASE = os.Getenv("MONGODB_DATABASE")
	MONGODB_USERNAME = os.Getenv("MONGODB_USERNAME")
	MONGODB_PASSWORD = os.Getenv("MONGODB_PASSWORD")
}

type Config struct {
	QueueName       string
	RabbitHost      string
	MongoUri        string
	MongoCollection string
	MongoDatabase   string
	MongoUsername   string
	MongoPassword   string
}

func NewConfig() *Config {
	return &Config{
		QueueName:       QUEUE_NAME,
		RabbitHost:      RABBITHOST,
		MongoUri:        MONGODB_URI,
		MongoDatabase:   MONGODB_DATABASE,
		MongoCollection: MONGODB_COLLECTION,
		MongoUsername:   MONGODB_USERNAME,
		MongoPassword:   MONGODB_PASSWORD,
	}
}
