package terminal

import (
	terminaldomain "LocalValet/internal/domain/terminal"
	"errors"
	"runtime"
	"testing"
)

type mockTerminalManager struct {
	failLaunch bool
	result     terminaldomain.LaunchResult
}

func (m *mockTerminalManager) Launch(options terminaldomain.LaunchOptions) (terminaldomain.LaunchResult, error) {
	if m.failLaunch {
		return terminaldomain.LaunchResult{}, errors.New("launch mock error")
	}
	return m.result, nil
}

func TestTerminalUseCase_Launch(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping linux-only terminal test")
	}

	mgr := &mockTerminalManager{
		result: terminaldomain.LaunchResult{
			Terminal: "xterm",
			WorkDir:  "/home/user/project",
		},
	}
	uc := New(mgr)

	msg := uc.LaunchTerminal("/home/user/project", "xterm")
	if msg.Level != "success" {
		t.Errorf("expected level success, got %s: %s", msg.Level, msg.Message)
	}

	// Failure case
	mgr.failLaunch = true
	failMsg := uc.LaunchTerminal("/home/user/project", "")
	if failMsg.Level != "error" {
		t.Errorf("expected level error on failure, got %s", failMsg.Level)
	}
}
