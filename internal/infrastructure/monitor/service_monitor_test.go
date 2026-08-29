package servicemonitor

import (
	"context"
	"sync"

	"testing"
	"time"
)

type mockMonitorManager struct {
	mu     sync.Mutex
	status map[string]bool
}

func (m *mockMonitorManager) StartService(serviceName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.status[serviceName] = true
	return nil
}

func (m *mockMonitorManager) StopService(serviceName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.status[serviceName] = false
	return nil
}

func (m *mockMonitorManager) GetServiceStatus(serviceName string) (bool, string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	running := m.status[serviceName]
	if running {
		return true, serviceName + " is running"
	}
	return false, serviceName + " is stopped"
}

func (m *mockMonitorManager) CheckHealth(serviceName string) (bool, string) {
	return true, "healthy"
}

func (m *mockMonitorManager) GetAllocatedPort(serviceName string) int {
	return 8080
}

func (m *mockMonitorManager) SetServiceVersion(serviceName, version string) error {
	return nil
}

func (m *mockMonitorManager) GetActiveServiceVersion(serviceName string) (string, error) {
	return "1.0", nil
}

func (m *mockMonitorManager) GetAvailableVersions(serviceName string) ([]string, error) {
	return []string{"1.0"}, nil
}

type mockEmitter struct {
	mu     sync.Mutex
	events []struct {
		event string
		data  interface{}
	}
}

func (e *mockEmitter) Emit(event string, data interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.events = append(e.events, struct {
		event string
		data  interface{}
	}{event: event, data: data})
}

func TestServiceMonitor_StatusChangeDetection(t *testing.T) {
	mgr := &mockMonitorManager{
		status: map[string]bool{"redis": false},
	}
	emitter := &mockEmitter{}

	monitor := NewServiceMonitor(mgr, []string{"redis"}, emitter)
	monitor.intervalFast = 50 * time.Millisecond
	monitor.intervalSlow = 50 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go monitor.Start(ctx)

	// Wait for monitor loop to start
	time.Sleep(100 * time.Millisecond)

	// Change service status
	mgr.mu.Lock()
	mgr.status["redis"] = true
	mgr.mu.Unlock()

	// Wait for monitor to detect change
	time.Sleep(200 * time.Millisecond)

	emitter.mu.Lock()
	eventCount := len(emitter.events)
	emitter.mu.Unlock()

	if eventCount == 0 {
		t.Errorf("expected ServiceMonitor to emit events on status change")
	}
}
