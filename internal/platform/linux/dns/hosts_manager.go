package dns

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
)


const (
	BlockStart = "# BEGIN LocalValet Managed Domains"
	BlockEnd   = "# END LocalValet Managed Domains"
)

type HostsManager struct {
	mu        sync.RWMutex
	hostsPath string
}

func NewHostsManager() *HostsManager {
	return &HostsManager{hostsPath: "/etc/hosts"}
}

func NewHostsManagerWithPath(hostsPath string) *HostsManager {
	return &HostsManager{hostsPath: hostsPath}
}

func (h *HostsManager) SyncDomains(domains []string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Normalize & deduplicate domains
	cleanDomains := make(map[string]bool)
	for _, d := range domains {
		trimmed := strings.TrimSpace(strings.ToLower(d))
		if trimmed != "" {
			cleanDomains[trimmed] = true
		}
	}

	sortedDomains := make([]string, 0, len(cleanDomains))
	for d := range cleanDomains {
		sortedDomains = append(sortedDomains, d)
	}
	sort.Strings(sortedDomains)

	// Build new block
	var blockBuf bytes.Buffer
	blockBuf.WriteString(BlockStart + "\n")
	for _, d := range sortedDomains {
		blockBuf.WriteString(fmt.Sprintf("127.0.0.1\t%s\n", d))
	}
	blockBuf.WriteString(BlockEnd + "\n")

	content, err := os.ReadFile(h.hostsPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var newContent string
	existing := string(content)

	startIdx := strings.Index(existing, BlockStart)
	endIdx := strings.Index(existing, BlockEnd)

	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		// Replace existing block
		before := existing[:startIdx]
		after := existing[endIdx+len(BlockEnd):]
		if strings.HasPrefix(after, "\n") {
			after = after[1:]
		}
		newContent = before + blockBuf.String() + after
	} else {
		// Append block
		if len(existing) > 0 && !strings.HasSuffix(existing, "\n") {
			existing += "\n"
		}
		newContent = existing + "\n" + blockBuf.String()
	}

	// Try direct write
	if err := os.WriteFile(h.hostsPath, []byte(newContent), 0o644); err == nil {
		return nil
	}

	// If direct write fails (due to root permission), write to a temp file and copy via pkexec/sudo
	tempFile, err := os.CreateTemp("", "localvalet_hosts_*")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write([]byte(newContent)); err != nil {
		tempFile.Close()
		return err
	}
	tempFile.Close()

	cmdStr := fmt.Sprintf("cp %s %s && chmod 644 %s", tempFile.Name(), h.hostsPath, h.hostsPath)

	if _, err := exec.LookPath("pkexec"); err == nil {
		cmd := exec.Command("pkexec", "sh", "-c", cmdStr)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	if _, err := exec.LookPath("sudo"); err == nil {
		cmd := exec.Command("sudo", "sh", "-c", cmdStr)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("root privileges required to update %s", h.hostsPath)
}

func (h *HostsManager) GetManagedDomains() ([]string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	content, err := os.ReadFile(h.hostsPath)
	if err != nil {
		return []string{}, err
	}

	existing := string(content)
	startIdx := strings.Index(existing, BlockStart)
	endIdx := strings.Index(existing, BlockEnd)

	if startIdx == -1 || endIdx == -1 || endIdx <= startIdx {
		return []string{}, nil
	}

	block := existing[startIdx+len(BlockStart) : endIdx]
	lines := strings.Split(block, "\n")

	var domains []string
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "127.0.0.1" {
			domains = append(domains, fields[1])
		}
	}

	return domains, nil
}
