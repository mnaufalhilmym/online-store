package log

import (
	"hilmy.dev/store/src/libs/db/mongo"
	applogger "hilmy.dev/store/src/libs/logger"
)

type Module struct {
	DBClient *mongo.Client
}

var logger = applogger.New("LogModule")

func Load(module *Module) {
	module.initRepository()
}
