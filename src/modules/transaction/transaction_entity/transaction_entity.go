package transactionentity

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"hilmy.dev/store/src/libs/db/pg"
	applogger "hilmy.dev/store/src/libs/logger"
	a "hilmy.dev/store/src/modules/account/account_entity"
)

type TransactionStatus string

const (
	STATUS_WAITING_PAYMENT TransactionStatus = "WAITING_PAYMENT"
	STATUS_COMPLETED       TransactionStatus = "COMPLETED"
	STATUS_CANCELLED       TransactionStatus = "CANCELLED"
)

type TransactionModel struct {
	pg.Model
	UserID *uuid.UUID         `gorm:"not null" json:"userId,omitempty"`
	User   *a.AccountModel    `json:"user,omitempty"`
	Status *TransactionStatus `gorm:"not null" json:"status,omitempty"`
	Price  *int               `gorm:"not null" json:"price,omitempty"`
	Data   datatypes.JSON     `gorm:"not null" json:"data,omitempty"`
}

func (TransactionModel) TableName() string {
	return "transactions"
}

type transactionDB = pg.Service[TransactionModel]

var transactionRepo *transactionDB
var logger = applogger.New("TransactionModule")

func InitRepository(db *pg.DB) {
	if db == nil {
		logger.Panic("db cannot be nil")
	}

	transactionRepo = pg.NewService[TransactionModel](db)
}

func TransactionRepository() *transactionDB {
	if transactionRepo == nil {
		logger.Panic("transactionRepo is nil")
	}

	return transactionRepo
}
