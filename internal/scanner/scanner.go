package scanner

import (
	"github.com/DioCGomes/infra-scout/internal/analyzers"
	"github.com/DioCGomes/infra-scout/internal/models"
	"github.com/DioCGomes/infra-scout/internal/rules"
)

// Scanner combines an analyzer and rule engine to scan IaC files
type Scanner struct {
	analyzer   analyzers.Analyzer
	ruleEngine rules.RuleEngine
}

func NewScanner(analyzer analyzers.Analyzer, ruleEngine rules.RuleEngine) *Scanner {
	return &Scanner{
		analyzer:   analyzer,
		ruleEngine: ruleEngine,
	}
}

func (s *Scanner) ScanFile(path string) (*models.ScanResult, error) {
	resources, err := s.analyzer.Analyze(path)
	if err != nil {
		return nil, err
	}

	findings, err := s.ruleEngine.Evaluate(resources)
	if err != nil {
		return nil, err
	}

	return &models.ScanResult{
		SourceFile: path,
		Provider:   s.analyzer.Provider(),
		Resources:  resources,
		Findings:   findings,
	}, nil
}
