package exporters

import "github.com/DioCGomes/infra-scout/internal/models"

// Exporter defines the interface for exporting scan results
type Exporter interface {
	// Export writes scan results to the configured output
	Export(results []*models.ScanResult) error
}
