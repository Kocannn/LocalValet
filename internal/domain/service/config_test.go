package service

import (
	"testing"
)

func TestDefaultConfigs(t *testing.T) {
	configs := DefaultConfigs()
	if len(configs) != 6 {
		t.Fatalf("expected 6 default configs, got %d", len(configs))
	}

	expectedServices := map[string]struct {
		category string
		port     int
	}{
		"apache":     {category: "Web", port: 8080},
		"nginx":      {category: "Web", port: 8080},
		"mysql":      {category: "Database", port: 3306},
		"postgresql": {category: "Database", port: 5432},
		"redis":      {category: "Database", port: 6379},
		"php-fpm":    {category: "Runtime", port: 9074},
	}

	for serviceName, expected := range expectedServices {
		cfg, ok := GetConfig(serviceName, configs)
		if !ok {
			t.Errorf("expected config for service %q to exist", serviceName)
			continue
		}
		if cfg.Category != expected.category {
			t.Errorf("service %q category = %q, expected %q", serviceName, cfg.Category, expected.category)
		}
		if cfg.DefaultPort != expected.port {
			t.Errorf("service %q default port = %d, expected %d", serviceName, cfg.DefaultPort, expected.port)
		}
	}
}

func TestCanonicalName(t *testing.T) {
	configs := DefaultConfigs()

	tests := []struct {
		input    string
		expected string
	}{
		{"Apache", "apache"},
		{"apache", "apache"},
		{"MySQL", "mysql"},
		{"PHP-FPM", "php-fpm"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		result := CanonicalName(tt.input, configs)
		if result != tt.expected {
			t.Errorf("CanonicalName(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestDependenciesAndDependents(t *testing.T) {
	configs := DefaultConfigs()

	// Apache depends on PHP-FPM
	deps := GetDependencies("apache", configs)
	if len(deps) != 1 || deps[0] != "php-fpm" {
		t.Errorf("GetDependencies(apache) = %v, expected [php-fpm]", deps)
	}

	// Unknown service dependencies
	depsUnknown := GetDependencies("unknown", configs)
	if depsUnknown != nil {
		t.Errorf("GetDependencies(unknown) = %v, expected nil", depsUnknown)
	}

	// Dependents of PHP-FPM should be Apache and Nginx
	dependents := GetDependents("php-fpm", configs)
	if len(dependents) != 2 {
		t.Fatalf("expected 2 dependents for php-fpm, got %d: %v", len(dependents), dependents)
	}

	hasApache := false
	hasNginx := false
	for _, dep := range dependents {
		if dep == "apache" {
			hasApache = true
		}
		if dep == "nginx" {
			hasNginx = true
		}
	}

	if !hasApache || !hasNginx {
		t.Errorf("expected dependents [apache, nginx], got %v", dependents)
	}
}
