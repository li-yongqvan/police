#!/usr/bin/env bash
# First-time deploy on Ubuntu cloud VM (single host, Docker Compose).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if [ ! -f .env ]; then
  cp .env.production.example .env
  echo "Created .env from .env.production.example — edit JWT_SECRET before public exposure."
fi

if grep -q "change-me-in-production" .env 2>/dev/null; then
  echo "WARNING: JWT_SECRET still default. Run: openssl rand -hex 32"
fi

bash scripts/init-db.sh

docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

echo "Waiting for service migrations..."
sleep 15

bash scripts/seed-content.sh
bash scripts/smoke-test.sh

echo ""
echo "Deployed. Open http://<your-server-ip>/ (port 80)."
echo "demo-login: only private networks unless DEMO_LOGIN_ALLOWLIST is set."
