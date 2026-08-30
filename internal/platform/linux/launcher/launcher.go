package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// OpenInEditor attempts to open the given directory in the specified or default editor.
func OpenInEditor(projectPath, preferredEditor string) error {
	if projectPath == "" {
		return fmt.Errorf("project path cannot be empty")
	}

	editors := []string{}
	if preferredEditor != "" {
		editors = append(editors, preferredEditor)
	}

	envEditor := os.Getenv("EDITOR")
	if envEditor != "" {
		editors = append(editors, envEditor)
	}

	// Common GUI code editors
	editors = append(editors, "code", "cursor", "phpstorm", "subl", "atom", "gedit", "kate")

	for _, editor := range editors {
		bin, err := exec.LookPath(editor)
		if err == nil && bin != "" {
			cmd := exec.Command(bin, projectPath)
			if err := cmd.Start(); err == nil {
				return nil
			}
		}
	}

	return fmt.Errorf("no suitable code editor found (tried: %s)", strings.Join(editors, ", "))
}

// OpenInBrowser opens a URL in the user's default web browser.
func OpenInBrowser(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	cmd := exec.Command("xdg-open", url)
	return cmd.Start()
}
