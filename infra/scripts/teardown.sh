#!/bin/bash

cd "$(dirname "$0")/../docker"
docker-compose down -v

echo "PostgreSQL stopped and cleaned up"
