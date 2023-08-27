package authmiddleware

import (
	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/contracts"
	"hilmy.dev/store/src/libs/parser"
	acc "hilmy.dev/store/src/modules/account/account_entity"
	a "hilmy.dev/store/src/modules/auth/auth_entity"
)

func AuthGuard(role ...acc.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(role) == 0 {
			return c.Next()
		}

		token := new(a.JWTPayload)
		if err := parser.ParseReqBearerToken(c, token); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(&contracts.Response{
				Error: &contracts.Error{
					Status:  fiber.ErrUnauthorized.Error(),
					Message: err.Error(),
				},
			})
		}

		isAuthorized := false
		for i := range role {
			if role[i] == *token.Role {
				isAuthorized = true
			}
		}

		if !isAuthorized {
			return c.Status(fiber.StatusForbidden).JSON(&contracts.Response{
				Error: &contracts.Error{
					Status:  fiber.ErrForbidden.Error(),
					Message: "you are prohibited from accessing this resource",
				},
			})
		}

		return c.Next()
	}
}
