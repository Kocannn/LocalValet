package terminal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureSandboxDirectories(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lv_test_sandbox_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sandboxDir, sandboxHome := ensureSandboxDirectories(tempDir)

	expectedDirs := []string{
		sandboxHome,
		filepath.Join(sandboxHome, ".config"),
		filepath.Join(sandboxHome, ".local", "share"),
		filepath.Join(sandboxHome, ".local", "state"),
		filepath.Join(sandboxHome, ".cache"),
		filepath.Join(sandboxDir, "composer", "vendor", "bin"),
		filepath.Join(sandboxDir, "composer", "cache"),
		filepath.Join(sandboxDir, "npm", "global", "bin"),
		filepath.Join(sandboxDir, "npm", "cache"),
		filepath.Join(sandboxDir, "conf.d"),
		filepath.Join(sandboxDir, "tmp"),
		filepath.Join(tempDir, "runtime", "data", "mysql"),
	}

	for _, dir := range expectedDirs {
		info, err := os.Stat(dir)
		if err != nil {
			t.Errorf("expected directory %s to exist, error: %v", dir, err)
		} else if !info.IsDir() {
			t.Errorf("expected %s to be a directory", dir)
		}
	}

	// Verify isolated .my.cnf and localvalet_mysql.ini were written
	myCnfPath := filepath.Join(sandboxHome, ".my.cnf")
	if !fileExists(myCnfPath) {
		t.Errorf("expected %s to be generated", myCnfPath)
	}

	phpIniPath := filepath.Join(sandboxDir, "conf.d", "localvalet_mysql.ini")
	if !fileExists(phpIniPath) {
		t.Errorf("expected %s to be generated", phpIniPath)
	}
}

func TestBuildAppEnvironment_SandboxIsolation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lv_test_env_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	workDir := filepath.Join(tempDir, "project")
	_ = os.MkdirAll(workDir, 0o755)

	injectedPaths := []string{
		filepath.Join(tempDir, "runtime", "linux", "php", "8.4", "bin"),
		filepath.Join(tempDir, "runtime", "linux", "node", "22", "bin"),
	}
	for _, p := range injectedPaths {
		_ = os.MkdirAll(p, 0o755)
	}

	envMap := map[string]string{
		"PHP_INI_SCAN_DIR": "/custom/php/ini",
	}

	envSlice, cleanPath := buildAppEnvironment(tempDir, workDir, injectedPaths, envMap)

	env := make(map[string]string)
	for _, item := range envSlice {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}

	expectedSandboxHome := filepath.Join(tempDir, "runtime", "sandbox", "home")
	if env["HOME"] != expectedSandboxHome {
		t.Errorf("expected HOME to be sandbox home %q, got %q", expectedSandboxHome, env["HOME"])
	}

	if env["LOCALVALET_SANDBOX"] != "1" {
		t.Errorf("expected LOCALVALET_SANDBOX=1, got %q", env["LOCALVALET_SANDBOX"])
	}

	if env["LOCALVALET_ENV"] != "isolated" {
		t.Errorf("expected LOCALVALET_ENV=isolated, got %q", env["LOCALVALET_ENV"])
	}

	expectedSock := filepath.Join(tempDir, "runtime", "data", "mysql", "mysql.sock")
	if env["MYSQL_UNIX_PORT"] != expectedSock {
		t.Errorf("expected MYSQL_UNIX_PORT=%s, got %s", expectedSock, env["MYSQL_UNIX_PORT"])
	}

	if env["DB_SOCKET"] != expectedSock {
		t.Errorf("expected DB_SOCKET=%s, got %s", expectedSock, env["DB_SOCKET"])
	}

	if env["COMPOSER_HOME"] != filepath.Join(tempDir, "runtime", "sandbox", "composer") {
		t.Errorf("expected COMPOSER_HOME to point to sandbox composer, got %q", env["COMPOSER_HOME"])
	}

	if env["NPM_CONFIG_PREFIX"] != filepath.Join(tempDir, "runtime", "sandbox", "npm", "global") {
		t.Errorf("expected NPM_CONFIG_PREFIX to point to sandbox npm global, got %q", env["NPM_CONFIG_PREFIX"])
	}

	if env["PHP_INI_SCAN_DIR"] != "/custom/php/ini" {
		t.Errorf("expected PHP_INI_SCAN_DIR to be preserved, got %q", env["PHP_INI_SCAN_DIR"])
	}

	if !strings.Contains(cleanPath, "/usr/bin") || !strings.Contains(cleanPath, "php/8.4/bin") {
		t.Errorf("expected cleanPath to contain injected php and system /usr/bin, got: %s", cleanPath)
	}
}

func TestWriteBashInitScript(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lv_test_bash_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cleanPath := "/runtime/php/bin:/usr/bin:/bin"
	versions := map[string]string{
		"php-fpm": "8.4",
		"node":    "22",
		"mysql":   "12.2.2",
		"nginx":   "1.26",
	}

	scriptPath, err := writeBashInitScript(tempDir, cleanPath, versions)
	if err != nil {
		t.Fatalf("failed to write bash init script: %v", err)
	}

	contentBytes, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("failed to read created script: %v", err)
	}
	_ = os.Remove(scriptPath)

	content := string(contentBytes)
	if !strings.Contains(content, "export LOCALVALET_SANDBOX=1") {
		t.Errorf("script missing LOCALVALET_SANDBOX=1: %s", content)
	}
	if !strings.Contains(content, "export MYSQL_UNIX_PORT=") {
		t.Errorf("script missing MYSQL_UNIX_PORT: %s", content)
	}
	if !strings.Contains(content, "alias mysql=") {
		t.Errorf("script missing mysql alias: %s", content)
	}
	if !strings.Contains(content, "export PS1=") || !strings.Contains(content, "[LocalValet:sandbox]") {
		t.Errorf("script missing sandbox PS1 prompt: %s", content)
	}
	if !strings.Contains(content, "alias artisan='php artisan'") {
		t.Errorf("script missing artisan alias: %s", content)
	}
	if !strings.Contains(content, cleanPath) {
		t.Errorf("script missing clean PATH export: %s", content)
	}
}

func TestRealProjectTerminalPaths(t *testing.T) {
	baseDir := findBaseDir()
	if !fileExists(filepath.Join(baseDir, "config", "runtime.json")) {
		t.Fatalf("findBaseDir() failed to locate valid project root: %s", baseDir)
	}

	cfg, err := loadRuntimeConfig(filepath.Join(baseDir, "config", "runtime.json"))
	if err != nil {
		t.Fatalf("loadRuntimeConfig failed: %v", err)
	}

	paths, envMap := buildInjectedEnv(baseDir, cfg)
	if len(paths) == 0 {
		t.Fatalf("expected injected paths to be non-empty")
	}

	envSlice, cleanPath := buildAppEnvironment(baseDir, baseDir, paths, envMap)
	_ = envSlice

	expectedMySQLBin := filepath.Join(baseDir, "runtime", "linux", "mysql", "12.2.2", "bin")
	if !strings.Contains(cleanPath, expectedMySQLBin) {
		t.Errorf("cleanPath missing MySQL runtime bin: %s", cleanPath)
	}

	// Verify MySQL path comes before /usr/bin
	mysqlIdx := strings.Index(cleanPath, expectedMySQLBin)
	usrBinIdx := strings.Index(cleanPath, "/usr/bin")
	if mysqlIdx > usrBinIdx {
		t.Errorf("MySQL bin (%d) must precede /usr/bin (%d) in PATH", mysqlIdx, usrBinIdx)
	}
}
