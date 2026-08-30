package vhost

import (
	projectdomain "LocalValet/internal/domain/project"
	vhostdomain "LocalValet/internal/domain/vhost"
	sslusecase "LocalValet/internal/usecase/ssl"
	"fmt"
	"path/filepath"
)

type UseCase struct {
	generator vhostdomain.Generator
	sslUC     *sslusecase.UseCase
}

func New(generator vhostdomain.Generator, sslUC *sslusecase.UseCase) *UseCase {
	return &UseCase{
		generator: generator,
		sslUC:     sslUC,
	}
}

func (u *UseCase) EnableVHost(proj projectdomain.Project) (string, error) {
	docRoot := proj.Path
	if proj.WebRoot != "" && proj.WebRoot != "." {
		docRoot = filepath.Join(proj.Path, proj.WebRoot)
	}

	cfg := vhostdomain.VHostConfig{
		Domain:       proj.Domain,
		ProjectName:  proj.Name,
		DocumentRoot: docRoot,
		HTTPPort:     80,
		HTTPSPort:    443,
	}

	if proj.TargetPort > 0 {
		cfg.ProxyPass = fmt.Sprintf("http://127.0.0.1:%d", proj.TargetPort)
	} else {
		cfg.PHPFpmAddress = "127.0.0.1:9074"
	}

	if proj.SSLEnabled && u.sslUC != nil {
		certPath, keyPath, exists := u.sslUC.GetCertPaths(proj.Domain)
		if !exists {
			pair, err := u.sslUC.CreateCertForProject(proj.Domain)
			if err == nil {
				certPath = pair.CertPath
				keyPath = pair.KeyPath
				exists = true
			}
		}

		if exists {
			cfg.SSLEnabled = true
			cfg.SSLCertPath = certPath
			cfg.SSLKeyPath = keyPath
		}
	}

	path, err := u.generator.Generate(cfg)
	if err != nil {
		return "", err
	}

	_ = u.generator.ReloadNginx()
	return path, nil
}

func (u *UseCase) DisableVHost(domain string) error {
	err := u.generator.Remove(domain)
	if err != nil {
		return err
	}
	_ = u.generator.ReloadNginx()
	return nil
}

func (u *UseCase) ListVHosts() ([]string, error) {
	return u.generator.List()
}
