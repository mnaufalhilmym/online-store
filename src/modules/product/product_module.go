package product

import (
	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/libs/db/pg"
	p "hilmy.dev/store/src/modules/product/product_entity"
)

type Module struct {
	App *fiber.App
	DB  *pg.DB
}

func Load(module *Module) {
	p.InitRepository(module.DB)
	module.controller()
}
