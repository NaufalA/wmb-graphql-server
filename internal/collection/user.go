package collection

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Email        string             `bson:"email" json:"email"`
	FullName     string             `bson:"fullName" json:"fullName"`
	Role         string             `bson:"role" json:"role"`
	PasswordHash string             `bson:"passwordHash" json:"passwordHash"`
	CreateTime   *time.Time         `bson:"createTime" json:"createTime"`
	UpdateTime   *time.Time         `bson:"updateTime" json:"updateTime"`
}

func (User) CollectionName() string {
	return "user"
}
