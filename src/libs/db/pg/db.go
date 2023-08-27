package pg

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"hilmy.dev/store/src/libs/gracefulshutdown"
	applogger "hilmy.dev/store/src/libs/logger"
	"hilmy.dev/store/src/libs/validator"
)

type DB = gorm.DB

type Config struct {
	Address      string `validate:"required"`
	User         string `validate:"required"`
	Password     string `validate:"required"`
	DatabaseName string `validate:"required"`
}

var logger = applogger.New("PostgreSQL")

func NewDB(config *Config) *DB {
	logger.Log("initializing PostgreSQL database")

	if err := validator.Struct(config); err != nil {
		logger.Panic(err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "postgresql://" + config.User + ":" + config.Password + "@" + config.Address + "/" + config.DatabaseName,
	}), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		logger.Panic(err)
	}

	gracefulshutdown.Add(gracefulshutdown.FnRunInShutdown{
		FnDescription: "close PostgreSQL database",
		Fn: func() {
			db, err := db.DB()
			if err != nil {
				logger.Error(err)
				return
			}
			if err := db.Close(); err != nil {
				logger.Error(err)
			}
		},
	})

	return db
}
