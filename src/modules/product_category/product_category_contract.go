package productcategory

import "github.com/google/uuid"

type getProductCategoryListReqQuery struct {
	Limit *int `query:"limit"`
	Page  *int `query:"page"`
}

type getProductCategoryDetailReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type addProductCategoryReq struct {
	Name *string `json:"name" validate:"required"`
}

type updateProductCategoryReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateProductCategoryReq struct {
	Name *string `json:"name"`
}

type deleteProductCategoryReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}
