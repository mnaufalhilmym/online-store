package balance

import (
	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/libs/db/pg"
	b "hilmy.dev/store/src/modules/balance/balance_entity"
)

type Module struct {
	App *fiber.App
	DB  *pg.DB
}

func Load(module *Module) {
	b.InitRepository(module.DB)
	module.controller()
}
