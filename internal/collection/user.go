package collection

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Email        string             `bson:"email"`
	Fullname     string             `bson:"fullname"`
	Role         string             `bson:"role"`
	PasswordHash string             `bson:"passwordHash"`
}
