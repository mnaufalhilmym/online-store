package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"time"

	applogger "hilmy.dev/store/src/libs/logger"
	"hilmy.dev/store/src/libs/validator"
)

var conf *jwtConf
var logger = applogger.New("JWT")

type jwtConf struct {
	duration   *time.Duration
	privateKey *rsa.PrivateKey
}

type Config struct {
	Bits     int
	Duration *time.Duration `validate:"required"`
}

func Init(config *Config) {
	logger.Log("initializing JWT")

	if err := validator.Struct(config); err != nil {
		logger.Panic(err)
	}

	if config.Bits == 0 {
		config.Bits = 2048
	}

	privKey, err := rsa.GenerateKey(rand.Reader, config.Bits)
	if err != nil {
		logger.Panic(err)
	}

	conf = &jwtConf{
		duration:   config.Duration,
		privateKey: privKey,
	}
}
