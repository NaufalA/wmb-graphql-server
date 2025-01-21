package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/NaufalA/wmb-graphql-server/api/controller"
	"github.com/NaufalA/wmb-graphql-server/api/middleware"
	"github.com/NaufalA/wmb-graphql-server/config"
	"github.com/NaufalA/wmb-graphql-server/graph/resolver"
	"github.com/NaufalA/wmb-graphql-server/internal/database"
	"github.com/NaufalA/wmb-graphql-server/internal/repository"
	"github.com/NaufalA/wmb-graphql-server/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const defaultPort = "8080"

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Error("error loading .env file ")
	}
}

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

	authService := service.NewAuthService(logger, userRepository)
	authController := controller.NewAuthController(logger, authService)

	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		auth.POST("/register", authController.Register)
		auth.POST("/resetPassword", authController.ResetPassword)
	}
	query := r.Group("/query")
	{
		query.Use(middleware.AuthTokenMiddleware())
		query.POST("", controller.GraphqlHandler(resolver.NewResolver(
			productRepository,
			userRepository,
		)))
	}
	r.GET("/playground", controller.PlaygroundHandler())

	logger.Printf("server listening at http://localhost:%s", port)
	logger.Fatal(http.ListenAndServe(":"+port, r))
}
