package controllers

import (
	"polarite/repository"

	"go.opentelemetry.io/otel"
)

// Dependency injection struct.
// Initialize once, use it everywhere.
type Dependency struct {
	repository.Paste
}

var tracer = otel.Tracer("polarite")
