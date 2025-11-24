package platform

import (
	"LocalValet/internal/platform/domain"
	"LocalValet/internal/platform/linux"
	"runtime"
)

func NewServiceManager() domain.ServiceManager {
	switch runtime.GOOS {
	case "linux":
		return linux.New()
	default:
		panic("unsupported OS")
	}
}
