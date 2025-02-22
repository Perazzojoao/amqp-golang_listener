package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PayloadWrapper struct {
	Pattern string       `json:"pattern"`
	Data    EventPayload `json:"data"`
}

type EventPayload struct {
	UserAgent     string `json:"userAgent" bson:"user_agent,omitempty"`
	UserId        int    `json:"userId" bson:"user_id" validate:"required,number"`
	Ip            string `json:"ip" bson:"ip,omitempty"`
	RequestMethod string `json:"requestMethod" bson:"request_method,omitempty"`
	Url           string `json:"url" bson:"url,omitempty"`
	TimeToProcess int    `json:"timeToProcess" bson:"time_to_process,omitempty"`
	Data          string `json:"data" bson:"data,omitempty"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func NewEventPayload(data []byte) EventPayload {
	var wrapper PayloadWrapper
	if err := json.Unmarshal(data, &wrapper); err != nil {
		fmt.Printf("Error unmarshaling wrapper: %v\n", err)
		return EventPayload{}
	}

	return wrapper.Data
}

func ValidatePayload(payload *EventPayload) error {
	if err := validate.Struct(payload); err != nil {
		return fmt.Errorf("error validating payload: %v", err)
	}

	return nil
}

type Mongo struct {
	uri    string
	Config *Config
	Client *mongo.Client
}

func NewMongo(uri string, config *Config) *Mongo {
	return &Mongo{uri: uri, Config: config}
}

func (m *Mongo) Connect(ctx context.Context) {
	c := m.Config
	clientOptions := options.Client().ApplyURI(m.uri)
	clientOptions.Auth = &options.Credential{
		Username: c.MongoUsername,
		Password: c.MongoPassword,
	}
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Conexão com MongoDB estabelecida com sucesso!")

	m.Client = client
}

func (m *Mongo) Disconnect(ctx context.Context) {
	if err := m.Client.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
}

func (m *Mongo) InsertOne(ctx context.Context, collection string, payload EventPayload) error {
	config := m.Config
	data := parsedData{
		Info:      payload,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	c := m.Client.Database(config.MongoDatabase).Collection(collection)
	_, err := c.InsertOne(ctx, data)
	if err != nil {
		return fmt.Errorf("error inserting document: %v", err)
	}

	return nil
}

type parsedData struct {
	Info      EventPayload `json:"info" bson:"info"`
	CreatedAt time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" bson:"updated_at"`
}
