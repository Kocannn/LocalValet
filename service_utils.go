package main

import (
	"os/exec"
	"runtime"
)

// CheckSudoAccess checks if the user has sudo access
func (a *App) CheckSudoAccess() bool {
	return a.controller.CheckSudoAccess()
}

// GetSystemInfo returns system information
func (a *App) GetSystemInfo() map[string]string {
	return a.controller.GetSystemInfo()
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
