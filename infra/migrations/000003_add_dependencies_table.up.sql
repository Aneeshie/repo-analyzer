-- Create enum for dependency ecosystems
DO $$ BEGIN
    CREATE TYPE dep_ecosystem AS ENUM (
        'npm',
        'go',
        'pip',
        'cargo',
        'bundler',
        'hex'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE dep_scope AS ENUM (
        'production',
        'development',
        'optional'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

-- Create dependencies table
CREATE TABLE IF NOT EXISTS dependencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id UUID NOT NULL REFERENCES repos(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(100),
    ecosystem dep_ecosystem NOT NULL,
    scope dep_scope,  -- 'production', 'development', 'optional'
    source_file TEXT,   -- Which file declared this (package.json, go.mod)
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(repo_id, name, ecosystem, source_file)
);

-- Create indexes for fast lookups
CREATE INDEX IF NOT EXISTS idx_dependencies_repo_id ON dependencies(repo_id);
CREATE INDEX IF NOT EXISTS idx_dependencies_ecosystem ON dependencies(ecosystem);
CREATE INDEX IF NOT EXISTS idx_dependencies_name ON dependencies(name);