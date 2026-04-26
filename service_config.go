package main

import servicedomain "LocalValet/internal/domain/service"

// ServiceConfig is exposed to Wails bindings.
type ServiceConfig = servicedomain.Config

// GetServiceConfigs returns all available service configurations
func GetServiceConfigs() []ServiceConfig {
	return servicedomain.DefaultConfigs()
}

// GetServiceName returns the appropriate service name for current OS
func GetServiceName(displayName string) string {
	return servicedomain.CanonicalName(displayName, GetServiceConfigs())
}

// GetAvailableServices returns list of services for frontend
func (a *App) GetAvailableServices() []ServiceConfig {
	return GetServiceConfigs()
}
