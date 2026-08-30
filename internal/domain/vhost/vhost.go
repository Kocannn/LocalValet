package vhost

// VHostConfig defines parameters needed to generate a virtual host configuration.
type VHostConfig struct {
	Domain         string `json:"domain"`
	ProjectName    string `json:"projectName"`
	DocumentRoot   string `json:"documentRoot"`
	PHPFpmAddress  string `json:"phpFpmAddress,omitempty"`
	ProxyPass      string `json:"proxyPass,omitempty"`
	SSLEnabled     bool   `json:"sslEnabled"`
	SSLCertPath    string `json:"sslCertPath,omitempty"`
	SSLKeyPath     string `json:"sslKeyPath,omitempty"`
	HTTPPort       int    `json:"httpPort,omitempty"`
	HTTPSPort      int    `json:"httpsPort,omitempty"`
}

// Generator manages generating and removing virtual host web server configurations.
type Generator interface {
	Generate(config VHostConfig) (string, error)
	Remove(domain string) error
	List() ([]string, error)
	GetVHostPath(domain string) string
	ReloadNginx() error
}
