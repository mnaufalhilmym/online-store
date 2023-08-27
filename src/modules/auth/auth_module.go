package auth

import (
	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App *fiber.App
}

func Load(module *Module) {
	module.controller()
}
