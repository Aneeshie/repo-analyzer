package models

import "time"

type DependencyEcosystem string

const (
	EcosystemNPM     DependencyEcosystem = "npm"
	EcosystemGo      DependencyEcosystem = "go"
	EcosystemPip     DependencyEcosystem = "pip"
	EcosystemCargo   DependencyEcosystem = "cargo"
	EcosystemBundler DependencyEcosystem = "bundler"
	EcosystemHex     DependencyEcosystem = "hex"
)

type DependencyScope string

const (
	ScopeProduction  DependencyScope = "production"
	ScopeDevelopment DependencyScope = "development"
	ScopeOptional    DependencyScope = "optional"
)

type Dependency struct {
	ID         string              `json:"id" db:"id"`
	RepoID     string              `json:"repo_id" db:"repo_id"`
	Name       string              `json:"name" db:"name"`
	Version    string              `json:"version" db:"version"`
	Ecosystem  DependencyEcosystem `json:"ecosystem" db:"ecosystem"`
	Scope      DependencyScope     `json:"scope" db:"scope"`
	SourceFile string              `json:"source_file" db:"source_file"`
	CreatedAt  time.Time           `json:"created_at" db:"created_at"`
}
