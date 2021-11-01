package handlers

import (
	"polarite/business/controllers"
)

// Dependency injection struct.
// Initialize once, use it everywhere.
type Dependency struct {
	controllers.PasteController
}
