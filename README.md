# Infra-Scout

> **⚠️ Work in Progress**: This project is under active development and is not yet functional. The core architecture is in place, but analyzers and security rules are still being implemented.
>
> Infra-Scout is based on [Scout](https://github.com/mlw157/scout), a Software Composition Analysis (SCA) tool for dependency vulnerabilities. This project extends the same architecture to infrastructure-as-code security scanning.

Infra-Scout is an Infrastructure as Code (IaC) security scanner. It analyzes your infrastructure configurations and checks them against security best practices and misconfigurations.

## Supported Providers

| Provider | Files | Status |
|----------|-------|--------|
| **Docker** | `Dockerfile`, `docker-compose.yml` | Planned |
| **Terraform** | `*.tf`, `*.tfvars` | Planned |
| **Kubernetes** | K8s manifests (YAML) | Planned |
| **Helm** | `Chart.yaml`, `values.yaml` | Planned |
| **CloudFormation** | AWS CFN templates | Planned |
| **Ansible** | Playbooks | Planned |

## Installation

### Build from source

```bash
git clone https://github.com/DioCGomes/infra-scout.git
cd infra-scout
make build
```

### Install to GOPATH

```bash
make install
```

## Usage

```bash
# Scan current directory
infra-scout .

# Scan specific providers only
infra-scout -p docker,terraform .

# Exclude directories
infra-scout -x .terraform,node_modules .

# Export as SARIF (for GitHub Security tab)
infra-scout -f sarif -o results.sarif .

# Filter by minimum severity
infra-scout -s high .
```

### Command-Line Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--providers` | `-p` | Providers to scan | all |
| `--exclude` | `-x` | Directories to exclude | - |
| `--format` | `-f` | Export format (json, sarif, html) | json |
| `--output` | `-o` | Output file path | infra-scout-report.[ext] |
| `--min-severity` | `-s` | Minimum severity (CRITICAL, HIGH, MEDIUM, LOW, INFO) | - |
| `--sequential` | | Process files sequentially | false |
| `--version` | `-v` | Print version | - |
| `--help` | `-h` | Show help | - |

## Architecture

Infra-Scout follows a modular, plugin-based architecture:

### Core Components

- **Engine**: Orchestrates the scanning process, coordinating detectors, analyzers, and exporters
- **Scanner**: Combines an analyzer and rule engine to scan individual files
- **Rule Engine**: Evaluates resources against security rules and generates findings

### Interfaces

- **Analyzer**: Parses IaC files and extracts resources (e.g., DockerAnalyzer, TerraformAnalyzer)
- **Detector**: Finds IaC files in the filesystem
- **Exporter**: Outputs scan results in various formats (JSON, SARIF, HTML)
- **Rule**: Defines a security check with severity, description, and remediation

### Models

- **Resource**: Represents an infrastructure component (S3 bucket, Docker image, K8s pod)
- **Finding**: A security issue found during scanning
- **ScanResult**: Results from scanning a single file

## Example Rules

### Docker

| Rule ID | Severity | Description |
|---------|----------|-------------|
| DOCKER-001 | HIGH | Container running as root |
| DOCKER-002 | MEDIUM | Using `latest` tag for base image |
| DOCKER-003 | LOW | Missing HEALTHCHECK instruction |

### Terraform (AWS)

| Rule ID | Severity | Description |
|---------|----------|-------------|
| TF-AWS-001 | CRITICAL | S3 bucket without encryption |
| TF-AWS-002 | HIGH | Security group with 0.0.0.0/0 ingress |
| TF-AWS-003 | HIGH | RDS instance publicly accessible |

### Kubernetes

| Rule ID | Severity | Description |
|---------|----------|-------------|
| K8S-001 | CRITICAL | Privileged container |
| K8S-002 | HIGH | Missing resource limits |
| K8S-003 | MEDIUM | Running as root |

## GitHub Actions

```yaml
name: "Infra-Scout"
on: [push, pull_request]

jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Infra-Scout
        run: |
          # Download and run infra-scout
          curl -LO "https://github.com/DioCGomes/infra-scout/releases/latest/download/infra-scout-linux-amd64"
          chmod +x infra-scout-linux-amd64
          ./infra-scout-linux-amd64 -f sarif -o results.sarif .
      
      - name: Upload SARIF
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: results.sarif
```

## Development

### Project Structure

```
infra-scout/
├── cmd/
│   └── infra-scout/
│       └── main.go           # CLI entry point
├── internal/
│   ├── analyzers/            # IaC file analyzers
│   │   ├── analyzer.go       # Analyzer interface
│   │   ├── docker/           # Docker analyzer (planned)
│   │   ├── terraform/        # Terraform analyzer (planned)
│   │   └── kubernetes/       # Kubernetes analyzer (planned)
│   ├── detectors/            # File detection
│   │   ├── detector.go       # Detector interface
│   │   ├── patterns.go       # File patterns
│   │   └── filesystem/       # Filesystem detector
│   ├── engine/               # Scan orchestration
│   │   └── engine.go
│   ├── exporters/            # Output exporters
│   │   ├── exporter.go       # Exporter interface
│   │   └── jsonexporter/
│   ├── models/               # Data models
│   │   ├── resource.go
│   │   ├── finding.go
│   │   └── scanresult.go
│   ├── rules/                # Security rules
│   │   ├── rule.go           # Rule interface
│   │   ├── engine.go         # Rule engine
│   │   └── builtin/          # Built-in rules (planned)
│   └── scanner/              # File scanning
│       └── scanner.go
├── Makefile
├── go.mod
└── README.md
```

### Running Tests

```bash
make test
```

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all
```
