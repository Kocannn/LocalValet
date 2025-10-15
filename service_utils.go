package main

import (
	"os/exec"
	"runtime"
	"strings"
)

// CheckSudoAccess checks if the user has sudo access
func (a *App) CheckSudoAccess() bool {
	if runtime.GOOS != "linux" {
		return true // Not needed on macOS/Windows
	}

	cmd := exec.Command("sudo", "-n", "true")
	err := cmd.Run()
	return err == nil
}

// GetSystemInfo returns system information
func (a *App) GetSystemInfo() map[string]string {
	info := make(map[string]string)
	info["os"] = runtime.GOOS
	info["arch"] = runtime.GOARCH

	// Get OS version
	if runtime.GOOS == "linux" {
		cmd := exec.Command("lsb_release", "-d")
		output, err := cmd.Output()
		if err == nil {
			info["version"] = strings.TrimSpace(strings.TrimPrefix(string(output), "Description:"))
		}
	}

	return info
}

// GetServiceCommand returns the appropriate service command for the OS
func getServiceCommand(action, serviceName string) *exec.Cmd {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("pkexec", "systemctl", action, serviceName)
	case "darwin":
		return exec.Command("brew", "services", action, serviceName)
	case "windows":
		if action == "start" {
			return exec.Command("net", "start", serviceName)
		} else if action == "stop" {
			return exec.Command("net", "stop", serviceName)
		}
	}
	return nil
}
