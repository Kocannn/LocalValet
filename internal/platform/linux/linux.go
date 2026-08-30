package linux

import (
	servicedomain "LocalValet/internal/domain/service"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// LinuxManager manages the lifecycle of local service processes on Linux.
type LinuxManager struct {
	registry *RuntimeRegistry
	portMgr  *PortManager
	health   *HealthChecker
	configs  []servicedomain.Config
}

// New creates a new LinuxManager instance.
func New() servicedomain.Manager {
	return &LinuxManager{
		registry: NewRuntimeRegistry(),
		portMgr:  NewPortManager(),
		health:   NewHealthChecker(),
		configs:  servicedomain.DefaultConfigs(),
	}
}

// NewWithConfigs creates a LinuxManager with specific service configurations.
func NewWithConfigs(configs []servicedomain.Config) servicedomain.Manager {
	return &LinuxManager{
		registry: NewRuntimeRegistry(),
		portMgr:  NewPortManager(),
		health:   NewHealthChecker(),
		configs:  append([]servicedomain.Config(nil), configs...),
	}
}

// GetServiceStatus checks if a service is currently running.
func (l *LinuxManager) GetServiceStatus(serviceName string) (bool, string) {
	pid, ok := l.readPID(serviceName)
	if !ok {
		l.portMgr.ReleasePort(serviceName)
		return false, fmt.Sprintf("%s is stopped", serviceName)
	}

	if processAlive(pid) {
		port := l.portMgr.GetAllocatedPort(serviceName)
		if port > 0 {
			return true, fmt.Sprintf("%s is running on port %d (PID %d)", serviceName, port, pid)
		}
		return true, fmt.Sprintf("%s is running (PID %d)", serviceName, pid)
	}

	// Clean up stale PID file and port
	_ = os.Remove(l.registry.pidFilePath(serviceName))
	l.portMgr.ReleasePort(serviceName)
	return false, fmt.Sprintf("%s is stopped", serviceName)
}

// CheckHealth evaluates the health and responsiveness of a running service.
func (l *LinuxManager) CheckHealth(serviceName string) (bool, string) {
	pid, ok := l.readPID(serviceName)
	if !ok || !processAlive(pid) {
		_ = os.Remove(l.registry.pidFilePath(serviceName))
		l.portMgr.ReleasePort(serviceName)
		return false, fmt.Sprintf("%s is stopped", serviceName)
	}

	cfg, _ := servicedomain.GetConfig(serviceName, l.configs)
	port := l.portMgr.GetAllocatedPort(serviceName)
	if port == 0 && cfg.DefaultPort > 0 {
		port = cfg.DefaultPort
	}

	return l.health.CheckService(serviceName, pid, port, cfg.HealthCheckType)
}

// GetAllocatedPort returns the currently allocated port for the given service.
func (l *LinuxManager) GetAllocatedPort(serviceName string) int {
	return l.portMgr.GetAllocatedPort(serviceName)
}

// StartService starts a service binary, auto-resolving port conflicts and configuring environment.
func (l *LinuxManager) StartService(serviceName string) error {
	runtimeCfg, err := l.registry.Resolve(serviceName)
	if err != nil {
		return err
	}

	if running, _ := l.GetServiceStatus(serviceName); running {
		return nil
	}

	// Resolve port and check for conflict
	cfg, hasCfg := servicedomain.GetConfig(serviceName, l.configs)
	defaultPort := 0
	if hasCfg {
		defaultPort = cfg.DefaultPort
	}

	allocatedPort := defaultPort
	if defaultPort > 0 {
		var err error
		allocatedPort, _, err = l.portMgr.ResolvePort(serviceName, defaultPort)
		if err != nil {
			return fmt.Errorf("port resolution failed for %s: %w", serviceName, err)
		}
	}

	args := l.prepareServiceArguments(serviceName, runtimeCfg, allocatedPort)

	if err := os.MkdirAll(filepath.Dir(l.registry.pidFilePath(serviceName)), 0o755); err != nil {
		return err
	}

	if err := os.MkdirAll(l.registry.logsDir(), 0o755); err != nil {
		return err
	}

	logPath := filepath.Join(l.registry.logsDir(), serviceName+".log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer logFile.Close()

	cmd := exec.Command(runtimeCfg.BinaryPath, args...)
	cmd.Env = buildServiceEnv(runtimeCfg.BinaryPath, runtimeCfg.WorkingDir)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if runtimeCfg.WorkingDir != "" {
		cmd.Dir = runtimeCfg.WorkingDir
	}

	if err := cmd.Start(); err != nil {
		l.portMgr.ReleasePort(serviceName)
		return err
	}

	pid := cmd.Process.Pid
	if err := os.WriteFile(l.registry.pidFilePath(serviceName), []byte(strconv.Itoa(pid)), 0o644); err != nil {
		_ = killProcessGroup(pid)
		l.portMgr.ReleasePort(serviceName)
		return err
	}

	_ = cmd.Process.Release()

	// Brief check to ensure process did not crash immediately
	time.Sleep(100 * time.Millisecond)
	if !processAlive(pid) {
		_ = os.Remove(l.registry.pidFilePath(serviceName))
		l.portMgr.ReleasePort(serviceName)
		return fmt.Errorf("%s failed to stay running (crashed on startup, see %s)", serviceName, logPath)
	}

	return nil
}

// StopService stops a running service process with graceful SIGTERM and SIGKILL timeout fallback.
func (l *LinuxManager) StopService(serviceName string) error {
	pid, ok := l.readPID(serviceName)
	if !ok {
		l.portMgr.ReleasePort(serviceName)
		return nil
	}

	if !processAlive(pid) {
		_ = os.Remove(l.registry.pidFilePath(serviceName))
		l.portMgr.ReleasePort(serviceName)
		return nil
	}

	// 1. Send SIGTERM to process group
	_ = signalProcessGroup(pid, syscall.SIGTERM)

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if !processAlive(pid) {
			_ = os.Remove(l.registry.pidFilePath(serviceName))
			l.portMgr.ReleasePort(serviceName)
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	// 2. Force kill with SIGKILL if still alive
	_ = signalProcessGroup(pid, syscall.SIGKILL)
	_ = os.Remove(l.registry.pidFilePath(serviceName))
	l.portMgr.ReleasePort(serviceName)
	return nil
}

// SetServiceVersion sets the active runtime version for a service.
func (l *LinuxManager) SetServiceVersion(serviceName, version string) error {
	return l.registry.SetActiveVersion(serviceName, version)
}

// GetActiveServiceVersion retrieves the active runtime version for a service.
func (l *LinuxManager) GetActiveServiceVersion(serviceName string) (string, error) {
	return l.registry.GetActiveVersion(serviceName)
}

// GetAvailableVersions retrieves all configured and available versions for a service.
func (l *LinuxManager) GetAvailableVersions(serviceName string) ([]string, error) {
	return l.registry.GetVersions(serviceName)
}

func (l *LinuxManager) readPID(serviceName string) (int, bool) {
	b, err := os.ReadFile(l.registry.pidFilePath(serviceName))
	if err != nil {
		return 0, false
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(b)))
	if err != nil {
		return 0, false
	}

	return pid, true
}

func processAlive(pid int) bool {
	if pid <= 0 {
		return false
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	return proc.Signal(syscall.Signal(0)) == nil
}

func signalProcessGroup(pid int, sig syscall.Signal) error {
	pgid, err := syscall.Getpgid(pid)
	if err == nil {
		if err := syscall.Kill(-pgid, sig); err == nil {
			return nil
		}
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Signal(sig)
}

func killProcessGroup(pid int) error {
	return signalProcessGroup(pid, syscall.SIGKILL)
}

func buildServiceEnv(binaryPath, workingDir string) []string {
	env := os.Environ()

	// Locate candidate lib directories
	var libDirs []string
	if binaryPath != "" {
		binDir := filepath.Dir(binaryPath)
		libDirs = append(libDirs,
			filepath.Join(binDir, "..", "lib"),
			filepath.Join(binDir, "lib"),
		)
	}
	if workingDir != "" {
		libDirs = append(libDirs,
			filepath.Join(workingDir, "lib"),
		)
	}

	var validLibs []string
	for _, libDir := range libDirs {
		clean := filepath.Clean(libDir)
		if info, err := os.Stat(clean); err == nil && info.IsDir() {
			validLibs = append(validLibs, clean)
		}
	}

	if len(validLibs) > 0 {
		currentLd := os.Getenv("LD_LIBRARY_PATH")
		newLd := strings.Join(validLibs, ":")
		if currentLd != "" {
			newLd = newLd + ":" + currentLd
		}
		env = appendOrReplaceEnv(env, "LD_LIBRARY_PATH", newLd)
	}

	return env
}

func appendOrReplaceEnv(env []string, key, value string) []string {
	prefix := key + "="
	for i, item := range env {
		if strings.HasPrefix(item, prefix) {
			env[i] = prefix + value
			return env
		}
	}
	return append(env, prefix+value)
}

func (l *LinuxManager) prepareServiceArguments(serviceName string, runtimeCfg *resolvedRuntime, port int) []string {
	baseDir := l.registry.baseDir()
	args := append([]string(nil), runtimeCfg.Args...)

	switch serviceName {
	case "mysql":
		dataDir := filepath.Join(baseDir, "runtime", "data", "mysql")
		socketPath := filepath.Join(dataDir, "mysql.sock")
		pidPath := l.registry.pidFilePath("mysql")
		binDir := filepath.Dir(runtimeCfg.BinaryPath)
		baseInstallDir := filepath.Dir(binDir)
		pluginDir := filepath.Join(baseInstallDir, "lib", "plugin")

		_ = os.MkdirAll(dataDir, 0o755)
		_ = initMySQLDataDir(baseDir, runtimeCfg.BinaryPath)
		_ = writeMySQLConfig(baseDir, port)
		_ = writePHPMysqlConfig(baseDir)

		mysqlArgs := make([]string, 0, len(args)+12)
		mysqlArgs = append(mysqlArgs, args...)
		mysqlArgs = append(mysqlArgs,
			fmt.Sprintf("--basedir=%s", baseInstallDir),
			fmt.Sprintf("--datadir=%s", dataDir),
			fmt.Sprintf("--plugin-dir=%s", pluginDir),
			fmt.Sprintf("--socket=%s", socketPath),
			fmt.Sprintf("--pid-file=%s", pidPath),
			fmt.Sprintf("--port=%d", port),
			"--bind-address=127.0.0.1",
		)
		return mysqlArgs

	case "redis":
		dataDir := filepath.Join(baseDir, "runtime", "data", "redis")
		_ = os.MkdirAll(dataDir, 0o755)
		if port > 0 {
			args = append(args, "--port", fmt.Sprintf("%d", port))
		}
		args = append(args,
			"--bind", "127.0.0.1",
			"--dir", dataDir,
			"--pidfile", l.registry.pidFilePath("redis"),
		)
		return args

	case "postgresql":
		dataDir := filepath.Join(baseDir, "runtime", "data", "postgresql")
		_ = os.MkdirAll(dataDir, 0o700)
		if port > 0 {
			args = append(args, "-p", fmt.Sprintf("%d", port))
		}
		args = append(args,
			"-D", dataDir,
			"-h", "127.0.0.1",
			"-k", dataDir,
		)
		return args

	case "apache":
		if port > 0 {
			args = append(args, "-D", fmt.Sprintf("PORT=%d", port))
		}
		return args

	default:
		if port > 0 {
			args = injectPortArgument(serviceName, args, port)
		}
		return args
	}
}

func initMySQLDataDir(baseDir, mysqlBinary string) error {
	dataDir := filepath.Join(baseDir, "runtime", "data", "mysql")
	if fileExists(filepath.Join(dataDir, "mysql")) {
		return nil
	}

	_ = os.MkdirAll(dataDir, 0o755)

	binDir := filepath.Dir(mysqlBinary)
	baseInstallDir := filepath.Dir(binDir)
	libDir := filepath.Join(baseInstallDir, "lib")

	installScripts := []string{
		filepath.Join(baseInstallDir, "scripts", "mariadb-install-db"),
		filepath.Join(baseInstallDir, "scripts", "mysql_install_db"),
		filepath.Join(baseInstallDir, "bin", "mariadb-install-db"),
		filepath.Join(baseInstallDir, "bin", "mysql_install_db"),
	}

	for _, script := range installScripts {
		if fileExists(script) {
			cmd := exec.Command(script,
				fmt.Sprintf("--basedir=%s", baseInstallDir),
				fmt.Sprintf("--datadir=%s", dataDir),
				"--auth-root-authentication-method=normal",
				"--skip-test-db",
			)
			cmd.Dir = baseInstallDir
			cmd.Env = append(os.Environ(), fmt.Sprintf("LD_LIBRARY_PATH=%s", libDir))
			if err := cmd.Run(); err == nil {
				return nil
			}
		}
	}

	return nil
}

func writeMySQLConfig(baseDir string, port int) error {
	dataDir := filepath.Join(baseDir, "runtime", "data", "mysql")
	_ = os.MkdirAll(dataDir, 0o755)

	socketPath := filepath.Join(dataDir, "mysql.sock")
	pidPath := filepath.Join(baseDir, "runtime", "pids", "mysql.pid")
	myCnfPath := filepath.Join(dataDir, "my.cnf")

	content := fmt.Sprintf(`[mysqld]
datadir = %s
socket = %s
pid-file = %s
port = %d
bind-address = 127.0.0.1

[client]
socket = %s
port = %d
user = root

[mysql]
socket = %s
port = %d
user = root

[mysqldump]
socket = %s
port = %d
user = root

[mariadb]
socket = %s
port = %d
user = root

[mysqladmin]
socket = %s
port = %d
user = root
`, dataDir, socketPath, pidPath, port, socketPath, port, socketPath, port, socketPath, port, socketPath, port, socketPath, port)

	if err := os.WriteFile(myCnfPath, []byte(content), 0o644); err != nil {
		return err
	}

	sandboxHome := filepath.Join(baseDir, "runtime", "sandbox", "home")
	_ = os.MkdirAll(sandboxHome, 0o755)
	_ = os.WriteFile(filepath.Join(sandboxHome, ".my.cnf"), []byte(content), 0o600)

	return nil
}

func writePHPMysqlConfig(baseDir string) error {
	socketPath := filepath.Join(baseDir, "runtime", "data", "mysql", "mysql.sock")
	confDir := filepath.Join(baseDir, "runtime", "sandbox", "conf.d")
	_ = os.MkdirAll(confDir, 0o755)

	content := fmt.Sprintf(`; LocalValet isolated MySQL configuration
pdo_mysql.default_socket = %s
mysqli.default_socket = %s
mysqlnd.default_socket = %s
`, socketPath, socketPath, socketPath)

	iniPath := filepath.Join(confDir, "localvalet_mysql.ini")
	return os.WriteFile(iniPath, []byte(content), 0o644)
}

func injectPortArgument(serviceName string, args []string, port int) []string {
	switch serviceName {
	case "mysql":
		return append(args, fmt.Sprintf("--port=%d", port))
	case "redis":
		return append(args, "--port", fmt.Sprintf("%d", port))
	case "postgresql":
		return append(args, "-p", fmt.Sprintf("%d", port))
	case "apache":
		return append(args, "-D", fmt.Sprintf("PORT=%d", port))
	default:
		return args
	}
}
