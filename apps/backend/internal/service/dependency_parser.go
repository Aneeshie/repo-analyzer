package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/Aneeshie/repo-analyzer/backend/pkg/models"
)

type DependencyParser struct{}

func NewDependencyParser() *DependencyParser {
	return &DependencyParser{}
}

func (p *DependencyParser) ParseRepo(repoPath, repoID string) ([]models.Dependency, error) {

	var allDeps []models.Dependency

	//walk through the repo

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			//skip .git
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		//check file name and parse accordingly
		switch filepath.Base(path) {
		case "package.json":
			deps, err := p.parsePackageJSON(path, repoID)
			if err == nil {
				allDeps = append(allDeps, deps...)
			}
		case "go.mod":
			deps, err := p.parseGoMod(path, repoID)
			if err == nil {
				allDeps = append(allDeps, deps...)
			}
		case "requirements.txt":
			deps, err := p.parseRequirementsTxt(path, repoID)
			if err == nil {
				allDeps = append(allDeps, deps...)
			}
		}
		return nil

	})

	return allDeps, err
}

func (p *DependencyParser) parsePackageJSON(path, repoID string) ([]models.Dependency, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pkgJSON struct {
		Dependencies         map[string]string `json:"dependencies"`
		DevDependencies      map[string]string `json:"devDependencies"`
		OptionalDependencies map[string]string `json:"optionalDependencies"`
	}
	if err := json.Unmarshal(file, &pkgJSON); err != nil {
		return nil, err
	}

	var deps []models.Dependency
	for name, version := range pkgJSON.Dependencies {
		deps = append(deps, models.Dependency{
			RepoID:     repoID,
			Name:       name,
			Version:    version,
			Ecosystem:  models.EcosystemNPM,
			Scope:      models.ScopeProduction,
			SourceFile: "package.json",
		})
	}
	for name, version := range pkgJSON.DevDependencies {
		deps = append(deps, models.Dependency{
			RepoID:     repoID,
			Name:       name,
			Version:    version,
			Ecosystem:  models.EcosystemNPM,
			Scope:      models.ScopeDevelopment,
			SourceFile: "package.json",
		})
	}
	for name, version := range pkgJSON.OptionalDependencies {
		deps = append(deps, models.Dependency{
			RepoID:     repoID,
			Name:       name,
			Version:    version,
			Ecosystem:  models.EcosystemNPM,
			Scope:      models.ScopeOptional,
			SourceFile: "package.json",
		})
	}
	return deps, nil
}

func (p *DependencyParser) parseGoMod(path, repoID string) ([]models.Dependency, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var deps []models.Dependency

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "require") && !strings.Contains(line, "(") {
			parts := strings.Fields(line)
			if len(parts) >= 3 && parts[0] == "require" {
				deps = append(deps, models.Dependency{
					RepoID:     repoID,
					Name:       parts[1],
					Version:    parts[2],
					Ecosystem:  models.EcosystemGo,
					Scope:      models.ScopeProduction,
					SourceFile: "go.mod",
				})
			}
		}
	}
	return deps, nil
}

func (p *DependencyParser) parseRequirementsTxt(path, repoID string) ([]models.Dependency, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var deps []models.Dependency

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "==")
		name := parts[0]
		version := ""
		if len(parts) > 1 {
			version = strings.Split(parts[1], " ")[0]
		}
		deps = append(deps, models.Dependency{
			RepoID:     repoID,
			Name:       name,
			Version:    version,
			Ecosystem:  models.EcosystemPip,
			Scope:      models.ScopeProduction,
			SourceFile: "requirements.txt",
		})
	}
	return deps, nil
}
