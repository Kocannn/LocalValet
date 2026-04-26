package linux

import (
	servicedomain "LocalValet/internal/domain/service"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

type LinuxManager struct {
	registry *RuntimeRegistry
}

func New() servicedomain.Manager {
	return &LinuxManager{registry: NewRuntimeRegistry()}
}

func (l *LinuxManager) GetServiceStatus(serviceName string) (bool, string) {
	pid, ok := l.readPID(serviceName)
	if !ok {
		return false, fmt.Sprintf("%s is stopped", serviceName)
	}

	if processAlive(pid) {
		return true, fmt.Sprintf("%s is running", serviceName)
	}

	_ = os.Remove(l.registry.pidFilePath(serviceName))
	return false, fmt.Sprintf("%s is stopped", serviceName)
}

func (l *LinuxManager) StartService(serviceName string) error {
	runtimeCfg, err := l.registry.Resolve(serviceName)
	if err != nil {
		return err
	}

	if running, _ := l.GetServiceStatus(serviceName); running {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(l.registry.pidFilePath(serviceName)), 0o755); err != nil {
		return err
	}

	if err := os.MkdirAll(l.registry.logsDir(), 0o755); err != nil {
		return err
	}

	logPath := filepath.Join(l.registry.logsDir(), serviceName+".log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer logFile.Close()

	cmd := exec.Command(runtimeCfg.BinaryPath, runtimeCfg.Args...)
	cmd.Env = os.Environ()
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if runtimeCfg.WorkingDir != "" {
		cmd.Dir = runtimeCfg.WorkingDir
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := os.WriteFile(l.registry.pidFilePath(serviceName), []byte(strconv.Itoa(cmd.Process.Pid)), 0o644); err != nil {
		_ = cmd.Process.Kill()
		return err
	}

	_ = cmd.Process.Release()
	return nil
}

func (l *LinuxManager) StopService(serviceName string) error {
	pid, ok := l.readPID(serviceName)
	if !ok {
		return nil
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		_ = os.Remove(l.registry.pidFilePath(serviceName))
		return nil
	}

	_ = proc.Signal(syscall.SIGTERM)
	deadline := time.Now().Add(5 * time.Second)

	for time.Now().Before(deadline) {
		if !processAlive(pid) {
			_ = os.Remove(l.registry.pidFilePath(serviceName))
			return nil
		}
		time.Sleep(150 * time.Millisecond)
	}

	_ = proc.Signal(syscall.SIGKILL)
	_ = os.Remove(l.registry.pidFilePath(serviceName))
	return nil
}

func (l *LinuxManager) SetServiceVersion(serviceName, version string) error {
	return l.registry.SetActiveVersion(serviceName, version)
}

func (l *LinuxManager) GetActiveServiceVersion(serviceName string) (string, error) {
	return l.registry.GetActiveVersion(serviceName)
}

func (l *LinuxManager) GetAvailableVersions(serviceName string) ([]string, error) {
	return l.registry.GetVersions(serviceName)
}

func (l *LinuxManager) readPID(serviceName string) (int, bool) {
	b, err := os.ReadFile(l.registry.pidFilePath(serviceName))
	if err != nil {
		return 0, false
	}

	pid, err := strconv.Atoi(string(b))
	if err != nil {
		return 0, false
	}

	return pid, true
}

func processAlive(pid int) bool {
	if pid <= 0 {
		return false
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	return proc.Signal(syscall.Signal(0)) == nil
}
