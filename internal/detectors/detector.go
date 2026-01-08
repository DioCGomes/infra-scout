package detectors

import "github.com/DioCGomes/infra-scout/internal/models"

// Detector defines the interface for finding IaC files in a directory
type Detector interface {
	DetectFiles(root string, excludeDirs []string, providers []string) ([]models.File, error)
	DetectFilesChannel(root string, excludeDirs []string, providers []string) (chan models.File, error)
}
