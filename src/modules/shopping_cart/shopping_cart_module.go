package shoppingcart

import (
	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/libs/db/pg"
	sc "hilmy.dev/store/src/modules/shopping_cart/shopping_cart_entity"
)

type Module struct {
	App *fiber.App
	DB  *pg.DB
}

func Load(module *Module) {
	sc.InitRepository(module.DB)
	module.controller()
}
