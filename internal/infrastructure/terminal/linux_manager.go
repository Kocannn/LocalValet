package terminal

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	terminaldomain "LocalValet/internal/domain/terminal"
)

type LinuxManager struct{}

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
	return &LinuxManager{}
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

	// Legacy support: versions can be flattened under service object.
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

func (m *LinuxManager) Launch(options terminaldomain.LaunchOptions) (terminaldomain.LaunchResult, error) {
	baseDir := findBaseDir()
	cfg, err := loadRuntimeConfig(filepath.Join(baseDir, "config", "runtime.json"))
	if err != nil {
		return terminaldomain.LaunchResult{}, err
	}

	workDir := resolveWorkDir(options.ProjectDir)
	preferred := strings.TrimSpace(options.PreferredTerminal)

	env, err := buildInjectedEnv(baseDir, cfg)
	if err != nil {
		return terminaldomain.LaunchResult{}, err
	}

	rcFile, err := writeRuntimeRC(baseDir)
	if err != nil {
		return terminaldomain.LaunchResult{}, err
	}

	termBin, termArgs, err := buildTerminalCommand(preferred, workDir, rcFile)
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
	return terminaldomain.LaunchResult{
		Terminal: termBin,
		WorkDir:  workDir,
	}, nil
}

func buildInjectedEnv(baseDir string, cfg *runtimeConfig) ([]string, error) {
	runtimeRoot := filepath.Join(baseDir, "runtime")
	paths := make([]string, 0, 8)

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

	phpSvc := cfg.Services["php-fpm"]
	if phpSvc != nil {
		phpVer := phpSvc.Versions[phpSvc.ActiveVersion]
		if phpVer != nil {
			phpRoot := filepath.Dir(filepath.Dir(resolvePath(baseDir, phpVer.Binary)))
			addPath(filepath.Join(phpRoot, "bin"))
			addPath(filepath.Join(phpRoot, "sbin"))
		}
	}

	mysqlSvc := cfg.Services["mysql"]
	if mysqlSvc != nil {
		mysqlVer := mysqlSvc.Versions[mysqlSvc.ActiveVersion]
		if mysqlVer != nil {
			addPath(filepath.Dir(resolvePath(baseDir, mysqlVer.Binary)))
		}
	}

	composerDir := filepath.Join(runtimeRoot, "linux", "composer")
	addPath(composerDir)
	addPath(filepath.Join(composerDir, "bin"))

	systemPath := os.Getenv("PATH")
	newPath := strings.Join(paths, string(os.PathListSeparator))
	if newPath != "" {
		newPath += string(os.PathListSeparator) + systemPath
	} else {
		newPath = systemPath
	}

	envMap := map[string]string{
		"PATH": newPath,
	}

	if phpSvc != nil {
		phpVer := phpSvc.Versions[phpSvc.ActiveVersion]
		if phpVer != nil {
			phpRoot := filepath.Dir(filepath.Dir(resolvePath(baseDir, phpVer.Binary)))
			iniScanDir := filepath.Join(phpRoot, "etc", "conf.d")
			if dirExists(iniScanDir) {
				envMap["PHP_INI_SCAN_DIR"] = iniScanDir
			}
		}
	}

	if mysqlSvc != nil {
		mysqlVer := mysqlSvc.Versions[mysqlSvc.ActiveVersion]
		if mysqlVer != nil {
			mysqlHome := resolvePath(baseDir, mysqlVer.WorkingDir)
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

	merged := append([]string{}, os.Environ()...)
	for key, value := range envMap {
		merged = appendOrReplaceEnv(merged, key, value)
	}
	return merged, nil
}

func buildTerminalCommand(preferred, workDir, rcFile string) (string, []string, error) {
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
			return bin, []string{"--directory", workDir, "bash", "--rcfile", rcFile}, nil
		case "alacritty":
			return bin, []string{"--working-directory", workDir, "-e", "bash", "--rcfile", rcFile}, nil
		case "foot":
			return bin, []string{"--working-directory", workDir, "bash", "--rcfile", rcFile}, nil
		case "gnome-terminal":
			return bin, []string{"--working-directory", workDir, "--", "bash", "--rcfile", rcFile}, nil
		case "konsole":
			return bin, []string{"--workdir", workDir, "-e", "bash", "--rcfile", rcFile}, nil
		case "xterm":
			return bin, []string{"-e", "bash", "--rcfile", rcFile}, nil
		default:
			return bin, []string{"-e", "bash", "--rcfile", rcFile}, nil
		}
	}

	return "", nil, fmt.Errorf("no supported terminal emulator found (tried: %v)", candidates)
}

func writeRuntimeRC(baseDir string) (string, error) {
	logPath := filepath.Join(baseDir, "runtime", "logs", "php-fpm.log")
	content := "" +
		"if [ -f ~/.bashrc ]; then . ~/.bashrc; fi\n" +
		"alias db-shell='mysql'\n" +
		"alias php-logs='tail -f " + shellQuote(logPath) + "'\n" +
		"alias valet-restart='echo \"Use LocalValet UI to restart services\"'\n"

	f, err := os.CreateTemp("", "localvalet-rc-*.sh")
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return "", err
	}

	return f.Name(), nil
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

	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		return home
	}

	cwd, err := os.Getwd()
	if err == nil {
		return cwd
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
	}

	if err == nil {
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
