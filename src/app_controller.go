package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/contracts"
	"hilmy.dev/store/src/libs/env"
)

func (m *module) controller() {
	m.app.Get("/", m.rootController)
}

func (*module) rootController(c *fiber.Ctx) error {
	return c.JSON(&contracts.Response{
		Data: fmt.Sprintf("%s is running", env.Get(env.APP_NAME)),
	})
}
