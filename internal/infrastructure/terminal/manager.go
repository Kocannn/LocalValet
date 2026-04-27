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
	appEnv := buildAppEnvironment(baseDir, workDir, paths, envMap)

	switch runtime.GOOS {
	case "linux":
		return launchLinux(preferred, workDir, baseDir, paths, appEnv, versions)
	case "windows":
		return launchWindows(preferred, workDir, baseDir, paths, appEnv, versions)
	default:
		return terminaldomain.LaunchResult{}, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func launchLinux(preferred, workDir, baseDir string, paths []string, env []string, versions map[string]string) (terminaldomain.LaunchResult, error) {
	shellPath := detectDefaultShell()
	shellName := strings.ToLower(filepath.Base(shellPath))

	var shellArgs []string
	switch shellName {
	case "bash":
		initScript, err := writeBashInitScript(baseDir, versions)
		if err != nil {
			return terminaldomain.LaunchResult{}, err
		}
		shellArgs = []string{"--login", "--rcfile", initScript, "-i"}
	case "zsh":
		startupDir, err := writeZshInitDir(baseDir, versions)
		if err != nil {
			return terminaldomain.LaunchResult{}, err
		}
		env = appendOrReplaceEnv(env, "ZDOTDIR", startupDir)
		shellArgs = []string{"-l", "-i"}
	default:
		initScript, err := writeFallbackInitScript(baseDir, shellPath, versions)
		if err != nil {
			return terminaldomain.LaunchResult{}, err
		}
		bootstrap := ". " + shellQuote(initScript) + "; exec " + shellQuote(shellPath) + " -i"
		shellPath = "/bin/sh"
		shellArgs = []string{"-lc", bootstrap}
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
	_ = paths
	return terminaldomain.LaunchResult{Terminal: termBin, WorkDir: workDir}, nil
}

func launchWindows(preferred, workDir, baseDir string, paths []string, env []string, versions map[string]string) (terminaldomain.LaunchResult, error) {
	_ = baseDir
	_ = paths
	banner := buildBanner(versions)
	psCmd := "$Host.UI.RawUI.WindowTitle = 'LocalValet'; " +
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
				"--title", "LocalValet",
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

func buildInjectedEnv(baseDir string, cfg *runtimeConfig) ([]string, map[string]string) {
	runtimeRoot := filepath.Join(baseDir, "runtime")
	paths := make([]string, 0, 12)

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

	composerRoot := filepath.Join(runtimeRoot, "linux", "composer")
	addPath(composerRoot)
	addPath(filepath.Join(composerRoot, "bin"))

	envMap := map[string]string{}

	if phpVersion != nil {
		phpRoot := filepath.Dir(filepath.Dir(resolvePath(baseDir, phpVersion.Binary)))
		iniScanDir := filepath.Join(phpRoot, "etc", "conf.d")
		if dirExists(iniScanDir) {
			envMap["PHP_INI_SCAN_DIR"] = iniScanDir
		}
	}

	mysqlSvc := cfg.Services["mysql"]
	if mysqlSvc != nil {
		mysqlVersion := mysqlSvc.Versions[mysqlSvc.ActiveVersion]
		if mysqlVersion != nil {
			mysqlHome := resolvePath(baseDir, mysqlVersion.WorkingDir)
			if dirExists(mysqlHome) {
				envMap["MYSQL_HOME"] = mysqlHome
			}
		}
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

func buildAppEnvironment(baseDir, workDir string, paths []string, envMap map[string]string) []string {
	currentUser, err := user.Current()
	username := "localvalet"
	if err == nil && currentUser.Username != "" {
		username = currentUser.Username
	}

	homeDir, err := os.UserHomeDir()
	if err != nil || homeDir == "" {
		homeDir = workDir
	}

	allowedSessionKeys := []string{
		"DISPLAY",
		"WAYLAND_DISPLAY",
		"XDG_RUNTIME_DIR",
		"XAUTHORITY",
		"DBUS_SESSION_BUS_ADDRESS",
		"XDG_SESSION_TYPE",
	}

	pathParts := append([]string{}, paths...)
	pathParts = append(pathParts,
		"/usr/local/sbin",
		"/usr/local/bin",
		"/usr/sbin",
		"/usr/bin",
		"/sbin",
		"/bin",
	)

	env := map[string]string{
		"HOME":                homeDir,
		"LOGNAME":             username,
		"PATH":                strings.Join(pathParts, string(os.PathListSeparator)),
		"PWD":                 workDir,
		"SHELL":               detectDefaultShell(),
		"TERM":                "xterm-256color",
		"USER":                username,
		"LANG":                "en_US.UTF-8",
		"LC_ALL":              "en_US.UTF-8",
		"LOCALVALET_BASE_DIR": baseDir,
		"LOCALVALET_HOME":     homeDir,
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

	return result
}

func writeBashInitScript(baseDir string, versions map[string]string) (string, error) {
	banner := buildBanner(versions)
	logPath := filepath.Join(baseDir, "runtime", "logs", "php-fpm.log")

	f, err := os.CreateTemp("", "lv_init_*.sh")
	if err != nil {
		return "", err
	}
	defer f.Close()

	content := strings.Join([]string{
		"# LocalValet terminal init",
		"export LOCALVALET_ACTIVE=1",
		"alias artisan='php artisan'",
		"alias db-shell='mysql'",
		"alias php-logs='tail -f " + shellQuote(logPath) + "'",
		"alias valet-restart='echo \"Use LocalValet UI to restart services\"'",
		"alias sudo='sudo env \"PATH=$PATH\"'",
		"if [ -f ./.localvaletrc ]; then . ./.localvaletrc || true; fi",
		"printf '\\n%s\\n' " + shellQuote(banner),
		"rm -f " + shellQuote(f.Name()) + " >/dev/null 2>&1 || true",
		"export SHELL=/bin/bash",
	}, "\n")

	if _, err := f.WriteString(content + "\n"); err != nil {
		return "", err
	}

	return f.Name(), nil
}

func writeZshInitDir(baseDir string, versions map[string]string) (string, error) {
	banner := buildBanner(versions)
	logPath := filepath.Join(baseDir, "runtime", "logs", "php-fpm.log")

	dir, err := os.MkdirTemp("", "lv_zsh_")
	if err != nil {
		return "", err
	}

	zprofile := strings.Join([]string{
		"# LocalValet zsh profile",
	}, "\n") + "\n"

	zshrc := strings.Join([]string{
		"# LocalValet terminal init",
		"export LOCALVALET_ACTIVE=1",
		"alias artisan='php artisan'",
		"alias db-shell='mysql'",
		"alias php-logs='tail -f " + shellQuote(logPath) + "'",
		"alias valet-restart='echo \"Use LocalValet UI to restart services\"'",
		"alias sudo='sudo env \"PATH=$PATH\"'",
		"if [ -f ./.localvaletrc ]; then . ./.localvaletrc || true; fi",
		"printf '\\n%s\\n' " + shellQuote(banner),
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

func writeFallbackInitScript(baseDir, shellPath string, versions map[string]string) (string, error) {
	banner := buildBanner(versions)
	logPath := filepath.Join(baseDir, "runtime", "logs", "php-fpm.log")

	f, err := os.CreateTemp("", "lv_init_*.sh")
	if err != nil {
		return "", err
	}
	defer f.Close()

	content := strings.Join([]string{
		"# LocalValet terminal init",
		"export LOCALVALET_ACTIVE=1",
		"alias artisan='php artisan'",
		"alias db-shell='mysql'",
		"alias php-logs='tail -f " + shellQuote(logPath) + "'",
		"alias valet-restart='echo \"Use LocalValet UI to restart services\"'",
		"alias sudo='sudo env \"PATH=$PATH\"'",
		"if [ -f ./.localvaletrc ]; then . ./.localvaletrc || true; fi",
		"printf '\\n%s\\n' " + shellQuote(banner),
		"rm -f " + shellQuote(f.Name()) + " >/dev/null 2>&1 || true",
		"export SHELL=" + shellQuote(shellPath),
	}, "\n")

	if _, err := f.WriteString(content + "\n"); err != nil {
		return "", err
	}

	return f.Name(), nil
}

func buildBanner(versions map[string]string) string {
	php := versions["php-fpm"]
	mysql := versions["mysql"]
	nginx := versions["nginx"]

	if php == "" {
		php = "-"
	}
	if mysql == "" {
		mysql = "-"
	}
	if nginx == "" {
		nginx = "-"
	}

	return fmt.Sprintf("LocalValet Active | PHP %s | MySQL %s | Nginx %s", php, mysql, nginx)
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

	if exePath != "" {
		return filepath.Dir(exePath)
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
