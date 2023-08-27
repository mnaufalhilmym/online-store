package main

import (
	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/libs/db/mongo"
	"hilmy.dev/store/src/libs/db/pg"
	"hilmy.dev/store/src/libs/env"
	"hilmy.dev/store/src/modules/account"
	"hilmy.dev/store/src/modules/log"
)

type module struct {
	app *fiber.App
}

func (m *module) load() {
	// PostgreSQL database
	pgDB := pg.NewDB(&pg.Config{
		Address:      env.Get(env.POSTGRES_ADDRESS),
		User:         env.Get(env.POSTGRES_USER),
		Password:     env.Get(env.POSTGRES_PASSWORD),
		DatabaseName: env.Get(env.POSTGRES_DB),
	})

	// MongoDB database
	mongoDBClient := mongo.NewClient(&mongo.Config{
		Address:  env.Get(env.MONGO_ADDRESS),
		User:     env.Get(env.MONGO_INITDB_ROOT_USERNAME),
		Password: env.Get(env.MONGO_INITDB_ROOT_PASSWORD),
	})

	m.controller()

	log.Load(&log.Module{
		DBClient: mongoDBClient,
	})

	account.Load(&account.Module{
		App: m.app,
		DB:  pgDB,
	})
}
