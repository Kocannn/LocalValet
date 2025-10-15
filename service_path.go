package main

import (
	"os"
	"path/filepath"
	"runtime"
)

// ServiceBinaryConfig holds binary path configuration for Windows
type ServiceBinaryConfig struct {
	ServiceName  string // Name of the service (e.g., "apache", "nginx")
	BinaryName   string // Executable name (e.g., "httpd.exe", "nginx.exe")
	RelativePath string // Relative path from bin/windows/{service}/ (optional)
}

// GetWindowsBinaryConfigs returns binary configurations for Windows services
func GetWindowsBinaryConfigs() map[string]ServiceBinaryConfig {
	return map[string]ServiceBinaryConfig{
		"Apache": {
			ServiceName:  "apache",
			BinaryName:   "httpd.exe",
			RelativePath: "bin",
		},
		"MySQL": {
			ServiceName:  "mysql",
			BinaryName:   "mysqld.exe",
			RelativePath: "bin",
		},
		"PostgreSQL": {
			ServiceName:  "postgresql",
			BinaryName:   "pg_ctl.exe",
			RelativePath: "bin",
		},
		"Redis": {
			ServiceName:  "redis",
			BinaryName:   "redis-server.exe",
			RelativePath: "",
		},
		"Nginx": {
			ServiceName:  "nginx",
			BinaryName:   "nginx.exe",
			RelativePath: "",
		},
		"PHP-FPM": {
			ServiceName:  "php",
			BinaryName:   "php-cgi.exe",
			RelativePath: "",
		},
	}
}

// GetExecutablePath returns the full path to service executable
// For Windows: returns path to binary in bin/windows/{service}/
// For Linux/macOS: returns the service name (uses system binaries)
func GetExecutablePath(displayName string) string {
	if runtime.GOOS == "windows" {
		return getWindowsExecutablePath(displayName)
	}

	// For Linux/macOS, return the service name (system will find it)
	return GetServiceName(displayName)
}

// getWindowsExecutablePath constructs the full path to Windows binary
func getWindowsExecutablePath(displayName string) string {
	binaryConfigs := GetWindowsBinaryConfigs()
	config, exists := binaryConfigs[displayName]

	if !exists {
		// Fallback to service name if config not found
		return GetServiceName(displayName)
	}

	// Get the executable directory (where the compiled binary is)
	exePath, err := os.Executable()
	if err != nil {
		return GetServiceName(displayName)
	}
	exeDir := filepath.Dir(exePath)

	// Construct path: {exe_dir}/bin/windows/{service}/{relative_path}/{binary}
	binPath := filepath.Join(exeDir, "bin", "windows", config.ServiceName)

	if config.RelativePath != "" {
		binPath = filepath.Join(binPath, config.RelativePath)
	}

	binPath = filepath.Join(binPath, config.BinaryName)

	// Check if the binary exists
	if _, err := os.Stat(binPath); err == nil {
		return binPath
	}

	// Fallback to service name if binary not found
	return GetServiceName(displayName)
}

// GetServiceWorkingDirectory returns the working directory for the service
// This is useful for services that need to be executed from their installation directory
func GetServiceWorkingDirectory(displayName string) string {
	if runtime.GOOS == "windows" {
		binaryConfigs := GetWindowsBinaryConfigs()
		config, exists := binaryConfigs[displayName]

		if !exists {
			return ""
		}

		exePath, err := os.Executable()
		if err != nil {
			return ""
		}
		exeDir := filepath.Dir(exePath)

		// Return the service directory: {exe_dir}/bin/windows/{service}/
		return filepath.Join(exeDir, "bin", "windows", config.ServiceName)
	}

	return ""
}

// IsUsingSystemBinaries returns true if the OS uses system binaries
func IsUsingSystemBinaries() bool {
	return runtime.GOOS != "windows"
}

// ValidateWindowsBinaries checks if all required Windows binaries exist
func ValidateWindowsBinaries() map[string]bool {
	if runtime.GOOS != "windows" {
		return nil
	}

	results := make(map[string]bool)
	binaryConfigs := GetWindowsBinaryConfigs()

	for displayName := range binaryConfigs {
		path := getWindowsExecutablePath(displayName)
		_, err := os.Stat(path)
		results[displayName] = err == nil
	}

	return results
}
