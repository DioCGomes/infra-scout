package rules

import "github.com/DioCGomes/infra-scout/internal/models"

// Rule defines a security check that can be applied to resources
type Rule struct {
	ID           string
	Provider     string
	ResourceType string
	Severity     string
	Title        string
	Description  string
	Remediation  string
	References   []string
	Check        CheckFunc
}

// evaluates a resource and returns true if there's a violation
type CheckFunc func(resource models.Resource) bool

// evaluates resources against security rules and returns findings
type RuleEngine interface {
	Evaluate(resources []models.Resource) ([]models.Finding, error)
	RegisterRule(rule Rule)
	GetRules() []Rule
	GetRulesForProvider(provider string) []Rule
}
