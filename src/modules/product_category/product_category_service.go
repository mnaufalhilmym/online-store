package productcategory

import (
	"github.com/google/uuid"
	"hilmy.dev/store/src/libs/db/pg"
	pc "hilmy.dev/store/src/modules/product_category/product_category_entity"
)

type paginationOptions struct {
	limit  *int
	offset *int
}

type paginationQuery struct {
	limit *int
	count *int
	total *int
}

func (*Module) getProductCategoryListService(pagination *paginationOptions) (*[]*pc.ProductCategoryModel, *paginationQuery, error) {
	limit := 0
	offset := 0

	if pagination != nil {
		if pagination.limit != nil && *pagination.limit > 0 {
			limit = *pagination.limit
		}
		if pagination.offset != nil && *pagination.offset > 0 {
			offset = *pagination.offset
		}
	}

	data, page, err := pc.ProductCategoryRepository().FindAll(&pg.FindAllOptions{
		Limit:  &limit,
		Offset: &offset,
		Order:  &[]string{"name asc"},
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

func (*Module) getProductCategoryDetailService(id *uuid.UUID) (*pc.ProductCategoryModel, error) {
	return pc.ProductCategoryRepository().FindOne(&pg.FindOneOptions{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
		IsUnscoped: true,
	})
}

func (*Module) addProductCategoryService(data *pc.ProductCategoryModel) (*pc.ProductCategoryModel, error) {
	return pc.ProductCategoryRepository().Create(data)
}

func (*Module) updateProductCategoryService(id *uuid.UUID, data *pc.ProductCategoryModel) (*pc.ProductCategoryModel, error) {
	if _, err := pc.ProductCategoryRepository().Update(data, &pg.UpdateOptions{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	}); err != nil {
		return nil, err
	}

	data, err := pc.ProductCategoryRepository().FindOne(&pg.FindOneOptions{
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

func (*Module) deleteProductCategoryService(id *uuid.UUID) error {
	return pc.ProductCategoryRepository().Destroy(&pc.ProductCategoryModel{
		Model: pg.Model{
			ID: id,
		},
	})
}
