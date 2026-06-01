#!/usr/bin/env bash
set -eu
cd /opt/ai-forum

COMPOSE_FILES="-f docker-compose.yml"
if [ -f docker-compose.server.yml ]; then
  COMPOSE_FILES="$COMPOSE_FILES -f docker-compose.server.yml"
elif [ -f docker-compose.host.yml ]; then
  COMPOSE_FILES="$COMPOSE_FILES -f docker-compose.host.yml"
fi

echo "==> Build user, forum, admin, frontend"
docker compose $COMPOSE_FILES build user-service forum-service admin-service frontend
docker compose $COMPOSE_FILES up -d user-service forum-service admin-service frontend nginx

echo "==> Reload nginx (refresh upstream DNS after container recreate)"
sleep 5
docker compose $COMPOSE_FILES restart nginx
sleep 5

echo "==> Verify"
curl -fsS http://127.0.0.1:8888/ >/dev/null && echo frontend OK
curl -fsS http://127.0.0.1:8888/forum-api/api/v1/stats/community | head -c 80 && echo
curl -fsS -X POST http://127.0.0.1:8888/user-api/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"demo_student","password":"demo123456"}' | grep -q access_token && echo login OK
