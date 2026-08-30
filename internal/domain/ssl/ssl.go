package ssl

// CertPair holds file paths and metadata for an issued certificate.
type CertPair struct {
	Domain    string `json:"domain"`
	CertPath  string `json:"certPath"`
	KeyPath   string `json:"keyPath"`
	CreatedAt string `json:"createdAt"`
	ExpiresAt string `json:"expiresAt"`
	IsTrusted bool   `json:"isTrusted"`
}

// Manager manages Root Certificate Authority generation and per-domain certificate signing.
type Manager interface {
	EnsureRootCA() (caCertPath string, caKeyPath string, err error)
	GenerateCert(domain string) (CertPair, error)
	GetCertPaths(domain string) (certPath, keyPath string, exists bool)
	GetCACertPath() string
	InstallRootCA() error
	IsRootCATrusted() bool
}

