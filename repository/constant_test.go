package repository_test

import (
	"os"
	"testing"

	"polarite/repository"
)

func TestBaseURL(t *testing.T) {
	err := os.Setenv("ENVIRONMENT", "production")
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	g := repository.BaseUrl()
	if g != "https://polarite.teknologiumum.com/" {
		t.Error("should be \"https://polarite.teknologiumum.com/\", got:", g)
	}

	err = os.Setenv("ENVIRONMENT", "development")
	if err != nil {
		t.Error("an error was thrown:", err)
	}
	err = os.Setenv("PORT", "3000")
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	b := repository.BaseUrl()
	if b != "http://localhost:3000/" {
		t.Error("should be \"http://localhost:3000/\", got:", g)
	}
}
