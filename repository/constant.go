package repository

import (
	"os"
	"strings"
)

const ID_NOT_FOUND = "ID specified was not found"

// this one can't be a constant, unless we agree on a static PORT number for dev
func BaseUrl() string {
	if strings.ToLower(os.Getenv("ENVIRONMENT")) == "production" {
		return "https://polarite.teknologiumum.com/"
	}
	return "http://localhost:" + os.Getenv("PORT") + "/"
}
