package collection

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	ProductName *string            `bson:"productName" json:"productName"`
	Price       *float64           `bson:"price" json:"price"`
	Stock       *float64           `bson:"stock" json:"stock"`
	StockUnit   *string            `bson:"stockUnit" json:"stockUnit"`
	CreateTime  *time.Time         `bson:"createTime" json:"createTime"`
	UpdateTime  *time.Time         `bson:"updateTime" json:"updateTime"`
}

func (Product) CollectionName() string {
	return "product"
}
