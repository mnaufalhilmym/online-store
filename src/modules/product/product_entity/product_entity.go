package productentity

import (
	"github.com/google/uuid"
	"hilmy.dev/store/src/libs/db/pg"
	applogger "hilmy.dev/store/src/libs/logger"
	pc "hilmy.dev/store/src/modules/product_category/product_category_entity"
)

type ProductModel struct {
	pg.Model
	CategoryID  *uuid.UUID               `gorm:"not null" json:"categoryId,omitempty"`
	Category    *pc.ProductCategoryModel `json:"category,omitempty"`
	Title       *string                  `gorm:"not null" json:"title,omitempty"`
	Description *string                  `gorm:"not null" json:"description,omitempty"`
	Price       *int                     `gorm:"not null" json:"price,omitempty"`
}

func (ProductModel) TableName() string {
	return "products"
}

type productDB = pg.Service[ProductModel]

var productRepo *productDB
var logger = applogger.New("ProductModule")

func InitRepository(db *pg.DB) {
	if db == nil {
		logger.Panic("db cannot be nil")
	}

	productRepo = pg.NewService[ProductModel](db)
}

func ProductRepository() *productDB {
	if productRepo == nil {
		logger.Panic("productRepo is nil")
	}

	return productRepo
}
