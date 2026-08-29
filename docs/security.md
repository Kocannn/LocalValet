# LocalValet v2 - Security Architecture

## Threat Model

### Attack Surface

```
┌─────────────────────────────────────────────────────────────┐
│                    Attack Surface                            │
├─────────────────────────────────────────────────────────────┤
│ 1. Wails IPC Interface                                      │
│    - All Go methods exposed to frontend                     │
│    - Input validation critical                              │
│                                                             │
│ 2. Managed Service Processes                                │
│    - Apache, Nginx, MySQL, Redis, PHP-FPM                   │
│    - Running as child processes                             │
│    - Port binding                                           │
│                                                             │
│ 3. File System Access                                       │
│    - Config files (read/write)                              │
│    - Runtime binaries (execute)                             │
│    - Project directories (read)                             │
│    - SSL certificates (sensitive)                           │
│                                                             │
│ 4. Network Exposure                                         │
│    - Localhost services                                     │
│    - Potential LAN exposure                                 │
│    - DNS resolution                                         │
│                                                             │
│ 5. Privilege Escalation                                     │
│    - Port binding < 1024                                    │
│    - System file modification                               │
│    - Service restart                                        │
└─────────────────────────────────────────────────────────────┘
```

### Threat Actors

| Actor | Risk Level | Motivation |
|-------|------------|------------|
| Local User | Medium | Accidental misconfiguration |
| Malicious Process | High | Lateral movement, data theft |
| Network Attacker | Low | Service exploitation |
| Supply Chain | Medium | Compromised dependencies |

## Process Isolation Strategy

### Process Group Management

```go
cmd.SysProcAttr = &syscall.SysProcAttr{
    Setpgid: true,   // Create new process group
    Pgid:    0,      // Set group ID to process ID
}
```

**Benefits:**
- Child processes don't receive parent's signals
- Can kill entire process group on shutdown
- Prevents orphaned processes

### Filesystem Isolation

**Option 1: Restricted Working Directory**
```go
cmd.Dir = filepath.Join(baseDir, "runtime", "linux", serviceName, version)
```

**Option 2: chroot (Advanced)**
```go
cmd.SysProcAttr = &syscall.SysProcAttr{
    Chroot: filepath.Join(baseDir, "runtime", "sandbox"),
}
```

**Option 3: Namespaces (Linux)**
```go
cmd.SysProcAttr = &syscall.SysProcAttr{
    Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWPID,
}
```

### Resource Limits

```go
cmd.SysProcAttr = &syscall.SysProcAttr{
    Setpgid: true,
    // CPU time limit (seconds)
    // Memory limit (bytes)
    // File size limit
}
```

**ulimit settings:**
```bash
# In service startup script
ulimit -n 1024      # Max open files
ulimit -v 524288    # Max virtual memory (512MB)
```

### Seccomp Profiles (Linux)

```go
// Restrict syscalls for service processes
seccomp.LoadDefaultProfile()
```

**Allowed syscalls:**
- read, write, open, close
- mmap, mprotect, brk
- socket, connect, bind, listen, accept
- fork, execve (for worker processes)

**Blocked syscalls:**
- mount, umount
- reboot
- kexec_load
- ptrace

## Privilege Escalation Handling

### When Root/Sudo is Needed

| Operation | Privilege Required | Strategy |
|-----------|-------------------|----------|
| Bind port < 1024 | root | pkexec / setcap |
| Modify /etc/hosts | root | pkexec |
| Trust SSL cert | root | pkexec |
| Start systemd service | root | systemctl |

### pkexec Strategy (Linux)

```go
func runPrivileged(name string, args ...string) error {
    cmd := exec.Command("pkexec", append([]string{name}, args...)...)
    return cmd.Run()
}
```

**Polkit Policy:**
```xml
<policyconfig>
  <action id="org.localvalet.service-start">
    <description>Start a local development service</description>
    <message>Authentication required to start services</message>
    <defaults>
      <allow_any>auth_admin</allow_any>
      <allow_inactive>auth_admin</allow_inactive>
      <allow_active>auth_self</allow_active>
    </defaults>
    <annotate key="org.freedesktop.policykit.exec.path">/usr/bin/systemctl</annotate>
  </action>
</policyconfig>
```

### Capability-Based Approach (Linux)

```bash
# Grant specific capabilities instead of full root
sudo setcap 'cap_net_bind_service=+ep' /usr/sbin/nginx
sudo setcap 'cap_net_bind_service=+ep' /usr/sbin/apache2
```

**Benefits:**
- No full root access needed
- Minimal privilege principle
- Auditable permissions

### UAC Strategy (Windows)

```go
func runAsAdmin(exe string, args ...string) error {
    verb := "runas"
    cmd := exec.Command("cmd", "/C", exe)
    cmd.SysProcAttr = &syscall.SysProcAttr{
        HideWindow:    true,
        CreationFlags: 0x00000010, // CREATE_NEW_CONSOLE
    }
    return cmd.Run()
}
```

## SSL Certificate Security

### CA Private Key Protection

```go
func protectCAKey(keyPath string) error {
    // Set restrictive permissions
    return os.Chmod(keyPath, 0600)
}
```

**File Permissions:**
```
runtime/certs/
├── ca.pem          (0644 - readable)
├── ca-key.pem      (0600 - owner only)
├── myproject.test.pem    (0644)
└── myproject.test-key.pem (0600)
```

### Certificate Storage

```go
type CertStore struct {
    baseDir string
}

func (s *CertStore) StoreCert(domain string, cert, key []byte) error {
    certPath := filepath.Join(s.baseDir, "runtime", "certs", domain+".pem")
    keyPath := filepath.Join(s.baseDir, "runtime", "certs", domain+"-key.pem")
    
    // Write cert (readable)
    os.WriteFile(certPath, cert, 0644)
    
    // Write key (owner only)
    os.WriteFile(keyPath, key, 0600)
}
```

### Trust Store Management

**Linux:**
```bash
# Copy CA to system trust store
sudo cp ca.pem /usr/local/share/ca-certificates/localvalet.crt
sudo update-ca-certificates
```

**macOS:**
```bash
sudo security add-trusted-cert -d -r trustRoot \
    -k /Library/Keychains/System.keychain ca.pem
```

**Windows:**
```powershell
Import-Certificate -FilePath ca.pem -CertStoreLocation Cert:\LocalMachine\Root
```

### Key Rotation Strategy

```go
func (m *SSLManager) RotateCA() error {
    // 1. Generate new CA
    newCA := generateCA()
    
    // 2. Sign all existing certs with new CA
    for _, domain := range m.ListDomains() {
        m.RegenerateCert(domain, newCA)
    }
    
    // 3. Trust new CA
    m.TrustCA(newCA)
    
    // 4. Remove old CA from trust store
    m.UntrustCA(m.oldCA)
    
    // 5. Delete old CA files
    m.DeleteOldCA()
}
```

## File System Security

### Config File Permissions

```go
func writeSecureConfig(path string, data []byte) error {
    // Write with restrictive permissions
    return os.WriteFile(path, data, 0600)
}
```

**Permission Matrix:**
```
config/
├── services.json    (0644 - readable)
├── runtime.json     (0644 - readable)
├── vhosts.json      (0644 - readable)
└── ssl.json         (0600 - sensitive)

runtime/
├── pids/            (0755)
│   └── *.pid        (0644)
├── logs/            (0755)
│   └── *.log        (0644)
└── certs/           (0700)
    ├── ca.pem       (0644)
    └── ca-key.pem   (0600)
```

### Temporary File Cleanup

```go
func cleanupTempFiles() error {
    tempDir := filepath.Join(baseDir, "runtime", "sandbox", "tmp")
    
    // Remove files older than 24 hours
    filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
        if time.Since(info.ModTime()) > 24*time.Hour {
            os.Remove(path)
        }
        return nil
    })
}
```

### Log File Rotation

```go
func rotateLogs() error {
    logDir := filepath.Join(baseDir, "runtime", "logs")
    
    filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
        // Rotate if > 10MB
        if info.Size() > 10*1024*1024 {
            os.Rename(path, path+".1")
        }
        return nil
    })
}
```

## Network Security

### Localhost-Only Binding

```go
// Default: bind to localhost only
cmd.Args = append(cmd.Args, "--bind", "127.0.0.1")

// Nginx: listen on localhost only
// listen 127.0.0.1:8080;

// Apache: bind to localhost
// Listen 127.0.0.1:8080
```

### Port Exposure Warnings

```go
func (m *Manager) StartService(name string) error {
    port := m.GetPort(name)
    
    if port < 1024 {
        m.EmitWarning("Service binding to privileged port", map[string]interface{}{
            "service": name,
            "port":    port,
            "risk":    "May require elevated privileges",
        })
    }
    
    if !isLocalhostOnly(name) {
        m.EmitWarning("Service accessible from network", map[string]interface{}{
            "service": name,
            "port":    port,
            "risk":    "Accessible from LAN",
        })
    }
}
```

### Firewall Integration

```go
func (m *Manager) AddFirewallRule(port int) error {
    // Linux: iptables
    cmd := exec.Command("iptables", "-A", "INPUT", "-p", "tcp",
        "--dport", strconv.Itoa(port), "-s", "127.0.0.1", "-j", "ACCEPT")
    cmd.Run()
    
    cmd = exec.Command("iptables", "-A", "INPUT", "-p", "tcp",
        "--dport", strconv.Itoa(port), "-j", "DROP")
    cmd.Run()
}
```

### DNS Rebinding Protection

```go
func validateDomain(domain string) error {
    // Ensure domain resolves to localhost
    ips, err := net.LookupIP(domain)
    if err != nil {
        return err
    }
    
    for _, ip := range ips {
        if !ip.IsLoopback() {
            return fmt.Errorf("domain %s resolves to non-loopback IP: %s", domain, ip)
        }
    }
    
    return nil
}
```

## IPC Security

### Input Validation

```go
func validateServiceName(name string) error {
    // Only allow alphanumeric and hyphens
    if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(name) {
        return fmt.Errorf("invalid service name: %s", name)
    }
    
    // Check against whitelist
    allowed := []string{"apache", "nginx", "mysql", "postgresql", "redis", "php-fpm"}
    if !contains(allowed, name) {
        return fmt.Errorf("unknown service: %s", name)
    }
    
    return nil
}
```

### Path Traversal Prevention

```go
func validatePath(path string, allowedRoot string) error {
    // Resolve to absolute path
    abs, err := filepath.Abs(path)
    if err != nil {
        return err
    }
    
    // Ensure path is under allowed root
    if !strings.HasPrefix(abs, allowedRoot) {
        return fmt.Errorf("path outside allowed directory: %s", path)
    }
    
    return nil
}
```

### Command Injection Prevention

```go
func validateArgs(args []string) error {
    for _, arg := range args {
        // Check for shell metacharacters
        if strings.ContainsAny(arg, ";|&$`") {
            return fmt.Errorf("invalid argument: %s", arg)
        }
    }
    return nil
}
```

### Rate Limiting

```go
type RateLimiter struct {
    mu       sync.Mutex
    counters map[string]int
    limits   map[string]int
}

func (rl *RateLimiter) Allow(operation string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    rl.counters[operation]++
    
    if rl.counters[operation] > rl.limits[operation] {
        return false
    }
    
    return true
}
```

**Rate Limits:**
| Operation | Limit |
|-----------|-------|
| StartService | 5/minute |
| StopService | 5/minute |
| GenerateCert | 10/hour |
| OpenTerminal | 10/minute |

## Dependency Security

### Binary Verification

```go
func verifyBinary(path string, expectedHash string) error {
    // Calculate SHA256
    hash := sha256.New()
    data, _ := os.ReadFile(path)
    hash.Write(data)
    actualHash := hex.EncodeToString(hash.Sum(nil))
    
    if actualHash != expectedHash {
        return fmt.Errorf("binary hash mismatch: expected %s, got %s", expectedHash, actualHash)
    }
    
    return nil
}
```

### Update Mechanism Security

```go
func (m *UpdateManager) CheckForUpdates() error {
    // 1. Fetch update manifest over HTTPS
    manifest, err := fetchManifest("https://updates.localvalet.com/manifest.json")
    
    // 2. Verify manifest signature
    if !verifySignature(manifest) {
        return fmt.Errorf("invalid manifest signature")
    }
    
    // 3. Download update over HTTPS
    update, err := downloadUpdate(manifest.URL)
    
    // 4. Verify update hash
    if !verifyHash(update, manifest.Hash) {
        return fmt.Errorf("update hash mismatch")
    }
    
    // 5. Apply update
    return applyUpdate(update)
}
```

### Plugin Sandboxing

```go
type PluginSandbox struct {
    allowedPaths []string
    allowedSyscalls []string
}

func (s *PluginSandbox) Validate(plugin Plugin) error {
    // Check plugin doesn't access restricted paths
    // Check plugin doesn't use restricted syscalls
    // Check plugin signature
}
```

## Security Checklist

### Pre-Release

- [ ] All IPC inputs validated
- [ ] Path traversal prevention implemented
- [ ] Command injection prevention implemented
- [ ] Rate limiting on sensitive operations
- [ ] File permissions set correctly
- [ ] SSL private keys protected
- [ ] Localhost-only binding by default
- [ ] Privilege escalation documented
- [ ] Temporary files cleaned up
- [ ] Logs don't contain sensitive data

### Runtime

- [ ] Process groups isolated
- [ ] Resource limits applied
- [ ] Health checks running
- [ ] Anomaly detection active
- [ ] Audit logging enabled

### Maintenance

- [ ] Dependencies updated regularly
- [ ] Security patches applied
- [ ] Certificates rotated
- [ ] Logs rotated
- [ ] Backups verified
