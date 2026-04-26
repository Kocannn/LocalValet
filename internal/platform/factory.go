package platform

import (
	servicedomain "LocalValet/internal/domain/service"
	"LocalValet/internal/platform/linux"
	"runtime"
)

func NewServiceManager() servicedomain.Manager {
	switch runtime.GOOS {
	case "linux":
		return linux.New()
	default:
		panic("unsupported OS")
	}
}
