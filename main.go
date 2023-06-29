package main

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"polarite/controllers"
	"polarite/repository"

	"github.com/gofiber/fiber/v2"

	"github.com/dgraph-io/badger/v3"
	sentry "github.com/getsentry/sentry-go"
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/template/html"
)

func main() {
	environment, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		environment = "development"
	}

	sentryDSN, ok := os.LookupEnv("SENTRY_DSN")
	if !ok {
		sentryDSN = ""
	}

	databaseDirectory, ok := os.LookupEnv("DATABASE_DIRECTORY")
	if !ok {
		databaseDirectory = os.TempDir()
	}

	httpHostname, ok := os.LookupEnv("HTTP_HOSTNAME")
	if !ok {
		httpHostname = "0.0.0.0"
	}

	httpPort, ok := os.LookupEnv("HTTP_PORT")
	if !ok {
		httpPort = "3000"
	}

	database, err := badger.Open(badger.DefaultOptions(databaseDirectory))
	if err != nil {
		log.Fatalf("Opening database: %s", err.Error())
	}
	defer func() {
		err := database.Close()
		if err != nil {
			log.Printf("Closing database: %s", err.Error())
		}
	}()

	// Setup Sentry
	err = sentry.Init(sentry.ClientOptions{
		Dsn:              sentryDSN,
		Debug:            environment != "production",
		AttachStacktrace: true,
		Environment:      environment,
	})
	if err != nil {
		log.Fatalf("Setting up Sentry client: %s", err.Error())
	}
	defer sentry.Flush(time.Minute)

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
		ErrorHandler:            controllers.ErrorHandler(),
		EnableTrustedProxyCheck: true,
		BodyLimit:               1024 * 1024 * 6,
		WriteTimeout:            30 * time.Second,
		ReadTimeout:             30 * time.Second,
		Views:                   viewEngine,
	})

	// Setup Dependency injection struct
	repositoryDependency := &repository.Dependency{
		DB: database,
	}

	r := controllers.Dependency{
		Paste: repositoryDependency,
	}

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{fiber.MethodGet, fiber.MethodPost, fiber.MethodHead}, ","),
		AllowHeaders: fiber.HeaderAuthorization,
	})

	app.Use(fibersentry.New(fibersentry.Config{
		Repanic:         true,
		WaitForDelivery: true,
	}))
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:         http.Dir("./views/assets"),
		Browse:       false,
		Index:        "404.html",
		NotFoundFile: "404.html",
		MaxAge:       60 * 60 * 24,
	}))

	app.Get("/", cache.New(cache.Config{Expiration: 1 * time.Hour, CacheControl: false}), r.HomePage)
	app.Get("/:id", corsMiddleware, r.Get)
	app.Post("/", corsMiddleware, limiter.New(limiter.Config{Max: 5, Expiration: 1 * time.Minute}), controllers.ValidateInput, r.AddPaste)

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt)

	go func() {
		<-exitSignal

		err := app.Shutdown()
		if err != nil {
			log.Printf("Shutting down: %s", err.Error())
		}
	}()

	err = app.Listen(net.JoinHostPort(httpHostname, httpPort))
	if err != nil {
		log.Printf("Server listening: %s", err.Error())
	}
}
