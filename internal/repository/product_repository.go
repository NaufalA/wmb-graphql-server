package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/NaufalA/wmb-graphql-server/graph/model"
	"github.com/NaufalA/wmb-graphql-server/internal/collection"
	"github.com/NaufalA/wmb-graphql-server/pkg/util"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ProductRepository struct {
	logger *logrus.Entry
	db     *mongo.Database
}

func NewProductRepository(
	logger *logrus.Logger,
	db *mongo.Database,
) *ProductRepository {
	return &ProductRepository{
		logger: logger.WithFields(logrus.Fields{
			"location": "ProductRepository",
		}),
		db: db,
	}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, input model.CreateProductInput) (*model.Product, error) {
	id := primitive.ObjectID(bson.NewObjectID())
	now := time.Now()
	col := collection.Product{
		ID:          id,
		ProductName: &input.ProductName,
		Price:       &input.Price,
		Stock:       &input.Stock,
		StockUnit:   &input.StockUnit,
		CreateTime:  &now,
	}
	_, err := r.db.Collection(col.CollectionName()).InsertOne(ctx, col)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	model := &model.Product{
		ID:          col.ID.Hex(),
		ProductName: &input.ProductName,
		Price:       &input.Price,
		Stock:       &input.Stock,
		StockUnit:   &input.StockUnit,
		CreateTime:  &now,
	}

	return model, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, input model.UpdateProductInput) (*model.Product, error) {
	id, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	product := collection.Product{
		ID: id,
	}
	filter := bson.M{"_id": product.ID}
	err = r.db.Collection(product.CollectionName()).FindOne(ctx, filter).Decode(&product)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	now := time.Now()
	if input.ProductName != nil {
		product.ProductName = input.ProductName
	}
	if input.Price != nil {
		product.Price = input.Price
	}
	if input.Stock != nil {
		product.Stock = input.Stock
	}
	if input.StockUnit != nil {
		product.StockUnit = input.StockUnit
	}
	product.UpdateTime = &now
	update := bson.M{"$set": product}

	collectionName := collection.Product{}.CollectionName()
	result, err := r.db.Collection(collectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if result.ModifiedCount == 0 {
		logrus.Error("no data modified")
		return nil, err
	}

	model := &model.Product{
		ID:          input.ID,
		ProductName: input.ProductName,
		Price:       input.Price,
		Stock:       input.Stock,
		StockUnit:   input.StockUnit,
		CreateTime:  product.CreateTime,
		UpdateTime:  &now,
	}

	return model, nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id string) (*string, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	col := collection.Product{
		ID: objectID,
	}

	filter := bson.M{"_id": col.ID}
	result, err := r.db.Collection(col.CollectionName()).DeleteOne(ctx, filter)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if result.DeletedCount == 0 {
		logrus.Error("no data deleted")
		return nil, err
	}

	return nil, nil
}

func (r *ProductRepository) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	col := collection.Product{
		ID: objectID,
	}
	filter := bson.M{"_id": col.ID}
	err = r.db.Collection(col.CollectionName()).FindOne(ctx, filter).Decode(&col)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &model.Product{
		ID:          col.ID.Hex(),
		ProductName: col.ProductName,
		Price:       col.Price,
		Stock:       col.Stock,
		StockUnit:   col.StockUnit,
		CreateTime:  col.CreateTime,
		UpdateTime:  col.UpdateTime,
	}, nil
}

func (r *ProductRepository) ListProducts(ctx context.Context, input model.ProductConnectionArgs) (*model.ProductConnection, error) {
	paginationUtil := util.PaginationUtil{}
	filter := bson.M{}
	var opts *options.FindOptionsBuilder
	direction := 1
	if input.First != nil {
		if *input.First < 0 {
			err := fmt.Errorf("invalid parameter, first cannot be negative")
			logrus.Error(err)
			return nil, err
		}
		// Forward Pagination
		opts = options.Find().SetLimit(int64(*input.First)).SetSort(bson.M{"createTime": direction})
		if input.After != nil {
			paginationCursor, _ := paginationUtil.DecodeCursor(*input.After)
			createtime, _ := time.Parse(time.RFC3339, paginationCursor)
			filter["createTime"] = bson.M{"$gt": createtime}
		}
	} else if input.Last != nil {
		if *input.Last < 0 {
			err := fmt.Errorf("invalid parameter, last cannot be negative")
			logrus.Error(err)
			return nil, err
		}
		// Backward Pagination
		direction = -1
		opts = options.Find().SetLimit(int64(*input.Last)).SetSort(bson.M{"createTime": direction})
		if input.Before != nil {
			paginationCursor, _ := paginationUtil.DecodeCursor(*input.After)
			filter["createTime"] = bson.M{"$lt": paginationCursor}
		}
	}

	collectionName := collection.Product{}.CollectionName()
	cursor, err := r.db.Collection(collectionName).Find(ctx, filter, opts)
	defer func() {
		if cursor != nil {
			cursor.Close(context.Background())
		}
	}()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	products := []collection.Product{}
	err = cursor.All(ctx, &products)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	result := &model.ProductConnection{
		Edges:    []*model.ProductEdge{},
		PageInfo: &model.PageInfo{},
	}
	for _, p := range products {
		edge := model.ProductEdge{
			Cursor: paginationUtil.EncodeCursor(p.CreateTime.Format(time.RFC3339)),
			Node: &model.Product{
				ID:          p.ID.Hex(),
				ProductName: p.ProductName,
				Price:       p.Price,
				Stock:       p.Stock,
				StockUnit:   p.StockUnit,
				CreateTime:  p.CreateTime,
				UpdateTime:  p.UpdateTime,
			},
		}
		if direction == 1 {
			result.Edges = append(result.Edges, &edge)
		} else {
			result.Edges = append([]*model.ProductEdge{&edge}, result.Edges...)
		}
	}

	if len(products) > 0 {
		result.PageInfo.StartCursor = &result.Edges[0].Cursor
		result.PageInfo.EndCursor = &result.Edges[len(result.Edges)-1].Cursor

		countOpts := options.Count().SetLimit(1)
		nextCount, err := r.db.Collection(collectionName).CountDocuments(
			ctx,
			bson.M{"createTime": bson.M{"$gt": products[len(products)-1].CreateTime}},
			countOpts,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result.PageInfo.HasNextPage = nextCount > 0
		prevCount, err := r.db.Collection(collectionName).CountDocuments(
			ctx,
			bson.M{"createTime": bson.M{"$lt": products[0].CreateTime}},
			countOpts,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result.PageInfo.HasNextPage = prevCount > 0
	}

	return result, nil
}
