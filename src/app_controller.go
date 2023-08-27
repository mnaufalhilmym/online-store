package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"hilmy.dev/store/src/contract"
	"hilmy.dev/store/src/libs/env"
)

func (m *module) controller() {
	m.app.Get("/", m.rootController)
	m.app.Get("/swagger/*", swagger.HandlerDefault)
}

func (*module) rootController(c *fiber.Ctx) error {
	return c.JSON(&contract.Response{
		Data: fmt.Sprintf("%s is running", env.Get(env.APP_NAME)),
	})
}
