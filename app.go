package main

import (
	"os"
	"polarite/business/controllers"
	"polarite/handlers"
	"time"

	sentryfiber "github.com/aldy505/sentry-fiber"
	"github.com/allegro/bigcache/v3"
	sentry "github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lib/pq"
)

func App() *fiber.App {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: false,
		ErrorHandler:  handlers.ErrorHandler,
	})

	// Setup Postgres/Cockroach
	dbURL, err := pq.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		panic(err)
	}

	// Setup Redis
	rdsConfig, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}

	rds := redis.NewClient(rdsConfig)

	// Setup In-Memory
	mem, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Hour * 12))
	if err != nil {
		panic(err)
	}

	// Setup Sentry
	err = sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})
	if err != nil {
		panic(err)
	}

	// Setup Dependency injection struct
	pasteController := &controllers.PasteControllerImpl{
		Cache:  rds,
		Memory: mem,
	}
	r := handlers.Dependency{
		DB:              db,
		PasteController: pasteController,
	}

	app.Use(cors.New())
	app.Use(sentryfiber.New(sentryfiber.Options{}))
	app.Get("/", cache.New(cache.Config{Expiration: 1 * time.Second, CacheControl: true}), r.HomePage)
	app.Get("/:id", cache.New(cache.Config{Expiration: 1 * time.Second, CacheControl: true}), r.Get)
	app.Post("/", limiter.New(limiter.Config{Max: 5, Expiration: 1 * time.Minute}), r.AddPaste)

	return app
}
