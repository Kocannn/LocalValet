package service

type Manager interface {
	StartService(serviceName string) error
	StopService(serviceName string) error
	GetServiceStatus(serviceName string) (bool, string)
	CheckHealth(serviceName string) (bool, string)
	GetAllocatedPort(serviceName string) int
	SetServiceVersion(serviceName, version string) error
	GetActiveServiceVersion(serviceName string) (string, error)
	GetAvailableVersions(serviceName string) ([]string, error)
}

// Status represents runtime status for a single service.
type Status struct {
	Name      string `json:"name"`
	IsRunning bool   `json:"isRunning"`
	Message   string `json:"message"`
	Port      int    `json:"port,omitempty"`
	Healthy   bool   `json:"healthy,omitempty"`
	Category  string `json:"category,omitempty"`
}

// RuntimeServiceInfo describes versioning details for a service or runtime environment.
type RuntimeServiceInfo struct {
	ServiceName       string   `json:"serviceName"`
	DisplayName       string   `json:"displayName"`
	ActiveVersion     string   `json:"activeVersion"`
	AvailableVersions []string `json:"availableVersions"`
	Category          string   `json:"category"`
	IsRunning         bool     `json:"isRunning"`
}


