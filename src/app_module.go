package main

import (
	"github.com/gofiber/fiber/v2"
)

type module struct {
	app *fiber.App
}

func (m *module) load() {
	m.controller()
}
