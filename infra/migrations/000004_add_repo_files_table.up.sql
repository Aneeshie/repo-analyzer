CREATE TABLE IF NOT EXISTS repo_files (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id     UUID NOT NULL REFERENCES repos(id) ON DELETE CASCADE,
    parent_path TEXT,
    path        TEXT NOT NULL,
    name        VARCHAR(255) NOT NULL,
    type        VARCHAR(10) NOT NULL CHECK (type IN ('file', 'directory')),
    size        BIGINT DEFAULT 0,
    language    VARCHAR(50),
    UNIQUE(repo_id, path)
);

CREATE INDEX IF NOT EXISTS idx_repo_files_repo_id ON repo_files(repo_id);
CREATE INDEX IF NOT EXISTS idx_repo_files_parent ON repo_files(repo_id, parent_path);
