package pg

import "github.com/google/uuid"

type Where struct {
	Query string
	Args  []interface{}
}

type FindAllWhere struct {
	Where          Where
	IncludeInCount bool
}

type IncludeTables struct {
	Query string
	Args  []interface{}
}

type CountOptions struct {
	Where      *[]Where
	IsUnscoped bool
}

type FindOneOptions struct {
	Where         *[]Where
	Order         *[]string
	IncludeTables *[]IncludeTables
	IsUnscoped    bool
}

type FindAllOptions struct {
	Where         *[]FindAllWhere
	Order         *[]string
	Limit         *int
	Offset        *int
	AfterID       *uuid.UUID
	IncludeTables *[]IncludeTables
	IsUnscoped    bool
}

type CreateOptions struct {
	IsUpsert bool
}

type UpdateOptions struct {
	Where      *[]Where
	IsUnscoped bool
}

type ReplaceOptions struct {
	Where      *[]Where
	IsUnscoped bool
}

type DestroyOptions struct {
	Where      *[]Where
	IsUnscoped bool
}

type Pagination struct {
	Limit int
	Count int
	Total int
}
