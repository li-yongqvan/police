#!/usr/bin/env bash
# Seed sample posts/comments after all services have run migrations once.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if [ -f .env ]; then
  set -a
  # shellcheck disable=SC1091
  source .env
  set +a
fi

DB_USER="${POSTGRES_USER:-ai_forum}"
DB_NAME="${POSTGRES_DB:-ai_forum}"

wait_for_table() {
  local table=$1
  for i in $(seq 1 90); do
    if docker compose exec -T postgres psql -U "${DB_USER}" -d "${DB_NAME}" -tAc \
      "SELECT to_regclass('${table}') IS NOT NULL;" 2>/dev/null | grep -q t; then
      return 0
    fi
    sleep 2
  done
  echo "Timeout waiting for table ${table}. Start services first: docker compose up -d --build"
  exit 1
}

echo "==> Waiting for forum posts table"
wait_for_table "schema_forum.posts"

echo "==> Seeding content"
docker compose exec -T postgres \
  psql -v ON_ERROR_STOP=1 -U "${DB_USER}" -d "${DB_NAME}" \
  < scripts/seed/001_content.sql

echo "Seed complete."
