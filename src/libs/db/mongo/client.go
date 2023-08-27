package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hilmy.dev/store/src/libs/gracefulshutdown"
	applogger "hilmy.dev/store/src/libs/logger"
	"hilmy.dev/store/src/libs/validator"
)

type Client = mongo.Client

type Config struct {
	Address  string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
}

var logger = applogger.New("MongoDB")

func NewClient(config *Config) *Client {
	logger.Log("initializing MongoDB client")

	if err := validator.Struct(config); err != nil {
		logger.Panic(err)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://"+config.User+":"+config.Password+"@"+config.Address))
	if err != nil {
		logger.Panic(err)
	}

	gracefulshutdown.Add(gracefulshutdown.FnRunInShutdown{
		FnDescription: "disconnect mongodb client",
		Fn: func() {
			if err := client.Disconnect(context.TODO()); err != nil {
				logger.Error(err)
			}
		},
	})

	return client
}
