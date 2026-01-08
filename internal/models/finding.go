package models

const (
	SeverityCritical = "CRITICAL"
	SeverityHigh     = "HIGH"
	SeverityMedium   = "MEDIUM"
	SeverityLow      = "LOW"
	SeverityInfo     = "INFO"
)

// represents a security misconfiguration or issue found during scanning
type Finding struct {
	RuleID      string // Unique rule identifier, e.g., "DOCKER-001", "TF-AWS-001"
	Severity    string
	Resource    Resource
	Title       string
	Description string
	Remediation string
	References  []string
}
