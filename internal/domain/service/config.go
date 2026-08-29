package service

// Config holds display and canonical service identifiers and runtime metadata.
type Config struct {
	DisplayName     string   `json:"displayName"`
	ServiceName     string   `json:"serviceName"`
	DefaultPort     int      `json:"defaultPort"`
	Category        string   `json:"category"`        // "Web", "Database", "Runtime"
	Dependencies    []string `json:"dependencies"`    // service names that must be running
	HealthCheckType string   `json:"healthCheckType"` // "tcp", "http", "process"
}

// DefaultConfigs returns all available service configurations.
func DefaultConfigs() []Config {
	return []Config{
		{
			DisplayName:     "Apache",
			ServiceName:     "apache",
			DefaultPort:     8080,
			Category:        "Web",
			Dependencies:    []string{"php-fpm"},
			HealthCheckType: "http",
		},
		{
			DisplayName:     "Nginx",
			ServiceName:     "nginx",
			DefaultPort:     8080,
			Category:        "Web",
			Dependencies:    []string{"php-fpm"},
			HealthCheckType: "http",
		},
		{
			DisplayName:     "MySQL",
			ServiceName:     "mysql",
			DefaultPort:     3306,
			Category:        "Database",
			Dependencies:    []string{},
			HealthCheckType: "tcp",
		},
		{
			DisplayName:     "PostgreSQL",
			ServiceName:     "postgresql",
			DefaultPort:     5432,
			Category:        "Database",
			Dependencies:    []string{},
			HealthCheckType: "tcp",
		},
		{
			DisplayName:     "Redis",
			ServiceName:     "redis",
			DefaultPort:     6379,
			Category:        "Database",
			Dependencies:    []string{},
			HealthCheckType: "tcp",
		},
		{
			DisplayName:     "PHP-FPM",
			ServiceName:     "php-fpm",
			DefaultPort:     9074,
			Category:        "Runtime",
			Dependencies:    []string{},
			HealthCheckType: "tcp",
		},
	}
}

// CanonicalName resolves display name to canonical service name.
func CanonicalName(displayName string, configs []Config) string {
	for _, config := range configs {
		if config.DisplayName == displayName || config.ServiceName == displayName {
			return config.ServiceName
		}
	}

	// Return display name as fallback.
	return displayName
}

// GetConfig finds a service configuration by its canonical name or display name.
func GetConfig(serviceName string, configs []Config) (Config, bool) {
	for _, cfg := range configs {
		if cfg.ServiceName == serviceName || cfg.DisplayName == serviceName {
			return cfg, true
		}
	}
	return Config{}, false
}

// GetDependencies returns the list of service names that serviceName depends on.
func GetDependencies(serviceName string, configs []Config) []string {
	if cfg, ok := GetConfig(serviceName, configs); ok {
		return append([]string(nil), cfg.Dependencies...)
	}
	return nil
}

// GetDependents returns the list of service names that depend on serviceName.
func GetDependents(serviceName string, configs []Config) []string {
	var dependents []string
	for _, cfg := range configs {
		for _, dep := range cfg.Dependencies {
			if dep == serviceName {
				dependents = append(dependents, cfg.ServiceName)
				break
			}
		}
	}
	return dependents
}

