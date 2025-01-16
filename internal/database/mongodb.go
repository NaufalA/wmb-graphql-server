package database

import (
	"context"
	"fmt"
	"time"

	"github.com/NaufalA/wmb-graphql-server/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func ConnectMongo(config config.MongoDBConfig) (*mongo.Client, error) {
	uri := fmt.Sprintf(
		"mongodb://%s%s",
		config.Host,
		config.Port,
	)
	opts := options.Client().ApplyURI(uri).SetAuth(options.Credential{
		Username: config.Username,
		Password: config.Password,
	})
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(timeoutCtx, readpref.Primary()); err != nil {
		return nil, err
	}
	return client, nil
}
