package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/DioCGomes/infra-scout/internal/detectors"
	"github.com/DioCGomes/infra-scout/internal/detectors/filesystem"
	"github.com/DioCGomes/infra-scout/internal/engine"
	"github.com/DioCGomes/infra-scout/internal/exporters/jsonexporter"
	"github.com/DioCGomes/infra-scout/internal/rules"
)

var version = "dev"

const art = `
   ____      ____              _____                  __ 
  /  _/___  / __/________ _   / ___/_________  __  __/ /_
  / // __ \/ /_/ ___/ __ '/   \__ \/ ___/ __ \/ / / / __/
_/ // / / / __/ /  / /_/ /   ___/ / /__/ /_/ / /_/ / /_  
/___/_/ /_/_/ /_/   \__,_/   /____/\___/\____/\__,_/\__/
`

func main() {
	var (
		providersFlag    string
		excludeDirsFlag  string
		exportFormatFlag string
		outputFileFlag   string
		minSeverityFlag  string
		sequentialFlag   bool
		versionFlag      bool
		helpFlag         bool
	)

	// Long flags
	flag.StringVar(&providersFlag, "providers", "", "Comma-separated list of providers to scan (e.g., docker,terraform,kubernetes)")
	flag.StringVar(&excludeDirsFlag, "exclude", "", "Comma-separated list of directories to exclude")
	flag.StringVar(&exportFormatFlag, "format", "json", "Export format: json, sarif, html")
	flag.StringVar(&outputFileFlag, "output", "", "Output file path (defaults to infra-scout-report.[format])")
	flag.StringVar(&minSeverityFlag, "min-severity", "", "Minimum severity to report (CRITICAL, HIGH, MEDIUM, LOW, INFO)")
	flag.BoolVar(&sequentialFlag, "sequential", false, "Process files sequentially instead of concurrently")
	flag.BoolVar(&versionFlag, "version", false, "Print version and exit")
	flag.BoolVar(&helpFlag, "help", false, "Show help message")

	// Short flag aliases
	flag.StringVar(&providersFlag, "p", "", "Alias for --providers")
	flag.StringVar(&excludeDirsFlag, "x", "", "Alias for --exclude")
	flag.StringVar(&exportFormatFlag, "f", "json", "Alias for --format")
	flag.StringVar(&outputFileFlag, "o", "", "Alias for --output")
	flag.StringVar(&minSeverityFlag, "s", "", "Alias for --min-severity")
	flag.BoolVar(&versionFlag, "v", false, "Alias for --version")
	flag.BoolVar(&helpFlag, "h", false, "Alias for --help")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, art)
		fmt.Fprintln(os.Stderr, "Infra-Scout - Infrastructure as Code Security Scanner")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  infra-scout [options] <directory>")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Supported Providers:")
		fmt.Fprintln(os.Stderr, "  docker       Dockerfile, docker-compose.yml")
		fmt.Fprintln(os.Stderr, "  terraform    *.tf, *.tfvars")
		fmt.Fprintln(os.Stderr, "  kubernetes   Kubernetes manifests (YAML)")
		fmt.Fprintln(os.Stderr, "  helm         Chart.yaml, values.yaml")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "  infra-scout .                           # Scan current directory")
		fmt.Fprintln(os.Stderr, "  infra-scout -p docker,terraform .       # Scan specific providers")
		fmt.Fprintln(os.Stderr, "  infra-scout -f sarif -o results.sarif . # Export as SARIF")
		fmt.Fprintln(os.Stderr, "  infra-scout -s high .                   # Only HIGH and CRITICAL")
		fmt.Fprintln(os.Stderr, "  infra-scout -x .terraform,vendor .      # Exclude directories")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if versionFlag {
		fmt.Printf("infra-scout v%s\n", version)
		os.Exit(0)
	}

	fmt.Print(art)

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Error: missing required argument <directory>")
		fmt.Fprintln(os.Stderr)
		flag.Usage()
		os.Exit(1)
	}

	rootDir := args[0]

	// Validate directory exists
	info, err := os.Stat(rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot access path: %s\n", err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Path is not a directory: %s\n", rootDir)
		os.Exit(1)
	}

	// Parse providers
	var providers []string
	if providersFlag != "" {
		for _, p := range strings.Split(providersFlag, ",") {
			providers = append(providers, strings.TrimSpace(p))
		}
	} else {
		providers = detectors.AllProviders()
	}

	// Parse exclude directories
	var excludeDirs []string
	if excludeDirsFlag != "" {
		excludeDirs = strings.Split(excludeDirsFlag, ",")
	}

	// Validate export format
	validFormats := map[string]bool{"json": true, "sarif": true, "html": true}
	if !validFormats[exportFormatFlag] {
		fmt.Fprintf(os.Stderr, "Invalid format '%s'. Valid options: json, sarif, html\n", exportFormatFlag)
		os.Exit(1)
	}

	// Validate severity if provided
	if minSeverityFlag != "" {
		minSeverityFlag = strings.ToUpper(minSeverityFlag)
		validSeverities := map[string]bool{
			"CRITICAL": true,
			"HIGH":     true,
			"MEDIUM":   true,
			"LOW":      true,
			"INFO":     true,
		}
		if !validSeverities[minSeverityFlag] {
			fmt.Fprintf(os.Stderr, "Invalid severity '%s'. Valid options: CRITICAL, HIGH, MEDIUM, LOW, INFO\n", minSeverityFlag)
			os.Exit(1)
		}
	}

	fmt.Println("Path:", rootDir)
	fmt.Println("Providers:", providers)
	if len(excludeDirs) > 0 {
		fmt.Println("Excluded:", excludeDirs)
	}

	// Set up output file
	formatExtensions := map[string]string{
		"json":  ".json",
		"sarif": ".sarif",
		"html":  ".html",
	}
	ext := formatExtensions[exportFormatFlag]

	outputFile := outputFileFlag
	if outputFile == "" {
		outputFile = "infra-scout-report" + ext
	} else if !strings.HasSuffix(outputFile, ext) {
		outputFile += ext
	}

	// Initialize components
	detector := filesystem.NewFSDetector()
	ruleEngine := rules.NewRuleEngine()

	// TODO: Register built-in rules here
	// registerDockerRules(ruleEngine)
	// registerTerraformRules(ruleEngine)
	// registerKubernetesRules(ruleEngine)

	config := engine.Config{
		Providers:      providers,
		ExcludeDirs:    excludeDirs,
		Exporter:       jsonexporter.NewJSONExporter(outputFile),
		SequentialMode: sequentialFlag,
		MinSeverity:    minSeverityFlag,
	}

	fmt.Println()
	fmt.Println("🔍 Scanning for security issues...")
	fmt.Printf("   Format: %s\n", exportFormatFlag)
	fmt.Printf("   Output: %s\n", outputFile)

	scanEngine := engine.NewEngine(detector, ruleEngine, config)

	// TODO: Register analyzers for each provider
	// scanEngine.RegisterAnalyzer(detectors.ProviderDocker, docker.NewDockerAnalyzer())
	// scanEngine.RegisterAnalyzer(detectors.ProviderTerraform, terraform.NewTerraformAnalyzer())

	scanResults, err := scanEngine.Scan(rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Scan failed: %v\n", err)
		os.Exit(1)
	}

	// Count findings
	totalFindings := 0
	totalResources := 0
	for _, result := range scanResults {
		totalResources += len(result.Resources)
		totalFindings += len(result.Findings)
	}

	fmt.Printf("\n   Found %d issues in %d resources across %d files.\n", totalFindings, totalResources, len(scanResults))

	// Print summary
	fmt.Println()
	fmt.Println("────────────────────────────────────────")

	if totalFindings > 0 {
		fmt.Println("⚠️  Security issues detected. Review the report for details.")
	} else {
		fmt.Println("✅ No security issues detected.")
	}

	fmt.Println("────────────────────────────────────────")

	if totalFindings > 0 {
		os.Exit(1)
	}
}
