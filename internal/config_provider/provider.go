package configprovider

import "runtime"

// 1. Definisikan Interface
type ConfigProvider interface {
    GetServiceConfigs() []ServiceConfig
    GetServiceName(displayName string) string
}

// 2. Definisikan Struct Konkret (Implementasi Nyata)
type DefaultConfigProvider struct{}

// 3. Buat Constructor
func NewConfigProvider() *DefaultConfigProvider {
    return &DefaultConfigProvider{}
}

// 4. Definisikan Struktur Data Config
type ServiceConfig struct {
    DisplayName string
    Linux       string
    Darwin      string
    Windows     string
}

// GetServiceConfigs returns all available service configurations
func (p *DefaultConfigProvider) GetServiceConfigs() []ServiceConfig {
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

// 6. Implementasi Method Interface: GetServiceName
// Logika ini dipindah ke sini agar App.go tidak pusing mikirin OS logic
func (p *DefaultConfigProvider) GetServiceName(displayName string) string {
    configs := p.GetServiceConfigs()

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
    return displayName
}
