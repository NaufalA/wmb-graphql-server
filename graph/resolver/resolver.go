package resolver

import (
	"context"

	"github.com/NaufalA/wmb-graphql-server/graph/model"
	"github.com/NaufalA/wmb-graphql-server/internal/collection"
	"github.com/NaufalA/wmb-graphql-server/internal/dto"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type ProductRepository interface {
	CreateProduct(ctx context.Context, input model.CreateProductInput) (*model.Product, error)
	UpdateProduct(ctx context.Context, input model.UpdateProductInput) (*model.Product, error)
	DeleteProduct(ctx context.Context, id string) (*string, error)
	GetProduct(ctx context.Context, id string) (*model.Product, error)
	ListProducts(ctx context.Context, input dto.ConnectionRequest) (*model.ProductConnection, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, input model.CreateUserInput) (*collection.User, error)
	UpdateUser(ctx context.Context, input collection.User) (*collection.User, error)
	DeleteUser(ctx context.Context, id string) (*string, error)
	GetUser(ctx context.Context, request dto.GetUserRequest) (*collection.User, error)
	ListUsers(ctx context.Context, input dto.ConnectionRequest) (*model.UserConnection, error)
}

type Resolver struct {
	productRepository ProductRepository
	userRepository    UserRepository
}

func NewResolver(
	productRepository ProductRepository,
	userRepository UserRepository,
) *Resolver {
	return &Resolver{
		productRepository: productRepository,
		userRepository:    userRepository,
	}
}
