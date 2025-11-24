package linux

import (
	"LocalValet/internal/platform/domain"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type LinuxManager struct{}

func New() domain.ServiceManager {
	return &LinuxManager{}
}

func (l *LinuxManager) GetServiceStatus(serviceName string) (bool, string) {
	cmd := exec.Command("systemctl", "is-active", serviceName)
	cmd.Env = os.Environ()
	output, err := cmd.Output()

	if err == nil && strings.TrimSpace(string(output)) == "active" {
		return true, fmt.Sprintf("%s is running", serviceName)
	}
	return false, fmt.Sprintf("%s is stopped", serviceName)
}

func (l *LinuxManager) StartService(serviceName string) error {
	cmd := exec.Command("pkexec", "systemctl", "start", serviceName)
	cmd.Env = os.Environ()
	return cmd.Run()
}

func (l *LinuxManager) StopService(serviceName string) error {
	cmd := exec.Command("pkexec", "systemctl", "stop", serviceName)
	cmd.Env = os.Environ()
	return cmd.Run()
}
