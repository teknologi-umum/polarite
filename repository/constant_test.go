package repository_test

import (
	"os"
	"polarite/repository"
	"testing"
)

func TestBaseURL(t *testing.T) {
	err := os.Setenv("ENVIRONMENT", "production")
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	g := repository.BASE_URL()
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

	b := repository.BASE_URL()
	if b != "http://localhost:3000/" {
		t.Error("should be \"http://localhost:3000/\", got:", g)
	}
}
