package linux

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// PortManager manages port availability checking and conflict resolution.
type PortManager struct {
	mu             sync.RWMutex
	allocatedPorts map[string]int
}

// NewPortManager creates a new PortManager instance.
func NewPortManager() *PortManager {
	return &PortManager{
		allocatedPorts: make(map[string]int),
	}
}

// IsPortAvailable checks if a TCP port is currently free to bind on localhost and 0.0.0.0.
func (pm *PortManager) IsPortAvailable(port int) bool {
	if port <= 0 || port > 65535 {
		return false
	}

	// Try binding 127.0.0.1
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	_ = ln.Close()

	// Try binding 0.0.0.0 (if dual-stack or wildcard is supported)
	lnWildcard, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return false
	}
	_ = lnWildcard.Close()

	return true
}

// FindAvailablePort scans starting from startPort up to maxAttempts to find the next available port.
func (pm *PortManager) FindAvailablePort(startPort int, maxAttempts int) (int, error) {
	if maxAttempts <= 0 {
		maxAttempts = 200
	}

	for i := 0; i < maxAttempts; i++ {
		candidate := startPort + i
		if candidate > 65535 {
			break
		}

		pm.mu.RLock()
		isAllocated := false
		for _, p := range pm.allocatedPorts {
			if p == candidate {
				isAllocated = true
				break
			}
		}
		pm.mu.RUnlock()

		if !isAllocated && pm.IsPortAvailable(candidate) {
			return candidate, nil
		}
	}

	return 0, fmt.Errorf("no available port found in range %d - %d", startPort, startPort+maxAttempts-1)
}

// ResolvePort resolves a port for a given service.
// If defaultPort is free, it allocates and returns (defaultPort, false, nil).
// If a conflict is detected, it auto-remaps to an available port within +200 range and returns (newPort, true, nil).
func (pm *PortManager) ResolvePort(serviceName string, defaultPort int) (int, bool, error) {
	if defaultPort <= 0 {
		return 0, false, nil
	}

	pm.mu.RLock()
	allocated, exists := pm.allocatedPorts[serviceName]
	pm.mu.RUnlock()

	if exists && allocated > 0 {
		return allocated, false, nil
	}

	// Check if defaultPort is available
	if pm.IsPortAvailable(defaultPort) {
		pm.SetAllocatedPort(serviceName, defaultPort)
		return defaultPort, false, nil
	}

	// Port conflict detected! Auto-remap to the next available port
	remappedPort, err := pm.FindAvailablePort(defaultPort+1, 200)
	if err != nil {
		return 0, true, fmt.Errorf("port conflict on %d for %s and could not find alternative: %w", defaultPort, serviceName, err)
	}

	pm.SetAllocatedPort(serviceName, remappedPort)
	return remappedPort, true, nil
}

// SetAllocatedPort records a specific port as allocated to a service.
func (pm *PortManager) SetAllocatedPort(serviceName string, port int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.allocatedPorts[serviceName] = port
}

// GetAllocatedPort returns the currently allocated port for a service, or 0 if none.
func (pm *PortManager) GetAllocatedPort(serviceName string) int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.allocatedPorts[serviceName]
}

// ReleasePort releases the port allocated to a service.
func (pm *PortManager) ReleasePort(serviceName string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.allocatedPorts, serviceName)
}

// IsServicePortListening verifies if the allocated port for a service is actively accepting connections.
func (pm *PortManager) IsServicePortListening(port int, timeout time.Duration) bool {
	if port <= 0 {
		return false
	}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
