package vhost

import (
	vhostdomain "LocalValet/internal/domain/vhost"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNginxGenerator_PHPAndSSL(t *testing.T) {
	tempDir := t.TempDir()
	vhostsDir := filepath.Join(tempDir, "vhosts")
	gen := NewNginxGeneratorWithPath(vhostsDir)

	cfg := vhostdomain.VHostConfig{
		Domain:        "my-laravel.test",
		ProjectName:   "My Laravel",
		DocumentRoot:  "/var/www/my-laravel/public",
		PHPFpmAddress: "127.0.0.1:9074",
		SSLEnabled:    true,
		SSLCertPath:   "/path/to/my-laravel.test.crt",
		SSLKeyPath:    "/path/to/my-laravel.test.key",
	}

	confPath, err := gen.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	content, err := os.ReadFile(confPath)
	if err != nil {
		t.Fatalf("failed to read generated conf: %v", err)
	}
	confStr := string(content)

	if !strings.Contains(confStr, "server_name my-laravel.test") {
		t.Errorf("expected server_name directive in conf")
	}
	if !strings.Contains(confStr, "ssl_certificate /path/to/my-laravel.test.crt") {
		t.Errorf("expected ssl_certificate directive in conf")
	}
	if !strings.Contains(confStr, "fastcgi_pass 127.0.0.1:9074") {
		t.Errorf("expected fastcgi_pass in conf")
	}
	if !strings.Contains(confStr, "try_files $uri $uri/ /index.php?$query_string") {
		t.Errorf("expected try_files in conf")
	}

	// List vhosts
	list, err := gen.List()
	if err != nil || len(list) != 1 || list[0] != "my-laravel.test" {
		t.Errorf("expected list to contain my-laravel.test, got %v", list)
	}

	// Remove vhost
	if err := gen.Remove("my-laravel.test"); err != nil {
		t.Errorf("Remove error: %v", err)
	}
	list, _ = gen.List()
	if len(list) != 0 {
		t.Errorf("expected 0 vhosts after remove, got %d", len(list))
	}
}

func TestNginxGenerator_ProxyPass(t *testing.T) {
	tempDir := t.TempDir()
	vhostsDir := filepath.Join(tempDir, "vhosts")
	gen := NewNginxGeneratorWithPath(vhostsDir)

	cfg := vhostdomain.VHostConfig{
		Domain:       "nextjs-app.test",
		ProjectName:  "Next.js App",
		DocumentRoot: "/var/www/nextjs-app",
		ProxyPass:    "http://127.0.0.1:3000",
		SSLEnabled:   false,
	}

	confPath, err := gen.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate proxy error: %v", err)
	}

	content, err := os.ReadFile(confPath)
	if err != nil {
		t.Fatalf("failed to read proxy conf: %v", err)
	}
	confStr := string(content)

	if !strings.Contains(confStr, "proxy_pass http://127.0.0.1:3000") {
		t.Errorf("expected proxy_pass in conf")
	}
	if !strings.Contains(confStr, "proxy_set_header Upgrade $http_upgrade") {
		t.Errorf("expected WebSocket upgrade headers in conf")
	}
}
