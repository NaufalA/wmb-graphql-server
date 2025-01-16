package collection

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `bson:"_id"`
	ProductName *string            `bson:"productName"`
	Price       *float64           `bson:"price"`
	Stock       *float64           `bson:"stock"`
	StockUnit   *string            `bson:"stockUnit"`
}

func (Product) CollectionName() string {
	return "product"
}
