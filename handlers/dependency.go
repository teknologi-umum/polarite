package handlers

import "polarite/packages/paste"

// Dependency injection struct.
// Initialize once, use it everywhere.
type Dependency struct {
	paste.PasteController
}
