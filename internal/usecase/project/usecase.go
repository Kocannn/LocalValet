package project

import (
	projectdomain "LocalValet/internal/domain/project"
	"LocalValet/internal/platform/linux/launcher"
	sslusecase "LocalValet/internal/usecase/ssl"
	vhostusecase "LocalValet/internal/usecase/vhost"
	"fmt"
	"path/filepath"
	"time"
)

type UseCase struct {
	scanner projectdomain.Scanner
	repo    projectdomain.Repository
	vhostUC *vhostusecase.UseCase
	sslUC   *sslusecase.UseCase
}

func New(
	scanner projectdomain.Scanner,
	repo projectdomain.Repository,
	vhostUC *vhostusecase.UseCase,
	sslUC *sslusecase.UseCase,
) *UseCase {
	return &UseCase{
		scanner: scanner,
		repo:    repo,
		vhostUC: vhostUC,
		sslUC:   sslUC,
	}
}

func (u *UseCase) ScanProjects() ([]projectdomain.Project, error) {
	roots, err := u.repo.GetRoots()
	if err != nil {
		return nil, err
	}

	scanned, err := u.scanner.ScanRoots(roots)
	if err != nil {
		return nil, err
	}

	// Auto-provision VHost & SSL for discovered projects
	for i := range scanned {
		p := &scanned[i]
		if u.sslUC != nil {
			_, _ = u.sslUC.CreateCertForProject(p.Domain)
			p.SSLEnabled = true
		}
		if u.vhostUC != nil && p.VHostEnabled {
			_, _ = u.vhostUC.EnableVHost(*p)
		}
	}

	_ = u.repo.SaveProjects(scanned)
	return scanned, nil
}

func (u *UseCase) GetProjects() ([]projectdomain.Project, error) {
	projects, err := u.repo.GetProjects()
	if err != nil || len(projects) == 0 {
		return u.ScanProjects()
	}
	return projects, nil
}

func (u *UseCase) GetProjectRoots() ([]string, error) {
	return u.repo.GetRoots()
}

func (u *UseCase) AddProjectRoot(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	roots, err := u.repo.GetRoots()
	if err != nil {
		roots = []string{}
	}

	for _, r := range roots {
		if r == absPath {
			return nil // already exists
		}
	}

	roots = append(roots, absPath)
	if err := u.repo.SaveRoots(roots); err != nil {
		return err
	}

	_, _ = u.ScanProjects()
	return nil
}

func (u *UseCase) RemoveProjectRoot(path string) error {
	cleanPath := filepath.Clean(path)
	roots, err := u.repo.GetRoots()
	if err != nil {
		return err
	}

	var updated []string
	for _, r := range roots {
		if filepath.Clean(r) != cleanPath {
			updated = append(updated, r)
		}
	}

	if err := u.repo.SaveRoots(updated); err != nil {
		return err
	}

	_, _ = u.ScanProjects()
	return nil
}

func (u *UseCase) ToggleProjectVHost(projectPath string, enable bool) error {
	projects, err := u.repo.GetProjects()
	if err != nil {
		return err
	}

	var found *projectdomain.Project
	for i := range projects {
		if projects[i].Path == projectPath {
			found = &projects[i]
			break
		}
	}

	if found == nil {
		return fmt.Errorf("project not found at %s", projectPath)
	}

	found.VHostEnabled = enable
	found.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	if enable && u.vhostUC != nil {
		_, err = u.vhostUC.EnableVHost(*found)
	} else if !enable && u.vhostUC != nil {
		err = u.vhostUC.DisableVHost(found.Domain)
	}

	if err != nil {
		return err
	}

	return u.repo.SaveProjects(projects)
}

func (u *UseCase) GenerateProjectSSL(projectPath string) error {
	projects, err := u.repo.GetProjects()
	if err != nil {
		return err
	}

	var found *projectdomain.Project
	for i := range projects {
		if projects[i].Path == projectPath {
			found = &projects[i]
			break
		}
	}

	if found == nil {
		return fmt.Errorf("project not found at %s", projectPath)
	}

	if u.sslUC != nil {
		_, err := u.sslUC.CreateCertForProject(found.Domain)
		if err != nil {
			return err
		}
		found.SSLEnabled = true
		found.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

		if found.VHostEnabled && u.vhostUC != nil {
			_, _ = u.vhostUC.EnableVHost(*found)
		}
	}

	return u.repo.SaveProjects(projects)
}

func (u *UseCase) OpenInEditor(projectPath, editor string) error {
	return launcher.OpenInEditor(projectPath, editor)
}

func (u *UseCase) OpenInBrowser(url string) error {
	return launcher.OpenInBrowser(url)
}

func (u *UseCase) GetAllDomains() ([]string, error) {
	projects, err := u.repo.GetProjects()
	if err != nil {
		return nil, err
	}

	var domains []string
	for _, p := range projects {
		if p.Domain != "" {
			domains = append(domains, p.Domain)
		}
	}
	return domains, nil
}

