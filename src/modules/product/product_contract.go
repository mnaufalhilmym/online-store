package product

import "github.com/google/uuid"

type getProductListReqQuery struct {
	SearchByCategoryID *uuid.UUID `query:"category_id"`
	Limit              *int       `query:"limit"`
	Page               *int       `query:"page"`
}

type getProductDetailReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type addProductReq struct {
	CategoryID  *uuid.UUID `json:"category_id" validate:"required"`
	Title       *string    `json:"title" validate:"required"`
	Description *string    `json:"description" validate:"required"`
	Price       *int       `json:"price" validate:"required"`
}

type updateProductReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateProductReq struct {
	CategoryID  *uuid.UUID `json:"category_id"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Price       *int       `json:"price"`
}

type deleteProductReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}
