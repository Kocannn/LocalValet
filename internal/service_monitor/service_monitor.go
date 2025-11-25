package servicemonitor

import (
	"context"
	"fmt"
	"time"

	"LocalValet/internal/platform/domain"
)

type ServiceMonitor struct {
	manager  domain.ServiceManager
	services []string
	emitter  EventEmitter

	intervalFast time.Duration
	intervalSlow time.Duration

	lastStatus map[string]bool
	lastEmit   map[string]time.Time
	lastLog    map[string]time.Time
}

func NewServiceMonitor(
	manager domain.ServiceManager,
	services []string,
	emitter EventEmitter,
) *ServiceMonitor {
	return &ServiceMonitor{
		manager:      manager,
		services:     services,
		emitter:      emitter,
		intervalFast: 300 * time.Millisecond,
		intervalSlow: 5 * time.Second,
		lastStatus:   make(map[string]bool),
		lastEmit:     make(map[string]time.Time),
		lastLog:      make(map[string]time.Time),
	}
}

func (m *ServiceMonitor) Start(ctx context.Context) {
	// initial check
	for _, s := range m.services {
		running, _ := m.manager.GetServiceStatus(s)
		m.lastStatus[s] = running
	}

	ticker := time.NewTicker(m.intervalSlow)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return

		case <-ticker.C:
			changed := m.checkServices()

			if changed {
				ticker.Reset(m.intervalFast)
			} else {
				ticker.Reset(m.intervalSlow)
			}
		}
	}
}

func (m *ServiceMonitor) checkServices() bool {
	changed := false
	ch := make(chan domain.ServiceStatus, len(m.services))

	for _, svc := range m.services {
		go func(name string) {
			running, msg := m.manager.GetServiceStatus(name)
			ch <- domain.ServiceStatus{
				Name:      name,
				IsRunning: running,
				Message:   msg,
			}
		}(svc)
	}

	for i := 0; i < len(m.services); i++ {
		status := <-ch

		if m.lastStatus[status.Name] != status.IsRunning {
			changed = true
			m.handleStatusChange(status)
			m.lastStatus[status.Name] = status.IsRunning
		}
	}

	return changed
}

func (m *ServiceMonitor) handleStatusChange(s domain.ServiceStatus) {
	now := time.Now()

	// Debounce UI event (max 1 per 400ms / service)
	if now.Sub(m.lastEmit[s.Name]) > 400*time.Millisecond {
		m.emitter.Emit("service:status-changed", s)
		m.lastEmit[s.Name] = now
	}

	// Log throttling (max 1 log per 2 detik)
	if now.Sub(m.lastLog[s.Name]) > 2*time.Second {
		level := "info"
		if s.IsRunning {
			level = "success"
		} else {
			level = "warning"
		}

		m.emitter.Emit("service:log", map[string]interface{}{
			"timestamp": now.Format("15:04:05"),
			"level":     level,
			"message":   fmt.Sprintf("%s status changed: %s", s.Name, s.Message),
		})

		m.lastLog[s.Name] = now
	}
}
