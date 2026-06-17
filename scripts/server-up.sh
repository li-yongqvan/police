#!/usr/bin/env bash
# Start / repair AI Forum on Ubuntu host with 1Panel on :80 (gateway on :8888).
set -eu
cd /opt/ai-forum

COMPOSE="docker compose -f docker-compose.yml -f docker-compose.server.yml"

echo "==> Building and starting stack (port 8888)..."
$COMPOSE up -d --build

echo "==> Waiting for gateway..."
for i in $(seq 1 30); do
  if curl -fsS -o /dev/null http://127.0.0.1:8888/; then
    echo "OK: http://127.0.0.1:8888/ returns 200"
    curl -fsS http://127.0.0.1:8888/ | grep -o '<title>[^<]*</title>' || true
    exit 0
  fi
  sleep 2
done

echo "FAIL: gateway not healthy on :8888" >&2
$COMPOSE ps
exit 1
