package balance

import (
	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/contract"
	"hilmy.dev/store/src/libs/db/pg"
	"hilmy.dev/store/src/libs/parser"
	acc "hilmy.dev/store/src/modules/account/account_entity"
	a "hilmy.dev/store/src/modules/auth/auth_entity"
	am "hilmy.dev/store/src/modules/auth/auth_middleware"
	balanceentity "hilmy.dev/store/src/modules/balance/balance_entity"
	"hilmy.dev/store/src/modules/log"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/balance", am.AuthGuard(acc.ROLE_USER), m.getBalance)
	m.App.Post("/api/v1/balance/add", am.AuthGuard(acc.ROLE_USER), m.addBalance)
}

func (m *Module) getBalance(c *fiber.Ctx) error {
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
		return c.Status(status).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  statusString,
				Message: err.Error(),
			},
		})
	}

	log.SaveLogService(c.OriginalURL(), "Ok", false)
	return c.Status(fiber.StatusOK).JSON(&contract.Response{
		Data: balanceDetailData,
	})
}

func (m *Module) addBalance(c *fiber.Ctx) error {
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

	req := new(addBalanceReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contract.Response{
			Error: &contract.Error{
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
		return c.Status(status).JSON(&contract.Response{
			Error: &contract.Error{
				Status:  statusString,
				Message: err.Error(),
			},
		})
	}

	balanceAmount := *balanceDetailData.Amount + *req.Amount
	balanceDetailData, err = m.updateBalanceByUserIDService(token.ID, &balanceentity.BalanceModel{
		Amount: &balanceAmount,
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
		Data: balanceDetailData,
	})
}
