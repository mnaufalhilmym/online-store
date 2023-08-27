package parser

import (
	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/libs/validator"
)

func ParseReqQuery[T any](c *fiber.Ctx, query T) error {
	if err := c.QueryParser(query); err != nil {
		logger.Error(err)
		return err
	}
	if err := validator.Struct(query); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
