package main

import (
	"log"
	"os"
	"os/signal"

	"net/http"
	"polarite/business/controllers"
	"polarite/handlers"
	"polarite/resources"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"

	sentryfiber "github.com/aldy505/sentry-fiber"
	"github.com/allegro/bigcache/v3"
	sentry "github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/template/html"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:                 "Teknologi Umum - Polarite",
		CaseSensitive:           true,
		StrictRouting:           false,
		ErrorHandler:            handlers.ErrorHandler,
		EnableTrustedProxyCheck: true,
		BodyLimit:               1024 * 1024 * 6,
		WriteTimeout:            30 * time.Second,
		ReadTimeout:             30 * time.Second,
		Views:                   html.New("./views", ".html"),
	})

	// Setup MySQL/Planetscale
	dbURL, err := resources.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	db, err := sqlx.Connect("mysql", dbURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Setup Redis
	rdsConfig, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}

	rds := redis.NewClient(rdsConfig)
	defer rds.Close()

	// Setup In-Memory
	mem, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Hour * 12))
	if err != nil {
		panic(err)
	}
	defer mem.Close()

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
		DB:     db,
	}
	r := handlers.Dependency{
		PasteController: pasteController,
	}

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{fiber.MethodGet, fiber.MethodPost, fiber.MethodHead}, ","),
		AllowHeaders: fiber.HeaderAuthorization,
	})
	app.Use(sentryfiber.New(sentryfiber.Options{}))
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:         http.Dir("./views/assets"),
		Browse:       false,
		Index:        "404.html",
		NotFoundFile: "404.html",
		MaxAge:       60 * 60 * 24,
	}))

	app.Get("/", cache.New(cache.Config{Expiration: 1 * time.Hour, CacheControl: true}), r.HomePage)
	app.Get("/:id", corsMiddleware, cache.New(cache.Config{Expiration: 1 * time.Hour, CacheControl: true}), r.Get)
	app.Post("/", corsMiddleware, limiter.New(limiter.Config{Max: 5, Expiration: 1 * time.Minute}), handlers.ValidateInput, r.AddPaste)

	if os.Getenv("ENVIRONMENT") == "development" {
		startServer(app)
	} else {
		startServerGraceful(app)
	}
}

func startServer(app *fiber.App) {
	if err := app.Listen(os.Getenv("HOST") + ":" + os.Getenv("PORT")); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}

func startServerGraceful(app *fiber.App) {
	// Create channel for idle connections.
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
		<-sigint

		// Received an interrupt signal, shutdown.
		if err := app.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnsClosed)
	}()

	// Run server.
	if err := app.Listen(os.Getenv("HOST") + ":" + os.Getenv("PORT")); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}
