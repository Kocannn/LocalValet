package service

import (
	servicedomain "LocalValet/internal/domain/service"
	"errors"
	"testing"
)

type mockManager struct {
	runningServices map[string]bool
	allocatedPorts  map[string]int
	versions        map[string][]string
	activeVersions  map[string]string
	failStart       map[string]bool
	failStop        map[string]bool
	healthy         map[string]bool
}

func newMockManager() *mockManager {
	return &mockManager{
		runningServices: make(map[string]bool),
		allocatedPorts:  make(map[string]int),
		versions:        make(map[string][]string),
		activeVersions:  make(map[string]string),
		failStart:       make(map[string]bool),
		failStop:        make(map[string]bool),
		healthy:         make(map[string]bool),
	}
}

func (m *mockManager) StartService(serviceName string) error {
	if m.failStart[serviceName] {
		return errors.New("mock start error")
	}
	m.runningServices[serviceName] = true
	if _, exists := m.allocatedPorts[serviceName]; !exists {
		m.allocatedPorts[serviceName] = 8080
	}
	return nil
}

func (m *mockManager) StopService(serviceName string) error {
	if m.failStop[serviceName] {
		return errors.New("mock stop error")
	}
	m.runningServices[serviceName] = false
	delete(m.allocatedPorts, serviceName)
	return nil
}

func (m *mockManager) GetServiceStatus(serviceName string) (bool, string) {
	running := m.runningServices[serviceName]
	if running {
		return true, serviceName + " is running"
	}
	return false, serviceName + " is stopped"
}

func (m *mockManager) CheckHealth(serviceName string) (bool, string) {
	if !m.runningServices[serviceName] {
		return false, serviceName + " is stopped"
	}
	if isHealthy, ok := m.healthy[serviceName]; ok {
		return isHealthy, serviceName + " health status"
	}
	return true, serviceName + " is healthy"
}

func (m *mockManager) GetAllocatedPort(serviceName string) int {
	return m.allocatedPorts[serviceName]
}

func (m *mockManager) SetServiceVersion(serviceName, version string) error {
	if serviceName == "error_svc" {
		return errors.New("mock version error")
	}
	m.activeVersions[serviceName] = version
	return nil
}

func (m *mockManager) GetActiveServiceVersion(serviceName string) (string, error) {
	return m.activeVersions[serviceName], nil
}

func (m *mockManager) GetAvailableVersions(serviceName string) ([]string, error) {
	return m.versions[serviceName], nil
}

func TestUseCase_ServiceNamesAndConfigs(t *testing.T) {
	configs := servicedomain.DefaultConfigs()
	mgr := newMockManager()
	uc := New(mgr, configs)

	names := uc.ServiceNames()
	if len(names) != len(configs) {
		t.Errorf("expected %d names, got %d", len(configs), len(names))
	}

	cfgs := uc.Configs()
	if len(cfgs) != len(configs) {
		t.Errorf("expected %d configs, got %d", len(configs), len(cfgs))
	}
}

func TestUseCase_StartupLog(t *testing.T) {
	uc := New(newMockManager(), servicedomain.DefaultConfigs())
	logMsg := uc.StartupLog("linux", "isolated runtime")

	if logMsg.Level != "info" {
		t.Errorf("expected level info, got %s", logMsg.Level)
	}
	if logMsg.Message == "" {
		t.Errorf("expected non-empty startup log message")
	}
}

func TestUseCase_StartWithDependency(t *testing.T) {
	mgr := newMockManager()
	configs := []servicedomain.Config{
		{ServiceName: "php-fpm", DefaultPort: 9074, Dependencies: []string{}},
		{ServiceName: "apache", DefaultPort: 8080, Dependencies: []string{"php-fpm"}},
	}
	uc := New(mgr, configs)

	// Start Apache -> PHP-FPM dependency should automatically start first
	logMsg := uc.StartService("apache")
	if logMsg.Level != "success" {
		t.Errorf("expected success starting apache, got %s: %s", logMsg.Level, logMsg.Message)
	}

	if !mgr.runningServices["php-fpm"] {
		t.Errorf("expected php-fpm dependency to be started automatically")
	}
	if !mgr.runningServices["apache"] {
		t.Errorf("expected apache to be running")
	}
}

func TestUseCase_StartFailure(t *testing.T) {
	mgr := newMockManager()
	mgr.failStart["fail_svc"] = true
	configs := []servicedomain.Config{
		{ServiceName: "fail_svc"},
	}
	uc := New(mgr, configs)

	logMsg := uc.StartService("fail_svc")
	if logMsg.Level != "error" {
		t.Errorf("expected error log level, got %s", logMsg.Level)
	}
}

func TestUseCase_StopAndToggle(t *testing.T) {
	mgr := newMockManager()
	configs := servicedomain.DefaultConfigs()
	uc := New(mgr, configs)

	// Toggle Start
	startMsg := uc.ToggleService("redis", true)
	if startMsg.Level != "success" || !mgr.runningServices["redis"] {
		t.Errorf("expected redis to start via toggle")
	}

	// Toggle Stop
	stopMsg := uc.ToggleService("redis", false)
	if stopMsg.Level != "success" || mgr.runningServices["redis"] {
		t.Errorf("expected redis to stop via toggle")
	}
}

func TestUseCase_StatusAndHealth(t *testing.T) {
	mgr := newMockManager()
	configs := servicedomain.DefaultConfigs()
	uc := New(mgr, configs)

	// Stopped status
	status := uc.GetServiceStatus("mysql")
	if status.IsRunning {
		t.Errorf("expected mysql to be stopped initially")
	}

	// Start and check status
	_ = uc.StartService("mysql")
	status = uc.GetServiceStatus("mysql")
	if !status.IsRunning {
		t.Errorf("expected mysql to be running")
	}
	if status.Category != "Database" {
		t.Errorf("expected Category Database, got %s", status.Category)
	}

	// Health check
	healthy, healthLog := uc.CheckHealth("mysql")
	if !healthy || healthLog.Level != "success" {
		t.Errorf("expected healthy status for running mysql")
	}

	// All services status
	allStatuses := uc.GetAllServicesStatus(uc.ServiceNames())
	if len(allStatuses) != len(configs) {
		t.Errorf("expected %d statuses, got %d", len(configs), len(allStatuses))
	}
}

func TestUseCase_VersionManagement(t *testing.T) {
	mgr := newMockManager()
	mgr.versions["php-fpm"] = []string{"8.3", "8.4"}
	mgr.activeVersions["php-fpm"] = "8.4"

	uc := New(mgr, servicedomain.DefaultConfigs())

	vers := uc.GetServiceVersions("php-fpm")
	if len(vers) != 2 {
		t.Errorf("expected 2 versions, got %d", len(vers))
	}

	active := uc.GetActiveServiceVersion("php-fpm")
	if active != "8.4" {
		t.Errorf("expected active 8.4, got %s", active)
	}

	switchMsg := uc.SetServiceVersion("php-fpm", "8.3")
	if switchMsg.Level != "success" {
		t.Errorf("expected success switching version, got %s: %s", switchMsg.Level, switchMsg.Message)
	}
	if mgr.activeVersions["php-fpm"] != "8.3" {
		t.Errorf("expected active version in manager to be 8.3")
	}
}

func TestUseCase_HotSwitching(t *testing.T) {
	mgr := newMockManager()
	mgr.versions["php-fpm"] = []string{"8.3", "8.4"}
	mgr.activeVersions["php-fpm"] = "8.4"
	mgr.runningServices["php-fpm"] = true // Currently running!

	uc := New(mgr, servicedomain.DefaultConfigs())

	// Hot-switch version from 8.4 to 8.3 while running
	msg := uc.SetServiceVersion("php-fpm", "8.3")
	if msg.Level != "success" {
		t.Errorf("expected success on hot-switch, got %s: %s", msg.Level, msg.Message)
	}

	// Service should still be running after hot-restart
	if !mgr.runningServices["php-fpm"] {
		t.Errorf("expected php-fpm to be running after hot-switch restart")
	}
	if mgr.activeVersions["php-fpm"] != "8.3" {
		t.Errorf("expected active version to be 8.3")
	}
}

func TestUseCase_GetAllRuntimeServices(t *testing.T) {
	mgr := newMockManager()
	mgr.activeVersions["php-fpm"] = "8.4"
	mgr.versions["php-fpm"] = []string{"8.3", "8.4"}
	mgr.activeVersions["node"] = "22"
	mgr.versions["node"] = []string{"18", "20", "22"}

	uc := New(mgr, servicedomain.DefaultConfigs())

	runtimes := uc.GetAllRuntimeServices()
	if len(runtimes) == 0 {
		t.Fatalf("expected runtime services list to not be empty")
	}

	foundPhp := false
	foundNode := false
	for _, r := range runtimes {
		if r.ServiceName == "php-fpm" {
			foundPhp = true
			if r.ActiveVersion != "8.4" {
				t.Errorf("expected PHP active version 8.4, got %s", r.ActiveVersion)
			}
		}
		if r.ServiceName == "node" {
			foundNode = true
			if r.ActiveVersion != "22" {
				t.Errorf("expected Node active version 22, got %s", r.ActiveVersion)
			}
		}
	}

	if !foundPhp || !foundNode {
		t.Errorf("expected php-fpm and node in runtime services, got %v", runtimes)
	}
}

