package linux

import (
	servicedomain "LocalValet/internal/domain/service"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestLinuxManager_LifecycleWithMockProcess(t *testing.T) {
	tempDir := t.TempDir()
	pidDir := filepath.Join(tempDir, "runtime", "pids")
	logsDir := filepath.Join(tempDir, "runtime", "logs")
	_ = os.MkdirAll(pidDir, 0o755)
	_ = os.MkdirAll(logsDir, 0o755)

	mgr := New().(*LinuxManager)

	// 1. Initial status - stopped
	running, msg := mgr.GetServiceStatus("mock_svc")
	if running {
		t.Errorf("expected mock_svc to be stopped initially")
	}

	// 2. Start a mock long-running process manually to test PID tracking
	cmd := exec.Command("sleep", "30")
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start mock process: %v", err)
	}
	pid := cmd.Process.Pid
	defer func() {
		_ = cmd.Process.Kill()
	}()

	pidPath := mgr.registry.pidFilePath("mock_svc")
	_ = os.MkdirAll(filepath.Dir(pidPath), 0o755)
	if err := os.WriteFile(pidPath, []byte(strconv.Itoa(pid)), 0o644); err != nil {
		t.Fatalf("failed to write pid file: %v", err)
	}

	// 3. Status should now be running
	running, msg = mgr.GetServiceStatus("mock_svc")
	if !running {
		t.Errorf("expected mock_svc to be running, msg: %s", msg)
	}

	// 4. Test StopService - graceful shutdown
	err := mgr.StopService("mock_svc")
	if err != nil {
		t.Errorf("StopService failed: %v", err)
	}

	// Wait briefly for process cleanup
	time.Sleep(100 * time.Millisecond)

	running, _ = mgr.GetServiceStatus("mock_svc")
	if running {
		t.Errorf("expected mock_svc to be stopped after StopService")
	}

	// PID file should be removed
	if _, err := os.Stat(pidPath); err == nil {
		t.Errorf("expected PID file to be removed after StopService")
	}
}

func TestLinuxManager_StalePIDCleanup(t *testing.T) {
	mgr := New().(*LinuxManager)
	pidPath := mgr.registry.pidFilePath("stale_svc")
	_ = os.MkdirAll(filepath.Dir(pidPath), 0o755)

	// Write an invalid/dead PID
	_ = os.WriteFile(pidPath, []byte("99999999"), 0o644)

	running, _ := mgr.GetServiceStatus("stale_svc")
	if running {
		t.Errorf("expected stale_svc to be reported as stopped")
	}

	// Verify PID file was cleaned up
	if _, err := os.Stat(pidPath); err == nil {
		t.Errorf("expected stale PID file to be automatically cleaned up")
	}
}

func TestLinuxManager_InjectPortArgument(t *testing.T) {
	tests := []struct {
		service string
		args    []string
		port    int
		check   string
	}{
		{"mysql", []string{}, 3307, "--port=3307"},
		{"redis", []string{}, 6380, "6380"},
		{"postgresql", []string{}, 5433, "5433"},
		{"apache", []string{}, 8081, "PORT=8081"},
	}

	for _, tt := range tests {
		result := injectPortArgument(tt.service, tt.args, tt.port)
		found := false
		for _, arg := range result {
			if arg == tt.check {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("service %s with port %d did not contain %q in args: %v", tt.service, tt.port, tt.check, result)
		}
	}
}

func TestNewWithConfigs(t *testing.T) {
	configs := []servicedomain.Config{
		{ServiceName: "custom", DefaultPort: 9000, Category: "Custom"},
	}
	mgr := NewWithConfigs(configs)
	if mgr == nil {
		t.Fatalf("expected NewWithConfigs to return a non-nil Manager")
	}
}
