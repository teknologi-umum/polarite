package main

import (
	"context"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"go.opentelemetry.io/otel"

	"polarite/controllers"
	slogotel "polarite/platform/slog-otel"
	"polarite/repository"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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

	slog.SetDefault(slog.New(slogotel.OtelHandler{
		Next: slog.NewJSONHandler(os.Stdout, nil),
	}))

	database, err := badger.Open(badger.DefaultOptions(databaseDirectory))
	if err != nil {
		slog.ErrorContext(context.Background(), "error during opening badger database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		err := database.Close()
		if err != nil {
			slog.WarnContext(context.Background(), "error during closing badger database", slog.String("error", err.Error()))
		}
	}()

	// Setup Sentry
	err = sentry.Init(sentry.ClientOptions{
		Dsn:           sentryDSN,
		Debug:         environment != "production",
		SampleRate:    1.0,
		EnableTracing: false,
		TracesSampler: func(ctx sentry.SamplingContext) float64 {
			if ctx.Span.Name == "GET /" || ctx.Span.Name == "GET /robots.txt" || ctx.Span.Name == "GET /favicon.ico" {
				return 0
			}
			return 0.2
		},
		Environment: environment,
	})
	if err != nil {
		slog.ErrorContext(context.Background(), "error during setting up sentry client", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer sentry.Flush(time.Minute)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())

	// Setup Dependency injection struct
	repositoryDependency := &repository.Dependency{
		DB: database,
	}

	r := controllers.Dependency{
		Paste: repositoryDependency,
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
		ErrorHandler:            controllers.ErrorHandler(),
		EnableTrustedProxyCheck: true,
		BodyLimit:               1024 * 1024 * 6,
		WriteTimeout:            30 * time.Second,
		ReadTimeout:             30 * time.Second,
		Views:                   viewEngine,
	})

	app.Use(fibersentry.New(fibersentry.Config{
		Repanic:         true,
		WaitForDelivery: true,
	}))

	app.Use(otelfiber.Middleware())

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
			slog.WarnContext(context.Background(), "error during shutting down server", slog.String("error", err.Error()))
		}
	}()

	err = app.Listen(net.JoinHostPort(httpHostname, httpPort))
	if err != nil {
		slog.ErrorContext(context.Background(), "error listening to port", slog.String("error", err.Error()))
	}
}
