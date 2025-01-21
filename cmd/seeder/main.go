package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"slices"
	"time"

	"github.com/NaufalA/wmb-graphql-server/config"
	"github.com/NaufalA/wmb-graphql-server/graph/model"
	"github.com/NaufalA/wmb-graphql-server/internal/collection"
	"github.com/NaufalA/wmb-graphql-server/internal/database"
	"github.com/NaufalA/wmb-graphql-server/internal/repository"
	"github.com/NaufalA/wmb-graphql-server/pkg/util"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

	ctx := context.Background()
	db := mongoClient.Database(mongoConfig.Database)
	collections, _ := db.ListCollectionNames(ctx, bson.D{})

	userCollectionName := collection.User{}.CollectionName()
	if !slices.Contains(collections, userCollectionName) {
		db.CreateCollection(ctx, userCollectionName)
		logrus.Info("Created user collection")
	}
	db.Collection(userCollectionName).Indexes().DropOne(ctx, "search")
	userIndexModel := mongo.IndexModel{
		Options: options.Index().SetName("search"),
		Keys:    bson.D{{Key: "email", Value: "text"}, {Key: "fullName", Value: "text"}},
	}
	_, err = db.Collection(userCollectionName).Indexes().CreateOne(ctx, userIndexModel)
	if err != nil {
		logrus.Errorf("Created user search index failed: %s", err.Error())
	} else {
		logrus.Info("Created user search index")
	}
	newUsers := []model.CreateUserInput{
		{
			Email:    "superadmin@wmb.com",
			FullName: "SuperAdmin",
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
		user := collection.User{}
		filter := bson.M{"email": u.Email}
		err := mongoClient.Database(mongoConfig.Database).Collection(user.CollectionName()).FindOne(ctx, filter).Decode(&user)
		if user.Email != nil {
			logrus.Error(err)
			continue
		}
		id := primitive.ObjectID(bson.NewObjectID())
		now := time.Now()
		passwordUtil := util.PasswordUtil{}
		passwordHash, err := passwordUtil.HashPassword(u.Password)
		if err != nil {
			logrus.Error(err)
		}
		col := collection.User{
			ID:           id,
			Email:        &u.Email,
			FullName:     &u.FullName,
			Role:         &u.FullName,
			PasswordHash: &passwordHash,
			CreateTime:   &now,
		}
		_, err = mongoClient.Database(mongoConfig.Database).Collection(col.CollectionName()).InsertOne(ctx, col)
		if err != nil {
			logrus.Error(err)
		}
		if err != nil {
			logrus.Errorf("failed create user: %s", err.Error())
		} else {
			logrus.Infof(
				"success create user with credentials: email: %s password: %s",
				u.Email,
				u.Password,
			)
		}
	}
	logrus.Info("Created user data")

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

	productCollectionName := collection.Product{}.CollectionName()
	if !slices.Contains(collections, productCollectionName) {
		db.CreateCollection(ctx, productCollectionName)
		logrus.Info("Created product collection")
	}
	db.Collection(productCollectionName).Indexes().DropOne(ctx, "search")
	productIndexModel := mongo.IndexModel{
		Options: options.Index().SetName("search"),
		Keys:    bson.D{{Key: "productName", Value: "text"}},
	}
	_, err = db.Collection(productCollectionName).Indexes().CreateOne(ctx, productIndexModel)
	if err != nil {
		logrus.Errorf("Created product search index failed: %s", err.Error())
	} else {
		logrus.Info("Created product search index")
	}
	for _, p := range productNames {
		productName := p
		price := math.Ceil(rand.Float64() * 10000)
		stock := math.Floor(rand.Float64() * 100)
		stockUnit := units[rand.Intn(4)]

		productRepository.CreateProduct(ctx, model.CreateProductInput{
			ProductName: productName,
			Price:       price,
			Stock:       stock,
			StockUnit:   stockUnit,
		})
		time.Sleep(2 * time.Second)
	}
	logrus.Info("Created product data")

	logrus.Info("seeding complete")
}
