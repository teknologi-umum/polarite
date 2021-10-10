package handlers

import (
	"polarite/business/controllers"

	"github.com/jmoiron/sqlx"
)

type Dependency struct {
	DB *sqlx.DB
	controllers.PasteController
}
