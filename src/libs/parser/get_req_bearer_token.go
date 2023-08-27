package parser

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetReqBearerToken(c *fiber.Ctx) (*string, error) {
	authorizationHeader := c.Get("authorization")
	if len(authorizationHeader) == 0 {
		err := errors.New("authorization header not found")
		logger.Error(err)
		return nil, err
	}

	authorization := strings.Split(authorizationHeader, " ")
	if strings.ToLower(authorization[0]) != "bearer" {
		err := errors.New("not a bearer token")
		logger.Error(err)
		return nil, err
	}

	return &authorization[1], nil
}
