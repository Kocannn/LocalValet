# LocalValet v2 - Process Manager Architecture

## Service Lifecycle Manager

### Process States

```
                    ┌──────────┐
                    │  Stopped │
                    └────┬─────┘
                         │ Start()
                         ▼
                    ┌──────────┐
              ┌─────│ Starting │─────┐
              │     └──────────┘     │
              │ Success              │ Error
              ▼                      ▼
        ┌──────────┐          ┌──────────┐
        │ Running  │          │  Error   │
        └────┬─────┘          └──────────┘
             │ Stop()
             ▼
        ┌──────────┐
        │ Stopping │
        └────┬─────┘
             │
             ▼
        ┌──────────┐
        │ Stopped  │
        └──────────┘
```

### Start Sequence

```go
func (m *Manager) StartService(name string) error {
    // 1. Resolve runtime config
    cfg, err := m.registry.Resolve(name)
    
    // 2. Check port availability
    selection, err := m.resolvePort(name, cfg)
    
    // 3. Apply port overrides if needed
    if selection.Remapped {
        applyPortOverride(name, cfg, selection.Port)
    }
    
    // 4. Create log file
    logFile := openLogFile(name)
    
    // 5. Start process with process group
    cmd := exec.Command(cfg.BinaryPath, cfg.Args...)
    cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
    cmd.Stdout = logFile
    cmd.Stderr = logFile
    
    // 6. Write PID file
    writePIDFile(name, cmd.Process.Pid)
    
    // 7. Write port file
    writePortFile(name, selection.Port)
    
    // 8. Health check with timeout
    return waitForHealth(name, 5*time.Second)
}
```

### Stop Sequence

```go
func (m *Manager) StopService(name string) error {
    pid := readPIDFile(name)
    
    // 1. Send SIGTERM
    proc, _ := os.FindProcess(pid)
    proc.Signal(syscall.SIGTERM)
    
    // 2. Wait for graceful shutdown (5s)
    deadline := time.Now().Add(5 * time.Second)
    for time.Now().Before(deadline) {
        if !processAlive(pid) {
            cleanup(name)
            return nil
        }
        time.Sleep(150 * time.Millisecond)
    }
    
    // 3. Force kill with SIGKILL
    proc.Signal(syscall.SIGKILL)
    cleanup(name)
    return nil
}
```

### Health Check Strategies

| Service   | Health Check Method                          |
|-----------|----------------------------------------------|
| Apache    | HTTP GET localhost:port                      |
| Nginx     | HTTP GET localhost:port                      |
| MySQL     | TCP connect to port                          |
| PostgreSQL| TCP connect to port                          |
| Redis     | PING command                                 |
| PHP-FPM   | TCP connect to FastCGI port                  |
| Node.js   | HTTP GET localhost:port (if web server)      |
| Python    | HTTP GET localhost:port (if web server)      |

## Port Management System

### Default Port Registry

```go
var defaultPorts = map[string]int{
    "apache":     8080,
    "nginx":      8080,
    "mysql":      3306,
    "postgresql": 5432,
    "redis":      6379,
    "php-fpm":    9074,
    "node":       3000,
    "python":     5000,
}
```

### Conflict Detection & Resolution

```go
func selectAvailablePort(default int) (int, bool, error) {
    // 1. Try default port
    if isTCPPortAvailable(default) {
        return default, false, nil
    }
    
    // 2. Scan +200 range
    for port := default + 1; port <= default + 200; port++ {
        if isTCPPortAvailable(port) {
            return port, true, nil  // remapped = true
        }
    }
    
    return 0, false, fmt.Errorf("no available port")
}
```

### Port Override Application

| Service    | Override Method                              |
|------------|----------------------------------------------|
| MySQL      | `--port=XXXX` CLI flag                       |
| PostgreSQL | `-p XXXX` CLI flag                           |
| Redis      | `--port XXXX` CLI flag                       |
| Apache     | `-C "Listen XXXX"` CLI flag                  |
| PHP-FPM    | Patch `listen = 127.0.0.1:XXXX` in config   |
| Nginx      | Patch `listen XXXX;` in config               |

## Auto Virtual Host Manager

### Project Discovery

```go
func (m *VHostManager) DiscoverProjects(roots []string) []Project {
    var projects []Project
    
    for _, root := range roots {
        filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
            // Detect framework by marker files
            if fileExists(path, "artisan") {
                projects = append(projects, Project{
                    Path:      path,
                    Framework: "laravel",
                    Domain:    filepath.Base(path) + ".test",
                })
            }
            if fileExists(path, "composer.json") {
                // PHP project
            }
            if fileExists(path, "package.json") {
                // Node.js project
            }
            return nil
        })
    }
    
    return projects
}
```

### VHost Config Generation

**Nginx Template:**
```nginx
server {
    listen 80;
    server_name {{.Domain}};
    root {{.Root}}/public;
    
    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }
    
    location ~ \.php$ {
        fastcgi_pass 127.0.0.1:{{.PHPFPMPort}};
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }
}
```

**Apache Template:**
```apache
<VirtualHost *:80>
    ServerName {{.Domain}}
    DocumentRoot {{.Root}}/public
    
    <Directory "{{.Root}}/public">
        AllowOverride All
        Require all granted
    </Directory>
</VirtualHost>
```

### DNS Resolution

**Option 1: /etc/hosts (simple)**
```
127.0.0.1 myproject.test
127.0.0.1 another-project.test
```

**Option 2: dnsmasq (wildcard)**
```
address=/.test/127.0.0.1
```

## SSL Certificate Manager

### Local CA Generation

```go
func (m *SSLManager) InitCA() error {
    // Generate CA key pair
    caKey, _ := rsa.GenerateKey(rand.Reader, 2048)
    caCert := &x509.Certificate{
        SerialNumber: big.NewInt(1),
        Subject: pkix.Name{
            Organization: []string{"LocalValet Local CA"},
        },
        NotBefore:             time.Now(),
        NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
        KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
        BasicConstraintsValid: true,
        IsCA:                  true,
    }
    
    // Self-sign
    certDER, _ := x509.CreateCertificate(rand.Reader, caCert, caCert, &caKey.PublicKey, caKey)
    
    // Save to runtime/certs/ca.pem and ca-key.pem
    // Trust via update-ca-certificates (Linux)
}
```

### Per-Domain Certificate

```go
func (m *SSLManager) GenerateCert(domain string) error {
    // Load CA
    caCert, caKey := loadCA()
    
    // Generate domain key pair
    domainKey, _ := rsa.GenerateKey(rand.Reader, 2048)
    domainCert := &x509.Certificate{
        SerialNumber: big.NewInt(2),
        Subject: pkix.Name{
            CommonName: domain,
        },
        DNSNames:  []string{domain, "*." + domain},
        NotBefore: time.Now(),
        NotAfter:  time.Now().Add(365 * 24 * time.Hour),
    }
    
    // Sign with CA
    certDER, _ := x509.CreateCertificate(rand.Reader, domainCert, caCert, &domainKey.PublicKey, caKey)
    
    // Save to runtime/certs/{domain}.pem and {domain}-key.pem
}
```

## Docker Integration

### Docker Compose Detection

```go
func (m *DockerManager) DetectComposeFiles(roots []string) []ComposeProject {
    var projects []ComposeProject
    
    for _, root := range roots {
        if fileExists(root, "docker-compose.yml") || fileExists(root, "compose.yml") {
            projects = append(projects, ComposeProject{
                Path: root,
                File: findComposeFile(root),
            })
        }
    }
    
    return projects
}
```

### Container Status Monitoring

```go
func (m *DockerManager) GetContainerStatus() []Container {
    cmd := exec.Command("docker", "ps", "-a", "--format", "{{.Names}}\t{{.Status}}\t{{.Ports}}")
    output, _ := cmd.Output()
    return parseContainerList(output)
}
```

### Port Coordination

```go
func (m *DockerManager) GetDockerPorts() map[int]string {
    // Parse docker ps output to get mapped ports
    // Avoid conflicts with LocalValet services
}
```

## Environment Variable Setup

### PATH Injection

```go
func buildPATH(baseDir string, activeServices map[string]string) string {
    paths := []string{}
    
    // Add active service binaries
    for service, version := range activeServices {
        binaryDir := filepath.Dir(resolveBinary(baseDir, service, version))
        paths = append(paths, binaryDir)
    }
    
    // Add system paths
    paths = append(paths, "/usr/local/bin", "/usr/bin", "/bin")
    
    return strings.Join(paths, ":")
}
```

### Shell Init Scripts

**Bash (~/.localvaletrc):**
```bash
export PATH="/opt/localvalet/runtime/linux/php/8.4/bin:$PATH"
export PATH="/opt/localvalet/runtime/linux/node/20/bin:$PATH"
export MYSQL_HOME="/opt/localvalet/runtime/linux/mysql/8.0"
alias artisan='php artisan'
alias sail='./vendor/bin/sail'
```

**Zsh:** Same content, loaded via ZDOTDIR override.

**Fish:**
```fish
set -gx PATH /opt/localvalet/runtime/linux/php/8.4/bin $PATH
```

## Runtime Version Management

### Version Detection

```go
func (m *Registry) GetVersions(service string) []string {
    base := filepath.Join(m.baseDir, "runtime", "linux", service)
    entries, _ := os.ReadDir(base)
    
    var versions []string
    for _, entry := range entries {
        if entry.IsDir() {
            versions = append(versions, entry.Name())
        }
    }
    
    sort.Strings(versions)
    return versions
}
```

### Hot-Switching

```go
func (m *Registry) SetActiveVersion(service, version string) error {
    // 1. Stop service if running
    if m.isRunning(service) {
        m.StopService(service)
    }
    
    // 2. Update config
    cfg := m.loadConfig()
    cfg.Services[service].ActiveVersion = version
    m.saveConfig(cfg)
    
    // 3. Restart service
    return m.StartService(service)
}
```

### Binary Verification

```go
func verifyBinary(path string) error {
    info, err := os.Stat(path)
    if err != nil {
        return fmt.Errorf("binary not found: %s", path)
    }
    
    // Check executable permission
    if info.Mode()&0111 == 0 {
        return fmt.Errorf("binary not executable: %s", path)
    }
    
    return nil
}
```
