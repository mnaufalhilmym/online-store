package main

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"hilmy.dev/store/src/libs/db/mongo"
	"hilmy.dev/store/src/libs/db/pg"
	"hilmy.dev/store/src/libs/env"
	"hilmy.dev/store/src/libs/hash/argon2"
	"hilmy.dev/store/src/libs/jwx/jwt"
	"hilmy.dev/store/src/modules/account"
	"hilmy.dev/store/src/modules/auth"
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

	// JWT
	jwt.Init(&jwt.Config{
		Duration: func() *time.Duration {
			duration, err := time.ParseDuration(env.Get(env.JWT_DURATION))
			if err != nil {
				logger.Panic(err)
			}
			return &duration
		}(),
	})

	// Argon2
	argon2.Init(&argon2.Config{
		Memory: func() uint32 {
			hashMemory, err := strconv.ParseUint(env.Get(env.HASH_MEMORY), 10, 32)
			if err != nil {
				logger.Panic(err)
			}
			return uint32(hashMemory)
		}(),
		Iterations: func() uint32 {
			hashIterations, err := strconv.ParseUint(env.Get(env.HASH_ITERATIONS), 10, 32)
			if err != nil {
				logger.Panic(err)
			}
			return uint32(hashIterations)
		}(),
		Parallelism: func() uint8 {
			hashParallelism, err := strconv.ParseUint(env.Get(env.HASH_PARALLELISM), 10, 8)
			if err != nil {
				logger.Panic(err)
			}
			return uint8(hashParallelism)
		}(),
		SaltLength: func() int {
			hashSaltLength, err := strconv.Atoi(env.Get(env.HASH_SALTLENGTH))
			if err != nil {
				logger.Panic(err)
			}
			return hashSaltLength
		}(),
		KeyLength: func() uint32 {
			hashKeyLength, err := strconv.ParseUint(env.Get(env.HASH_KEYLENGTH), 10, 32)
			if err != nil {
				logger.Panic(err)
			}
			return uint32(hashKeyLength)
		}(),
	})

	m.controller()

	log.Load(&log.Module{
		DBClient: mongoDBClient,
	})

	account.Load(&account.Module{
		App: m.app,
		DB:  pgDB,
	})

	auth.Load(&auth.Module{
		App: m.app,
	})
}
