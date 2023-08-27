package productcategory

import (
	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/libs/db/pg"
	pc "hilmy.dev/store/src/modules/product_category/product_category_entity"
)

type Module struct {
	App *fiber.App
	DB  *pg.DB
}

func Load(module *Module) {
	pc.InitRepository(module.DB)
	module.controller()
}
