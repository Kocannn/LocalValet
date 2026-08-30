package project

import (
	projectdomain "LocalValet/internal/domain/project"
	"os"
	"path/filepath"
	"testing"
)

func TestScanner_DetectFrameworks(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Setup Laravel project
	laravelDir := filepath.Join(tempDir, "my-laravel-app")
	_ = os.MkdirAll(filepath.Join(laravelDir, "public"), 0o755)
	_ = os.WriteFile(filepath.Join(laravelDir, "artisan"), []byte("#!/usr/bin/env php\n"), 0o755)
	_ = os.WriteFile(filepath.Join(laravelDir, "public", "index.php"), []byte("<?php\n"), 0o644)

	// 2. Setup WordPress project
	wpDir := filepath.Join(tempDir, "my-wp-blog")
	_ = os.MkdirAll(wpDir, 0o755)
	_ = os.WriteFile(filepath.Join(wpDir, "wp-config.php"), []byte("<?php\n"), 0o644)

	// 3. Setup Next.js project
	nextDir := filepath.Join(tempDir, "my-next-app")
	_ = os.MkdirAll(nextDir, 0o755)
	_ = os.WriteFile(filepath.Join(nextDir, "next.config.js"), []byte("module.exports = {};\n"), 0o644)

	// 4. Setup Vite React project
	reactDir := filepath.Join(tempDir, "my-react-app")
	_ = os.MkdirAll(filepath.Join(reactDir, "dist"), 0o755)
	_ = os.WriteFile(filepath.Join(reactDir, "vite.config.ts"), []byte("export default {};\n"), 0o644)
	_ = os.WriteFile(filepath.Join(reactDir, "package.json"), []byte(`{"dependencies":{"react":"^18.0.0"}}`), 0o644)

	// 5. Setup Static project
	staticDir := filepath.Join(tempDir, "my-static-site")
	_ = os.MkdirAll(staticDir, 0o755)
	_ = os.WriteFile(filepath.Join(staticDir, "index.html"), []byte("<h1>Hello</h1>"), 0o644)

	scanner := NewScanner()

	// Test Laravel
	p1, err := scanner.DetectProject(laravelDir)
	if err != nil || p1 == nil {
		t.Fatalf("failed to detect Laravel: %v", err)
	}
	if p1.Framework != string(projectdomain.FrameworkLaravel) {
		t.Errorf("expected framework laravel, got %s", p1.Framework)
	}
	if p1.WebRoot != "public" {
		t.Errorf("expected webroot public, got %s", p1.WebRoot)
	}
	if p1.Domain != "my-laravel-app.test" {
		t.Errorf("expected domain my-laravel-app.test, got %s", p1.Domain)
	}

	// Test WordPress
	p2, err := scanner.DetectProject(wpDir)
	if err != nil || p2 == nil {
		t.Fatalf("failed to detect WordPress: %v", err)
	}
	if p2.Framework != string(projectdomain.FrameworkWordPress) {
		t.Errorf("expected framework wordpress, got %s", p2.Framework)
	}

	// Test Next.js
	p3, err := scanner.DetectProject(nextDir)
	if err != nil || p3 == nil {
		t.Fatalf("failed to detect Next.js: %v", err)
	}
	if p3.Framework != string(projectdomain.FrameworkNextJS) {
		t.Errorf("expected framework nextjs, got %s", p3.Framework)
	}
	if p3.TargetPort != 3000 {
		t.Errorf("expected targetPort 3000, got %d", p3.TargetPort)
	}

	// Test React
	p4, err := scanner.DetectProject(reactDir)
	if err != nil || p4 == nil {
		t.Fatalf("failed to detect React: %v", err)
	}
	if p4.Framework != string(projectdomain.FrameworkReact) {
		t.Errorf("expected framework react, got %s", p4.Framework)
	}

	// Test Static
	p5, err := scanner.DetectProject(staticDir)
	if err != nil || p5 == nil {
		t.Fatalf("failed to detect Static: %v", err)
	}
	if p5.Framework != string(projectdomain.FrameworkStatic) {
		t.Errorf("expected framework static, got %s", p5.Framework)
	}

	// Test ScanRoots
	projects, err := scanner.ScanRoots([]string{tempDir})
	if err != nil {
		t.Fatalf("ScanRoots error: %v", err)
	}
	if len(projects) != 5 {
		t.Errorf("expected 5 projects discovered, got %d", len(projects))
	}
}
