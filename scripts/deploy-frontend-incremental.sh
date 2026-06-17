#!/usr/bin/env bash
set -eu
cd /opt/ai-forum

COMPOSE_FILES="-f docker-compose.yml"
if [ -f docker-compose.server.yml ]; then
  COMPOSE_FILES="$COMPOSE_FILES -f docker-compose.server.yml"
elif [ -f docker-compose.host.yml ]; then
  COMPOSE_FILES="$COMPOSE_FILES -f docker-compose.host.yml"
fi

echo "==> Rebuild frontend"
docker compose $COMPOSE_FILES build frontend
docker compose $COMPOSE_FILES up -d frontend nginx

echo "==> Verify"
sleep 3
curl -fsS http://127.0.0.1:8888/ >/dev/null && echo "frontend OK"
curl -fsS http://127.0.0.1:8888/ | grep -o 'index-[^"]*\.js' | head -1
