package main

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx               context.Context
	monitoringActive  bool
	servicesToMonitor []string
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
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.monitoringActive = true

	// Define services to monitor
	a.servicesToMonitor = []string{
		"apache2",
		"mysql",
		"postgresql",
		"redis-server",
		"nginx",
		"php8.1-fpm",
	}

	// Start background service monitor
	go a.StartServiceMonitor()

	// Send initial log
	wailsRuntime.EventsEmit(ctx, "service:log", LogMessage{
		Timestamp: time.Now().Format("15:04:05"),
		Level:     "info",
		Message:   "LocalValet event-driven monitoring started",
	})
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetServiceStatus checks if a service is running
func (a *App) GetServiceStatus(serviceName string) ServiceStatus {
	isRunning := false
	message := ""

	switch runtime.GOOS {
	case "linux":
		// Use systemctl to check service status
		cmd := exec.Command("systemctl", "is-active", serviceName)
		output, err := cmd.Output()
		if err == nil && strings.TrimSpace(string(output)) == "active" {
			isRunning = true
			message = fmt.Sprintf("%s is running", serviceName)
		} else {
			message = fmt.Sprintf("%s is stopped", serviceName)
		}
	case "darwin":
		// macOS - check with brew services
		cmd := exec.Command("brew", "services", "list")
		output, err := cmd.Output()
		if err == nil {
			if strings.Contains(string(output), serviceName) && strings.Contains(string(output), "started") {
				isRunning = true
				message = fmt.Sprintf("%s is running", serviceName)
			} else {
				message = fmt.Sprintf("%s is stopped", serviceName)
			}
		}
	case "windows":
		// Windows - check with sc query
		cmd := exec.Command("sc", "query", serviceName)
		output, err := cmd.Output()
		if err == nil && strings.Contains(string(output), "RUNNING") {
			isRunning = true
			message = fmt.Sprintf("%s is running", serviceName)
		} else {
			message = fmt.Sprintf("%s is stopped", serviceName)
		}
	}

	return ServiceStatus{
		Name:      serviceName,
		IsRunning: isRunning,
		Message:   message,
	}
}

// StartService starts a service
func (a *App) StartService(serviceName string) LogMessage {
	var cmd *exec.Cmd
	var err error

	switch runtime.GOOS {
	case "linux":
		// Use systemctl to start service
		cmd = exec.Command("sudo", "systemctl", "start", serviceName)
		err = cmd.Run()
	case "darwin":
		// macOS - use brew services
		cmd = exec.Command("brew", "services", "start", serviceName)
		err = cmd.Run()
	case "windows":
		// Windows - use net start
		cmd = exec.Command("net", "start", serviceName)
		err = cmd.Run()
	}

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
	var cmd *exec.Cmd
	var err error

	switch runtime.GOOS {
	case "linux":
		// Use systemctl to stop service
		cmd = exec.Command("sudo", "systemctl", "stop", serviceName)
		err = cmd.Run()
	case "darwin":
		// macOS - use brew services
		cmd = exec.Command("brew", "services", "stop", serviceName)
		err = cmd.Run()
	case "windows":
		// Windows - use net stop
		cmd = exec.Command("net", "stop", serviceName)
		err = cmd.Run()
	}

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
