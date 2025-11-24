package domain

type ServiceManager interface {
	StartService(serviceName string) error
	StopService(serviceName string) error
	GetServiceStatus(serviceName string) (bool, string)
}
