package shoppingcart

import "github.com/google/uuid"

type getShoppingCartItemListReqQuery struct {
	Limit *int `query:"limit"`
	Page  *int `query:"page"`
}

type addShoppingCartItemReq struct {
	ProductID *uuid.UUID `json:"productId" validate:"required"`
	Amount    *int       `json:"amount" validate:"required"`
}

type updateShoppingCartItemReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateShoppingCartItemReq struct {
	Amount *int `json:"amount"`
}

type deleteShoppingCartItemReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}
