package accountentity

import (
	"hilmy.dev/store/src/libs/db/pg"
	"hilmy.dev/store/src/libs/env"
	"hilmy.dev/store/src/libs/hash/argon2"
	applogger "hilmy.dev/store/src/libs/logger"
)

type Role string

const (
	ROLE_ADMIN Role = "ADMIN"
	ROLE_USER  Role = "USER"
)

type AccountModel struct {
	pg.Model
	Name     *string `gorm:"not null" json:"name,omitempty"`
	Username *string `gorm:"uniqueIndex;not null" json:"username,omitempty"`
	Password *string `gorm:"not null" json:"-"`
	Role     *Role   `gorm:"not null" json:"role,omitempty"`
}

func (AccountModel) TableName() string {
	return "account"
}

type accountDB = pg.Service[AccountModel]

var accountRepo *accountDB
var logger = applogger.New("AccountModule")

func InitRepository(db *pg.DB) {
	if db == nil {
		logger.Panic("db cannot be nil")
	}

	accountRepo = pg.NewService[AccountModel](db)
}

func AccountRepository() *accountDB {
	if accountRepo == nil {
		logger.Panic("accountRepo is nil")
	}

	return accountRepo
}

func CreateInitialAccount() {
	accountRole := Role(env.Get(env.INITIAL_ACCOUNT_ROLE, env.Option{MustExist: true}))

	count, err := AccountRepository().Count(&pg.CountOptions{
		Where: &[]pg.Where{
			{
				Query: "role = ?",
				Args:  []interface{}{accountRole},
			},
		},
	})
	if err != nil {
		logger.Panic(err)
	}
	if *count > 0 {
		return
	}

	accountName := env.Get(env.INITIAL_ACCOUNT_NAME, env.Option{MustExist: true})
	accountUsername := env.Get(env.INITIAL_ACCOUNT_USERNAME, env.Option{MustExist: true})
	accountPassword := env.Get(env.INITIAL_ACCOUNT_PASSWORD, env.Option{MustExist: true})
	encodedHash, err := argon2.GetEncodedHash(&accountPassword)
	if err != nil {
		logger.Panic(err)
	}
	if _, err := AccountRepository().Create(&AccountModel{
		Name:     &accountName,
		Username: &accountUsername,
		Password: encodedHash,
		Role:     &accountRole,
	}); err != nil {
		logger.Panic(err)
	}
}
