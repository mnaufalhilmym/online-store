package shoppingcartentity

import (
	"github.com/google/uuid"
	"hilmy.dev/store/src/libs/db/pg"
	applogger "hilmy.dev/store/src/libs/logger"
	a "hilmy.dev/store/src/modules/account/account_entity"
	p "hilmy.dev/store/src/modules/product/product_entity"
)

type ShoppingCartItemModel struct {
	pg.Model
	UserID    *uuid.UUID      `gorm:"not null" json:"userId,omitempty"`
	User      *a.AccountModel `json:"user,omitempty"`
	ProductID *uuid.UUID      `gorm:"not null" json:"productId,omitempty"`
	Product   *p.ProductModel `json:"product,omitempty"`
	Amount    *int            `gorm:"not null" json:"amount,omitempty"`
}

func (ShoppingCartItemModel) TableName() string {
	return "shopping_cart_items"
}

type shoppingCartItemDB = pg.Service[ShoppingCartItemModel]

var shoppingCartItemRepo *shoppingCartItemDB
var logger = applogger.New("ShoppingCartModule")

func InitRepository(db *pg.DB) {
	if db == nil {
		logger.Panic("db cannot be nil")
	}

	shoppingCartItemRepo = pg.NewService[ShoppingCartItemModel](db)
}

func ShoppingCartItemRepository() *shoppingCartItemDB {
	if shoppingCartItemRepo == nil {
		logger.Panic("shoppingCartItemRepo is nil")
	}

	return shoppingCartItemRepo
}
