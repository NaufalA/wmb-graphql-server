package resolver

import (
	"context"

	"github.com/NaufalA/wmb-graphql-server/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type ProductRepository interface {
	CreateProduct(ctx context.Context, input model.CreateProductInput) (*model.Product, error)
	UpdateProduct(ctx context.Context, input model.UpdateProductInput) (*model.Product, error)
	DeleteProduct(ctx context.Context, id string) (*string, error)
	GetProduct(ctx context.Context, id string) (*model.Product, error)
	ListProducts(ctx context.Context, input model.ProductConnectionArgs) (*model.ProductConnection, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error)
	UpdateUser(ctx context.Context, input model.UpdateUserInput) (*model.User, error)
	DeleteUser(ctx context.Context, id string) (*string, error)
	GetUser(ctx context.Context, id string) (*model.User, error)
	ListUsers(ctx context.Context, input model.UserConnectionArgs) (*model.UserConnection, error)
}

type Resolver struct {
	productRepository ProductRepository
	userRepository UserRepository
}

func NewResolver(
	productRepository ProductRepository,
	userRepository UserRepository,
) *Resolver {
	return &Resolver{
		productRepository: productRepository,
		userRepository: userRepository,
	}
}
