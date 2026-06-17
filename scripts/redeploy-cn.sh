#!/usr/bin/env bash
set -euo pipefail
cd /home/liyongquan/projects/ai-forum
export PATH="$HOME/.docker/cli-plugins:$PATH"
CF="-f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml"
docker compose $CF pull nginx postgres redis || true
docker compose $CF up -d --build
sleep 40
docker compose $CF restart nginx || true
sleep 10
[ -f scripts/seed-content.sh ] && bash scripts/seed-content.sh || true
curl -fsS -o /dev/null -w "HTTP8888=%{http_code}\n" http://127.0.0.1:8888/ || true
docker compose $CF ps
echo REDEPLOY_DONE
