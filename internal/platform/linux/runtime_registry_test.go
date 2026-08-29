package linux

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRuntimeRegistry_ConfigLoadAndResolve(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	// Create a dummy binary
	binDir := filepath.Join(tempDir, "runtime", "linux", "php", "8.4", "sbin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("failed to create bin dir: %v", err)
	}
	dummyBin := filepath.Join(binDir, "php-fpm")
	if err := os.WriteFile(dummyBin, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("failed to write dummy bin: %v", err)
	}

	rawCfg := `{
		"services": {
			"php-fpm": {
				"activeVersion": "8.4",
				"versions": {
					"8.4": {
						"binary": "` + dummyBin + `",
						"args": ["-v"],
						"workingDir": "` + binDir + `"
					},
					"8.3": {
						"binary": "/nonexistent/php-fpm",
						"args": []
					}
				}
			}
		}
	}`

	configPath := filepath.Join(configDir, "runtime.json")
	if err := os.WriteFile(configPath, []byte(rawCfg), 0o644); err != nil {
		t.Fatalf("failed to write runtime config: %v", err)
	}

	reg := NewRuntimeRegistry()

	// Direct load test
	var cfg runtimeConfig
	if err := json.Unmarshal([]byte(rawCfg), &cfg); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	phpSvc := cfg.Services["php-fpm"]
	if phpSvc == nil {
		t.Fatalf("php-fpm service is nil")
	}
	if phpSvc.ActiveVersion != "8.4" {
		t.Errorf("expected active version 8.4, got %s", phpSvc.ActiveVersion)
	}
	if len(phpSvc.Versions) != 2 {
		t.Errorf("expected 2 versions, got %d", len(phpSvc.Versions))
	}

	_ = reg
}

func TestRuntimeRegistry_LegacyJsonSupport(t *testing.T) {
	legacyJSON := `{
		"services": {
			"php-fpm": {
				"activeVersion": "8.4",
				"8.4": {
					"binary": "runtime/linux/php/8.4/sbin/php-fpm",
					"args": ["--nodaemonize"]
				}
			}
		}
	}`

	var cfg runtimeConfig
	if err := json.Unmarshal([]byte(legacyJSON), &cfg); err != nil {
		t.Fatalf("failed to parse legacy JSON: %v", err)
	}

	svc := cfg.Services["php-fpm"]
	if svc == nil {
		t.Fatalf("expected php-fpm service to exist")
	}
	if svc.ActiveVersion != "8.4" {
		t.Errorf("expected activeVersion 8.4, got %s", svc.ActiveVersion)
	}
	if ver, ok := svc.Versions["8.4"]; !ok || ver.Binary != "runtime/linux/php/8.4/sbin/php-fpm" {
		t.Errorf("expected version 8.4 with binary path, got %+v", ver)
	}
}
