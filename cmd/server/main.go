package main

import (
	"context"
	"fmt"

	"github.com/NaufalA/wmb-graphql-server/config"
	"github.com/NaufalA/wmb-graphql-server/internal/database"

	"github.com/sirupsen/logrus"
)

func main() {
	mongoConfig := config.MongoDBConfig{
		Host:     "localhost",
		Port:     ":27017",
		Username: "root",
		Password: "D1fficultPAssw0rd",
	}
	mongoClient, err := database.ConnectMongo(mongoConfig)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Info(fmt.Sprintf("Successfully Connected to MongoDB Server %s", mongoConfig.Host))

	defer func() {
		err := mongoClient.Disconnect(context.Background())
		if err != nil {
			logrus.Panic(err)
		}
	}()
}
