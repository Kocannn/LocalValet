package terminal

import (
	terminaldomain "LocalValet/internal/domain/terminal"
	"fmt"
	"runtime"
	"time"
)

type LogMessage struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

type UseCase struct {
	manager terminaldomain.Manager
}

func New(manager terminaldomain.Manager) *UseCase {
	return &UseCase{manager: manager}
}

func (u *UseCase) OpenContextTerminal(projectDir string, preferredTerminal string) LogMessage {
	now := time.Now().Format("15:04:05")
	if runtime.GOOS != "linux" {
		return LogMessage{
			Timestamp: now,
			Level:     "error",
			Message:   "Terminal launch is currently supported on Linux only",
		}
	}

	result, err := u.manager.Launch(terminaldomain.LaunchOptions{
		ProjectDir:        projectDir,
		PreferredTerminal: preferredTerminal,
	})
	if err != nil {
		return LogMessage{
			Timestamp: now,
			Level:     "error",
			Message:   fmt.Sprintf("Failed to open terminal: %v", err),
		}
	}

	return LogMessage{
		Timestamp: now,
		Level:     "success",
		Message:   fmt.Sprintf("Terminal opened with %s in %s", result.Terminal, result.WorkDir),
	}
}
