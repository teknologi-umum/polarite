package handlers

import (
	"polarite/business/controllers"

	"github.com/jmoiron/sqlx"
)

// Dependency injection struct.
// Initialize once, use it everywhere.
type Dependency struct {
	DB *sqlx.DB
	controllers.PasteController
}
