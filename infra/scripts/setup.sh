#!/bin/bash

set -e

cd "$(dirname "$0")/../docker"
docker-compose up -d

<<<<<<< HEAD
echo "Waiting for PostgreSQL..."
=======
echo "Waiting for PostgreSQL.."
>>>>>>> bb2c496 (feat(scripts): add setup and teardown scripts)
sleep 5

cd ../..
migrate -path infra/migrations -database "postgres://repo_user:repo_pass@localhost:5432/repo_analyzer?sslmode=disable" up

<<<<<<< HEAD
echo "PostgreSQL ready and migrations applied!!"
=======
echo "PostgreSQL ready and migrations applied"
>>>>>>> bb2c496 (feat(scripts): add setup and teardown scripts)
