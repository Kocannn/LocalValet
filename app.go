package main

import (
	projectdomain "LocalValet/internal/domain/project"
	servicedomain "LocalValet/internal/domain/service"
	"LocalValet/internal/infrastructure/events"
	terminalinfra "LocalValet/internal/infrastructure/terminal"
	"LocalValet/internal/platform"
	projectplatform "LocalValet/internal/platform/linux/project"
	sslplatform "LocalValet/internal/platform/linux/ssl"
	vhostplatform "LocalValet/internal/platform/linux/vhost"
	dnsplatform "LocalValet/internal/platform/linux/dns"
	projectusecase "LocalValet/internal/usecase/project"
	serviceusecase "LocalValet/internal/usecase/service"
	sslusecase "LocalValet/internal/usecase/ssl"
	systemusecase "LocalValet/internal/usecase/system"
	terminalusecase "LocalValet/internal/usecase/terminal"
	vhostusecase "LocalValet/internal/usecase/vhost"
	"context"
	"fmt"
	"log"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)



// App struct
type App struct {
	ctx               context.Context
	serviceManager    servicedomain.Manager
	serviceUC         *serviceusecase.UseCase
	systemUC          *systemusecase.UseCase
	terminalUC        *terminalusecase.UseCase
	sslUC             *sslusecase.UseCase
	vhostUC           *vhostusecase.UseCase
	projectUC         *projectusecase.UseCase
	hostsManager      *dnsplatform.HostsManager
	servicesToMonitor []string
	monitor           *servicemonitor.ServiceMonitor
	emitter           servicemonitor.EventEmitter
}

// LogMessage represents a log entry.
type LogMessage = serviceusecase.LogMessage

// NewApp creates a new App application struct.
func NewApp() *App {
	configs := servicedomain.DefaultConfigs()
	manager := platform.NewServiceManager()

	sslManager := sslplatform.NewCAManager()
	sslUC := sslusecase.New(sslManager)

	vhostGen := vhostplatform.NewNginxGenerator()
	vhostUC := vhostusecase.New(vhostGen, sslUC)

	scanner := projectplatform.NewScanner()
	repo := projectplatform.NewRepository()
	projectUC := projectusecase.New(scanner, repo, vhostUC, sslUC)
	hostsManager := dnsplatform.NewHostsManager()

	return &App{
		serviceManager: manager,
		serviceUC:      serviceusecase.New(manager, configs),
		systemUC:       systemusecase.New(),
		terminalUC:     terminalusecase.New(terminalinfra.NewLinuxManager()),
		sslUC:          sslUC,
		vhostUC:        vhostUC,
		projectUC:      projectUC,
		hostsManager:   hostsManager,
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

	// Start monitor
	go a.monitor.Start(ctx)

	// Initial scan for projects in background
	go func() {
		time.Sleep(1 * time.Second)
		if a.projectUC != nil {
			_, _ = a.projectUC.ScanProjects()
		}
	}()
}

// shutdown is called when the app terminates
func (a *App) shutdown(ctx context.Context) {
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return "Hello " + name + ", It's show time!"
}

// GetServiceStatus returns status for a single service
func (a *App) GetServiceStatus(serviceName string) servicedomain.Status {
	return a.controller.GetServiceStatus(serviceName)
}

// CheckServiceHealth returns health status for a single service
func (a *App) CheckServiceHealth(serviceName string) (bool, LogMessage) {
	return a.serviceUC.CheckHealth(serviceName)
}

// GetAllServicesStatus returns status for all services
func (a *App) GetAllServicesStatus() []servicedomain.Status {
	return a.serviceUC.GetAllServicesStatus(a.servicesToMonitor)
}

// StartService starts a service
func (a *App) StartService(serviceName string) LogMessage {
	msg := a.serviceUC.StartService(serviceName)
	if a.emitter != nil {
		a.emitter.Emit("service:log", msg)
	}
	return msg
}

// StopService stops a service.
func (a *App) StopService(serviceName string) LogMessage {
	msg := a.serviceUC.StopService(serviceName)
	if a.emitter != nil {
		a.emitter.Emit("service:log", msg)
	}
	return msg
}

// ToggleService toggles a service state
func (a *App) ToggleService(serviceName string, shouldStart bool) LogMessage {
	msg := a.serviceUC.ToggleService(serviceName, shouldStart)
	if a.emitter != nil {
		a.emitter.Emit("service:log", msg)
	}
	return msg
}


// GetServiceVersions returns available versions for a service
func (a *App) GetServiceVersions(serviceName string) []string {
	versions, err := a.serviceUC.GetServiceVersionsWithError(serviceName)
	if err != nil {
		log.Printf("[versions] error for %s: %v", serviceName, err)
		return []string{}
	}

	return versions
}

// GetActiveServiceVersion returns active runtime version for a service.
func (a *App) GetActiveServiceVersion(serviceName string) string {
	return a.controller.GetActiveServiceVersion(serviceName)
}

// SetServiceVersion updates active runtime version for a service.
func (a *App) SetServiceVersion(serviceName, version string) LogMessage {
	msg := a.serviceUC.SetServiceVersion(serviceName, version)
	if msg.Level == "success" && a.ctx != nil {
		wailsRuntime.EventsEmit(a.ctx, "service:log", msg)
	}
	return msg
}

// GetAllRuntimeServices returns version and status information for all runtime services.
func (a *App) GetAllRuntimeServices() []servicedomain.RuntimeServiceInfo {
	return a.serviceUC.GetAllRuntimeServices()
}

// OpenContextTerminal launches a new terminal window with runtime-injected environment.
func (a *App) OpenContextTerminal(projectDir string) LogMessage {
	return a.controller.OpenContextTerminal(projectDir)
}

// -------------------------------------------------------------
// Phase 4: Projects, Virtual Hosts & SSL Wails Bindings
// -------------------------------------------------------------

// GetProjects returns all discovered projects.
func (a *App) GetProjects() []projectdomain.Project {
	projects, err := a.projectUC.GetProjects()
	if err != nil {
		log.Printf("[projects] GetProjects error: %v", err)
		return []projectdomain.Project{}
	}
	return projects
}

// ScanProjects triggers a full rescan of project root directories.
func (a *App) ScanProjects() []projectdomain.Project {
	projects, err := a.projectUC.ScanProjects()
	if err != nil {
		log.Printf("[projects] ScanProjects error: %v", err)
		return []projectdomain.Project{}
	}
	if a.ctx != nil {
		wailsRuntime.EventsEmit(a.ctx, "project:scanned", projects)
	}
	return projects
}

// GetProjectRoots returns configured scan directories.
func (a *App) GetProjectRoots() []string {
	roots, err := a.projectUC.GetProjectRoots()
	if err != nil {
		return []string{}
	}
	return roots
}

// AddProjectRoot adds a new scan directory.
func (a *App) AddProjectRoot(path string) LogMessage {
	timestamp := time.Now().Format("15:04:05")
	if err := a.projectUC.AddProjectRoot(path); err != nil {
		return LogMessage{
			Timestamp: timestamp,
			Level:     "error",
			Message:   fmt.Sprintf("Failed to add project root %s: %v", path, err),
		}
	}
	return LogMessage{
		Timestamp: timestamp,
		Level:     "success",
		Message:   fmt.Sprintf("Added project root: %s", path),
	}
}

// RemoveProjectRoot removes a scan directory.
func (a *App) RemoveProjectRoot(path string) LogMessage {
	timestamp := time.Now().Format("15:04:05")
	if err := a.projectUC.RemoveProjectRoot(path); err != nil {
		return LogMessage{
			Timestamp: timestamp,
			Level:     "error",
			Message:   fmt.Sprintf("Failed to remove project root %s: %v", path, err),
		}
	}
	return LogMessage{
		Timestamp: timestamp,
		Level:     "success",
		Message:   fmt.Sprintf("Removed project root: %s", path),
	}
}

// ToggleProjectVHost toggles virtual host generation for a project.
func (a *App) ToggleProjectVHost(projectPath string, enable bool) LogMessage {
	timestamp := time.Now().Format("15:04:05")
	if err := a.projectUC.ToggleProjectVHost(projectPath, enable); err != nil {
		return LogMessage{
			Timestamp: timestamp,
			Level:     "error",
			Message:   fmt.Sprintf("Failed to toggle VHost: %v", err),
		}
	}
	statusStr := "enabled"
	if !enable {
		statusStr = "disabled"
	}
	return LogMessage{
		Timestamp: timestamp,
		Level:     "success",
		Message:   fmt.Sprintf("VHost %s for project", statusStr),
	}
}

// GenerateProjectSSL generates an SSL certificate for the given project.
func (a *App) GenerateProjectSSL(projectPath string) LogMessage {
	timestamp := time.Now().Format("15:04:05")
	if err := a.projectUC.GenerateProjectSSL(projectPath); err != nil {
		return LogMessage{
			Timestamp: timestamp,
			Level:     "error",
			Message:   fmt.Sprintf("Failed to generate SSL: %v", err),
		}
	}
	return LogMessage{
		Timestamp: timestamp,
		Level:     "success",
		Message:   "Generated SSL certificate successfully",
	}
}

// OpenProjectInEditor opens project in code editor.
func (a *App) OpenProjectInEditor(projectPath, editor string) LogMessage {
	timestamp := time.Now().Format("15:04:05")
	if err := a.projectUC.OpenInEditor(projectPath, editor); err != nil {
		return LogMessage{
			Timestamp: timestamp,
			Level:     "error",
			Message:   fmt.Sprintf("Could not open in editor: %v", err),
		}
	}
	return LogMessage{
		Timestamp: timestamp,
		Level:     "success",
		Message:   fmt.Sprintf("Opened %s in editor", projectPath),
	}
}

// OpenProjectInBrowser opens URL in web browser.
func (a *App) OpenProjectInBrowser(url string) LogMessage {
	timestamp := time.Now().Format("15:04:05")
	if err := a.projectUC.OpenInBrowser(url); err != nil {
		return LogMessage{
			Timestamp: timestamp,
			Level:     "error",
			Message:   fmt.Sprintf("Could not open in browser: %v", err),
		}
	}
	return LogMessage{
		Timestamp: timestamp,
		Level:     "info",
		Message:   fmt.Sprintf("Opened %s in browser", url),
	}
}

// TrustRootCA installs the LocalValet Root CA into the system trust store.
func (a *App) TrustRootCA() LogMessage {
	timestamp := time.Now().Format("15:04:05")
	if a.sslUC == nil {
		return LogMessage{
			Timestamp: timestamp,
			Level:     "error",
			Message:   "SSL manager not initialized",
		}
	}

	if err := a.sslUC.InstallRootCA(); err != nil {
		return LogMessage{
			Timestamp: timestamp,
			Level:     "error",
			Message:   fmt.Sprintf("Failed to install Root CA: %v", err),
		}
	}

	return LogMessage{
		Timestamp: timestamp,
		Level:     "success",
		Message:   "LocalValet Root CA trusted in system successfully",
	}
}

// IsRootCATrusted checks if the LocalValet Root CA is trusted in the system store.
func (a *App) IsRootCATrusted() bool {
	if a.sslUC == nil {
		return false
	}
	return a.sslUC.IsRootCATrusted()
}

// SyncHostsDomains synchronizes all discovered project domains to /etc/hosts.
func (a *App) SyncHostsDomains() LogMessage {
	timestamp := time.Now().Format("15:04:05")
	if a.projectUC == nil || a.hostsManager == nil {
		return LogMessage{
			Timestamp: timestamp,
			Level:     "error",
			Message:   "Hosts manager not initialized",
		}
	}

	domains, err := a.projectUC.GetAllDomains()
	if err != nil {
		return LogMessage{
			Timestamp: timestamp,
			Level:     "error",
			Message:   fmt.Sprintf("Failed to get domains: %v", err),
		}
	}

	if err := a.hostsManager.SyncDomains(domains); err != nil {
		return LogMessage{
			Timestamp: timestamp,
			Level:     "error",
			Message:   fmt.Sprintf("Failed to sync /etc/hosts: %v", err),
		}
	}

	return LogMessage{
		Timestamp: timestamp,
		Level:     "success",
		Message:   fmt.Sprintf("Synchronized %d domains to /etc/hosts", len(domains)),
	}
}

