package rules

import "github.com/DioCGomes/infra-scout/internal/models"

// DefaultRuleEngine implements the RuleEngine interface
type DefaultRuleEngine struct {
	rules []Rule
}

// NewRuleEngine creates a new rule engine instance
func NewRuleEngine() *DefaultRuleEngine {
	return &DefaultRuleEngine{
		rules: make([]Rule, 0),
	}
}

// RegisterRule adds a rule to the engine
func (e *DefaultRuleEngine) RegisterRule(rule Rule) {
	e.rules = append(e.rules, rule)
}

// GetRules returns all registered rules
func (e *DefaultRuleEngine) GetRules() []Rule {
	return e.rules
}

// GetRulesForProvider returns rules for a specific provider
func (e *DefaultRuleEngine) GetRulesForProvider(provider string) []Rule {
	var result []Rule
	for _, rule := range e.rules {
		if rule.Provider == provider || rule.Provider == "*" {
			result = append(result, rule)
		}
	}
	return result
}

// Evaluate checks resources against applicable rules and returns findings
func (e *DefaultRuleEngine) Evaluate(resources []models.Resource) ([]models.Finding, error) {
	var findings []models.Finding

	for _, resource := range resources {
		rules := e.GetRulesForProvider(resource.Provider)

		for _, rule := range rules {
			// Check if rule applies to this resource type
			if rule.ResourceType != "*" && rule.ResourceType != resource.Type {
				continue
			}

			// Run the check function
			if rule.Check != nil && rule.Check(resource) {
				finding := models.Finding{
					RuleID:      rule.ID,
					Severity:    rule.Severity,
					Resource:    resource,
					Title:       rule.Title,
					Description: rule.Description,
					Remediation: rule.Remediation,
					References:  rule.References,
				}
				findings = append(findings, finding)
			}
		}
	}

	return findings, nil
}
