package main

import (
	"runtime"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"hilmy.dev/store/src/constants"
	"hilmy.dev/store/src/contracts"
	"hilmy.dev/store/src/libs/env"
	"hilmy.dev/store/src/libs/gracefulshutdown"
	applogger "hilmy.dev/store/src/libs/logger"
)

var logger = applogger.New("App")

func main() {
	appName := env.Get(env.APP_NAME)
	appMode := env.Get(env.APP_MODE)
	appAddress := env.Get(env.APP_ADDRESS)
	webAddress := env.Get(env.WEB_ADDRESS)

	logger.Log("starting " + appName + " in " + appMode + " on " + runtime.Version())

	app := fiber.New(fiber.Config{
		AppName:     appName,
		Network:     fiber.NetworkTCP,
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		ReadTimeout: 30 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			status := fiber.StatusInternalServerError
			statusString := fiber.ErrInternalServerError.Error()
			if err != nil {
				if fiberError, ok := err.(*fiber.Error); ok {
					status = fiberError.Code
					statusString = fiberError.Error()
				}
				logger.Error(err, &applogger.Options{IsPrintStack: false})
				return c.Status(status).JSON(&contracts.Response{
					Error: &contracts.Error{
						Status:  statusString,
						Message: err.Error(),
					},
				})
			}
			return c.Status(status).JSON(&contracts.Response{
				Error: &contracts.Error{
					Status:  statusString,
					Message: "unexpected error occurred",
				},
			})
		},
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: appMode != constants.APP_MODE_RELEASE,
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: func() string {
			if appMode == constants.APP_MODE_RELEASE && len(webAddress) > 0 {
				return webAddress
			}
			return "*"
		}(),
	}))

	if env.Get(env.APP_MODE) != constants.APP_MODE_RELEASE {
		app.Use(fiberlogger.New())
	}

	module := module{app: app}
	module.load()

	gracefulshutdown.Add(gracefulshutdown.FnRunInShutdown{
		FnDescription: "shutting down app",
		Fn: func() {
			if err := app.Shutdown(); err != nil {
				logger.Error(err)
			}
		},
	})
	gracefulshutdown.Run()

	if err := app.Listen(appAddress); err != nil {
		logger.Panic(err)
	}
}
