package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type EventPayload struct {
	UserAgent     string      `json:"userAgent" bson:"user_agent,omitempty"`
	UserId        string      `json:"userId" bson:"user_id,omitempty"`
	Ip            string      `json:"ip" bson:"ip,omitempty"`
	RequestMethod string      `json:"requestMethod" bson:"request_method,omitempty"`
	Url           string      `json:"url" bson:"url,omitempty"`
	TimeToProcess int         `json:"timeToProcess" bson:"time_to_process,omitempty"`
	Data          interface{} `json:"data" bson:"data,omitempty"`
}

func NewEventPayload(data []byte) EventPayload {
	var event EventPayload
	if err := json.Unmarshal(data, &event); err != nil {
		return EventPayload{}
	}
	return event
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
	c := m.Client.Database(config.MongoDatabase).Collection(collection)
	_, err := c.InsertOne(ctx, payload)
	if err != nil {
		return err
	}
	return nil
}
