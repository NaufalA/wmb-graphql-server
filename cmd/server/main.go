package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/NaufalA/wmb-graphql-server/config"
	"github.com/NaufalA/wmb-graphql-server/graph"
	"github.com/NaufalA/wmb-graphql-server/graph/resolver"
	"github.com/NaufalA/wmb-graphql-server/internal/database"
	"github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

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
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers:	&resolver.Resolver{},
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	logrus.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	logrus.Fatal(http.ListenAndServe(":"+port, nil))
}