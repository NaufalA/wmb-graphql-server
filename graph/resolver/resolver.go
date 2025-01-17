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

type Resolver struct {
	productRepository ProductRepository
}

func NewResolver(
	productRepository ProductRepository,
) *Resolver {
	return &Resolver{
		productRepository: productRepository,
	}
}
