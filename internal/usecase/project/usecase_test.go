package project

import (
	projectdomain "LocalValet/internal/domain/project"
	projectplatform "LocalValet/internal/platform/linux/project"
	sslplatform "LocalValet/internal/platform/linux/ssl"
	vhostplatform "LocalValet/internal/platform/linux/vhost"
	sslusecase "LocalValet/internal/usecase/ssl"
	vhostusecase "LocalValet/internal/usecase/vhost"
	"os"
	"path/filepath"
	"testing"
)

func TestProjectUseCase_Workflow(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Setup mock workspace
	projectsRoot := filepath.Join(tempDir, "Projects")
	_ = os.MkdirAll(filepath.Join(projectsRoot, "demo-app", "public"), 0o755)
	_ = os.WriteFile(filepath.Join(projectsRoot, "demo-app", "artisan"), []byte("#!/usr/bin/env php\n"), 0o755)

	configPath := filepath.Join(tempDir, "projects.json")
	certsDir := filepath.Join(tempDir, "certs")
	vhostsDir := filepath.Join(tempDir, "vhosts")

	scanner := projectplatform.NewScanner()
	repo := projectplatform.NewRepositoryWithPath(configPath)
	_ = repo.SaveRoots([]string{projectsRoot})

	sslManager := sslplatform.NewCAManagerWithPath(certsDir)
	sslUC := sslusecase.New(sslManager)
	vhostGen := vhostplatform.NewNginxGeneratorWithPath(vhostsDir)
	vhostUC := vhostusecase.New(vhostGen, sslUC)

	uc := New(scanner, repo, vhostUC, sslUC)

	roots, _ := uc.GetProjectRoots()
	if len(roots) != 1 || roots[0] != projectsRoot {
		t.Errorf("expected 1 root %s, got %v", projectsRoot, roots)
	}


	// Scan projects
	projects, err := uc.ScanProjects()
	if err != nil {
		t.Fatalf("ScanProjects error: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}

	p := projects[0]
	if p.Framework != string(projectdomain.FrameworkLaravel) {
		t.Errorf("expected laravel, got %s", p.Framework)
	}
	if p.Domain != "demo-app.test" {
		t.Errorf("expected demo-app.test, got %s", p.Domain)
	}
	if !p.SSLEnabled {
		t.Errorf("expected SSL enabled")
	}

	// Verify vhost file created
	vhostPath := filepath.Join(vhostsDir, "demo-app.test.conf")
	if _, err := os.Stat(vhostPath); err != nil {
		t.Errorf("expected vhost file to exist at %s", vhostPath)
	}

	// Toggle VHost
	if err := uc.ToggleProjectVHost(p.Path, false); err != nil {
		t.Errorf("ToggleProjectVHost error: %v", err)
	}
	if _, err := os.Stat(vhostPath); err == nil {
		t.Errorf("expected vhost file to be removed when disabled")
	}
}
