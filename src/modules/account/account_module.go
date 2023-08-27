package account

import (
	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/libs/db/pg"
	a "hilmy.dev/store/src/modules/account/account_entity"
)

type Module struct {
	App *fiber.App
	DB  *pg.DB
}

func Load(module *Module) {
	a.InitRepository(module.DB)
	a.CreateInitialAccount()
	module.controller()
}
