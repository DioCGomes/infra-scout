package detectors

import "regexp"

type FilePattern struct {
	Regex    *regexp.Regexp
	Provider string
}

// Provider constants
const (
	ProviderDocker         = "docker"
	ProviderTerraform      = "terraform"
	ProviderKubernetes     = "kubernetes"
	ProviderCloudFormation = "cloudformation"
	ProviderHelm           = "helm"
	ProviderAnsible        = "ansible"
)

// File patterns for IaC providers
var (
	// Docker patterns
	DockerfilePattern    = FilePattern{Regex: regexp.MustCompile(`^Dockerfile(\..+)?$`), Provider: ProviderDocker}
	DockerComposePattern = FilePattern{Regex: regexp.MustCompile(`^(docker-)?compose.*\.ya?ml$`), Provider: ProviderDocker}
	DockerBakePattern    = FilePattern{Regex: regexp.MustCompile(`^docker-bake\.(hcl|json)$`), Provider: ProviderDocker}

	// Terraform patterns
	TerraformPattern     = FilePattern{Regex: regexp.MustCompile(`\.tf$`), Provider: ProviderTerraform}
	TerraformVarsPattern = FilePattern{Regex: regexp.MustCompile(`\.tfvars$`), Provider: ProviderTerraform}

	// Kubernetes patterns
	KubernetesPattern = FilePattern{Regex: regexp.MustCompile(`\.(ya?ml)$`), Provider: ProviderKubernetes}

	// Helm patterns
	HelmChartPattern  = FilePattern{Regex: regexp.MustCompile(`^Chart\.ya?ml$`), Provider: ProviderHelm}
	HelmValuesPattern = FilePattern{Regex: regexp.MustCompile(`^values.*\.ya?ml$`), Provider: ProviderHelm}

	// CloudFormation patterns
	CloudFormationPattern = FilePattern{Regex: regexp.MustCompile(`\.(ya?ml|json)$`), Provider: ProviderCloudFormation}

	// Ansible patterns
	AnsiblePlaybookPattern = FilePattern{Regex: regexp.MustCompile(`^(playbook|site|main).*\.ya?ml$`), Provider: ProviderAnsible}
)

// provider names to their file patterns
// Note: Some patterns (like Kubernetes, CloudFormation) may need content-based detection
var DefaultFilePatterns = map[string][]FilePattern{
	ProviderDocker: {
		DockerfilePattern,
		DockerComposePattern,
		DockerBakePattern,
	},
	ProviderTerraform: {
		TerraformPattern,
		TerraformVarsPattern,
	},
	ProviderKubernetes: {
		KubernetesPattern,
	},
	ProviderHelm: {
		HelmChartPattern,
		HelmValuesPattern,
	},
	ProviderCloudFormation: {
		CloudFormationPattern,
	},
	ProviderAnsible: {
		AnsiblePlaybookPattern,
	},
}

// patterns that identify a provider without needing content-based detection
var HighConfidencePatterns = map[string][]FilePattern{
	ProviderDocker: {
		DockerfilePattern,
		DockerComposePattern,
	},
	ProviderTerraform: {
		TerraformPattern,
		TerraformVarsPattern,
	},
	ProviderHelm: {
		HelmChartPattern,
	},
}

// list of all supported provider names
func AllProviders() []string {
	return []string{
		ProviderDocker,
		ProviderTerraform,
		ProviderKubernetes,
		ProviderCloudFormation,
		ProviderHelm,
		ProviderAnsible,
	}
}
