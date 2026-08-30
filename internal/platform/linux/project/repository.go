package project

import (
	projectdomain "LocalValet/internal/domain/project"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type projectsData struct {
	Roots    []string                `json:"roots"`
	Projects []projectdomain.Project `json:"projects"`
}

type JSONRepository struct {
	mu         sync.RWMutex
	configPath string
}

func NewRepository() projectdomain.Repository {
	baseDir := findBaseDir()
	configPath := filepath.Join(baseDir, "config", "projects.json")
	return &JSONRepository{configPath: configPath}
}

func NewRepositoryWithPath(configPath string) projectdomain.Repository {
	return &JSONRepository{configPath: configPath}
}

func (r *JSONRepository) GetProjects() ([]projectdomain.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := r.load()
	if err != nil {
		return []projectdomain.Project{}, nil
	}
	return data.Projects, nil
}

func (r *JSONRepository) SaveProjects(projects []projectdomain.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, _ := r.load()
	data.Projects = projects
	return r.save(data)
}

func (r *JSONRepository) GetRoots() ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := r.load()
	if err != nil || len(data.Roots) == 0 {
		return defaultRoots(), nil
	}
	return data.Roots, nil
}

func (r *JSONRepository) SaveRoots(roots []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, _ := r.load()
	data.Roots = roots
	return r.save(data)
}

func (r *JSONRepository) load() (projectsData, error) {
	b, err := os.ReadFile(r.configPath)
	if err != nil {
		return projectsData{Roots: defaultRoots(), Projects: []projectdomain.Project{}}, err
	}

	var data projectsData
	if err := json.Unmarshal(b, &data); err != nil {
		return projectsData{Roots: defaultRoots(), Projects: []projectdomain.Project{}}, err
	}
	return data, nil
}

func (r *JSONRepository) save(data projectsData) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(r.configPath), 0o755); err != nil {
		return err
	}

	return os.WriteFile(r.configPath, b, 0o644)
}

func defaultRoots() []string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return []string{"/var/www"}
	}

	candidates := []string{
		filepath.Join(home, "Projects"),
		filepath.Join(home, "Coding"),
		filepath.Join(home, "www"),
		filepath.Join(home, "public_html"),
	}

	var valid []string
	for _, c := range candidates {
		if dirExists(c) {
			valid = append(valid, c)
		}
	}

	if len(valid) == 0 {
		valid = append(valid, filepath.Join(home, "Projects"))
	}
	return valid
}

func findBaseDir() string {
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		if fileExists(filepath.Join(exeDir, "config", "runtime.json")) {
			return exeDir
		}
	}

	cwd, err := os.Getwd()
	if err == nil {
		if fileExists(filepath.Join(cwd, "config", "runtime.json")) {
			return cwd
		}
		parent := filepath.Dir(cwd)
		if fileExists(filepath.Join(parent, "config", "runtime.json")) {
			return parent
		}
		return cwd
	}

	return "."
}
