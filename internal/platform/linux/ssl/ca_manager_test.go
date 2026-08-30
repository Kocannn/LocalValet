package ssl

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"
)

func TestCAManager_FullLifecycle(t *testing.T) {
	tempDir := t.TempDir()
	caMgr := NewCAManagerWithPath(tempDir)

	// 1. Ensure Root CA
	caCertPath, caKeyPath, err := caMgr.EnsureRootCA()
	if err != nil {
		t.Fatalf("EnsureRootCA error: %v", err)
	}

	if !fileExists(caCertPath) || !fileExists(caKeyPath) {
		t.Fatalf("CA cert or key does not exist on disk")
	}

	// 2. Generate Domain Cert
	domain := "my-awesome-project.test"
	pair, err := caMgr.GenerateCert(domain)
	if err != nil {
		t.Fatalf("GenerateCert error: %v", err)
	}

	if pair.Domain != domain {
		t.Errorf("expected domain %s, got %s", domain, pair.Domain)
	}
	if !fileExists(pair.CertPath) || !fileExists(pair.KeyPath) {
		t.Fatalf("leaf cert or key does not exist on disk")
	}

	// 3. Cryptographic Verification
	caCertBytes, _ := os.ReadFile(caCertPath)
	caBlock, _ := pem.Decode(caCertBytes)
	rootCert, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to parse Root CA cert: %v", err)
	}

	roots := x509.NewCertPool()
	roots.AddCert(rootCert)

	leafCertBytes, _ := os.ReadFile(pair.CertPath)
	leafBlock, _ := pem.Decode(leafCertBytes)
	leafCert, err := x509.ParseCertificate(leafBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to parse Leaf cert: %v", err)
	}

	opts := x509.VerifyOptions{
		Roots:   roots,
		DNSName: domain,
	}

	if _, err := leafCert.Verify(opts); err != nil {
		t.Fatalf("leaf certificate verification against Root CA failed: %v", err)
	}

	// 4. Check SAN wildcard verification
	optsWildcard := x509.VerifyOptions{
		Roots:   roots,
		DNSName: "sub." + domain,
	}
	if _, err := leafCert.Verify(optsWildcard); err != nil {
		t.Fatalf("wildcard SAN verification failed: %v", err)
	}
}
