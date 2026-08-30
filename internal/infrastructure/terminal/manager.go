package terminal

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	terminaldomain "LocalValet/internal/domain/terminal"
)

const terminalClass = "LocalValetTerm"

type Manager struct{}

type runtimeVersion struct {
	Binary     string   `json:"binary"`
	Args       []string `json:"args"`
	WorkingDir string   `json:"workingDir"`
}

type runtimeService struct {
	ActiveVersion string                     `json:"activeVersion"`
	Versions      map[string]*runtimeVersion `json:"versions"`
}

type runtimeConfig struct {
	Services map[string]*runtimeService `json:"services"`
}

func NewLinuxManager() terminaldomain.Manager {
	return &Manager{}
}

func (s *runtimeService) UnmarshalJSON(data []byte) error {
	type serviceAlias struct {
		ActiveVersion string                     `json:"activeVersion"`
		Versions      map[string]*runtimeVersion `json:"versions"`
	}

	var alias serviceAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	if alias.Versions == nil {
		alias.Versions = make(map[string]*runtimeVersion)
	}

	// Backward compatibility: allow flattened version entries under the service object.
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for key, value := range raw {
		if key == "activeVersion" || key == "versions" {
			continue
		}

		if _, exists := alias.Versions[key]; exists {
			continue
		}

		var version runtimeVersion
		if err := json.Unmarshal(value, &version); err != nil {
			continue
		}

		if version.Binary == "" {
			continue
		}

		v := version
		alias.Versions[key] = &v
	}

	s.ActiveVersion = alias.ActiveVersion
	s.Versions = alias.Versions
	return nil
}

func (m *Manager) Launch(options terminaldomain.LaunchOptions) (terminaldomain.LaunchResult, error) {
	baseDir := findBaseDir()
	cfg, err := loadRuntimeConfig(filepath.Join(baseDir, "config", "runtime.json"))
	if err != nil {
		return terminaldomain.LaunchResult{}, err
	}

	workDir := resolveWorkDir(options.ProjectDir)
	paths, envMap := buildInjectedEnv(baseDir, cfg)
	preferred := strings.TrimSpace(options.PreferredTerminal)
	versions := activeVersions(cfg)
	appEnv, cleanPath := buildAppEnvironment(baseDir, workDir, paths, envMap)

	switch runtime.GOOS {
	case "linux":
		return launchLinux(preferred, workDir, baseDir, cleanPath, appEnv, versions)
	case "windows":
		return launchWindows(preferred, workDir, baseDir, cleanPath, appEnv, versions)
	default:
		return terminaldomain.LaunchResult{}, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func launchLinux(preferred, workDir, baseDir, cleanPath string, env []string, versions map[string]string) (terminaldomain.LaunchResult, error) {
	shellPath := detectDefaultShell()
	shellName := strings.ToLower(filepath.Base(shellPath))

	var shellArgs []string
	switch shellName {
	case "bash":
		initScript, err := writeBashInitScript(baseDir, cleanPath, versions)
		if err != nil {
			return terminaldomain.LaunchResult{}, err
		}
		// Use --noprofile so host /etc/profile and ~/.bash_profile are not loaded,
		// preventing host version managers (nvm, asdf, etc.) from polluting PATH.
		shellArgs = []string{"--noprofile", "--rcfile", initScript, "-i"}
	case "zsh":
		startupDir, err := writeZshInitDir(baseDir, cleanPath, versions)
		if err != nil {
			return terminaldomain.LaunchResult{}, err
		}
		env = appendOrReplaceEnv(env, "ZDOTDIR", startupDir)
		// Run interactive shell without loading login global profiles that could overwrite PATH
		shellArgs = []string{"-i"}
	default:
		initScript, err := writeFallbackInitScript(baseDir, shellPath, cleanPath, versions)
		if err != nil {
			return terminaldomain.LaunchResult{}, err
		}
		bootstrap := ". " + shellQuote(initScript) + "; exec " + shellQuote(shellPath) + " -i"
		shellPath = "/bin/sh"
		shellArgs = []string{"-c", bootstrap}
	}

	termBin, termArgs, err := buildLinuxTerminalCommand(preferred, workDir, shellPath, shellArgs)
	if err != nil {
		return terminaldomain.LaunchResult{}, err
	}

	cmd := exec.Command(termBin, termArgs...)
	cmd.Dir = workDir
	cmd.Env = env

	if err := cmd.Start(); err != nil {
		return terminaldomain.LaunchResult{}, err
	}

	_ = cmd.Process.Release()
	return terminaldomain.LaunchResult{Terminal: termBin, WorkDir: workDir}, nil
}

func launchWindows(preferred, workDir, baseDir, cleanPath string, env []string, versions map[string]string) (terminaldomain.LaunchResult, error) {
	sandboxHome := filepath.Join(baseDir, "runtime", "sandbox", "home")
	banner := buildBanner(versions, sandboxHome)
	psCmd := "$Host.UI.RawUI.WindowTitle = 'LocalValet [Sandbox]'; " +
		"$env:PATH = '" + escapeSingleQuotes(cleanPath) + "'; " +
		"Set-Location -LiteralPath '" + escapeSingleQuotes(workDir) + "'; " +
		"Write-Host '" + escapeSingleQuotes(banner) + "' -ForegroundColor Green"

	candidates := []string{}
	if preferred != "" {
		candidates = append(candidates, preferred)
	}
	candidates = append(candidates, "wt", "powershell", "cmd")

	for _, term := range candidates {
		bin, err := exec.LookPath(term)
		if err != nil {
			continue
		}

		var args []string
		switch term {
		case "wt":
			args = []string{
				"-w", "0", "new-tab",
				"--title", "LocalValet [Sandbox]",
				"-d", workDir,
				"powershell", "-NoExit", "-Command", psCmd,
			}
		case "powershell":
			args = []string{"-NoExit", "-Command", psCmd}
		case "cmd":
			args = []string{"/K", "cd /d " + workDir}
		default:
			args = []string{}
		}

		cmd := exec.Command(bin, args...)
		cmd.Dir = workDir
		cmd.Env = env
		if err := cmd.Start(); err != nil {
			continue
		}

		_ = cmd.Process.Release()
		return terminaldomain.LaunchResult{Terminal: bin, WorkDir: workDir}, nil
	}

	return terminaldomain.LaunchResult{}, fmt.Errorf("no supported terminal emulator found")
}

func buildLinuxTerminalCommand(preferred, workDir, shellBin string, shellArgs []string) (string, []string, error) {
	candidates := []string{}
	if preferred != "" {
		candidates = append(candidates, preferred)
	}
	candidates = append(candidates, "kitty", "alacritty", "foot", "gnome-terminal", "konsole", "xterm")

	for _, term := range candidates {
		bin, err := exec.LookPath(term)
		if err != nil {
			continue
		}

		switch term {
		case "kitty":
			args := []string{
				"--class", terminalClass,
				"--directory", workDir,
				shellBin,
			}
			args = append(args, shellArgs...)
			return bin, args, nil
		case "alacritty":
			args := []string{
				"--class", terminalClass,
				"--working-directory", workDir,
				"-e", shellBin,
			}
			args = append(args, shellArgs...)
			return bin, args, nil
		case "foot":
			args := []string{
				"--app-id", terminalClass,
				"--working-directory", workDir,
				shellBin,
			}
			args = append(args, shellArgs...)
			return bin, args, nil
		case "gnome-terminal":
			args := []string{
				"--class", terminalClass,
				"--working-directory", workDir,
				"--", shellBin,
			}
			args = append(args, shellArgs...)
			return bin, args, nil
		case "konsole":
			args := []string{
				"--name", terminalClass,
				"--workdir", workDir,
				"-e", shellBin,
			}
			args = append(args, shellArgs...)
			return bin, args, nil
		case "xterm":
			args := []string{"-class", terminalClass, "-e", shellBin}
			args = append(args, shellArgs...)
			return bin, args, nil
		}
	}

	return "", nil, fmt.Errorf("no supported terminal emulator found (tried: %v)", candidates)
}

// writeMySQLConfig creates an isolated my.cnf for LocalValet MySQL server and clients.
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

	// Also write to sandbox home
	sandboxHome := filepath.Join(baseDir, "runtime", "sandbox", "home")
	_ = os.MkdirAll(sandboxHome, 0o755)
	_ = os.WriteFile(filepath.Join(sandboxHome, ".my.cnf"), []byte(content), 0o600)

	return nil
}

// writePHPMysqlConfig creates an isolated PHP configuration snippet pointing PDO/MySQLi to LocalValet socket.
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

// ensureSandboxDirectories ensures the full isolated sandbox directory tree is created.
func ensureSandboxDirectories(baseDir string) (string, string) {
	sandboxDir := filepath.Join(baseDir, "runtime", "sandbox")
	sandboxHome := filepath.Join(sandboxDir, "home")
	mysqlDataDir := filepath.Join(baseDir, "runtime", "data", "mysql")

	dirs := []string{
		sandboxHome,
		filepath.Join(sandboxHome, ".config"),
		filepath.Join(sandboxHome, ".local", "share"),
		filepath.Join(sandboxHome, ".local", "state"),
		filepath.Join(sandboxHome, ".cache"),
		filepath.Join(sandboxDir, "composer"),
		filepath.Join(sandboxDir, "composer", "vendor", "bin"),
		filepath.Join(sandboxDir, "composer", "cache"),
		filepath.Join(sandboxDir, "npm"),
		filepath.Join(sandboxDir, "npm", "global", "bin"),
		filepath.Join(sandboxDir, "npm", "cache"),
		filepath.Join(sandboxDir, "conf.d"),
		filepath.Join(sandboxDir, "tmp"),
		mysqlDataDir,
	}

	for _, dir := range dirs {
		_ = os.MkdirAll(dir, 0o755)
	}

	_ = writeMySQLConfig(baseDir, 3306)
	_ = writePHPMysqlConfig(baseDir)

	return sandboxDir, sandboxHome
}

func buildInjectedEnv(baseDir string, cfg *runtimeConfig) ([]string, map[string]string) {
	runtimeRoot := filepath.Join(baseDir, "runtime")
	paths := make([]string, 0, 16)

	addPath := func(path string) {
		clean := filepath.Clean(path)
		if clean == "" {
			return
		}
		if !isUnderPath(clean, runtimeRoot) {
			return
		}
		if _, err := os.Stat(clean); err != nil {
			return
		}
		for _, existing := range paths {
			if existing == clean {
				return
			}
		}
		paths = append(paths, clean)
	}

	addActiveBinPaths := func(serviceName string) *runtimeVersion {
		service := cfg.Services[serviceName]
		if service == nil {
			return nil
		}
		version := service.Versions[service.ActiveVersion]
		if version == nil {
			return nil
		}
		binaryAbs := resolvePath(baseDir, version.Binary)
		addPath(filepath.Dir(binaryAbs))
		if serviceName == "php-fpm" {
			addPath(filepath.Join(filepath.Dir(filepath.Dir(binaryAbs)), "bin"))
		}
		return version
	}

	phpVersion := addActiveBinPaths("php-fpm")
	_ = addActiveBinPaths("mysql")
	_ = addActiveBinPaths("nginx")
	_ = addActiveBinPaths("node")
	_ = addActiveBinPaths("apache")
	_ = addActiveBinPaths("postgresql")
	_ = addActiveBinPaths("redis")

	composerRoot := filepath.Join(runtimeRoot, "linux", "composer")
	addPath(composerRoot)
	addPath(filepath.Join(composerRoot, "bin"))

	// Sandbox global binary paths for composer & npm global tools
	sandboxComposerBin := filepath.Join(runtimeRoot, "sandbox", "composer", "vendor", "bin")
	sandboxNpmBin := filepath.Join(runtimeRoot, "sandbox", "npm", "global", "bin")
	addPath(sandboxComposerBin)
	addPath(sandboxNpmBin)

	envMap := map[string]string{}

	// Configure isolated PHP ini scanner
	phpConfDir := filepath.Join(runtimeRoot, "sandbox", "conf.d")
	_ = writePHPMysqlConfig(baseDir)
	if phpVersion != nil {
		phpRoot := filepath.Dir(filepath.Dir(resolvePath(baseDir, phpVersion.Binary)))
		iniScanDir := filepath.Join(phpRoot, "etc", "conf.d")
		if dirExists(iniScanDir) {
			envMap["PHP_INI_SCAN_DIR"] = fmt.Sprintf("%s:%s", iniScanDir, phpConfDir)
		} else {
			envMap["PHP_INI_SCAN_DIR"] = phpConfDir
		}
	} else {
		envMap["PHP_INI_SCAN_DIR"] = phpConfDir
	}

	// Configure isolated MySQL socket and client settings
	mysqlSock := filepath.Join(runtimeRoot, "data", "mysql", "mysql.sock")
	mysqlDataDir := filepath.Join(runtimeRoot, "data", "mysql")
	_ = writeMySQLConfig(baseDir, 3306)

	envMap["MYSQL_UNIX_PORT"] = mysqlSock
	envMap["MYSQL_HOME"] = mysqlDataDir
	envMap["MYSQL_HOST"] = "127.0.0.1"
	envMap["MYSQL_TCP_PORT"] = "3306"
	envMap["DB_CONNECTION"] = "mysql"
	envMap["DB_HOST"] = "127.0.0.1"
	envMap["DB_PORT"] = "3306"
	envMap["DB_SOCKET"] = mysqlSock
	envMap["DB_DATABASE"] = "localvalet"
	envMap["DB_USERNAME"] = "root"

	// Add runtime lib directories for dynamic linker (e.g. libcrypt.so.1 for MariaDB)
	var libDirs []string
	mysqlSvc := cfg.Services["mysql"]
	if mysqlSvc != nil {
		mysqlVersion := mysqlSvc.Versions[mysqlSvc.ActiveVersion]
		if mysqlVersion != nil {
			libDir := filepath.Join(resolvePath(baseDir, mysqlVersion.WorkingDir), "lib")
			if dirExists(libDir) {
				libDirs = append(libDirs, libDir)
			}
		}
	}
	if len(libDirs) > 0 {
		envMap["LD_LIBRARY_PATH"] = strings.Join(libDirs, ":")
	}

	sslCandidates := []string{
		filepath.Join(runtimeRoot, "certs", "cacert.pem"),
		filepath.Join(runtimeRoot, "certs", "ca-bundle.crt"),
		"/etc/ssl/certs/ca-certificates.crt",
	}
	for _, candidate := range sslCandidates {
		if fileExists(candidate) {
			envMap["SSL_CERT_FILE"] = candidate
			break
		}
	}

	return paths, envMap
}

func buildAppEnvironment(baseDir, workDir string, paths []string, envMap map[string]string) ([]string, string) {
	currentUser, err := user.Current()
	username := "localvalet"
	if err == nil && currentUser.Username != "" {
		username = currentUser.Username
	}

	hostHome, err := os.UserHomeDir()
	if err != nil || hostHome == "" {
		hostHome = workDir
	}

	sandboxDir, sandboxHome := ensureSandboxDirectories(baseDir)
	sandboxTmp := filepath.Join(sandboxDir, "tmp")
	sandboxComposer := filepath.Join(sandboxDir, "composer")
	sandboxNpm := filepath.Join(sandboxDir, "npm")
	mysqlSock := filepath.Join(baseDir, "runtime", "data", "mysql", "mysql.sock")
	mysqlDataDir := filepath.Join(baseDir, "runtime", "data", "mysql")

	allowedSessionKeys := []string{
		"DISPLAY",
		"WAYLAND_DISPLAY",
		"XDG_RUNTIME_DIR",
		"XAUTHORITY",
		"DBUS_SESSION_BUS_ADDRESS",
		"XDG_SESSION_TYPE",
		"SSH_AUTH_SOCK",
		"COLORTERM",
		"TERM_PROGRAM",
	}

	pathParts := make([]string, 0, len(paths)+8)
	pathParts = append(pathParts, paths...)
	pathParts = append(pathParts,
		"/usr/local/bin",
		"/usr/bin",
		"/bin",
		"/usr/local/sbin",
		"/usr/sbin",
		"/sbin",
	)

	cleanPath := strings.Join(pathParts, string(os.PathListSeparator))

	env := map[string]string{
		"HOME":                   sandboxHome,
		"LOCALVALET_HOME":        sandboxHome,
		"LOCALVALET_HOST_HOME":   hostHome,
		"LOCALVALET_SANDBOX":     "1",
		"LOCALVALET_ENV":         "isolated",
		"LOCALVALET_BASE_DIR":    baseDir,
		"LOCALVALET_SANDBOX_DIR": sandboxDir,
		"LOGNAME":                username,
		"USER":                   username,
		"PATH":                   cleanPath,
		"PWD":                    workDir,
		"SHELL":                  detectDefaultShell(),
		"TERM":                   "xterm-256color",
		"LANG":                   "en_US.UTF-8",
		"LC_ALL":                 "en_US.UTF-8",
		"COMPOSER_HOME":          sandboxComposer,
		"COMPOSER_CACHE_DIR":     filepath.Join(sandboxComposer, "cache"),
		"NPM_CONFIG_PREFIX":      filepath.Join(sandboxNpm, "global"),
		"NPM_CONFIG_CACHE":       filepath.Join(sandboxNpm, "cache"),
		"NPM_CONFIG_USERCONFIG":  filepath.Join(sandboxHome, ".npmrc"),
		"XDG_CONFIG_HOME":        filepath.Join(sandboxHome, ".config"),
		"XDG_DATA_HOME":          filepath.Join(sandboxHome, ".local", "share"),
		"XDG_STATE_HOME":         filepath.Join(sandboxHome, ".local", "state"),
		"XDG_CACHE_HOME":         filepath.Join(sandboxHome, ".cache"),
		"TMPDIR":                 sandboxTmp,
		"TEMP":                   sandboxTmp,
		"TMP":                    sandboxTmp,
		"HISTFILE":               filepath.Join(sandboxHome, ".bash_history"),
		"MYSQL_UNIX_PORT":        mysqlSock,
		"MYSQL_HOME":             mysqlDataDir,
		"MYSQL_HOST":             "127.0.0.1",
		"MYSQL_TCP_PORT":         "3306",
		"DB_CONNECTION":          "mysql",
		"DB_HOST":                "127.0.0.1",
		"DB_PORT":                "3306",
		"DB_SOCKET":              mysqlSock,
		"DB_DATABASE":            "localvalet",
		"DB_USERNAME":            "root",
		"PGHOST":                 "127.0.0.1",
		"PGPORT":                 "5432",
		"REDIS_HOST":             "127.0.0.1",
		"REDIS_PORT":             "6379",
	}

	// Share host ~/.gitconfig globally if available so git commit author info works smoothly
	hostGitConfig := filepath.Join(hostHome, ".gitconfig")
	if fileExists(hostGitConfig) {
		env["GIT_CONFIG_GLOBAL"] = hostGitConfig
	}

	for _, key := range allowedSessionKeys {
		if value := os.Getenv(key); value != "" {
			env[key] = value
		}
	}

	for key, value := range envMap {
		env[key] = value
	}

	result := make([]string, 0, len(env))
	keys := make([]string, 0, len(env))
	for key := range env {
		keys = append(keys, key)
	}
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[j] < keys[i] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	for _, key := range keys {
		result = append(result, key+"="+env[key])
	}

	return result, cleanPath
}

func resolveServiceBinaries(baseDir string, versions map[string]string) (mysqlBin, mariadbBin, mysqldumpBin, mysqladminBin string) {
	mysqlVersion := versions["mysql"]
	if mysqlVersion != "" {
		binDir := filepath.Join(baseDir, "runtime", "linux", "mysql", mysqlVersion, "bin")
		mBin := filepath.Join(binDir, "mysql")
		if fileExists(mBin) {
			mysqlBin = mBin
		}
		mariaBin := filepath.Join(binDir, "mariadb")
		if fileExists(mariaBin) {
			mariadbBin = mariaBin
		}
		dumpBin := filepath.Join(binDir, "mysqldump")
		if fileExists(dumpBin) {
			mysqldumpBin = dumpBin
		}
		adminBin := filepath.Join(binDir, "mysqladmin")
		if fileExists(adminBin) {
			mysqladminBin = adminBin
		}
	}

	if mysqlBin == "" {
		mysqlBin = "mysql"
	}
	if mariadbBin == "" {
		mariadbBin = "mariadb"
	}
	if mysqldumpBin == "" {
		mysqldumpBin = "mysqldump"
	}
	if mysqladminBin == "" {
		mysqladminBin = "mysqladmin"
	}
	return
}

func writeBashInitScript(baseDir, cleanPath string, versions map[string]string) (string, error) {
	sandboxHome := filepath.Join(baseDir, "runtime", "sandbox", "home")
	mysqlSock := filepath.Join(baseDir, "runtime", "data", "mysql", "mysql.sock")
	banner := buildBanner(versions, sandboxHome)
	logPath := filepath.Join(baseDir, "runtime", "logs", "php-fpm.log")
	mysqlBin, mariadbBin, mysqldumpBin, mysqladminBin := resolveServiceBinaries(baseDir, versions)

	f, err := os.CreateTemp("", "lv_init_*.sh")
	if err != nil {
		return "", err
	}
	defer f.Close()

	content := strings.Join([]string{
		"# LocalValet isolated terminal init",
		"export LOCALVALET_ACTIVE=1",
		"export LOCALVALET_SANDBOX=1",
		"export PATH=" + shellQuote(cleanPath),
		"export MYSQL_UNIX_PORT=" + shellQuote(mysqlSock),
		"export MYSQL_TCP_PORT=3306",
		"export MYSQL_HOST=127.0.0.1",
		"export DB_CONNECTION=mysql",
		"export DB_HOST=127.0.0.1",
		"export DB_PORT=3306",
		"export DB_SOCKET=" + shellQuote(mysqlSock),
		"export DB_USERNAME=root",
		"export DB_DATABASE=localvalet",
		"export PS1='\\[\\033[01;32m\\][LocalValet:sandbox]\\[\\033[00m\\] \\[\\033[01;34m\\]\\w\\[\\033[00m\\]\\$ '",
		"alias mysql='" + mysqlBin + " --socket=" + mysqlSock + " -u root'",
		"alias mariadb='" + mariadbBin + " --socket=" + mysqlSock + " -u root'",
		"alias mysqldump='" + mysqldumpBin + " --socket=" + mysqlSock + " -u root'",
		"alias mysqladmin='" + mysqladminBin + " --socket=" + mysqlSock + " -u root'",
		"alias db-shell='" + mysqlBin + " --socket=" + mysqlSock + " -u root'",
		"alias artisan='php artisan'",
		"alias php-logs='tail -f " + shellQuote(logPath) + "'",
		"alias valet-restart='echo \"Use LocalValet UI to restart services\"'",
		"alias sudo='sudo env \"PATH=$PATH\"'",
		"alias ll='ls -la --color=auto'",
		"alias l='ls -CF --color=auto'",
		"alias la='ls -A --color=auto'",
		"if [ -f ./.localvaletrc ]; then . ./.localvaletrc || true; fi",
		"printf '%s\\n' " + shellQuote(banner),
		"rm -f " + shellQuote(f.Name()) + " >/dev/null 2>&1 || true",
		"export SHELL=/bin/bash",
	}, "\n")

	if _, err := f.WriteString(content + "\n"); err != nil {
		return "", err
	}

	return f.Name(), nil
}

func writeZshInitDir(baseDir, cleanPath string, versions map[string]string) (string, error) {
	sandboxHome := filepath.Join(baseDir, "runtime", "sandbox", "home")
	mysqlSock := filepath.Join(baseDir, "runtime", "data", "mysql", "mysql.sock")
	banner := buildBanner(versions, sandboxHome)
	logPath := filepath.Join(baseDir, "runtime", "logs", "php-fpm.log")
	mysqlBin, mariadbBin, mysqldumpBin, mysqladminBin := resolveServiceBinaries(baseDir, versions)

	dir, err := os.MkdirTemp("", "lv_zsh_")
	if err != nil {
		return "", err
	}

	zprofile := strings.Join([]string{
		"# LocalValet isolated zsh profile",
		"export PATH=" + shellQuote(cleanPath),
	}, "\n") + "\n"

	zshrc := strings.Join([]string{
		"# LocalValet isolated zsh init",
		"export LOCALVALET_ACTIVE=1",
		"export LOCALVALET_SANDBOX=1",
		"export PATH=" + shellQuote(cleanPath),
		"export MYSQL_UNIX_PORT=" + shellQuote(mysqlSock),
		"export MYSQL_TCP_PORT=3306",
		"export MYSQL_HOST=127.0.0.1",
		"export DB_CONNECTION=mysql",
		"export DB_HOST=127.0.0.1",
		"export DB_PORT=3306",
		"export DB_SOCKET=" + shellQuote(mysqlSock),
		"export DB_USERNAME=root",
		"export DB_DATABASE=localvalet",
		"export PROMPT='%F{green}[LocalValet:sandbox]%f %F{blue}%~%f %# '",
		"alias mysql='" + mysqlBin + " --socket=" + mysqlSock + " -u root'",
		"alias mariadb='" + mariadbBin + " --socket=" + mysqlSock + " -u root'",
		"alias mysqldump='" + mysqldumpBin + " --socket=" + mysqlSock + " -u root'",
		"alias mysqladmin='" + mysqladminBin + " --socket=" + mysqlSock + " -u root'",
		"alias db-shell='" + mysqlBin + " --socket=" + mysqlSock + " -u root'",
		"alias artisan='php artisan'",
		"alias php-logs='tail -f " + shellQuote(logPath) + "'",
		"alias valet-restart='echo \"Use LocalValet UI to restart services\"'",
		"alias sudo='sudo env \"PATH=$PATH\"'",
		"alias ll='ls -la --color=auto'",
		"alias l='ls -CF --color=auto'",
		"alias la='ls -A --color=auto'",
		"if [ -f ./.localvaletrc ]; then . ./.localvaletrc || true; fi",
		"printf '%s\\n' " + shellQuote(banner),
		"rm -f " + shellQuote(filepath.Join(dir, ".zprofile")) + " >/dev/null 2>&1 || true",
		"rm -f " + shellQuote(filepath.Join(dir, ".zshrc")) + " >/dev/null 2>&1 || true",
		"rmdir " + shellQuote(dir) + " >/dev/null 2>&1 || true",
	}, "\n") + "\n"

	if err := os.WriteFile(filepath.Join(dir, ".zprofile"), []byte(zprofile), 0o600); err != nil {
		return "", err
	}

	if err := os.WriteFile(filepath.Join(dir, ".zshrc"), []byte(zshrc), 0o600); err != nil {
		return "", err
	}

	return dir, nil
}

func writeFallbackInitScript(baseDir, shellPath, cleanPath string, versions map[string]string) (string, error) {
	sandboxHome := filepath.Join(baseDir, "runtime", "sandbox", "home")
	mysqlSock := filepath.Join(baseDir, "runtime", "data", "mysql", "mysql.sock")
	banner := buildBanner(versions, sandboxHome)
	logPath := filepath.Join(baseDir, "runtime", "logs", "php-fpm.log")
	mysqlBin, mariadbBin, mysqldumpBin, mysqladminBin := resolveServiceBinaries(baseDir, versions)

	f, err := os.CreateTemp("", "lv_init_*.sh")
	if err != nil {
		return "", err
	}
	defer f.Close()

	content := strings.Join([]string{
		"# LocalValet terminal init",
		"export LOCALVALET_ACTIVE=1",
		"export LOCALVALET_SANDBOX=1",
		"export PATH=" + shellQuote(cleanPath),
		"export MYSQL_UNIX_PORT=" + shellQuote(mysqlSock),
		"export MYSQL_TCP_PORT=3306",
		"export MYSQL_HOST=127.0.0.1",
		"export DB_CONNECTION=mysql",
		"export DB_HOST=127.0.0.1",
		"export DB_PORT=3306",
		"export DB_SOCKET=" + shellQuote(mysqlSock),
		"export DB_USERNAME=root",
		"export DB_DATABASE=localvalet",
		"alias mysql='" + mysqlBin + " --socket=" + mysqlSock + " -u root'",
		"alias mariadb='" + mariadbBin + " --socket=" + mysqlSock + " -u root'",
		"alias mysqldump='" + mysqldumpBin + " --socket=" + mysqlSock + " -u root'",
		"alias mysqladmin='" + mysqladminBin + " --socket=" + mysqlSock + " -u root'",
		"alias db-shell='" + mysqlBin + " --socket=" + mysqlSock + " -u root'",
		"alias artisan='php artisan'",
		"alias php-logs='tail -f " + shellQuote(logPath) + "'",
		"alias valet-restart='echo \"Use LocalValet UI to restart services\"'",
		"alias sudo='sudo env \"PATH=$PATH\"'",
		"if [ -f ./.localvaletrc ]; then . ./.localvaletrc || true; fi",
		"printf '%s\\n' " + shellQuote(banner),
		"rm -f " + shellQuote(f.Name()) + " >/dev/null 2>&1 || true",
		"export SHELL=" + shellQuote(shellPath),
	}, "\n")

	if _, err := f.WriteString(content + "\n"); err != nil {
		return "", err
	}

	return f.Name(), nil
}

func buildBanner(versions map[string]string, sandboxHome string) string {
	php := versions["php-fpm"]
	node := versions["node"]
	mysql := versions["mysql"]
	nginx := versions["nginx"]

	if php == "" {
		php = "-"
	}
	if node == "" {
		node = "-"
	}
	if mysql == "" {
		mysql = "-"
	}
	if nginx == "" {
		nginx = "-"
	}

	line1 := fmt.Sprintf("LocalValet Isolated Environment (Sandbox Mode)")
	line2 := fmt.Sprintf("PHP: %s | Node: %s | MySQL: %s | Nginx: %s", php, node, mysql, nginx)
	line3 := fmt.Sprintf("Sandbox HOME: %s", sandboxHome)

	return fmt.Sprintf("\n=== %s ===\n %s\n %s\n=======================================================\n", line1, line2, line3)
}


func activeVersions(cfg *runtimeConfig) map[string]string {
	result := make(map[string]string)
	for serviceName, service := range cfg.Services {
		if service == nil {
			continue
		}
		result[serviceName] = service.ActiveVersion
	}
	return result
}

func detectDefaultShell() string {
	candidates := []string{"/bin/bash", "/usr/bin/bash", "/bin/zsh", "/usr/bin/zsh", "/bin/sh", "/usr/bin/sh"}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return "/bin/bash"
}

func loadRuntimeConfig(path string) (*runtimeConfig, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read runtime config at %s: %w", path, err)
	}

	var cfg runtimeConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("invalid runtime config: %w", err)
	}
	return &cfg, nil
}

func resolveWorkDir(projectDir string) string {
	if projectDir != "" {
		abs, err := filepath.Abs(projectDir)
		if err == nil && dirExists(abs) {
			if root, ok := findProjectRoot(abs); ok {
				return root
			}
			return abs
		}
	}

	cwd, err := os.Getwd()
	if err == nil {
		if root, ok := findProjectRoot(cwd); ok {
			return root
		}
		return cwd
	}

	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		return home
	}

	return "."
}

func findProjectRoot(start string) (string, bool) {
	current := filepath.Clean(start)
	for {
		if fileExists(filepath.Join(current, "composer.json")) || fileExists(filepath.Join(current, "package.json")) {
			return current, true
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", false
		}
		current = parent
	}
}

func findBaseDir() string {
	// 1. Check exeDir and its ancestors
	if exePath, err := os.Executable(); err == nil && exePath != "" {
		dir := filepath.Dir(exePath)
		for {
			if fileExists(filepath.Join(dir, "config", "runtime.json")) {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir || parent == "." || parent == "/" {
				break
			}
			dir = parent
		}
	}

	// 2. Check cwd and its ancestors
	if cwd, err := os.Getwd(); err == nil && cwd != "" {
		dir := cwd
		for {
			if fileExists(filepath.Join(dir, "config", "runtime.json")) {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir || parent == "." || parent == "/" {
				break
			}
			dir = parent
		}
	}

	return "."
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

func resolvePath(baseDir, value string) string {
	if filepath.IsAbs(value) {
		return value
	}
	return filepath.Join(baseDir, filepath.FromSlash(value))
}

func isUnderPath(path, root string) bool {
	rel, err := filepath.Rel(filepath.Clean(root), filepath.Clean(path))
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func shellQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\"'\"'")
	return "'" + escaped + "'"
}

func escapeSingleQuotes(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
