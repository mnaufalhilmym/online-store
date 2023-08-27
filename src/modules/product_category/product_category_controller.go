package productcategory

import (
	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/contract"
	"hilmy.dev/store/src/libs/db/pg"
	"hilmy.dev/store/src/libs/parser"
	acc "hilmy.dev/store/src/modules/account/account_entity"
	am "hilmy.dev/store/src/modules/auth/auth_middleware"
	"hilmy.dev/store/src/modules/log"
	pc "hilmy.dev/store/src/modules/product_category/product_category_entity"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/product-categories", m.getProductCategoryList)
	m.App.Get("/api/v1/product-category/:id", m.getProductCategoryDetail)
	m.App.Post("/api/v1/product-category", am.AuthGuard(acc.ROLE_ADMIN), m.addProductCategory)
	m.App.Patch("/api/v1/product-category/:id", am.AuthGuard(acc.ROLE_ADMIN), m.updateProductCategory)
	m.App.Delete("/api/v1/product-category/:id", am.AuthGuard(acc.ROLE_ADMIN), m.deleteProductCategory)
}

func (m *Module) getProductCategoryList(c *fiber.Ctx) error {
	query := new(getProductCategoryListReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	offset := 0
	if query.Page != nil && query.Limit != nil && *query.Page > 0 && *query.Limit > 0 {
		offset = (*query.Page - 1) * *query.Limit
	}

	productCategoryListData, page, err := m.getProductCategoryListService(&paginationOptions{
		limit:  query.Limit,
		offset: &offset,
	})
	if err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), true)
		return c.Status(fiber.StatusInternalServerError).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrInternalServerError.Error(),
				Message: err.Error(),
			},
		})
	}

	log.SaveLogService(c.OriginalURL(), "Ok", false)
	return c.Status(fiber.StatusOK).JSON(&contract.Response{
		Pagination: &contract.Pagination{
			Limit: page.limit,
			Count: page.count,
			Page:  query.Page,
			Total: page.total,
		},
		Data: productCategoryListData,
	})
}

func (m *Module) getProductCategoryDetail(c *fiber.Ctx) error {
	param := new(getProductCategoryDetailReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	productCategoryDetailData, err := m.getProductCategoryDetailService(param.ID)
	if err != nil {
		status := fiber.StatusInternalServerError
		statusString := fiber.ErrInternalServerError.Error()
		printStack := true
		if pg.IsErrRecordNotFound(err) {
			status = fiber.StatusNotFound
			statusString = fiber.ErrNotFound.Error()
			printStack = false
		}
		log.SaveLogService(c.OriginalURL(), err.Error(), printStack)
		return c.Status(status).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  statusString,
				Message: err.Error(),
			},
		})
	}

	log.SaveLogService(c.OriginalURL(), "Ok", false)
	return c.Status(fiber.StatusOK).JSON(&contract.Response{
		Data: productCategoryDetailData,
	})
}

func (m *Module) addProductCategory(c *fiber.Ctx) error {
	req := new(addProductCategoryReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	productCategoryDetailData, err := m.addProductCategoryService(&pc.ProductCategoryModel{
		Name: req.Name,
	})
	if err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), true)
		return c.Status(fiber.StatusInternalServerError).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrInternalServerError.Error(),
				Message: err.Error(),
			},
		})
	}

	log.SaveLogService(c.OriginalURL(), "Ok", false)
	return c.Status(fiber.StatusCreated).JSON(&contract.Response{
		Data: productCategoryDetailData,
	})
}

func (m *Module) updateProductCategory(c *fiber.Ctx) error {
	param := new(updateProductCategoryReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	req := new(updateProductCategoryReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	productCategoryDetailData, err := m.updateProductCategoryService(param.ID, &pc.ProductCategoryModel{
		Name: req.Name,
	})
	if err != nil {
		status := fiber.StatusInternalServerError
		statusString := fiber.ErrInternalServerError.Error()
		printStack := true
		if pg.IsErrRecordNotFound(err) {
			status = fiber.StatusNotFound
			statusString = fiber.ErrNotFound.Error()
			printStack = false
		}
		log.SaveLogService(c.OriginalURL(), err.Error(), printStack)
		return c.Status(status).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  statusString,
				Message: err.Error(),
			},
		})
	}

	log.SaveLogService(c.OriginalURL(), "Ok", false)
	return c.Status(fiber.StatusOK).JSON(&contract.Response{
		Data: productCategoryDetailData,
	})
}

func (m *Module) deleteProductCategory(c *fiber.Ctx) error {
	param := new(deleteProductCategoryReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	if err := m.deleteProductCategoryService(param.ID); err != nil {
		status := fiber.StatusInternalServerError
		statusString := fiber.ErrInternalServerError.Error()
		printStack := true
		if pg.IsErrRecordNotFound(err) {
			status = fiber.StatusNotFound
			statusString = fiber.ErrNotFound.Error()
			printStack = false
		}
		log.SaveLogService(c.OriginalURL(), err.Error(), printStack)
		return c.Status(status).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  statusString,
				Message: err.Error(),
			},
		})
	}

	log.SaveLogService(c.OriginalURL(), "Ok", false)
	return c.Status(fiber.StatusOK).JSON(&contract.Response{
		Data: param.ID,
	})
}
