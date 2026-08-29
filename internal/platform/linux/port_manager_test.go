package linux

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestPortManager_IsPortAvailable(t *testing.T) {
	pm := NewPortManager()

	// Invalid ports
	if pm.IsPortAvailable(0) {
		t.Errorf("expected port 0 to be invalid")
	}
	if pm.IsPortAvailable(-1) {
		t.Errorf("expected negative port to be invalid")
	}
	if pm.IsPortAvailable(70000) {
		t.Errorf("expected port > 65535 to be invalid")
	}

	// Find an ephemeral port to bind
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen on ephemeral port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port

	// Since ln is listening, IsPortAvailable must return false
	if pm.IsPortAvailable(port) {
		t.Errorf("expected busy port %d to not be available", port)
	}

	// Close listener
	_ = ln.Close()
	time.Sleep(10 * time.Millisecond)

	// Now port should be available
	if !pm.IsPortAvailable(port) {
		t.Errorf("expected closed port %d to be available", port)
	}
}

func TestPortManager_ResolvePort(t *testing.T) {
	pm := NewPortManager()

	// 1. Resolve for default port <= 0
	p0, remapped0, err := pm.ResolvePort("dummy", 0)
	if err != nil || p0 != 0 || remapped0 {
		t.Errorf("ResolvePort(0) = (%d, %v, %v), expected (0, false, nil)", p0, remapped0, err)
	}

	// 2. Resolve on an available port
	// Find free port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	freePort := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()
	time.Sleep(10 * time.Millisecond)

	resolvedPort, remapped, err := pm.ResolvePort("test_svc", freePort)
	if err != nil {
		t.Fatalf("ResolvePort failed: %v", err)
	}
	if remapped {
		t.Errorf("expected remapped to be false for free port %d", freePort)
	}
	if resolvedPort != freePort {
		t.Errorf("expected resolvedPort %d, got %d", freePort, resolvedPort)
	}
	if pm.GetAllocatedPort("test_svc") != freePort {
		t.Errorf("expected allocated port to be %d, got %d", freePort, pm.GetAllocatedPort("test_svc"))
	}

	// 3. Resolve with conflict
	// Intentionally occupy a port
	busyLn, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", freePort))
	if err != nil {
		t.Fatalf("failed to occupy port %d: %v", freePort, err)
	}
	defer busyLn.Close()

	pm2 := NewPortManager()
	remappedPort, remapped2, err := pm2.ResolvePort("conflict_svc", freePort)
	if err != nil {
		t.Fatalf("ResolvePort failed during conflict: %v", err)
	}
	if !remapped2 {
		t.Errorf("expected remapped to be true for busy port %d", freePort)
	}
	if remappedPort == freePort {
		t.Errorf("expected remapped port to differ from busy port %d", freePort)
	}

	// Release port
	pm2.ReleasePort("conflict_svc")
	if pm2.GetAllocatedPort("conflict_svc") != 0 {
		t.Errorf("expected allocated port to be 0 after release")
	}
}

func TestPortManager_IsServicePortListening(t *testing.T) {
	pm := NewPortManager()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	if !pm.IsServicePortListening(port, 500*time.Millisecond) {
		t.Errorf("expected port %d to be listening", port)
	}

	if pm.IsServicePortListening(port+1000, 100*time.Millisecond) {
		t.Errorf("expected unassigned port %d to not be listening", port+1000)
	}
}
