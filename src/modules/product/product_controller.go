package product

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/contract"
	"hilmy.dev/store/src/libs/db/pg"
	"hilmy.dev/store/src/libs/parser"
	acc "hilmy.dev/store/src/modules/account/account_entity"
	am "hilmy.dev/store/src/modules/auth/auth_middleware"
	"hilmy.dev/store/src/modules/log"
	p "hilmy.dev/store/src/modules/product/product_entity"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/products", m.getProductList)
	m.App.Get("/api/v1/product/:id", m.getProductDetail)
	m.App.Post("/api/v1/product", am.AuthGuard(acc.ROLE_ADMIN), m.addProduct)
	m.App.Patch("/api/v1/product/:id", am.AuthGuard(acc.ROLE_ADMIN), m.updateProduct)
	m.App.Delete("/api/v1/product/:id", am.AuthGuard(acc.ROLE_ADMIN), m.deleteProduct)
}

func (m *Module) getProductList(c *fiber.Ctx) error {
	query := new(getProductListReqQuery)
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

	productListData, page, err := m.getProductListService(&paginationOptions{
		limit:  query.Limit,
		offset: &offset,
	}, &searchOptions{
		byCategoryID: query.SearchByCategoryID,
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
		Data: productListData,
	})
}

func (m *Module) getProductDetail(c *fiber.Ctx) error {
	param := new(getProductDetailReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	productDetailData, err := m.getProductDetailService(param.ID)
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
		Data: productDetailData,
	})
}

func (m *Module) addProduct(c *fiber.Ctx) error {
	req := new(addProductReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	pcCount, err := m.getProductCategoryCountByProductID(req.CategoryID)
	if err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), true)
		return c.Status(fiber.StatusInternalServerError).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrInternalServerError.Error(),
				Message: err.Error(),
			},
		})
	}
	if *pcCount == 0 {
		err := errors.New("category does not exist")
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	productDetailData, err := m.addProductService(&p.ProductModel{
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
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
		Data: productDetailData,
	})
}

func (m *Module) updateProduct(c *fiber.Ctx) error {
	param := new(updateProductReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	req := new(updateProductReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	pcCount, err := m.getProductCategoryCountByProductID(req.CategoryID)
	if err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), true)
		return c.Status(fiber.StatusInternalServerError).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrInternalServerError.Error(),
				Message: err.Error(),
			},
		})
	}
	if *pcCount == 0 {
		err := errors.New("category does not exist")
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	productDetailData, err := m.updateProductService(param.ID, &p.ProductModel{
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
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
		Data: productDetailData,
	})
}

func (m *Module) deleteProduct(c *fiber.Ctx) error {
	param := new(deleteProductReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	if err := m.deleteProductService(param.ID); err != nil {
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
