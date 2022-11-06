package repository_test

import (
	"log"
	"os"
	"polarite/repository"
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/getsentry/sentry-go"
)

var dependency *repository.Dependency

func TestMain(m *testing.M) {
	databaseDirectory, ok := os.LookupEnv("DATABASE_DIRECTORY")
	if !ok {
		databaseDirectory = ""
	}

	database, err := badger.Open(badger.DefaultOptions(databaseDirectory).WithInMemory(databaseDirectory == ""))
	if err != nil {
		log.Fatalf("Opening database: %s", err.Error())
	}
	defer func() {
		err := database.Close()
		if err != nil {
			log.Printf("Closing database: %s", err.Error())
		}
	}()

	logger, err := sentry.NewClient(sentry.ClientOptions{})
	if err != nil {
		log.Fatalf("Initializing sentry client: %s", err.Error())
	}

	dependency = &repository.Dependency{
		Logger: logger,
		DB:     database,
	}

	exitCode := m.Run()

	err = database.Close()
	if err != nil {
		log.Printf("Closing database: %s", err.Error())
	}

	os.Exit(exitCode)
}
