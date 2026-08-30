package project

import (
	projectdomain "LocalValet/internal/domain/project"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type LinuxScanner struct{}

func NewScanner() projectdomain.Scanner {
	return &LinuxScanner{}
}

func (s *LinuxScanner) DetectProject(dirPath string) (*projectdomain.Project, error) {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(absPath)
	if err != nil || !info.IsDir() {
		return nil, err
	}

	folderName := filepath.Base(absPath)
	projectName := formatProjectName(folderName)
	domain := formatDomain(folderName)

	framework, webRoot, targetPort := detectFramework(absPath)
	if framework == projectdomain.FrameworkUnknown {
		return nil, nil // Not recognized as a web project
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	return &projectdomain.Project{
		ID:           domainSlug(folderName),
		Name:         projectName,
		Path:         absPath,
		Framework:    string(framework),
		WebRoot:      webRoot,
		Domain:       domain,
		VHostEnabled: true,
		SSLEnabled:   true,
		TargetPort:   targetPort,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (s *LinuxScanner) ScanRoots(roots []string) ([]projectdomain.Project, error) {
	var projects []projectdomain.Project
	seenPaths := make(map[string]bool)

	for _, root := range roots {
		cleanRoot := filepath.Clean(root)
		if cleanRoot == "" {
			continue
		}

		info, err := os.Stat(cleanRoot)
		if err != nil || !info.IsDir() {
			continue
		}

		// 1. Check if the root directory itself is a project
		if proj, err := s.DetectProject(cleanRoot); err == nil && proj != nil {
			if !seenPaths[proj.Path] {
				seenPaths[proj.Path] = true
				projects = append(projects, *proj)
			}
			continue
		}

		// 2. Scan child directories of the root (1 level deep)
		entries, err := os.ReadDir(cleanRoot)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			name := entry.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "runtime" {
				continue
			}

			subPath := filepath.Join(cleanRoot, name)
			if proj, err := s.DetectProject(subPath); err == nil && proj != nil {
				if !seenPaths[proj.Path] {
					seenPaths[proj.Path] = true
					projects = append(projects, *proj)
				}
			}
		}
	}

	return projects, nil
}

func detectFramework(dirPath string) (projectdomain.Framework, string, int) {
	// 1. Check Laravel
	if fileExists(filepath.Join(dirPath, "artisan")) || packageInComposer(dirPath, "laravel/framework") {
		webRoot := "public"
		if !dirExists(filepath.Join(dirPath, "public")) {
			webRoot = "."
		}
		return projectdomain.FrameworkLaravel, webRoot, 0
	}

	// 2. Check WordPress
	if fileExists(filepath.Join(dirPath, "wp-config.php")) || fileExists(filepath.Join(dirPath, "wp-config-sample.php")) || fileExists(filepath.Join(dirPath, "wp-load.php")) {
		return projectdomain.FrameworkWordPress, ".", 0
	}

	// 3. Check Next.js
	if fileExists(filepath.Join(dirPath, "next.config.js")) || fileExists(filepath.Join(dirPath, "next.config.mjs")) || fileExists(filepath.Join(dirPath, "next.config.ts")) || packageInJSON(dirPath, "next") {
		return projectdomain.FrameworkNextJS, ".", 3000
	}

	// 4. Check Nuxt
	if fileExists(filepath.Join(dirPath, "nuxt.config.ts")) || fileExists(filepath.Join(dirPath, "nuxt.config.js")) || packageInJSON(dirPath, "nuxt") {
		return projectdomain.FrameworkNuxt, ".", 3000
	}

	// 5. Check React / Vue / Vite
	if fileExists(filepath.Join(dirPath, "vite.config.ts")) || fileExists(filepath.Join(dirPath, "vite.config.js")) {
		if packageInJSON(dirPath, "vue") {
			return projectdomain.FrameworkVue, resolveDistRoot(dirPath), 5173
		}
		return projectdomain.FrameworkReact, resolveDistRoot(dirPath), 5173
	}

	// 6. Generic PHP
	if fileExists(filepath.Join(dirPath, "index.php")) || fileExists(filepath.Join(dirPath, "composer.json")) {
		if dirExists(filepath.Join(dirPath, "public")) && fileExists(filepath.Join(dirPath, "public", "index.php")) {
			return projectdomain.FrameworkPHP, "public", 0
		}
		return projectdomain.FrameworkPHP, ".", 0
	}

	// 7. Static HTML
	if fileExists(filepath.Join(dirPath, "index.html")) {
		return projectdomain.FrameworkStatic, ".", 0
	}

	return projectdomain.FrameworkUnknown, ".", 0
}

func resolveDistRoot(dirPath string) string {
	if dirExists(filepath.Join(dirPath, "dist")) {
		return "dist"
	}
	if dirExists(filepath.Join(dirPath, "build")) {
		return "build"
	}
	if dirExists(filepath.Join(dirPath, "public")) {
		return "public"
	}
	return "."
}

func packageInComposer(dirPath, pkgName string) bool {
	composerPath := filepath.Join(dirPath, "composer.json")
	data, err := os.ReadFile(composerPath)
	if err != nil {
		return false
	}
	var content map[string]interface{}
	if err := json.Unmarshal(data, &content); err != nil {
		return false
	}

	checkDeps := func(key string) bool {
		if deps, ok := content[key].(map[string]interface{}); ok {
			_, exists := deps[pkgName]
			return exists
		}
		return false
	}

	return checkDeps("require") || checkDeps("require-dev")
}

func packageInJSON(dirPath, pkgName string) bool {
	pkgPath := filepath.Join(dirPath, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return false
	}
	var content map[string]interface{}
	if err := json.Unmarshal(data, &content); err != nil {
		return false
	}

	checkDeps := func(key string) bool {
		if deps, ok := content[key].(map[string]interface{}); ok {
			_, exists := deps[pkgName]
			return exists
		}
		return false
	}

	return checkDeps("dependencies") || checkDeps("devDependencies")
}

func formatProjectName(folderName string) string {
	parts := strings.FieldsFunc(folderName, func(r rune) bool {
		return r == '-' || r == '_'
	})
	if len(parts) == 0 {
		return folderName
	}
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}

func formatDomain(folderName string) string {
	return domainSlug(folderName) + ".test"
}

func domainSlug(name string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	slug := strings.ToLower(reg.ReplaceAllString(name, "-"))
	return strings.Trim(slug, "-")
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
