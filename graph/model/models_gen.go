// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type CreateProductInput struct {
	ProductName string  `json:"productName"`
	Price       float64 `json:"price"`
	Stock       float64 `json:"stock"`
	StockUnit   string  `json:"stockUnit"`
}

type CreateUserInput struct {
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Password string `json:"password"`
}

type Mutation struct {
}

type PageInfo struct {
	HasPrevPage bool    `json:"hasPrevPage"`
	HasNextPage bool    `json:"hasNextPage"`
	StartCursor *string `json:"startCursor,omitempty"`
	EndCursor   *string `json:"endCursor,omitempty"`
}

type Product struct {
	ID          string     `json:"id"`
	ProductName *string    `json:"productName,omitempty"`
	Price       *float64   `json:"price,omitempty"`
	Stock       *float64   `json:"stock,omitempty"`
	StockUnit   *string    `json:"stockUnit,omitempty"`
	CreateTime  *time.Time `json:"createTime,omitempty"`
	UpdateTime  *time.Time `json:"updateTime,omitempty"`
}

type ProductConnection struct {
	Edges    []*ProductEdge `json:"edges"`
	PageInfo *PageInfo      `json:"pageInfo"`
}

type ProductConnectionArgs struct {
	First    *int32  `json:"first,omitempty"`
	Last     *int32  `json:"last,omitempty"`
	Before   *string `json:"before,omitempty"`
	After    *string `json:"after,omitempty"`
	Email    *string `json:"email,omitempty"`
	FullName *string `json:"fullName,omitempty"`
	Role     *string `json:"role,omitempty"`
}

type ProductEdge struct {
	Cursor string   `json:"cursor"`
	Node   *Product `json:"node,omitempty"`
}

type Query struct {
}

type UpdateProductInput struct {
	ID          string   `json:"id"`
	ProductName *string  `json:"productName,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Stock       *float64 `json:"stock,omitempty"`
	StockUnit   *string  `json:"stockUnit,omitempty"`
}

type UpdateUserInput struct {
	ID       string  `json:"id"`
	Email    *string `json:"email,omitempty"`
	FullName *string `json:"fullName,omitempty"`
	Role     *string `json:"role,omitempty"`
}

type User struct {
	ID         string     `json:"id"`
	Email      *string    `json:"email,omitempty"`
	FullName   *string    `json:"fullName,omitempty"`
	Role       *string    `json:"role,omitempty"`
	CreateTime *time.Time `json:"createTime,omitempty"`
	UpdateTime *time.Time `json:"updateTime,omitempty"`
}

type UserConnection struct {
	Edges    []*UserEdge `json:"edges"`
	PageInfo *PageInfo   `json:"pageInfo"`
}

type UserConnectionArgs struct {
	First    *int32  `json:"first,omitempty"`
	Last     *int32  `json:"last,omitempty"`
	Before   *string `json:"before,omitempty"`
	After    *string `json:"after,omitempty"`
	Email    *string `json:"email,omitempty"`
	FullName *string `json:"fullName,omitempty"`
	Role     *string `json:"role,omitempty"`
}

type UserEdge struct {
	Cursor string `json:"cursor"`
	Node   *User  `json:"node,omitempty"`
}
