#!/bin/bash

set -e

cd "$(dirname "$0")/../docker"

# Load environment variables
set -a
source .env
set +a

docker-compose up --build -d

echo "Waiting for PostgreSQL.."
sleep 5

cd ../..
migrate -path infra/migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable" up

echo "PostgreSQL ready and migrations applied!!"
