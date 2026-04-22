CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- status enums
DO $$ BEGIN
    CREATE TYPE repos_status AS ENUM (
        'pending',
        'cloning',
        'parsing',
        'analyzing',
        'completed',
        'failed'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS repos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    url VARCHAR(500) UNIQUE NOT NULL,
    status repos_status DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_repos_url ON repos(url);
