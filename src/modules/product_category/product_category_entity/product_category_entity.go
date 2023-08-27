package productcategoryentity

import (
	"hilmy.dev/store/src/libs/db/pg"
	applogger "hilmy.dev/store/src/libs/logger"
)

type ProductCategoryModel struct {
	pg.Model
	Name *string `gorm:"uniqueIndex;not null" json:"name,omitempty"`
}

func (ProductCategoryModel) TableName() string {
	return "product_categories"
}

type productCategoryDB = pg.Service[ProductCategoryModel]

var productCategoryRepo *productCategoryDB
var logger = applogger.New("ProductCategoryModule")

func InitRepository(db *pg.DB) {
	if db == nil {
		logger.Panic("db cannot be nil")
	}

	productCategoryRepo = pg.NewService[ProductCategoryModel](db)
}

func ProductCategoryRepository() *productCategoryDB {
	if productCategoryRepo == nil {
		logger.Panic("productCategoryRepo is nil")
	}

	return productCategoryRepo
}
