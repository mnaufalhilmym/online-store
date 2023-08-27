package shoppingcart

import (
	"github.com/google/uuid"
	"hilmy.dev/store/src/libs/db/pg"
	p "hilmy.dev/store/src/modules/product/product_entity"
	sc "hilmy.dev/store/src/modules/shopping_cart/shopping_cart_entity"
)

type searchOptions struct {
	byUserID *uuid.UUID
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

func (*Module) getShoppingCartItemListService(pagination *paginationOptions, search *searchOptions) (*[]*sc.ShoppingCartItemModel, *paginationQuery, error) {
	where := []pg.FindAllWhere{}
	limit := 0
	offset := 0

	if search != nil {
		if search.byUserID != nil {
			where = append(where, pg.FindAllWhere{
				Where: pg.Where{
					Query: "user_id = ?",
					Args:  []interface{}{search.byUserID},
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

	data, page, err := sc.ShoppingCartItemRepository().FindAll(&pg.FindAllOptions{
		Where:  &where,
		Limit:  &limit,
		Offset: &offset,
		Order:  &[]string{"created_at asc"},
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

func (*Module) getShoppingCartItemByProductIDService(userID *uuid.UUID, productID *uuid.UUID) (*sc.ShoppingCartItemModel, error) {
	return sc.ShoppingCartItemRepository().FindOne(&pg.FindOneOptions{
		Where: &[]pg.Where{
			{
				Query: "user_id = ? AND product_id = ?",
				Args:  []interface{}{userID, productID},
			},
		},
	})
}

func (*Module) addShoppingCartItemService(data *sc.ShoppingCartItemModel) (*sc.ShoppingCartItemModel, error) {
	return sc.ShoppingCartItemRepository().Create(data)
}

func (*Module) updateShoppingCartItemService(userID *uuid.UUID, id *uuid.UUID, data *sc.ShoppingCartItemModel) (*sc.ShoppingCartItemModel, error) {
	if _, err := sc.ShoppingCartItemRepository().Update(data, &pg.UpdateOptions{
		Where: &[]pg.Where{
			{
				Query: "user_id = ? AND id = ?",
				Args:  []interface{}{userID, id},
			},
		},
	}); err != nil {
		return nil, err
	}

	data, err := sc.ShoppingCartItemRepository().FindOne(&pg.FindOneOptions{
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

	return data, nil
}

func (*Module) deleteShoppingCartItemService(userID *uuid.UUID, id *uuid.UUID) error {
	return sc.ShoppingCartItemRepository().Destroy(&sc.ShoppingCartItemModel{}, &pg.DestroyOptions{
		Where: &[]pg.Where{
			{
				Query: "user_id = ? AND id = ?",
				Args:  []interface{}{userID, id},
			},
		},
		IsUnscoped: true,
	})
}

func (*Module) countProductService(id *uuid.UUID) (*int64, error) {
	return p.ProductRepository().Count(&pg.CountOptions{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}
