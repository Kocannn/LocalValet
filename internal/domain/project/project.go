package project

// Framework represents the detected project framework.
type Framework string

const (
	FrameworkLaravel   Framework = "laravel"
	FrameworkWordPress Framework = "wordpress"
	FrameworkNextJS    Framework = "nextjs"
	FrameworkNuxt      Framework = "nuxt"
	FrameworkReact     Framework = "react"
	FrameworkVue       Framework = "vue"
	FrameworkPHP       Framework = "php"
	FrameworkStatic    Framework = "static"
	FrameworkUnknown   Framework = "unknown"
)

// Project represents a discovered local development project.
type Project struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	Framework    string `json:"framework"`
	WebRoot      string `json:"webRoot"`
	Domain       string `json:"domain"`
	VHostEnabled bool   `json:"vhostEnabled"`
	SSLEnabled   bool   `json:"sslEnabled"`
	TargetPort   int    `json:"targetPort,omitempty"`
	PHPVersion   string `json:"phpVersion,omitempty"`
	NodeVersion  string `json:"nodeVersion,omitempty"`
	CreatedAt    string `json:"createdAt,omitempty"`
	UpdatedAt    string `json:"updatedAt,omitempty"`
}

// Scanner scans directories to discover and identify projects.
type Scanner interface {
	ScanRoots(roots []string) ([]Project, error)
	DetectProject(dirPath string) (*Project, error)
}

// Repository handles persistence of projects and scanned root directories.
type Repository interface {
	GetProjects() ([]Project, error)
	SaveProjects(projects []Project) error
	GetRoots() ([]string, error)
	SaveRoots(roots []string) error
}
