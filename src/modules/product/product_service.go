package product

import (
	"github.com/google/uuid"
	"hilmy.dev/store/src/libs/db/pg"
	p "hilmy.dev/store/src/modules/product/product_entity"
	pc "hilmy.dev/store/src/modules/product_category/product_category_entity"
)

type searchOptions struct {
	byCategoryID *uuid.UUID
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

func (*Module) getProductListService(pagination *paginationOptions, search *searchOptions) (*[]*p.ProductModel, *paginationQuery, error) {
	where := []pg.FindAllWhere{}
	limit := 0
	offset := 0

	if search != nil {
		if search.byCategoryID != nil && len(*search.byCategoryID) > 0 {
			where = append(where, pg.FindAllWhere{
				Where: pg.Where{
					Query: "category_id = ?",
					Args:  []interface{}{search.byCategoryID},
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

	data, page, err := p.ProductRepository().FindAll(&pg.FindAllOptions{
		Where:  &where,
		Limit:  &limit,
		Offset: &offset,
		Order:  &[]string{"title asc"},
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

func (*Module) getProductDetailService(id *uuid.UUID) (*p.ProductModel, error) {
	return p.ProductRepository().FindOne(&pg.FindOneOptions{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
		IncludeTables: &[]pg.IncludeTables{
			{
				Query: "Product",
			},
		},
		IsUnscoped: true,
	})
}

func (*Module) addProductService(data *p.ProductModel) (*p.ProductModel, error) {
	return p.ProductRepository().Create(data)
}

func (*Module) updateProductService(id *uuid.UUID, data *p.ProductModel) (*p.ProductModel, error) {
	if _, err := p.ProductRepository().Update(data, &pg.UpdateOptions{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	}); err != nil {
		return nil, err
	}

	data, err := p.ProductRepository().FindOne(&pg.FindOneOptions{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (*Module) deleteProductService(id *uuid.UUID) error {
	return p.ProductRepository().Destroy(&p.ProductModel{
		Model: pg.Model{
			ID: id,
		},
	})
}

func (*Module) getProductCategoryCountByProductID(id *uuid.UUID) (*int64, error) {
	return pc.ProductCategoryRepository().Count(&pg.CountOptions{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}
