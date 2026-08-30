package vhost

import (
	vhostdomain "LocalValet/internal/domain/vhost"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

type NginxGenerator struct {
	mu        sync.RWMutex
	vhostsDir string
}

func NewNginxGenerator() vhostdomain.Generator {
	baseDir := findBaseDir()
	vhostsDir := filepath.Join(baseDir, "runtime", "linux", "nginx", "vhosts")
	_ = os.MkdirAll(vhostsDir, 0o755)
	return &NginxGenerator{vhostsDir: vhostsDir}
}

func NewNginxGeneratorWithPath(vhostsDir string) vhostdomain.Generator {
	_ = os.MkdirAll(vhostsDir, 0o755)
	return &NginxGenerator{vhostsDir: vhostsDir}
}

const nginxTemplate = `# LocalValet Auto-Generated Virtual Host for {{ .Domain }}
# Generated at {{ .ProjectName }}

{{ if .SSLEnabled }}
server {
    listen {{ if .HTTPPort }}{{ .HTTPPort }}{{ else }}80{{ end }};
    server_name {{ .Domain }} *.{{ .Domain }};
    return 301 https://$host$request_uri;
}

server {
    listen {{ if .HTTPSPort }}{{ .HTTPSPort }}{{ else }}443{{ end }} ssl;
    server_name {{ .Domain }} *.{{ .Domain }};

    ssl_certificate {{ .SSLCertPath }};
    ssl_certificate_key {{ .SSLKeyPath }};
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    {{ if .ProxyPass }}
    location / {
        proxy_pass {{ .ProxyPass }};
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    {{ else }}
    root {{ .DocumentRoot }};
    index index.php index.html index.htm;
    charset utf-8;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location = /favicon.ico { access_log off; log_not_found off; }
    location = /robots.txt  { access_log off; log_not_found off; }

    error_page 404 /index.php;

    location ~ \.php$ {
        fastcgi_pass {{ if .PHPFpmAddress }}{{ .PHPFpmAddress }}{{ else }}127.0.0.1:9074{{ end }};
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $realpath_root$fastcgi_script_name;
        fastcgi_param DOCUMENT_ROOT $realpath_root;
        include fastcgi_params;
    }

    location ~ /\.(?!well-known).* {
        deny all;
    }
    {{ end }}
}
{{ else }}
server {
    listen {{ if .HTTPPort }}{{ .HTTPPort }}{{ else }}80{{ end }};
    server_name {{ .Domain }} *.{{ .Domain }};

    {{ if .ProxyPass }}
    location / {
        proxy_pass {{ .ProxyPass }};
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
    {{ else }}
    root {{ .DocumentRoot }};
    index index.php index.html index.htm;
    charset utf-8;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \.php$ {
        fastcgi_pass {{ if .PHPFpmAddress }}{{ .PHPFpmAddress }}{{ else }}127.0.0.1:9074{{ end }};
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $realpath_root$fastcgi_script_name;
        include fastcgi_params;
    }
    {{ end }}
}
{{ end }}
`

func (g *NginxGenerator) Generate(config vhostdomain.VHostConfig) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if config.Domain == "" {
		return "", fmt.Errorf("domain cannot be empty")
	}

	tmpl, err := template.New("nginx_vhost").Parse(nginxTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	vhostPath := g.GetVHostPath(config.Domain)
	if err := os.WriteFile(vhostPath, buf.Bytes(), 0o644); err != nil {
		return "", fmt.Errorf("failed to write vhost config at %s: %w", vhostPath, err)
	}

	return vhostPath, nil
}

func (g *NginxGenerator) Remove(domain string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	vhostPath := g.GetVHostPath(domain)
	if fileExists(vhostPath) {
		return os.Remove(vhostPath)
	}
	return nil
}

func (g *NginxGenerator) List() ([]string, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	entries, err := os.ReadDir(g.vhostsDir)
	if err != nil {
		return []string{}, nil
	}

	var vhosts []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".conf") {
			domain := strings.TrimSuffix(entry.Name(), ".conf")
			vhosts = append(vhosts, domain)
		}
	}
	return vhosts, nil
}

func (g *NginxGenerator) GetVHostPath(domain string) string {
	return filepath.Join(g.vhostsDir, domain+".conf")
}

func (g *NginxGenerator) ReloadNginx() error {
	cmd := exec.Command("pkill", "-HUP", "-f", "nginx")
	_ = cmd.Run()
	return nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func findBaseDir() string {
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		if fileExists(filepath.Join(exeDir, "config", "runtime.json")) {
			return exeDir
		}
	}

	cwd, err := os.Getwd()
	if err == nil {
		if fileExists(filepath.Join(cwd, "config", "runtime.json")) {
			return cwd
		}
		parent := filepath.Dir(cwd)
		if fileExists(filepath.Join(parent, "config", "runtime.json")) {
			return parent
		}
		return cwd
	}

	return "."
}
