package main

import "runtime"

// ServiceConfig holds service configuration for different platforms
type ServiceConfig struct {
	DisplayName string
	Linux       string // systemctl service name
	Darwin      string // brew service name
	Windows     string // Windows service name
}

// GetServiceConfigs returns all available service configurations
func GetServiceConfigs() []ServiceConfig {
	return []ServiceConfig{
		{
			DisplayName: "Apache",
			Linux:       "apache2",
			Darwin:      "httpd",
			Windows:     "Apache2.4",
		},
		{
			DisplayName: "MySQL",
			Linux:       "mysql",
			Darwin:      "mysql",
			Windows:     "MySQL80",
		},
		{
			DisplayName: "PostgreSQL",
			Linux:       "postgresql",
			Darwin:      "postgresql",
			Windows:     "PostgreSQL",
		},
		{
			DisplayName: "Redis",
			Linux:       "redis-server",
			Darwin:      "redis",
			Windows:     "Redis",
		},
		{
			DisplayName: "Nginx",
			Linux:       "nginx",
			Darwin:      "nginx",
			Windows:     "nginx",
		},
		{
			DisplayName: "PHP-FPM",
			Linux:       "php8.1-fpm", // or php7.4-fpm, adjust version
			Darwin:      "php",
			Windows:     "php-cgi",
		},
	}
}

// GetServiceName returns the appropriate service name for current OS
func GetServiceName(displayName string) string {
	configs := GetServiceConfigs()

	for _, config := range configs {
		if config.DisplayName == displayName {
			switch runtime.GOOS {
			case "linux":
				return config.Linux
			case "darwin":
				return config.Darwin
			case "windows":
				return config.Windows
			}
		}
	}

	// Return display name as fallback
	return displayName
}

// GetAvailableServices returns list of services for frontend
func (a *App) GetAvailableServices() []ServiceConfig {
	return GetServiceConfigs()
}
