package transaction

import (
	"errors"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/contracts"
	"hilmy.dev/store/src/libs/db/pg"
	"hilmy.dev/store/src/libs/parser"
	acc "hilmy.dev/store/src/modules/account/account_entity"
	a "hilmy.dev/store/src/modules/auth/auth_entity"
	am "hilmy.dev/store/src/modules/auth/auth_middleware"
	"hilmy.dev/store/src/modules/log"
	sc "hilmy.dev/store/src/modules/shopping_cart/shopping_cart_entity"
	t "hilmy.dev/store/src/modules/transaction/transaction_entity"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/transactions", am.AuthGuard(acc.ROLE_USER), m.getTransactionList)
	m.App.Get("/api/v1/transaction/:id", am.AuthGuard(acc.ROLE_USER), m.getTransactionDetail)
	m.App.Post("/api/v1/transaction", am.AuthGuard(acc.ROLE_USER), m.addTransaction)
	m.App.Post("/api/v1/transaction/:id/pay", am.AuthGuard(acc.ROLE_USER), m.payTransaction)
	m.App.Post("/api/v1/transaction/:id/cancel", am.AuthGuard(acc.ROLE_USER), m.cancelTransaction)
}

func (m *Module) getTransactionList(c *fiber.Ctx) error {
	token := new(a.JWTPayload)
	if err := parser.ParseReqBearerToken(c, token); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusUnauthorized).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrUnauthorized.Error(),
				Message: err.Error(),
			},
		})
	}

	query := new(getTransactionListReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	offset := 0
	if query.Page != nil && query.Limit != nil && *query.Page > 0 && *query.Limit > 0 {
		offset = (*query.Page - 1) * *query.Limit
	}

	transactionDataList, page, err := m.getTransactionListService(&paginationOptions{
		limit:  query.Limit,
		offset: &offset,
	}, &searchOptions{
		byUserID:            token.ID,
		byTransactionStatus: query.SearchByStatus,
	})
	if err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), true)
		return c.Status(fiber.StatusInternalServerError).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrInternalServerError.Error(),
				Message: err.Error(),
			},
		})
	}

	log.SaveLogService(c.OriginalURL(), "Ok", false)
	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Pagination: &contracts.Pagination{
			Limit: page.limit,
			Count: page.count,
			Page:  query.Page,
			Total: page.total,
		},
		Data: transactionDataList,
	})
}

func (m *Module) getTransactionDetail(c *fiber.Ctx) error {
	token := new(a.JWTPayload)
	if err := parser.ParseReqBearerToken(c, token); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusUnauthorized).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrUnauthorized.Error(),
				Message: err.Error(),
			},
		})
	}

	param := new(getTransactionDetailReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	transactionDetailData, err := m.getTransactionDetailService(token.ID, param.ID)
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
		return c.Status(status).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  statusString,
				Message: err.Error(),
			},
		})
	}

	log.SaveLogService(c.OriginalURL(), "Ok", false)
	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: transactionDetailData,
	})
}

func (m *Module) addTransaction(c *fiber.Ctx) error {
	token := new(a.JWTPayload)
	if err := parser.ParseReqBearerToken(c, token); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusUnauthorized).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrUnauthorized.Error(),
				Message: err.Error(),
			},
		})
	}

	req := new(addTransactionReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	transactionPrice := 0
	shoppingCartItemListData := []*sc.ShoppingCartItemModel{}
	for i := range *req.ShoppingCartItemIDs {
		shoppingCartItemDetailData, err := m.getShoppingCartItemDetailService((*req.ShoppingCartItemIDs)[i])
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
			return c.Status(status).JSON(&contracts.Response{
				Error: &contracts.Error{
					Status:  statusString,
					Message: err.Error(),
				},
			})
		}
		productDetailData, err := m.getProductDetailService(shoppingCartItemDetailData.ProductID)
		if err != nil {
			if pg.IsErrRecordNotFound(err) {
				if err := m.deleteShoppingCartItemDetailService((*req.ShoppingCartItemIDs)[i]); err != nil {
					log.SaveLogService(c.OriginalURL(), err.Error(), true)
					return c.Status(fiber.StatusInternalServerError).JSON(&contracts.Response{
						Error: &contracts.Error{
							Status:  fiber.ErrInternalServerError.Error(),
							Message: err.Error(),
						},
					})
				}
			}
			log.SaveLogService(c.OriginalURL(), err.Error(), true)
			return c.Status(fiber.StatusInternalServerError).JSON(&contracts.Response{
				Error: &contracts.Error{
					Status:  fiber.ErrInternalServerError.Error(),
					Message: err.Error(),
				},
			})
		}
		transactionPrice += *shoppingCartItemDetailData.Amount * *productDetailData.Price
		shoppingCartItemListData = append(shoppingCartItemListData, shoppingCartItemDetailData)
	}

	dataBytes, err := sonic.Marshal(shoppingCartItemListData)
	if err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), true)
		return c.Status(fiber.StatusInternalServerError).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrInternalServerError.Error(),
				Message: err.Error(),
			},
		})
	}

	transactionStatus := t.STATUS_WAITING_PAYMENT
	transactionDetailData := t.TransactionModel{
		UserID: token.ID,
		Status: &transactionStatus,
		Price:  &transactionPrice,
		Data:   dataBytes,
	}
	if err := m.addTransactionService(&transactionDetailData, req.ShoppingCartItemIDs); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), true)
		return c.Status(fiber.StatusInternalServerError).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrInternalServerError.Error(),
				Message: err.Error(),
			},
		})
	}

	log.SaveLogService(c.OriginalURL(), "Ok", false)
	return c.Status(fiber.StatusCreated).JSON(&contracts.Response{
		Data: &transactionDetailData,
	})
}

func (m *Module) payTransaction(c *fiber.Ctx) error {
	token := new(a.JWTPayload)
	if err := parser.ParseReqBearerToken(c, token); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusUnauthorized).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrUnauthorized.Error(),
				Message: err.Error(),
			},
		})
	}

	param := new(payTransactionReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	transactionDetailData, err := m.getTransactionDetailService(token.ID, param.ID)
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
		return c.Status(status).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  statusString,
				Message: err.Error(),
			},
		})
	}
	if *transactionDetailData.Status != t.STATUS_WAITING_PAYMENT {
		err := fmt.Errorf("cannot pay for a transaction that is not in %s status", t.STATUS_WAITING_PAYMENT)
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	balanceDetailData, err := m.getBalanceByUserIDService(token.ID)
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
		return c.Status(status).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  statusString,
				Message: err.Error(),
			},
		})
	}

	if *balanceDetailData.Amount < *transactionDetailData.Price {
		err := errors.New("insufficient balance")
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	transactionStatus := t.STATUS_COMPLETED
	transactionDetailData.Status = &transactionStatus
	*balanceDetailData.Amount -= *transactionDetailData.Price
	if err := m.payTransactionService(token.ID, transactionDetailData.ID, transactionDetailData, balanceDetailData.ID, balanceDetailData); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), true)
		return c.Status(fiber.StatusInternalServerError).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrInternalServerError.Error(),
				Message: err.Error(),
			},
		})
	}

	log.SaveLogService(c.OriginalURL(), "Ok", false)
	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: transactionDetailData,
	})
}

func (m *Module) cancelTransaction(c *fiber.Ctx) error {
	token := new(a.JWTPayload)
	if err := parser.ParseReqBearerToken(c, token); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusUnauthorized).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrUnauthorized.Error(),
				Message: err.Error(),
			},
		})
	}

	param := new(cancelTransactionReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	transactionDetailData, err := m.getTransactionDetailService(token.ID, param.ID)
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
		return c.Status(status).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  statusString,
				Message: err.Error(),
			},
		})
	}

	if *transactionDetailData.Status == t.STATUS_COMPLETED {
		err := errors.New("cannot cancel a transaction that has already been completed")
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	transactionDetailData, err = m.cancelTransactionService(token.ID, param.ID)
	if err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), true)
		return c.Status(fiber.StatusInternalServerError).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrInternalServerError.Error(),
				Message: err.Error(),
			},
		})
	}

	log.SaveLogService(c.OriginalURL(), "Ok", false)
	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: transactionDetailData,
	})
}
