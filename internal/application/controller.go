package application

import (
	servicedomain "LocalValet/internal/domain/service"
	terminaldomain "LocalValet/internal/domain/terminal"
	servicemonitor "LocalValet/internal/infrastructure/monitor"
	serviceusecase "LocalValet/internal/usecase/service"
	systemusecase "LocalValet/internal/usecase/system"
	terminalusecase "LocalValet/internal/usecase/terminal"
	"context"
	"log"
	"runtime"
	"time"
)

type EventPublisher interface {
	Emit(event string, data interface{})
}

type Controller struct {
	serviceManager    servicedomain.Manager
	serviceUC         *serviceusecase.UseCase
	systemUC          *systemusecase.UseCase
	terminalUC        *terminalusecase.UseCase
	servicesToMonitor []string
	monitor           *servicemonitor.ServiceMonitor
	publisher         EventPublisher
	ctx               context.Context
}

func New(
	serviceManager servicedomain.Manager,
	serviceConfigs []servicedomain.Config,
	terminalManager terminaldomain.Manager,
) *Controller {
	return &Controller{
		serviceManager: serviceManager,
		serviceUC:      serviceusecase.New(serviceManager, serviceConfigs),
		systemUC:       systemusecase.New(),
		terminalUC:     terminalusecase.New(terminalManager),
	}
}

func (c *Controller) Startup(ctx context.Context, publisher EventPublisher) {
	c.ctx = ctx
	c.publisher = publisher
	c.servicesToMonitor = c.serviceUC.ServiceNames()
	c.monitor = servicemonitor.NewServiceMonitor(c.serviceManager, c.servicesToMonitor, publisher)

	go c.monitor.Start(ctx)

	c.emit("service:log", c.serviceUC.StartupLog(runtime.GOOS, "isolated runtime (runtime/)"))
}

func (c *Controller) GetServiceStatus(serviceName string) servicedomain.Status {
	return c.serviceUC.GetServiceStatus(serviceName)
}

func (c *Controller) StartService(serviceName string) serviceusecase.LogMessage {
	return c.serviceUC.StartService(serviceName)
}

func (c *Controller) StopService(serviceName string) serviceusecase.LogMessage {
	return c.serviceUC.StopService(serviceName)
}

func (c *Controller) ToggleService(serviceName string, shouldStart bool) serviceusecase.LogMessage {
	logMsg := c.serviceUC.ToggleService(serviceName, shouldStart)
	c.emit("service:log", logMsg)

	go c.waitForServiceStatus(serviceName, shouldStart)

	return logMsg
}

func (c *Controller) GetAllServicesStatus() []servicedomain.Status {
	return c.serviceUC.GetAllServicesStatus(c.servicesToMonitor)
}

func (c *Controller) GetBinarySourceInfo(usingSystemBinaries bool, binaryLocation string) map[string]interface{} {
	return c.systemUC.GetBinarySourceInfo(usingSystemBinaries, binaryLocation)
}

func (c *Controller) GetServiceVersions(serviceName string) []string {
	versions, err := c.serviceUC.GetServiceVersionsWithError(serviceName)
	if err != nil {
		log.Printf("[versions] failed to fetch versions for %s: %v", serviceName, err)
		c.emit("service:log", serviceusecase.LogMessage{
			Timestamp: time.Now().Format("15:04:05"),
			Level:     "warning",
			Message:   "Failed to fetch versions for " + serviceName + ": " + err.Error(),
		})
		return []string{}
	}

	log.Printf("[versions] fetched %d versions for %s: %v", len(versions), serviceName, versions)
	c.emit("service:log", serviceusecase.LogMessage{
		Timestamp: time.Now().Format("15:04:05"),
		Level:     "info",
		Message:   "Fetched versions for " + serviceName + ": " + formatVersions(versions),
	})

	return versions
}

func (c *Controller) GetActiveServiceVersion(serviceName string) string {
	return c.serviceUC.GetActiveServiceVersion(serviceName)
}

func (c *Controller) SetServiceVersion(serviceName, version string) serviceusecase.LogMessage {
	msg := c.serviceUC.SetServiceVersion(serviceName, version)
	if msg.Level == "success" {
		c.emit("service:log", msg)
	}
	return msg
}

func (c *Controller) OpenContextTerminal(projectDir string) serviceusecase.LogMessage {
	msg := c.terminalUC.LaunchTerminal(projectDir, "")
	c.emit("service:log", msg)
	return serviceusecase.LogMessage(msg)
}

func (c *Controller) LaunchTerminal(projectPath string) serviceusecase.LogMessage {
	msg := c.terminalUC.LaunchTerminal(projectPath, "")
	c.emit("service:log", msg)
	return serviceusecase.LogMessage(msg)
}

func (c *Controller) CheckSudoAccess() bool {
	return c.systemUC.CheckSudoAccess()
}

func (c *Controller) GetSystemInfo() map[string]string {
	return c.systemUC.GetSystemInfo()
}

func (c *Controller) GetAvailableServices() []servicedomain.Config {
	return c.serviceUC.Configs()
}

func (c *Controller) Services() []string {
	return c.serviceUC.ServiceNames()
}

func (c *Controller) emit(event string, data interface{}) {
	if c.publisher == nil {
		return
	}

	c.publisher.Emit(event, data)
}

func (c *Controller) waitForServiceStatus(serviceName string, shouldStart bool) {
	if c.ctx == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			finalStatus := c.GetServiceStatus(serviceName)
			c.emit("service:status-changed", finalStatus)
			return

		case <-ticker.C:
			currentStatus := c.GetServiceStatus(serviceName)
			if currentStatus.IsRunning == shouldStart {
				c.emit("service:status-changed", currentStatus)
				return
			}
		}
	}
}

func formatVersions(versions []string) string {
	if len(versions) == 0 {
		return "(none)"
	}

	result := versions[0]
	for i := 1; i < len(versions); i++ {
		result += ", " + versions[i]
	}

	return result
}
