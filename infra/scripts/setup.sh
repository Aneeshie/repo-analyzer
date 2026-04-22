#!/bin/bash

set -e

cd "$(dirname "$0")/../docker"
docker-compose up -d

echo "Waiting for PostgreSQL.."
sleep 5

cd ../..
migrate -path infra/migrations -database "postgres://repo_user:repo_pass@localhost:5432/repo_analyzer?sslmode=disable" up

echo "PostgreSQL ready and migrations applied!!"
