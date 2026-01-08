package filesystem

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/DioCGomes/infra-scout/internal/detectors"
	"github.com/DioCGomes/infra-scout/internal/models"
)

type FSDetector struct{}

func NewFSDetector() *FSDetector {
	return &FSDetector{}
}

// returns a channel immediately, running the anonymous function in another goroutine
func (d *FSDetector) DetectFilesChannel(root string, excludeDirs []string, providers []string) (chan models.File, error) {
	filesChan := make(chan models.File)
	providerSet := make(map[string]bool)

	if len(providers) == 0 {
		for _, p := range detectors.AllProviders() {
			providerSet[p] = true
		}
	} else {
		for _, p := range providers {
			providerSet[p] = true
		}
	}

	excludeSet := make(map[string]bool)
	for _, dir := range excludeDirs {
		excludeSet[dir] = true
	}

	go func() {
		defer close(filesChan)

		_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				if d.shouldExclude(info.Name(), excludeSet) {
					return filepath.SkipDir
				}
				return nil
			}

			if file, ok := d.matchFile(path, info.Name(), providerSet); ok {
				filesChan <- file
			}

			return nil
		})
	}()

	return filesChan, nil
}

// walks the directory tree and returns all IaC files found
func (d *FSDetector) DetectFiles(root string, excludeDirs []string, providers []string) ([]models.File, error) {
	var files []models.File

	providerSet := make(map[string]bool)
	if len(providers) == 0 {
		for _, p := range detectors.AllProviders() {
			providerSet[p] = true
		}
	} else {
		for _, p := range providers {
			providerSet[p] = true
		}
	}

	excludeSet := make(map[string]bool)
	for _, dir := range excludeDirs {
		excludeSet[dir] = true
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if d.shouldExclude(info.Name(), excludeSet) {
				return filepath.SkipDir
			}
			return nil
		}

		if file, ok := d.matchFile(path, info.Name(), providerSet); ok {
			files = append(files, file)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func (d *FSDetector) shouldExclude(name string, excludeSet map[string]bool) bool {
	// Always exclude common non-IaC directories
	defaultExcludes := []string{
		".git",
		".terraform",
		"node_modules",
		"vendor",
		".venv",
		"__pycache__",
	}

	for _, exclude := range defaultExcludes {
		if name == exclude {
			return true
		}
	}

	return excludeSet[name]
}

// checks if a file matches any IaC pattern and returns the matched file
func (d *FSDetector) matchFile(path, name string, providerSet map[string]bool) (models.File, bool) {
	// Check high-confidence patterns first (these don't need content detection)
	for provider, patterns := range detectors.HighConfidencePatterns {
		if !providerSet[provider] {
			continue
		}

		for _, pattern := range patterns {
			if pattern.Regex.MatchString(name) {
				return models.File{
					Path:     path,
					Provider: provider,
				}, true
			}
		}
	}

	// For Terraform files, also check the extension
	if providerSet[detectors.ProviderTerraform] {
		if strings.HasSuffix(name, ".tf") || strings.HasSuffix(name, ".tfvars") {
			return models.File{
				Path:     path,
				Provider: detectors.ProviderTerraform,
			}, true
		}
	}

	// TODO: Add content-based detection for YAML files (Kubernetes, CloudFormation, etc.)
	// This would require reading the file and checking for specific keys like:
	// - Kubernetes: apiVersion, kind
	// - CloudFormation: AWSTemplateFormatVersion, Resources
	// - Ansible: hosts, tasks

	return models.File{}, false
}
