package main

import (
	"LocalValet/internal/platform"
	"LocalValet/internal/platform/domain"
	"context"
	"fmt"
	"runtime"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx               context.Context
	monitoringActive  bool
	servicesToMonitor []string

	serviceManager domain.ServiceManager
}

// ServiceStatus represents the status of a service
type ServiceStatus struct {
	Name      string `json:"name"`
	IsRunning bool   `json:"isRunning"`
	Message   string `json:"message"`
}

// LogMessage represents a log entry
type LogMessage struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		serviceManager: platform.NewServiceManager(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.monitoringActive = true

	// Define services to monitor based on configured services
	configs := GetServiceConfigs()
	a.servicesToMonitor = make([]string, 0, len(configs))

	for _, config := range configs {
		serviceName := GetServiceName(config.DisplayName)
		a.servicesToMonitor = append(a.servicesToMonitor, serviceName)
	}

	// Start background service monitor
	go a.StartServiceMonitor()

	// Send initial log with binary source info
	binarySource := "system binaries"
	if runtime.GOOS == "windows" {
		binarySource = "custom binaries (bin/windows/)"
	}

	wailsRuntime.EventsEmit(ctx, "service:log", LogMessage{
		Timestamp: time.Now().Format("15:04:05"),
		Level:     "info",
		Message:   fmt.Sprintf("LocalValet started on %s using %s", runtime.GOOS, binarySource),
	})
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetServiceStatus checks if a service is running
func (a *App) GetServiceStatus(serviceName string) ServiceStatus {
	isRunning, msg := a.serviceManager.GetServiceStatus(serviceName)

	return ServiceStatus{
		Name:      serviceName,
		IsRunning: isRunning,
		Message:   msg,
	}
}

// StartService starts a service
func (a *App) StartService(serviceName string) LogMessage {

	err := a.serviceManager.StartService(serviceName)
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

// StopService stops a service
func (a *App) StopService(serviceName string) LogMessage {

	timestamp := time.Now().Format("15:04:05")
	err := a.serviceManager.StopService(serviceName)

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

// ToggleService toggles a service on/off
func (a *App) ToggleService(serviceName string, shouldStart bool) LogMessage {
	var logMsg LogMessage

	if shouldStart {
		logMsg = a.StartService(serviceName)
	} else {
		logMsg = a.StopService(serviceName)
	}

	// Emit log event to frontend
	wailsRuntime.EventsEmit(a.ctx, "service:log", logMsg)

	// Trigger immediate status check for this service
	go func() {
		time.Sleep(500 * time.Millisecond) // Wait for service to start/stop
		status := a.GetServiceStatus(serviceName)
		wailsRuntime.EventsEmit(a.ctx, "service:status-changed", status)
	}()

	return logMsg
}

// StartServiceMonitor runs in background and monitors service status changes
func (a *App) StartServiceMonitor() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// Store last known status for each service
	lastStatus := make(map[string]bool)

	// Initialize with current status
	for _, serviceName := range a.servicesToMonitor {
		status := a.GetServiceStatus(serviceName)
		lastStatus[serviceName] = status.IsRunning
	}

	for range ticker.C {
		if !a.monitoringActive {
			break
		}

		// Check each service
		for _, serviceName := range a.servicesToMonitor {
			status := a.GetServiceStatus(serviceName)

			// Emit event ONLY if status has changed
			if lastStatus[serviceName] != status.IsRunning {
				wailsRuntime.EventsEmit(a.ctx, "service:status-changed", status)

				// Log the change
				logLevel := "info"
				if status.IsRunning {
					logLevel = "success"
				} else {
					logLevel = "warning"
				}

				wailsRuntime.EventsEmit(a.ctx, "service:log", LogMessage{
					Timestamp: time.Now().Format("15:04:05"),
					Level:     logLevel,
					Message:   fmt.Sprintf("%s status changed: %s", serviceName, status.Message),
				})

				lastStatus[serviceName] = status.IsRunning
			}
		}
	}
}

// GetAllServicesStatus returns status for all monitored services (for initial load)
func (a *App) GetAllServicesStatus() []ServiceStatus {
	statuses := make([]ServiceStatus, 0, len(a.servicesToMonitor))

	for _, serviceName := range a.servicesToMonitor {
		status := a.GetServiceStatus(serviceName)
		statuses = append(statuses, status)
	}

	return statuses
}

// GetBinarySourceInfo returns information about where binaries are executed from
func (a *App) GetBinarySourceInfo() map[string]interface{} {
	info := make(map[string]interface{})

	info["os"] = runtime.GOOS
	info["using_system_binaries"] = IsUsingSystemBinaries()

	if runtime.GOOS == "windows" {
		info["binary_location"] = "bin/windows/"
		info["binary_validation"] = ValidateWindowsBinaries()
	} else {
		info["binary_location"] = "system"
		info["binary_validation"] = nil
	}

	return info
}
