package parser

import (
	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/libs/jwx/jwt"
)

func ParseReqBearerToken[T any](c *fiber.Ctx, tokenData *T) error {
	token, err := GetReqBearerToken(c)
	if err != nil {
		return err
	}

	if err := jwt.Parse(*token, tokenData); err != nil {
		return err
	}

	return nil
}
