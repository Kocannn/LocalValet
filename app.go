package main

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// App struct
type App struct {
	ctx context.Context
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
	if shouldStart {
		return a.StartService(serviceName)
	}
	return a.StopService(serviceName)
}
