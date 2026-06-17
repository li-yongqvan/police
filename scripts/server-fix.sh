#!/usr/bin/env bash
set -euo pipefail
cd /opt/ai-forum

python3 <<'PY'
from pathlib import Path
import re
text = Path("docker-compose.yml").read_text()
for port in ["8001:8001", "8002:8002", "8003:8003", "5432:5432", "6379:6379", "3000:80"]:
    text = text.replace(f'      - "{port}"\n', "")
if '- "80:80"' in text:
    text = text.replace('- "80:80"', '- "8888:80"')
Path("docker-compose.yml").write_text(text)
PY

docker compose down --remove-orphans
docker compose -f docker-compose.yml -f docker-compose.host.yml up -d

for i in $(seq 1 60); do
  if docker compose exec -T postgres pg_isready -U ai_forum -d ai_forum >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

sleep 15
docker ps

curl -fsS http://127.0.0.1:8888/api/v1/boards >/dev/null && echo boards_ok
curl -fsS -X POST http://127.0.0.1:8888/user-api/api/v1/demo-login \
  -H 'Content-Type: application/json' \
  -d '{"role":"student"}' | tee /tmp/demo.json | grep -q token && echo demo_login_ok
curl -fsS -X POST http://127.0.0.1:8888/user-api/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"demo_student","password":"demo123456"}' | grep -q token && echo password_login_ok

echo DONE
