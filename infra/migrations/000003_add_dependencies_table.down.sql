DROP INDEX IF EXISTS idx_dependencies_repo_id;
DROP INDEX IF EXISTS idx_dependencies_ecosystem;
DROP INDEX IF EXISTS idx_dependencies_name;

DROP TABLE IF EXISTS dependencies;

DROP TYPE IF EXISTS dep_scope;
DROP TYPE IF EXISTS dep_ecosystem;