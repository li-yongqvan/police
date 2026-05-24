#!/usr/bin/env bash
# Initialize PostgreSQL schemas (first-time). Run from repo root on Ubuntu/macOS/WSL.
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
COMPOSE_FILE="${COMPOSE_FILE:-infra/docker-compose.yml}"

echo "==> Starting Postgres + Redis (${COMPOSE_FILE})"
docker compose -f "${COMPOSE_FILE}" up -d

echo "==> Waiting for Postgres"
for i in $(seq 1 60); do
  if docker compose -f "${COMPOSE_FILE}" exec -T postgres pg_isready -U "${DB_USER}" -d "${DB_NAME}" >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

echo "==> Applying schema init (idempotent)"
docker compose -f "${COMPOSE_FILE}" exec -T postgres \
  psql -v ON_ERROR_STOP=1 -U "${DB_USER}" -d "${DB_NAME}" \
  < migrations/init/000_init_schemas.up.sql

echo ""
echo "Done. Next steps:"
echo "  1) cp .env.example .env   # set JWT_SECRET for production"
echo "  2) docker compose up -d --build"
echo "  3) bash scripts/seed-content.sh"
