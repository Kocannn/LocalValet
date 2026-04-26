package service

type Manager interface {
	StartService(serviceName string) error
	StopService(serviceName string) error
	GetServiceStatus(serviceName string) (bool, string)
	SetServiceVersion(serviceName, version string) error
	GetActiveServiceVersion(serviceName string) (string, error)
	GetAvailableVersions(serviceName string) ([]string, error)
}

// Status represents runtime status for a single service.
type Status struct {
	Name      string `json:"name"`
	IsRunning bool   `json:"isRunning"`
	Message   string `json:"message"`
}
