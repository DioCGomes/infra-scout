package jsonexporter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/DioCGomes/infra-scout/internal/models"
)

type JSONExporter struct {
	OutputFile string
}

func NewJSONExporter(outputFile string) *JSONExporter {
	return &JSONExporter{OutputFile: outputFile}
}

func (j *JSONExporter) Export(results []*models.ScanResult) error {
	file, err := os.Create(j.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", j.OutputFile, err)
	}
	defer file.Close()

	criticalCount := 0
	highCount := 0
	mediumCount := 0
	lowCount := 0

	var findings []map[string]interface{}

	for _, result := range results {
		for _, finding := range result.Findings {
			findingWithFile := map[string]interface{}{
				"finding": finding,
				"file":    result.SourceFile,
			}

			findings = append(findings, findingWithFile)

			switch finding.Severity {
			case models.SeverityCritical:
				criticalCount++
			case models.SeverityHigh:
				highCount++
			case models.SeverityMedium:
				mediumCount++
			case models.SeverityLow:
				lowCount++
			}
		}
	}

	summary := map[string]int{
		"Total Findings": len(findings),
		"Critical":       criticalCount,
		"High":           highCount,
		"Medium":         mediumCount,
		"Low":            lowCount,
	}

	output := map[string]interface{}{
		"summary":  summary,
		"findings": findings,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		return fmt.Errorf("failed to encode results to JSON: %v", err)
	}

	log.Printf("Findings exported to %s\n", j.OutputFile)
	return nil
}
