package engine

import (
	"fmt"
	"log"
	"sync"

	"github.com/DioCGomes/infra-scout/internal/analyzers"
	"github.com/DioCGomes/infra-scout/internal/detectors"
	"github.com/DioCGomes/infra-scout/internal/exporters"
	"github.com/DioCGomes/infra-scout/internal/models"
	"github.com/DioCGomes/infra-scout/internal/rules"
	"github.com/DioCGomes/infra-scout/internal/scanner"
)

// Config holds the engine configuration
type Config struct {
	Providers      []string           // Providers to scan (docker, terraform, etc.)
	ExcludeDirs    []string           // Directories to exclude
	Exporter       exporters.Exporter // Output exporter
	SequentialMode bool               // Process files sequentially vs concurrently
	MinSeverity    string             // Minimum severity to report
}

// Engine orchestrates the scanning process
type Engine struct {
	detector   detectors.Detector
	analyzers  map[string]analyzers.Analyzer
	ruleEngine rules.RuleEngine
	config     Config
}

func NewEngine(detector detectors.Detector, ruleEngine rules.RuleEngine, config Config) *Engine {
	return &Engine{
		detector:   detector,
		analyzers:  make(map[string]analyzers.Analyzer),
		ruleEngine: ruleEngine,
		config:     config,
	}
}

// adds an analyzer for a specific provider
func (e *Engine) RegisterAnalyzer(provider string, analyzer analyzers.Analyzer) {
	e.analyzers[provider] = analyzer
}

// detects IaC files and scans them for security issues
func (e *Engine) Scan(root string) ([]*models.ScanResult, error) {
	var scanResults []*models.ScanResult

	if e.config.SequentialMode {
		return e.scanSequential(root)
	}

	return e.scanConcurrent(root, scanResults)
}

// processes files one at a time
func (e *Engine) scanSequential(root string) ([]*models.ScanResult, error) {
	var scanResults []*models.ScanResult

	files, err := e.detector.DetectFiles(root, e.config.ExcludeDirs, e.config.Providers)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		result, err := e.scanFile(file)
		if err != nil {
			log.Printf("Could not scan file %s: %v\n", file.Path, err)
			continue
		}
		scanResults = append(scanResults, result)
	}

	if e.config.Exporter != nil {
		if err := e.config.Exporter.Export(scanResults); err != nil {
			return nil, fmt.Errorf("failed to export results: %v", err)
		}
	}

	return scanResults, nil
}

// processes files in parallel
func (e *Engine) scanConcurrent(root string, scanResults []*models.ScanResult) ([]*models.ScanResult, error) {
	var mu sync.Mutex
	var wg sync.WaitGroup

	filesChan, err := e.detector.DetectFilesChannel(root, e.config.ExcludeDirs, e.config.Providers)
	if err != nil {
		return nil, err
	}

	for file := range filesChan {
		wg.Add(1)
		go func(f models.File) {
			defer wg.Done()

			result, err := e.scanFile(f)
			if err != nil {
				log.Printf("Could not scan file %s: %v\n", f.Path, err)
				return
			}

			mu.Lock()
			scanResults = append(scanResults, result)
			mu.Unlock()
		}(file)
	}

	wg.Wait()

	if e.config.Exporter != nil {
		if err := e.config.Exporter.Export(scanResults); err != nil {
			return nil, fmt.Errorf("failed to export results: %v", err)
		}
	}

	return scanResults, nil
}

// scans a single file using the appropriate analyzer
func (e *Engine) scanFile(file models.File) (*models.ScanResult, error) {
	analyzer, ok := e.analyzers[file.Provider]
	if !ok {
		return nil, fmt.Errorf("no analyzer registered for provider: %s", file.Provider)
	}

	log.Printf("Scanning %s\n", file.Path)

	s := scanner.NewScanner(analyzer, e.ruleEngine)
	return s.ScanFile(file.Path)
}

// filters findings by minimum severity
func (e *Engine) filterBySeverity(findings []models.Finding) []models.Finding {
	if e.config.MinSeverity == "" {
		return findings
	}

	severityOrder := map[string]int{
		models.SeverityCritical: 5,
		models.SeverityHigh:     4,
		models.SeverityMedium:   3,
		models.SeverityLow:      2,
		models.SeverityInfo:     1,
	}

	minLevel := severityOrder[e.config.MinSeverity]
	if minLevel == 0 {
		return findings
	}

	var filtered []models.Finding
	for _, f := range findings {
		if severityOrder[f.Severity] >= minLevel {
			filtered = append(filtered, f)
		}
	}

	return filtered
}
