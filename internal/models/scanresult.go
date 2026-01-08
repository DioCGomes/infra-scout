package models

type ScanResult struct {
	SourceFile string
	Provider   string
	Resources  []Resource
	Findings   []Finding
}
