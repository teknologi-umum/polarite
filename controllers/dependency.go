package controllers

import "polarite/repository"

// Dependency injection struct.
// Initialize once, use it everywhere.
type Dependency struct {
	repository.Paste
}
