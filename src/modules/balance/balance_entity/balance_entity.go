package balanceentity

import (
	"github.com/google/uuid"
	"hilmy.dev/store/src/libs/db/pg"
	applogger "hilmy.dev/store/src/libs/logger"
	a "hilmy.dev/store/src/modules/account/account_entity"
)

type BalanceModel struct {
	pg.Model
	UserID *uuid.UUID      `gorm:"uniqueIndex;not null" json:"userId,omitempty"`
	User   *a.AccountModel `json:"user,omitempty"`
	Amount *int            `gorm:"not null" json:"amount,omitempty"`
}

func (BalanceModel) TableName() string {
	return "balance"
}

type balanceDB = pg.Service[BalanceModel]

var balanceRepo *balanceDB
var logger = applogger.New("BalanceModule")

func InitRepository(db *pg.DB) {
	if db == nil {
		logger.Panic("db cannot be nil")
	}

	balanceRepo = pg.NewService[BalanceModel](db)
}

func BalanceRepository() *balanceDB {
	if balanceRepo == nil {
		logger.Panic("balanceRepo is nil")
	}

	return balanceRepo
}
