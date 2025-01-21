package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"

	"github.com/NaufalA/wmb-graphql-server/config"
	"github.com/NaufalA/wmb-graphql-server/graph/model"
	"github.com/NaufalA/wmb-graphql-server/internal/database"
	"github.com/NaufalA/wmb-graphql-server/internal/repository"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Error("error loading .env file ")
	}
}

const defaultPort = "8080"

func main() {
	logger := logrus.New()

	mongoConfig := config.MongoDBConfig{
		Host:     os.Getenv("MONGODB_HOST"),
		Port:     os.Getenv("MONGODB_PORT"),
		Username: os.Getenv("MONGODB_USERNAME"),
		Password: os.Getenv("MONGODB_PASSWORD"),
		Database: os.Getenv("MONGODB_DATABASE"),
	}
	mongoClient, err := database.ConnectMongo(mongoConfig)
	if err != nil {
		logger.Panic(err)
	}
	logger.Info(fmt.Sprintf("Successfully Connected to MongoDB Server %s", mongoConfig.Host))

	defer func() {
		err := mongoClient.Disconnect(context.Background())
		if err != nil {
			logger.Panic(err)
		}
	}()
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = defaultPort
	}

	productRepository := repository.NewProductRepository(logger, mongoClient.Database(mongoConfig.Database))
	userRepository := repository.NewUserRepository(logger, mongoClient.Database(mongoConfig.Database))

	ctx := context.Background()
	newUsers := []model.CreateUserInput{
		{
			Email:    "superadmin@wmb.com",
			FullName: "Super Admin",
			Password: bson.NewObjectID().Hex(),
		},
		{
			Email:    "admin@wmb.com",
			FullName: "Admin",
			Password: bson.NewObjectID().Hex(),
		},
		{
			Email:    "guest@wmb.com",
			FullName: "Guest",
			Password: bson.NewObjectID().Hex(),
		},
	}
	for _, u := range newUsers {
		createdUser, err := userRepository.CreateUser(ctx, u)
		if err != nil {
			logrus.Errorf("failed create user: %s", err.Error())
		} else {
			logrus.Infof(
				"success create user with credentials: email: %s password: %s",
				*createdUser.Email,
				u.Password,
			)
		}
	}

	var productNames = []string{
		"Rendang",
		"Gado-gado",
		"Nasi Goreng",
		"Sate",
		"Bakso",
		"Soto Ayam",
		"Pempek",
		"Mie Goreng",
		"Rawon",
		"Gudeg",
		"Rendang Jengkol",
		"Sambal Goreng Ati",
		"Tongseng",
		"Sate Lilit",
		"Nasi Kuning",
		"Semur Jengkol",
		"Soto Betawi",
		"Gulai Ikan",
		"Sate Padang",
		"Opor Ayam",
	}

	units := []string{
		"portions",
		"plates",
		"orders",
		"servings",
	}

	for _, p := range productNames {
		productName := p
		price := rand.Float64() * 10000
		stock := rand.Float64() * 100
		stockUnit := units[rand.Intn(4)]

		productRepository.CreateProduct(ctx, model.CreateProductInput{
			ProductName: productName,
			Price:       price,
			Stock:       stock,
			StockUnit:   stockUnit,
		})
	}

	logrus.Info("seeding complete")
}
