package dns

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHostsManager_SyncAndRead(t *testing.T) {
	tempDir := t.TempDir()
	hostsFile := filepath.Join(tempDir, "hosts")

	initialContent := "127.0.0.1 localhost\n::1 localhost ip6-localhost\n"
	if err := os.WriteFile(hostsFile, []byte(initialContent), 0o644); err != nil {
		t.Fatalf("failed to write initial hosts: %v", err)
	}

	mgr := NewHostsManagerWithPath(hostsFile)

	// Sync 2 domains
	domains := []string{"my-laravel.test", "my-wp.test"}
	if err := mgr.SyncDomains(domains); err != nil {
		t.Fatalf("SyncDomains error: %v", err)
	}

	readDomains, err := mgr.GetManagedDomains()
	if err != nil {
		t.Fatalf("GetManagedDomains error: %v", err)
	}
	if len(readDomains) != 2 {
		t.Errorf("expected 2 domains, got %d", len(readDomains))
	}

	content, _ := os.ReadFile(hostsFile)
	contentStr := string(content)

	if !strings.Contains(contentStr, "127.0.0.1 localhost") {
		t.Errorf("expected original content to be preserved")
	}
	if !strings.Contains(contentStr, "127.0.0.1\tmy-laravel.test") {
		t.Errorf("expected my-laravel.test in hosts")
	}

	// Re-sync with updated domain list
	newDomains := []string{"my-laravel.test", "nextjs-app.test", "api.test"}
	if err := mgr.SyncDomains(newDomains); err != nil {
		t.Fatalf("Re-sync error: %v", err)
	}

	readDomains2, _ := mgr.GetManagedDomains()
	if len(readDomains2) != 3 {
		t.Errorf("expected 3 domains after re-sync, got %d", len(readDomains2))
	}
}
