CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ BEGIN
    CREATE TYPE llm_model AS ENUM (
        'gemini',
        'openAI',
        'anthropic'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS repo_details (
    repo_id UUID PRIMARY KEY REFERENCES repos(id) ON DELETE CASCADE,
    name VARCHAR(255),
    owner VARCHAR(255),
    description TEXT,
    stars INT DEFAULT 0,
    forks INT DEFAULT 0,
    watchers INT DEFAULT 0,
    open_issues INT DEFAULT 0,
    language VARCHAR(100),
    topics TEXT[],
    license VARCHAR(100),
    last_github_fetch TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS analysis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id UUID NOT NULL REFERENCES repos(id) ON DELETE CASCADE,
    summary TEXT NOT NULL,
    local_setup_instructions TEXT,
    how_it_works TEXT,
    tech_stack TEXT[],
    llm_model llm_model DEFAULT 'gemini',
    generated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_analysis_repo_id ON analysis(repo_id);
CREATE INDEX IF NOT EXISTS idx_repo_details_language ON repo_details(language);
