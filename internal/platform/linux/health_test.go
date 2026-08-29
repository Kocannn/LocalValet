package linux

import (
	"net"
	"net/http"
	"net/http/httptest"

	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestHealthChecker_CheckProcess(t *testing.T) {
	hc := NewHealthChecker()

	// Current process is alive
	if !hc.CheckProcess(os.Getpid()) {
		t.Errorf("expected current PID %d to be alive", os.Getpid())
	}

	// Non-existent PID
	if hc.CheckProcess(99999999) {
		t.Errorf("expected invalid PID to not be alive")
	}

	// PID <= 0
	if hc.CheckProcess(0) || hc.CheckProcess(-1) {
		t.Errorf("expected PID <= 0 to not be alive")
	}
}

func TestHealthChecker_CheckTCP(t *testing.T) {
	hc := NewHealthChecker()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	if !hc.CheckTCP(port, 500*time.Millisecond) {
		t.Errorf("expected TCP port %d to be reachable", port)
	}

	if hc.CheckTCP(port+1000, 100*time.Millisecond) {
		t.Errorf("expected inactive port %d to not be reachable", port+1000)
	}
}

func TestHealthChecker_CheckHTTP(t *testing.T) {
	hc := NewHealthChecker()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))
	defer server.Close()

	parts := strings.Split(server.URL, ":")
	port, _ := strconv.Atoi(parts[len(parts)-1])

	if !hc.CheckHTTP(port, "/", 1*time.Second) {
		t.Errorf("expected HTTP server on port %d to be healthy", port)
	}
}

func TestHealthChecker_CheckService(t *testing.T) {
	hc := NewHealthChecker()
	myPid := os.Getpid()

	// Dead process
	healthy, msg := hc.CheckService("test_dead", 99999999, 8080, "tcp")
	if healthy {
		t.Errorf("expected dead process to be unhealthy, msg: %s", msg)
	}

	// Alive process with "process" check
	healthy, _ = hc.CheckService("test_proc", myPid, 0, "process")
	if !healthy {
		t.Errorf("expected process check to be healthy for current PID")
	}

	// Alive process with TCP check
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	healthy, _ = hc.CheckService("test_tcp", myPid, port, "tcp")
	if !healthy {
		t.Errorf("expected TCP check to be healthy on open port %d", port)
	}

	// Alive process with unresponsive TCP port
	healthy, msg = hc.CheckService("test_tcp_fail", myPid, port+2000, "tcp")
	if healthy {
		t.Errorf("expected TCP check to fail on closed port %d, got msg: %s", port+2000, msg)
	}
}
