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

