package main

import (
	"html/template"
	"log"
	"os"
	"os/signal"

	"net/http"
	"polarite/handlers"
	"polarite/packages/paste"
	"polarite/resources"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"

	"github.com/allegro/bigcache/v3"
	sentry "github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/template/html"
	"github.com/jmoiron/sqlx"
)

func main() {
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
	logger, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		Debug:            true,
		AttachStacktrace: true,
	})
	if err != nil {
		panic(err)
	}

	viewEngine := html.New("./views", ".html")
	viewEngine.AddFunc(
		// add unescape function
		"unescape", func(s string) template.HTML {
			return template.HTML(s)
		},
	)

	app := fiber.New(fiber.Config{
		AppName:                 "Teknologi Umum - Polarite",
		CaseSensitive:           true,
		StrictRouting:           false,
		ErrorHandler:            handlers.ErrorHandler(logger),
		EnableTrustedProxyCheck: true,
		BodyLimit:               1024 * 1024 * 6,
		WriteTimeout:            30 * time.Second,
		ReadTimeout:             30 * time.Second,
		Views:                   viewEngine,
	})

	// Setup Dependency injection struct
	pasteController := &paste.Dependency{
		Cache:  rds,
		Memory: mem,
		DB:     db,
		Logger: logger,
	}
	r := handlers.Dependency{
		PasteController: pasteController,
	}

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{fiber.MethodGet, fiber.MethodPost, fiber.MethodHead}, ","),
		AllowHeaders: fiber.HeaderAuthorization,
	})

	app.Use(helmet.New(helmet.Config{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSPreloadEnabled:    true,
		HSTSMaxAge:            63072000,
		HSTSExcludeSubdomains: false,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	}))

	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:         http.Dir("./views/assets"),
		Browse:       false,
		Index:        "404.html",
		NotFoundFile: "404.html",
		MaxAge:       60 * 60 * 24,
	}))

	app.Get("/", handlers.GoSecure, cache.New(cache.Config{Expiration: 1 * time.Hour, CacheControl: false}), r.HomePage)
	app.Get("/:id", corsMiddleware, handlers.GoSecure, r.Get)
	app.Post("/", corsMiddleware, handlers.GoSecure, limiter.New(limiter.Config{Max: 5, Expiration: 1 * time.Minute}), handlers.ValidateInput, r.AddPaste)

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
