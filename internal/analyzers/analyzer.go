package analyzers

import "github.com/DioCGomes/infra-scout/internal/models"

// Analyzer defines the interface for parsing IaC files and extracting resources
type Analyzer interface {
	// Analyze parses an IaC file and returns the resources defined in it
	Analyze(path string) ([]models.Resource, error)

	// Provider returns the provider name this analyzer handles (e.g., "terraform", "docker")
	Provider() string
}
