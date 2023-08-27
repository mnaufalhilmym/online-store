package transaction

import (
	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/libs/db/pg"
	transactionentity "hilmy.dev/store/src/modules/transaction/transaction_entity"
)

type Module struct {
	App *fiber.App
	DB  *pg.DB
}

func Load(module *Module) {
	transactionentity.InitRepository(module.DB)
	module.controller()
}
