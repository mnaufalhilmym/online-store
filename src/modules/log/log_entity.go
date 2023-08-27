package log

import (
	"hilmy.dev/store/src/libs/db/mongo"
	"hilmy.dev/store/src/libs/env"
)

type logModel struct {
	mongo.Model
	Location *string `gorm:"not null"`
	Message  *string `gorm:"not null"`
	Stack    *string
}

func (logModel) DatabaseName() string {
	return env.Get(env.MONGO_DATABASE_NAME)
}

func (logModel) CollectionName() string {
	return "log"
}

type logDB = mongo.Service[logModel]

var logRepo *logDB

func (m *Module) initRepository() {
	if m.DBClient == nil {
		logger.Panic("dbClient cannot be nil")
	}

	logRepo = mongo.NewService[logModel](m.DBClient)
}

func LogRepository() *logDB {
	if logRepo == nil {
		logger.Panic("logRepo is nil")
	}

	return logRepo
}
