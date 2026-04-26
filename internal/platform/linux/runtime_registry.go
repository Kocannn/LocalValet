package linux

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type RuntimeRegistry struct{}

type runtimeVersion struct {
	Binary     string   `json:"binary"`
	Args       []string `json:"args"`
	WorkingDir string   `json:"workingDir"`
}

type runtimeService struct {
	ActiveVersion string                     `json:"activeVersion"`
	Versions      map[string]*runtimeVersion `json:"versions"`
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

	// Support legacy layout where versions are placed directly under the service object.
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

type runtimeConfig struct {
	Services map[string]*runtimeService `json:"services"`
}

type resolvedRuntime struct {
	BinaryPath string
	Args       []string
	WorkingDir string
}

func NewRuntimeRegistry() *RuntimeRegistry {
	return &RuntimeRegistry{}
}

func (r *RuntimeRegistry) Resolve(serviceName string) (*resolvedRuntime, error) {
	cfg, err := r.loadConfig()
	if err != nil {
		return nil, err
	}

	svc, ok := cfg.Services[serviceName]
	if !ok {
		return nil, fmt.Errorf("service %q is not configured in config/runtime.json", serviceName)
	}

	ver, ok := svc.Versions[svc.ActiveVersion]
	if !ok {
		return nil, fmt.Errorf("active version %q for service %q is not defined", svc.ActiveVersion, serviceName)
	}

	baseDir := r.baseDir()
	binaryPath := resolvePath(baseDir, ver.Binary)
	if _, err := os.Stat(binaryPath); err != nil {
		return nil, fmt.Errorf("binary for %s (%s) not found at %s", serviceName, svc.ActiveVersion, binaryPath)
	}

	args := make([]string, 0, len(ver.Args))
	for _, arg := range ver.Args {
		if looksLikePath(arg) {
			args = append(args, resolvePath(baseDir, arg))
			continue
		}
		args = append(args, arg)
	}

	workingDir := ""
	if ver.WorkingDir != "" {
		workingDir = resolvePath(baseDir, ver.WorkingDir)
	}

	return &resolvedRuntime{
		BinaryPath: binaryPath,
		Args:       args,
		WorkingDir: workingDir,
	}, nil
}

func (r *RuntimeRegistry) SetActiveVersion(serviceName, version string) error {
	cfg, err := r.loadConfig()
	if err != nil {
		return err
	}

	svc, ok := cfg.Services[serviceName]
	if !ok {
		return fmt.Errorf("service %q is not configured", serviceName)
	}

	if _, ok := svc.Versions[version]; !ok {
		return fmt.Errorf("version %q for service %q is not available", version, serviceName)
	}

	svc.ActiveVersion = version
	return r.saveConfig(cfg)
}

func (r *RuntimeRegistry) GetActiveVersion(serviceName string) (string, error) {
	cfg, err := r.loadConfig()
	if err != nil {
		return "", err
	}

	svc, ok := cfg.Services[serviceName]
	if !ok {
		return "", fmt.Errorf("service %q is not configured", serviceName)
	}

	return svc.ActiveVersion, nil
}

func (r *RuntimeRegistry) GetVersions(serviceName string) ([]string, error) {
	cfg, err := r.loadConfig()
	if err != nil {
		return nil, err
	}

	svc, ok := cfg.Services[serviceName]
	if !ok {
		return nil, fmt.Errorf("service %q is not configured", serviceName)
	}

	versions := make([]string, 0, len(svc.Versions))
	for version := range svc.Versions {
		versions = append(versions, version)
	}
	sort.Strings(versions)
	return versions, nil
}

func (r *RuntimeRegistry) pidFilePath(serviceName string) string {
	return filepath.Join(r.baseDir(), "runtime", "pids", serviceName+".pid")
}

func (r *RuntimeRegistry) logsDir() string {
	return filepath.Join(r.baseDir(), "runtime", "logs")
}

func (r *RuntimeRegistry) configPath() string {
	return filepath.Join(r.baseDir(), "config", "runtime.json")
}

func (r *RuntimeRegistry) loadConfig() (*runtimeConfig, error) {
	configPath := r.configPath()
	b, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read runtime config at %s: %w", configPath, err)
	}

	var cfg runtimeConfig
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("invalid runtime config: %w", err)
	}

	return &cfg, nil
}

func (r *RuntimeRegistry) saveConfig(cfg *runtimeConfig) error {
	configPath := r.configPath()
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}

	return os.WriteFile(configPath, b, 0o644)
}

func (r *RuntimeRegistry) baseDir() string {
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
	}

	if err == nil {
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func looksLikePath(value string) bool {
	if strings.HasPrefix(value, "/") || strings.HasPrefix(value, "./") || strings.HasPrefix(value, "../") {
		return true
	}

	return strings.Contains(value, "/")
}

func resolvePath(baseDir, value string) string {
	if filepath.IsAbs(value) {
		return value
	}

	return filepath.Join(baseDir, filepath.FromSlash(value))
}
