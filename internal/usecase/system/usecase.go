package system

import (
	"os/exec"
	"runtime"
	"strings"
)

type UseCase struct{}

func New() *UseCase {
	return &UseCase{}
}

// CheckSudoAccess checks if the user has sudo access.
func (u *UseCase) CheckSudoAccess() bool {
	if runtime.GOOS != "linux" {
		return true
	}

	cmd := exec.Command("sudo", "-n", "true")
	err := cmd.Run()
	return err == nil
}

// GetSystemInfo returns system information.
func (u *UseCase) GetSystemInfo() map[string]string {
	info := make(map[string]string)
	info["os"] = runtime.GOOS
	info["arch"] = runtime.GOARCH

	if runtime.GOOS == "linux" {
		cmd := exec.Command("lsb_release", "-d")
		output, err := cmd.Output()
		if err == nil {
			info["version"] = strings.TrimSpace(strings.TrimPrefix(string(output), "Description:"))
		}
	}

	return info
}

// GetBinarySourceInfo returns binary source details.
func (u *UseCase) GetBinarySourceInfo(usingSystemBinaries bool, binaryLocation string) map[string]interface{} {
	info := make(map[string]interface{})
	info["os"] = runtime.GOOS
	info["using_system_binaries"] = usingSystemBinaries
	info["binary_location"] = binaryLocation
	info["binary_validation"] = nil
	return info
}
