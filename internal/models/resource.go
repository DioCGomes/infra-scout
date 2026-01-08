package models

// represents where a resource is defined in a file
type Location struct {
	File      string
	StartLine int
	EndLine   int
}

// represents an infrastructure component (e.g., S3 bucket, Docker image, K8s pod)
type Resource struct {
	Type       string
	Name       string
	Provider   string
	Attributes map[string]any
	Location   Location
}
