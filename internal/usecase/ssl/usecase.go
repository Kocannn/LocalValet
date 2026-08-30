package ssl

import (
	ssldomain "LocalValet/internal/domain/ssl"
)

type UseCase struct {
	manager ssldomain.Manager
}

func New(manager ssldomain.Manager) *UseCase {
	return &UseCase{manager: manager}
}

func (u *UseCase) EnsureRootCA() (string, string, error) {
	return u.manager.EnsureRootCA()
}

func (u *UseCase) CreateCertForProject(domain string) (ssldomain.CertPair, error) {
	return u.manager.GenerateCert(domain)
}

func (u *UseCase) GetCertPaths(domain string) (string, string, bool) {
	return u.manager.GetCertPaths(domain)
}

func (u *UseCase) GetCACertPath() string {
	return u.manager.GetCACertPath()
}

func (u *UseCase) InstallRootCA() error {
	return u.manager.InstallRootCA()
}

func (u *UseCase) IsRootCATrusted() bool {
	return u.manager.IsRootCATrusted()
}

