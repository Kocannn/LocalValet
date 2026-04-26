package main

import (
	servicedomain "LocalValet/internal/domain/service"
	servicemonitor "LocalValet/internal/infrastructure/monitor"
	"LocalValet/internal/platform"
	serviceusecase "LocalValet/internal/usecase/service"
	systemusecase "LocalValet/internal/usecase/system"
	"context"
	"runtime"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx               context.Context
	serviceManager    servicedomain.Manager
	serviceUC         *serviceusecase.UseCase
	systemUC          *systemusecase.UseCase
	servicesToMonitor []string
	monitor           *servicemonitor.ServiceMonitor
	emitter           servicemonitor.EventEmitter
}

// LogMessage represents a log entry
type LogMessage = serviceusecase.LogMessage

// NewApp creates a new App application struct
func NewApp() *App {
	configs := servicedomain.DefaultConfigs()
	manager := platform.NewServiceManager()

	return &App{
		serviceManager: manager,
		serviceUC:      serviceusecase.New(manager, configs),
		systemUC:       systemusecase.New(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.servicesToMonitor = a.serviceUC.ServiceNames()

	// Create emitter
	a.emitter = servicemonitor.NewEventEmitter(ctx)

	// Create monitor instance
	a.monitor = servicemonitor.NewServiceMonitor(
		a.serviceManager,
		a.servicesToMonitor,
		a.emitter,
	)

	// Start monitoring using context
	go a.monitor.Start(ctx)

	// Startup log
	binarySource := "isolated runtime (runtime/)"
	a.emitter.Emit("service:log", a.serviceUC.StartupLog(runtime.GOOS, binarySource))
}

// GetServiceStatus checks if a service is running
func (a *App) GetServiceStatus(serviceName string) servicedomain.Status {
	return a.serviceUC.GetServiceStatus(serviceName)
}

// StartService starts a service
func (a *App) StartService(serviceName string) LogMessage {
	return a.serviceUC.StartService(serviceName)
}

// StopService stops a service
func (a *App) StopService(serviceName string) LogMessage {
	return a.serviceUC.StopService(serviceName)
}

// ToggleService toggles a service on/off
func (a *App) ToggleService(serviceName string, shouldStart bool) LogMessage {
	logMsg := a.serviceUC.ToggleService(serviceName, shouldStart)

	wailsRuntime.EventsEmit(a.ctx, "service:log", logMsg)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				finalStatus := a.GetServiceStatus(serviceName)
				wailsRuntime.EventsEmit(a.ctx, "service:status-changed", finalStatus)
				return

			case <-ticker.C:
				currentStatus := a.GetServiceStatus(serviceName)

				if currentStatus.IsRunning == shouldStart {
					wailsRuntime.EventsEmit(a.ctx, "service:status-changed", currentStatus)
					return
				}
			}
		}
	}()

	return logMsg
}

// GetAllServicesStatus returns status for all monitored services (for initial load)
func (a *App) GetAllServicesStatus() []servicedomain.Status {
	return a.serviceUC.GetAllServicesStatus(a.servicesToMonitor)
}

// GetBinarySourceInfo returns information about where binaries are executed from
func (a *App) GetBinarySourceInfo() map[string]interface{} {
	return a.systemUC.GetBinarySourceInfo(IsUsingSystemBinaries(), "runtime/")
}

// GetServiceVersions returns available versions for a service runtime.
func (a *App) GetServiceVersions(serviceName string) []string {
	return a.serviceUC.GetServiceVersions(serviceName)
}

// GetActiveServiceVersion returns active runtime version for a service.
func (a *App) GetActiveServiceVersion(serviceName string) string {
	return a.serviceUC.GetActiveServiceVersion(serviceName)
}

// SetServiceVersion updates active runtime version for a service.
func (a *App) SetServiceVersion(serviceName, version string) LogMessage {
	msg := a.serviceUC.SetServiceVersion(serviceName, version)
	if msg.Level == "success" {
		wailsRuntime.EventsEmit(a.ctx, "service:log", msg)
	}
	return msg
}
