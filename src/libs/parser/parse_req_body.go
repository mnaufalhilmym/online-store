package parser

import (
	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/libs/validator"
)

func ParseReqBody[T any](c *fiber.Ctx, req T) error {
	if err := c.BodyParser(req); err != nil {
		logger.Error(err)
		return err
	}
	if err := validator.Struct(req); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
