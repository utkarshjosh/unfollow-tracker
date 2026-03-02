#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

export DATABASE_URL="${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/unfollow_tracker?sslmode=disable}"

echo "Using DATABASE_URL=$DATABASE_URL"
echo "Checking migration status..."

status_output="$(go run ./cmd/migrator status)"
echo "$status_output"

if echo "$status_output" | grep -q "Pending"; then
  echo "Pending migrations found. Applying..."
  go run ./cmd/migrator up
  echo "Done."
else
  echo "No pending migrations."
fi
