package main

import (
	"LocalValet/internal/application"
	servicedomain "LocalValet/internal/domain/service"
	"LocalValet/internal/infrastructure/events"
	terminalinfra "LocalValet/internal/infrastructure/terminal"
	"LocalValet/internal/platform"
	serviceusecase "LocalValet/internal/usecase/service"
	"context"
)

// App struct
type App struct {
	controller *application.Controller
}

// LogMessage represents a log entry.
type LogMessage = serviceusecase.LogMessage

// NewApp creates a new App application struct.
func NewApp() *App {
	configs := servicedomain.DefaultConfigs()
	manager := platform.NewServiceManager()

	return &App{
		controller: application.New(manager, configs, terminalinfra.NewLinuxManager()),
	}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.controller.Startup(ctx, events.NewWailsPublisher(ctx))
}

// GetServiceStatus checks if a service is running.
func (a *App) GetServiceStatus(serviceName string) servicedomain.Status {
	return a.controller.GetServiceStatus(serviceName)
}

// StartService starts a service.
func (a *App) StartService(serviceName string) LogMessage {
	return a.controller.StartService(serviceName)
}

// StopService stops a service.
func (a *App) StopService(serviceName string) LogMessage {
	return a.controller.StopService(serviceName)
}

// CheckServiceHealth checks the health of a service
func (a *App) CheckServiceHealth(serviceName string) LogMessage {
	_, msg := a.serviceUC.CheckHealth(serviceName)
	return msg
}


// ToggleService toggles a service on/off
func (a *App) ToggleService(serviceName string, shouldStart bool) LogMessage {
	return a.controller.ToggleService(serviceName, shouldStart)
}

// GetAllServicesStatus returns status for all monitored services.
func (a *App) GetAllServicesStatus() []servicedomain.Status {
	return a.controller.GetAllServicesStatus()
}

// GetBinarySourceInfo returns information about where binaries are executed from.
func (a *App) GetBinarySourceInfo() map[string]interface{} {
	return a.controller.GetBinarySourceInfo(IsUsingSystemBinaries(), "runtime/")
}

// GetServiceVersions returns available versions for a service runtime.
func (a *App) GetServiceVersions(serviceName string) []string {
	return a.controller.GetServiceVersions(serviceName)
}

// GetActiveServiceVersion returns active runtime version for a service.
func (a *App) GetActiveServiceVersion(serviceName string) string {
	return a.controller.GetActiveServiceVersion(serviceName)
}

// SetServiceVersion updates active runtime version for a service.
func (a *App) SetServiceVersion(serviceName, version string) LogMessage {
	return a.controller.SetServiceVersion(serviceName, version)
}

// OpenContextTerminal launches a new terminal window with runtime-injected environment.
func (a *App) OpenContextTerminal(projectDir string) LogMessage {
	return a.controller.OpenContextTerminal(projectDir)
}

// LaunchTerminal launches external terminal with LocalValet context.
func (a *App) LaunchTerminal(projectPath string) LogMessage {
	return a.controller.LaunchTerminal(projectPath)
}
