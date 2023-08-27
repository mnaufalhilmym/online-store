package shoppingcart

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/contract"
	"hilmy.dev/store/src/libs/db/pg"
	"hilmy.dev/store/src/libs/parser"
	acc "hilmy.dev/store/src/modules/account/account_entity"
	a "hilmy.dev/store/src/modules/auth/auth_entity"
	am "hilmy.dev/store/src/modules/auth/auth_middleware"
	"hilmy.dev/store/src/modules/log"
	sc "hilmy.dev/store/src/modules/shopping_cart/shopping_cart_entity"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/shopping-cart-items", am.AuthGuard(acc.ROLE_USER), m.getShoppingCartItemList)
	m.App.Post("/api/v1/shopping-cart-item", am.AuthGuard(acc.ROLE_USER), m.addShoppingCartItem)
	m.App.Patch("/api/v1/shopping-cart-item/:id", am.AuthGuard(acc.ROLE_USER), m.updateShoppingCartItem)
	m.App.Delete("/api/v1/shopping-cart-item/:id", am.AuthGuard(acc.ROLE_USER), m.deleteShoppingCartItem)
}

func (m *Module) getShoppingCartItemList(c *fiber.Ctx) error {
	token := new(a.JWTPayload)
	if err := parser.ParseReqBearerToken(c, token); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusUnauthorized).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrUnauthorized.Error(),
				Message: err.Error(),
			},
		})
	}

	query := new(getShoppingCartItemListReqQuery)
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

	shoppingCartItemListData, page, err := m.getShoppingCartItemListService(&paginationOptions{
		limit:  query.Limit,
		offset: &offset,
	}, &searchOptions{
		byUserID: token.ID,
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
		Data: shoppingCartItemListData,
	})
}

func (m *Module) addShoppingCartItem(c *fiber.Ctx) error {
	token := new(a.JWTPayload)
	if err := parser.ParseReqBearerToken(c, token); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusUnauthorized).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrUnauthorized.Error(),
				Message: err.Error(),
			},
		})
	}

	req := new(addShoppingCartItemReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	productCount, err := m.countProductService(req.ProductID)
	if err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), true)
		return c.Status(fiber.StatusInternalServerError).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrInternalServerError.Error(),
				Message: err.Error(),
			},
		})
	}
	if *productCount == 0 {
		err := errors.New("unregistered product")
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	shoppingCartItemDetailData, err := m.getShoppingCartItemByProductIDService(token.ID, req.ProductID)
	if err != nil {
		if pg.IsErrRecordNotFound(err) {
			_shoppingCartItemDetailData, err := m.addShoppingCartItemService(&sc.ShoppingCartItemModel{
				UserID:    token.ID,
				ProductID: req.ProductID,
				Amount:    req.Amount,
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
			shoppingCartItemDetailData = _shoppingCartItemDetailData
		} else {
			log.SaveLogService(c.OriginalURL(), err.Error(), true)
			return c.Status(fiber.StatusInternalServerError).JSON(&contract.Response{
				Error: &contract.Error{
					Status:  fiber.ErrInternalServerError.Error(),
					Message: err.Error(),
				},
			})
		}
	} else if shoppingCartItemDetailData != nil {
		*shoppingCartItemDetailData.Amount += *req.Amount
		_shoppingCartItemDetailData, err := m.updateShoppingCartItemService(token.ID, shoppingCartItemDetailData.ID, &sc.ShoppingCartItemModel{
			Amount: shoppingCartItemDetailData.Amount,
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
		shoppingCartItemDetailData = _shoppingCartItemDetailData
	} else {
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
		Data: shoppingCartItemDetailData,
	})
}

func (m *Module) updateShoppingCartItem(c *fiber.Ctx) error {
	token := new(a.JWTPayload)
	if err := parser.ParseReqBearerToken(c, token); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusUnauthorized).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrUnauthorized.Error(),
				Message: err.Error(),
			},
		})
	}

	param := new(updateShoppingCartItemReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	req := new(updateShoppingCartItemReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	var shoppingCartItemDetailData *sc.ShoppingCartItemModel
	var err error
	if req.Amount != nil {
		if *req.Amount > 0 {
			shoppingCartItemDetailData, err = m.updateShoppingCartItemService(token.ID, param.ID, &sc.ShoppingCartItemModel{
				Amount: req.Amount,
			})
		} else {
			err = m.deleteShoppingCartItemService(token.ID, param.ID)
		}
	}
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
		Data: func() interface{} {
			if shoppingCartItemDetailData != nil {
				return shoppingCartItemDetailData
			} else {
				return param.ID
			}
		}(),
	})
}

func (m *Module) deleteShoppingCartItem(c *fiber.Ctx) error {
	token := new(a.JWTPayload)
	if err := parser.ParseReqBearerToken(c, token); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusUnauthorized).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrUnauthorized.Error(),
				Message: err.Error(),
			},
		})
	}

	param := new(deleteShoppingCartItemReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	if err := m.deleteShoppingCartItemService(token.ID, param.ID); err != nil {
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
