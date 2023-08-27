package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Where = bson.E

type FindAllWhere struct {
	Where          Where
	IncludeInCount bool
}

type Order = bson.E

type CountOptions struct {
	Where *[]Where
}

type FindOneOptions struct {
	Where *[]Where
	Order *[]Order
}

type FindAllOptions struct {
	Where   *[]FindAllWhere
	Order   *[]Order
	Limit   *int
	Offset  *int
	AfterID *primitive.ObjectID
}

type UpdateOptions struct {
	Where *[]Where
}

type DestroyOptions struct {
	Where *[]Where
}

type Pagination struct {
	Limit int
	Count int
	Total int
}
