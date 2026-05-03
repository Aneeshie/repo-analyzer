ALTER TABLE repos ADD COLUMN entry_points JSONB DEFAULT '[]'::jsonb;

CREATE TABLE IF NOT EXISTS file_explanations (
    repo_id UUID NOT NULL REFERENCES repos(id) ON DELETE CASCADE,
    path TEXT NOT NULL,
    explanation TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (repo_id, path)
);
