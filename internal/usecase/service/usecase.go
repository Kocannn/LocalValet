package service

import (
	servicedomain "LocalValet/internal/domain/service"
	"fmt"
	"time"
)

// LogMessage represents a user-facing service operation log.
type LogMessage struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

type UseCase struct {
	manager servicedomain.Manager
	configs []servicedomain.Config
}

func New(manager servicedomain.Manager, configs []servicedomain.Config) *UseCase {
	return &UseCase{
		manager: manager,
		configs: append([]servicedomain.Config(nil), configs...),
	}
}

func (u *UseCase) Configs() []servicedomain.Config {
	return append([]servicedomain.Config(nil), u.configs...)
}

func (u *UseCase) ServiceNames() []string {
	names := make([]string, 0, len(u.configs))
	for _, cfg := range u.configs {
		names = append(names, cfg.ServiceName)
	}
	return names
}

func (u *UseCase) StartupLog(goos string, binarySource string) LogMessage {
	return LogMessage{
		Timestamp: time.Now().Format("15:04:05"),
		Level:     "info",
		Message:   fmt.Sprintf("LocalValet started on %s using %s", goos, binarySource),
	}
}

func (u *UseCase) GetServiceStatus(serviceName string) servicedomain.Status {
	isRunning, msg := u.manager.GetServiceStatus(serviceName)
	return servicedomain.Status{
		Name:      serviceName,
		IsRunning: isRunning,
		Message:   msg,
	}
}

func (u *UseCase) GetAllServicesStatus(serviceNames []string) []servicedomain.Status {
	statuses := make([]servicedomain.Status, 0, len(serviceNames))
	for _, serviceName := range serviceNames {
		statuses = append(statuses, u.GetServiceStatus(serviceName))
	}
	return statuses
}

func (u *UseCase) StartService(serviceName string) LogMessage {
	err := u.manager.StartService(serviceName)
	timestamp := time.Now().Format("15:04:05")

	if err != nil {
		return LogMessage{
			Timestamp: timestamp,
			Level:     "error",
			Message:   fmt.Sprintf("Failed to start %s: %v", serviceName, err),
		}
	}

	return LogMessage{
		Timestamp: timestamp,
		Level:     "success",
		Message:   fmt.Sprintf("%s started successfully", serviceName),
	}
}

func (u *UseCase) StopService(serviceName string) LogMessage {
	err := u.manager.StopService(serviceName)
	timestamp := time.Now().Format("15:04:05")

	if err != nil {
		return LogMessage{
			Timestamp: timestamp,
			Level:     "error",
			Message:   fmt.Sprintf("Failed to stop %s: %v", serviceName, err),
		}
	}

	return LogMessage{
		Timestamp: timestamp,
		Level:     "success",
		Message:   fmt.Sprintf("%s stopped successfully", serviceName),
	}
}

func (u *UseCase) ToggleService(serviceName string, shouldStart bool) LogMessage {
	if shouldStart {
		return u.StartService(serviceName)
	}

	return u.StopService(serviceName)
}

func (u *UseCase) GetServiceVersions(serviceName string) []string {
	versions, err := u.manager.GetAvailableVersions(serviceName)
	if err != nil {
		return []string{}
	}

	return versions
}

func (u *UseCase) GetServiceVersionsWithError(serviceName string) ([]string, error) {
	return u.manager.GetAvailableVersions(serviceName)
}

func (u *UseCase) GetActiveServiceVersion(serviceName string) string {
	version, err := u.manager.GetActiveServiceVersion(serviceName)
	if err != nil {
		return ""
	}

	return version
}

func (u *UseCase) SetServiceVersion(serviceName, version string) LogMessage {
	err := u.manager.SetServiceVersion(serviceName, version)
	timestamp := time.Now().Format("15:04:05")

	if err != nil {
		return LogMessage{
			Timestamp: timestamp,
			Level:     "error",
			Message:   fmt.Sprintf("Failed to set %s version to %s: %v", serviceName, version, err),
		}
	}

	return LogMessage{
		Timestamp: timestamp,
		Level:     "success",
		Message:   fmt.Sprintf("%s version switched to %s", serviceName, version),
	}
}
