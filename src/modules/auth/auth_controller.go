package auth

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/contracts"
	"hilmy.dev/store/src/libs/db/pg"
	"hilmy.dev/store/src/libs/hash/argon2"
	"hilmy.dev/store/src/libs/jwx/jwt"
	"hilmy.dev/store/src/libs/parser"
	acc "hilmy.dev/store/src/modules/account/account_entity"
	a "hilmy.dev/store/src/modules/auth/auth_entity"
	balanceentity "hilmy.dev/store/src/modules/balance/balance_entity"
	"hilmy.dev/store/src/modules/log"
)

func (m *Module) controller() {
	m.App.Post("/api/v1/signup", m.signup)
	m.App.Post("/api/v1/signin", m.signin)
	m.App.Get("/api/v1/auth", m.auth)
}

func (m *Module) signup(c *fiber.Ctx) error {
	req := new(signupReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	encodedHash, err := argon2.GetEncodedHash(req.Password)
	if err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), true)
		return c.Status(fiber.StatusInternalServerError).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrInternalServerError.Error(),
				Message: err.Error(),
			},
		})
	}

	accountRole := acc.ROLE_USER
	accountDetailData, err := m.addAccountService(&acc.AccountModel{
		Name:     req.Name,
		Username: req.Username,
		Password: encodedHash,
		Role:     &accountRole,
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
	balanceAmount := 0
	if _, err := m.createBalanceService(&balanceentity.BalanceModel{
		UserID: accountDetailData.ID,
		Amount: &balanceAmount,
	}); err != nil {
		if err := m.deleteAccountService(accountDetailData.ID); err != nil {
			log.SaveLogService(c.OriginalURL(), err.Error(), true)
			return c.Status(fiber.StatusInternalServerError).JSON(&contracts.Response{
				Error: &contracts.Error{
					Status:  fiber.ErrInternalServerError.Error(),
					Message: err.Error(),
				},
			})
		}
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
		Data: accountDetailData,
	})
}

func (m *Module) signin(c *fiber.Ctx) error {
	req := new(signinReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	accountDetailData, err := m.getAccountDetailByUsernameService(req.Username)
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

	isAuthorized, err := argon2.CompareStringAndEncodedHash(req.Password, accountDetailData.Password)
	if err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), true)
		return c.Status(fiber.StatusInternalServerError).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrInternalServerError.Error(),
				Message: err.Error(),
			},
		})
	}
	if !isAuthorized {
		err := errors.New("incorrect username or password")
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}

	jwtToken, err := jwt.Create(&a.JWTPayload{
		ID:   accountDetailData.ID,
		Role: accountDetailData.Role,
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
		Data: &signinRes{
			Token: jwtToken,
			ID:    accountDetailData.ID,
			Name:  accountDetailData.Name,
			Role:  accountDetailData.Role,
		},
	})
}

func (m *Module) auth(c *fiber.Ctx) error {
	tokenString, err := parser.GetReqBearerToken(c)
	if err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusUnauthorized).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrUnauthorized.Error(),
				Message: err.Error(),
			},
		})
	}

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

	tokenExp := time.Unix(*token.Expiration, 0)
	renewToken, err := jwt.Renew[a.JWTPayload](tokenString, &tokenExp)
	if err != nil {
		log.SaveLogService(c.OriginalURL(), err.Error(), false)
		return c.Status(fiber.StatusBadRequest).JSON(&contracts.Response{
			Error: &contracts.Error{
				Status:  fiber.ErrBadRequest.Error(),
				Message: err.Error(),
			},
		})
	}
	tokenString = renewToken

	accountDetailData, err := m.getAccountDetailService(token.ID)
	if err != nil {
		if pg.IsErrRecordNotFound(err) {
			log.SaveLogService(c.OriginalURL(), err.Error(), false)
			return c.Status(fiber.StatusNotFound).JSON(&contracts.Response{
				Error: &contracts.Error{
					Status:  fiber.ErrNotFound.Error(),
					Message: err.Error(),
				},
			})
		}
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
		Data: &accountRes{
			Token: tokenString,
			ID:    accountDetailData.ID,
			Name:  accountDetailData.Name,
			Role:  accountDetailData.Role,
		},
	})
}
