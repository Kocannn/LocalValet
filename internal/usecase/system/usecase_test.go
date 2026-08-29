package system

import (
	"runtime"
	"testing"
)

func TestUseCase_GetSystemInfo(t *testing.T) {
	uc := New()
	info := uc.GetSystemInfo()

	if info["os"] != runtime.GOOS {
		t.Errorf("expected os %s, got %s", runtime.GOOS, info["os"])
	}
	if info["arch"] != runtime.GOARCH {
		t.Errorf("expected arch %s, got %s", runtime.GOARCH, info["arch"])
	}
}

func TestUseCase_GetBinarySourceInfo(t *testing.T) {
	uc := New()
	info := uc.GetBinarySourceInfo(false, "runtime/")

	if info["using_system_binaries"] != false {
		t.Errorf("expected using_system_binaries false")
	}
	if info["binary_location"] != "runtime/" {
		t.Errorf("expected binary_location runtime/, got %v", info["binary_location"])
	}
}

func TestUseCase_CheckSudoAccess(t *testing.T) {
	uc := New()
	// Should not panic
	_ = uc.CheckSudoAccess()
}
