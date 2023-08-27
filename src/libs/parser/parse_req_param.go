package parser

import (
	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/libs/validator"
)

func ParseReqParam[T any](c *fiber.Ctx, param T) error {
	if err := c.ParamsParser(param); err != nil {
		logger.Error(err)
		return err
	}
	if err := validator.Struct(param); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
