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
	healthy, _ := u.manager.CheckHealth(serviceName)
	port := u.manager.GetAllocatedPort(serviceName)

	cfg, hasCfg := servicedomain.GetConfig(serviceName, u.configs)
	category := ""
	if hasCfg {
		category = cfg.Category
		if port == 0 && isRunning {
			port = cfg.DefaultPort
		}
	}

	return servicedomain.Status{
		Name:      serviceName,
		IsRunning: isRunning,
		Message:   msg,
		Port:      port,
		Healthy:   healthy,
		Category:  category,
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
	timestamp := time.Now().Format("15:04:05")

	// 1. Check and start dependencies first
	deps := servicedomain.GetDependencies(serviceName, u.configs)
	for _, dep := range deps {
		if running, _ := u.manager.GetServiceStatus(dep); !running {
			if err := u.manager.StartService(dep); err != nil {
				return LogMessage{
					Timestamp: timestamp,
					Level:     "error",
					Message:   fmt.Sprintf("Failed to start dependency %s for %s: %v", dep, serviceName, err),
				}
			}
		}
	}

	// 2. Start the primary service
	err := u.manager.StartService(serviceName)
	if err != nil {
		return LogMessage{
			Timestamp: timestamp,
			Level:     "error",
			Message:   fmt.Sprintf("Failed to start %s: %v", serviceName, err),
		}
	}

	port := u.manager.GetAllocatedPort(serviceName)
	msg := fmt.Sprintf("%s started successfully", serviceName)
	if port > 0 {
		msg = fmt.Sprintf("%s started successfully on port %d", serviceName, port)
	}

	return LogMessage{
		Timestamp: timestamp,
		Level:     "success",
		Message:   msg,
	}
}

func (u *UseCase) StopService(serviceName string) LogMessage {
	timestamp := time.Now().Format("15:04:05")

	err := u.manager.StopService(serviceName)
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

func (u *UseCase) CheckHealth(serviceName string) (bool, LogMessage) {
	timestamp := time.Now().Format("15:04:05")
	healthy, msg := u.manager.CheckHealth(serviceName)

	level := "info"
	if healthy {
		level = "success"
	} else {
		level = "warning"
	}

	return healthy, LogMessage{
		Timestamp: timestamp,
		Level:     level,
		Message:   msg,
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
	timestamp := time.Now().Format("15:04:05")

	isRunning, _ := u.manager.GetServiceStatus(serviceName)

	if isRunning {
		// Hot-switch: Gracefully restart service with the new version
		if err := u.manager.StopService(serviceName); err != nil {
			return LogMessage{
				Timestamp: timestamp,
				Level:     "error",
				Message:   fmt.Sprintf("Failed to stop %s during version switch: %v", serviceName, err),
			}
		}

		if err := u.manager.SetServiceVersion(serviceName, version); err != nil {
			return LogMessage{
				Timestamp: timestamp,
				Level:     "error",
				Message:   fmt.Sprintf("Failed to set %s version to %s: %v", serviceName, version, err),
			}
		}

		if err := u.manager.StartService(serviceName); err != nil {
			return LogMessage{
				Timestamp: timestamp,
				Level:     "error",
				Message:   fmt.Sprintf("Switched to version %s, but failed to restart %s: %v", version, serviceName, err),
			}
		}

		return LogMessage{
			Timestamp: timestamp,
			Level:     "success",
			Message:   fmt.Sprintf("%s switched to version %s and restarted successfully", serviceName, version),
		}
	}

	err := u.manager.SetServiceVersion(serviceName, version)
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

// GetAllRuntimeServices returns version information and status for all configurable runtimes.
func (u *UseCase) GetAllRuntimeServices() []servicedomain.RuntimeServiceInfo {
	serviceList := []string{"php-fpm", "node", "mysql", "nginx", "apache", "redis", "postgresql"}
	result := make([]servicedomain.RuntimeServiceInfo, 0, len(serviceList))

	for _, svcName := range serviceList {
		cfg, hasCfg := servicedomain.GetConfig(svcName, u.configs)
		displayName := svcName
		category := "Runtime"
		if hasCfg {
			displayName = cfg.DisplayName
			category = cfg.Category
		} else if svcName == "node" {
			displayName = "Node.js"
			category = "Runtime"
		}

		active, _ := u.manager.GetActiveServiceVersion(svcName)
		versions, _ := u.manager.GetAvailableVersions(svcName)
		isRunning, _ := u.manager.GetServiceStatus(svcName)

		result = append(result, servicedomain.RuntimeServiceInfo{
			ServiceName:       svcName,
			DisplayName:       displayName,
			ActiveVersion:     active,
			AvailableVersions: versions,
			Category:          category,
			IsRunning:         isRunning,
		})
	}

	return result
}

