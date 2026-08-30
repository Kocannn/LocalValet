package ssl

import (
	ssldomain "LocalValet/internal/domain/ssl"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)


type CAManager struct {
	mu       sync.RWMutex
	certsDir string
}

func NewCAManager() ssldomain.Manager {
	baseDir := findBaseDir()
	certsDir := filepath.Join(baseDir, "runtime", "certs")
	_ = os.MkdirAll(certsDir, 0o755)
	return &CAManager{certsDir: certsDir}
}

func NewCAManagerWithPath(certsDir string) ssldomain.Manager {
	_ = os.MkdirAll(certsDir, 0o755)
	return &CAManager{certsDir: certsDir}
}

func (m *CAManager) GetCACertPath() string {
	return filepath.Join(m.certsDir, "ca.crt")
}

func (m *CAManager) caKeyPath() string {
	return filepath.Join(m.certsDir, "ca.key")
}

func (m *CAManager) EnsureRootCA() (string, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	caCertPath := m.GetCACertPath()
	caKeyPath := m.caKeyPath()

	if fileExists(caCertPath) && fileExists(caKeyPath) {
		return caCertPath, caKeyPath, nil
	}

	// Generate Root CA Private Key
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate CA private key: %w", err)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate CA serial number: %w", err)
	}

	caTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:  []string{"LocalValet Development CA"},
			Country:       []string{"ID"},
			CommonName:    "LocalValet Local Root CA",
		},
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().AddDate(10, 0, 0), // 10 years validity
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		MaxPathLenZero:        false,
	}

	// Self-sign Root CA
	caCertDER, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to create Root CA certificate: %w", err)
	}

	// Write CA Cert PEM
	caCertFile, err := os.OpenFile(caCertPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return "", "", err
	}
	defer caCertFile.Close()
	if err := pem.Encode(caCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: caCertDER}); err != nil {
		return "", "", err
	}

	// Write CA Key PEM
	caKeyFile, err := os.OpenFile(caKeyPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return "", "", err
	}
	defer caKeyFile.Close()
	if err := pem.Encode(caKeyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(caPrivateKey)}); err != nil {
		return "", "", err
	}

	return caCertPath, caKeyPath, nil
}

func (m *CAManager) GenerateCert(domain string) (ssldomain.CertPair, error) {
	if domain == "" {
		return ssldomain.CertPair{}, fmt.Errorf("domain cannot be empty")
	}

	caCertPath, caKeyPath, err := m.EnsureRootCA()
	if err != nil {
		return ssldomain.CertPair{}, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Load CA Cert & Key
	caCertBytes, err := os.ReadFile(caCertPath)
	if err != nil {
		return ssldomain.CertPair{}, err
	}
	caCertBlock, _ := pem.Decode(caCertBytes)
	if caCertBlock == nil {
		return ssldomain.CertPair{}, fmt.Errorf("failed to parse CA cert PEM")
	}
	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return ssldomain.CertPair{}, err
	}

	caKeyBytes, err := os.ReadFile(caKeyPath)
	if err != nil {
		return ssldomain.CertPair{}, err
	}
	caKeyBlock, _ := pem.Decode(caKeyBytes)
	if caKeyBlock == nil {
		return ssldomain.CertPair{}, fmt.Errorf("failed to parse CA key PEM")
	}
	caKey, err := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return ssldomain.CertPair{}, err
	}

	// Generate Leaf Private Key
	leafKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return ssldomain.CertPair{}, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return ssldomain.CertPair{}, err
	}

	// Prepare Leaf Certificate
	dnsNames := []string{domain, "*." + domain, "localhost"}
	ipAddresses := []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")}

	leafTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"LocalValet Development"},
			CommonName:   domain,
		},
		DNSNames:              dnsNames,
		IPAddresses:           ipAddresses,
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().AddDate(2, 0, 0), // 2 years validity
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	leafDER, err := x509.CreateCertificate(rand.Reader, &leafTemplate, caCert, &leafKey.PublicKey, caKey)
	if err != nil {
		return ssldomain.CertPair{}, fmt.Errorf("failed to issue certificate for %s: %w", domain, err)
	}

	certPath := filepath.Join(m.certsDir, domain+".crt")
	keyPath := filepath.Join(m.certsDir, domain+".key")

	// Write Leaf Cert + CA Cert Chain
	certFile, err := os.OpenFile(certPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return ssldomain.CertPair{}, err
	}
	defer certFile.Close()

	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: leafDER}); err != nil {
		return ssldomain.CertPair{}, err
	}
	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: caCertBlock.Bytes}); err != nil {
		return ssldomain.CertPair{}, err
	}

	// Write Leaf Key PEM
	keyFile, err := os.OpenFile(keyPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return ssldomain.CertPair{}, err
	}
	defer keyFile.Close()

	if err := pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(leafKey)}); err != nil {
		return ssldomain.CertPair{}, err
	}

	return ssldomain.CertPair{
		Domain:    domain,
		CertPath:  certPath,
		KeyPath:   keyPath,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		ExpiresAt: time.Now().AddDate(2, 0, 0).Format("2006-01-02 15:04:05"),
		IsTrusted: true,
	}, nil
}

func (m *CAManager) GetCertPaths(domain string) (string, string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	certPath := filepath.Join(m.certsDir, domain+".crt")
	keyPath := filepath.Join(m.certsDir, domain+".key")

	if fileExists(certPath) && fileExists(keyPath) {
		return certPath, keyPath, true
	}
	return "", "", false
}

// InstallRootCA installs the LocalValet Root CA into the system trust store.
func (m *CAManager) InstallRootCA() error {
	caCertPath, _, err := m.EnsureRootCA()
	if err != nil {
		return err
	}

	targetDir := "/usr/local/share/ca-certificates"
	targetFile := filepath.Join(targetDir, "localvalet_ca.crt")

	caBytes, err := os.ReadFile(caCertPath)
	if err != nil {
		return err
	}

	// Try direct write
	if err := os.MkdirAll(targetDir, 0o755); err == nil {
		if err := os.WriteFile(targetFile, caBytes, 0o644); err == nil {
			cmd := exec.Command("update-ca-certificates")
			_ = cmd.Run()
			return nil
		}
	}

	// Try via pkexec
	cmdStr := fmt.Sprintf("mkdir -p %s && cp %s %s && update-ca-certificates", targetDir, caCertPath, targetFile)
	if _, err := exec.LookPath("pkexec"); err == nil {
		cmd := exec.Command("pkexec", "sh", "-c", cmdStr)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	// Try via sudo
	if _, err := exec.LookPath("sudo"); err == nil {
		cmd := exec.Command("sudo", "sh", "-c", cmdStr)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("root privileges required to install CA certificate into %s", targetFile)
}

// IsRootCATrusted checks if the LocalValet Root CA is installed in the system trust store.
func (m *CAManager) IsRootCATrusted() bool {
	targetFile := "/usr/local/share/ca-certificates/localvalet_ca.crt"
	if fileExists(targetFile) {
		return true
	}

	caCertPath := m.GetCACertPath()
	if !fileExists(caCertPath) {
		return false
	}

	caBytes, err := os.ReadFile(caCertPath)
	if err != nil {
		return false
	}
	block, _ := pem.Decode(caBytes)
	if block == nil {
		return false
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false
	}

	pool, err := x509.SystemCertPool()
	if err != nil || pool == nil {
		return false
	}

	opts := x509.VerifyOptions{
		Roots: pool,
	}
	_, err = cert.Verify(opts)
	return err == nil
}


func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func findBaseDir() string {
	if exePath, err := os.Executable(); err == nil && exePath != "" {
		dir := filepath.Dir(exePath)
		for {
			if fileExists(filepath.Join(dir, "config", "runtime.json")) {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir || parent == "." || parent == "/" {
				break
			}
			dir = parent
		}
	}

	if cwd, err := os.Getwd(); err == nil && cwd != "" {
		dir := cwd
		for {
			if fileExists(filepath.Join(dir, "config", "runtime.json")) {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir || parent == "." || parent == "/" {
				break
			}
			dir = parent
		}
	}

	return "."
}
