package service

// Config holds display and canonical service identifiers.
type Config struct {
	DisplayName string
	ServiceName string
}

// DefaultConfigs returns all available service configurations.
func DefaultConfigs() []Config {
	return []Config{
		{DisplayName: "Apache", ServiceName: "apache"},
		{DisplayName: "MySQL", ServiceName: "mysql"},
		{DisplayName: "PostgreSQL", ServiceName: "postgresql"},
		{DisplayName: "Redis", ServiceName: "redis"},
		{DisplayName: "Nginx", ServiceName: "nginx"},
		{DisplayName: "PHP-FPM", ServiceName: "php-fpm"},
	}
}

// CanonicalName resolves display name to canonical service name.
func CanonicalName(displayName string, configs []Config) string {
	for _, config := range configs {
		if config.DisplayName == displayName {
			return config.ServiceName
		}
	}

	// Return display name as fallback.
	return displayName
}
