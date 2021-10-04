package main

import (
	"context"
	"os"
	"polarite/handlers"
	"time"

	sentryfiber "github.com/aldy505/sentry-fiber"
	"github.com/allegro/bigcache/v3"
	sentry "github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v4/pgxpool"
)

func App() *fiber.App {
	app := fiber.New(fiber.Config{
		ETag: true,
		CaseSensitive: true,
		StrictRouting: false,
		ErrorHandler: ,
	})

	// Setup Postgres/Cockroach
	dbConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	db, err := pgxpool.ConnectConfig(context.Background(), dbConfig)
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

	// Setup Dependency injection struct
	r := handlers.Dependency{
		DB:     db,
		Cache:  rds,
		Memory: mem,
	}

	app.Use(cors.New())
	app.Use(sentryfiber.New(sentryfiber.Options{}))
	app.Get("/", cache.New(cache.Config{Expiration: 1 * time.Hour, CacheControl: true}), r.HomePage)
	app.Get("/:id", cache.New(cache.Config{Expiration: 1 * time.Hour, CacheControl: true}), r.Get)
	app.Post("/", r.AddPaste)

	return app
}
