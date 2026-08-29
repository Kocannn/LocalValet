package linux

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"
)

// HealthChecker checks the liveness and readiness of service processes and ports.
type HealthChecker struct {
	httpTimeout time.Duration
	tcpTimeout  time.Duration
}

// NewHealthChecker creates a new HealthChecker.
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		httpTimeout: 1 * time.Second,
		tcpTimeout:  800 * time.Millisecond,
	}
}

// CheckProcess checks if the process with the given PID is alive.
func (h *HealthChecker) CheckProcess(pid int) bool {
	if pid <= 0 {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}

// CheckTCP checks if a TCP port is open and accepting connections.
func (h *HealthChecker) CheckTCP(port int, timeout time.Duration) bool {
	if port <= 0 {
		return false
	}
	if timeout <= 0 {
		timeout = h.tcpTimeout
	}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// CheckHTTP checks if an HTTP server on localhost:port responds.
func (h *HealthChecker) CheckHTTP(port int, path string, timeout time.Duration) bool {
	if port <= 0 {
		return false
	}
	if timeout <= 0 {
		timeout = h.httpTimeout
	}
	if path == "" {
		path = "/"
	}

	client := &http.Client{
		Timeout: timeout,
	}
	url := fmt.Sprintf("http://127.0.0.1:%d%s", port, path)
	resp, err := client.Get(url)
	if err != nil {
		// If TCP connection worked but returned 4xx/5xx or TLS error, it's still alive
		return h.CheckTCP(port, timeout)
	}
	defer resp.Body.Close()
	return true
}

// CheckService evaluates the health of a service based on its PID, port, and health check type.
func (h *HealthChecker) CheckService(serviceName string, pid int, port int, checkType string) (bool, string) {
	if !h.CheckProcess(pid) {
		return false, fmt.Sprintf("%s process (PID %d) is not running", serviceName, pid)
	}

	switch checkType {
	case "tcp":
		if port > 0 {
			if h.CheckTCP(port, h.tcpTimeout) {
				return true, fmt.Sprintf("%s is healthy on port %d (PID %d)", serviceName, port, pid)
			}
			return false, fmt.Sprintf("%s process is running (PID %d) but port %d is not responding", serviceName, pid, port)
		}
	case "http":
		if port > 0 {
			if h.CheckHTTP(port, "/", h.httpTimeout) {
				return true, fmt.Sprintf("%s is healthy on port %d (PID %d)", serviceName, port, pid)
			}
			return false, fmt.Sprintf("%s process is running (PID %d) but HTTP port %d is not responding", serviceName, pid, port)
		}
	case "process":
		return true, fmt.Sprintf("%s is running (PID %d)", serviceName, pid)
	}

	// Default fallback: if port is defined, check TCP; otherwise just process
	if port > 0 {
		if h.CheckTCP(port, h.tcpTimeout) {
			return true, fmt.Sprintf("%s is healthy on port %d (PID %d)", serviceName, port, pid)
		}
	}

	return true, fmt.Sprintf("%s is running (PID %d)", serviceName, pid)
}
