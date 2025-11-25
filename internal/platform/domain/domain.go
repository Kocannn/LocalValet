package domain

type ServiceManager interface {
	StartService(serviceName string) error
	StopService(serviceName string) error
	GetServiceStatus(serviceName string) (bool, string)
}

// ServiceStatus represents the status of a service
type ServiceStatus struct {
	Name      string `json:"name"`
	IsRunning bool   `json:"isRunning"`
	Message   string `json:"message"`
}

type BinarySourceInfo struct {
	OS                  string      `json:"os"`
	UsingSystemBinaries bool        `json:"using_system_binaries"`
	BinaryLocation      string      `json:"binary_location"`
	BinaryValidation    interface{} `json:"binary_validation"` // Bisa dispesifikkan lagi jika strukturnya jelas
}
