package transaction

import (
	"github.com/google/uuid"
	"hilmy.dev/store/src/libs/db/pg"
	b "hilmy.dev/store/src/modules/balance/balance_entity"
	p "hilmy.dev/store/src/modules/product/product_entity"
	sc "hilmy.dev/store/src/modules/shopping_cart/shopping_cart_entity"
	t "hilmy.dev/store/src/modules/transaction/transaction_entity"
)

type searchOptions struct {
	byUserID            *uuid.UUID
	byTransactionStatus *t.TransactionStatus
}

type paginationOptions struct {
	limit  *int
	offset *int
}

type paginationQuery struct {
	limit *int
	count *int
	total *int
}

func (*Module) getTransactionListService(pagination *paginationOptions, search *searchOptions) (*[]*t.TransactionModel, *paginationQuery, error) {
	where := []pg.FindAllWhere{}
	limit := 0
	offset := 0

	if search != nil {
		if search.byUserID != nil && len(*search.byUserID) > 0 {
			where = append(where, pg.FindAllWhere{
				Where: pg.Where{
					Query: "user_id = ?",
					Args:  []interface{}{search.byUserID},
				},
				IncludeInCount: true,
			})
		}
		if search.byTransactionStatus != nil && len(*search.byTransactionStatus) > 0 {
			where = append(where, pg.FindAllWhere{
				Where: pg.Where{
					Query: "status = ?",
					Args:  []interface{}{search.byTransactionStatus},
				},
				IncludeInCount: true,
			})
		}
	}

	if pagination != nil {
		if pagination.limit != nil && *pagination.limit > 0 {
			limit = *pagination.limit
		}
		if pagination.offset != nil && *pagination.offset > 0 {
			offset = *pagination.offset
		}
	}

	data, page, err := t.TransactionRepository().FindAll(&pg.FindAllOptions{
		Where:  &where,
		Limit:  &limit,
		Offset: &offset,
		Order:  &[]string{"created_at desc"},
	})
	if err != nil {
		return nil, nil, err
	}

	return data, &paginationQuery{
		limit: &page.Limit,
		count: &page.Count,
		total: &page.Total,
	}, nil
}

func (*Module) getTransactionDetailService(userID *uuid.UUID, id *uuid.UUID) (*t.TransactionModel, error) {
	return t.TransactionRepository().FindOne(&pg.FindOneOptions{
		Where: &[]pg.Where{
			{
				Query: "user_id = ? AND id = ?",
				Args:  []interface{}{userID, id},
			},
		},
	})
}

func (m *Module) addTransactionService(data *t.TransactionModel, shoppingCartItemIDs *[]*uuid.UUID) error {
	return pg.Transaction(m.DB, func(tx *pg.DB) *pg.DB {
		txz := t.TransactionRepository().CreateTx(tx, data)
		return txz
	}, func(tx *pg.DB) *pg.DB {
		data := []*sc.ShoppingCartItemModel{}
		for i := range *shoppingCartItemIDs {
			data = append(data, &sc.ShoppingCartItemModel{
				Model: pg.Model{
					ID: (*shoppingCartItemIDs)[i],
				},
			})
		}
		txz := sc.ShoppingCartItemRepository().BulkDestroyTx(tx, &data)
		return txz
	})
}

func (m *Module) payTransactionService(userID *uuid.UUID, tID *uuid.UUID, tData *t.TransactionModel, bID *uuid.UUID, bData *b.BalanceModel) error {
	return pg.Transaction(m.DB, func(tx *pg.DB) *pg.DB {
		return t.TransactionRepository().UpdateTx(tx, tData, &pg.UpdateOptions{
			Where: &[]pg.Where{
				{
					Query: "user_id = ? AND id = ?",
					Args:  []interface{}{userID, tID},
				},
			},
		})
	}, func(tx *pg.DB) *pg.DB {
		return b.BalanceRepository().UpdateTx(tx, bData, &pg.UpdateOptions{
			Where: &[]pg.Where{
				{
					Query: "id = ?",
					Args:  []interface{}{bID},
				},
			},
		})
	})
}

func (*Module) cancelTransactionService(userID *uuid.UUID, id *uuid.UUID) (*t.TransactionModel, error) {
	updateTransactionStatus := t.STATUS_CANCELLED
	data, err := t.TransactionRepository().Update(&t.TransactionModel{
		Status: &updateTransactionStatus,
	}, &pg.UpdateOptions{
		Where: &[]pg.Where{
			{
				Query: "user_id = ? AND id = ?",
				Args:  []interface{}{userID, id},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	data.ID = id
	return data, nil
}

func (*Module) getShoppingCartItemDetailService(id *uuid.UUID) (*sc.ShoppingCartItemModel, error) {
	return sc.ShoppingCartItemRepository().FindOne(&pg.FindOneOptions{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (*Module) deleteShoppingCartItemDetailService(id *uuid.UUID) error {
	return sc.ShoppingCartItemRepository().Destroy(&sc.ShoppingCartItemModel{
		Model: pg.Model{
			ID: id,
		},
	})
}

func (*Module) getProductDetailService(id *uuid.UUID) (*p.ProductModel, error) {
	return p.ProductRepository().FindOne(&pg.FindOneOptions{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (*Module) getBalanceByUserIDService(userID *uuid.UUID) (*b.BalanceModel, error) {
	return b.BalanceRepository().FindOne(&pg.FindOneOptions{
		Where: &[]pg.Where{
			{
				Query: "user_id = ?",
				Args:  []interface{}{userID},
			},
		},
	})
}
