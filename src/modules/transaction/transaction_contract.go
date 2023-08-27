package transaction

import (
	"github.com/google/uuid"
	t "hilmy.dev/store/src/modules/transaction/transaction_entity"
)

type getTransactionListReqQuery struct {
	SearchByStatus *t.TransactionStatus `query:"status"`
	Limit          *int                 `query:"limit"`
	Page           *int                 `query:"page"`
}

type getTransactionDetailReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type addTransactionReq struct {
	ShoppingCartItemIDs *[]*uuid.UUID `json:"shoppingCartItemIds" validate:"required,unique"`
}

type payTransactionReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type cancelTransactionReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}
