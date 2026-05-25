#!/usr/bin/env bash
set -euo pipefail

DEPLOY_DIR="/opt/ai-forum"
OLD_DIR="/opt/ai-forum.old.$(date +%Y%m%d%H%M%S)"
REPO="https://github.com/li-yongqvan/police.git"
BRANCH="feature/ai-forum-production"

echo "==> Stop old stack and remove volumes for clean demo accounts"
if [ -d "$DEPLOY_DIR" ] && [ -f "$DEPLOY_DIR/docker-compose.yml" ]; then
  (cd "$DEPLOY_DIR" && docker compose down -v --remove-orphans 2>/dev/null) || true
fi
docker volume rm ai-forum_pgdata 2>/dev/null || true

echo "==> Backup and clone fresh code"
if [ -d "$DEPLOY_DIR" ]; then
  mv "$DEPLOY_DIR" "$OLD_DIR"
fi
git clone -b "$BRANCH" --depth 1 "$REPO" "$DEPLOY_DIR"
cd "$DEPLOY_DIR"
chmod +x scripts/*.sh

echo "==> Add frontend Dockerfile (missing in branch)"
cat > frontend/Dockerfile <<'DOCKER'
FROM node:20-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
DOCKER

echo "==> Patch nginx for college frontend API prefixes"
python3 <<'PY'
from pathlib import Path
path = Path("nginx/nginx.conf")
text = path.read_text()
block = """
        location /user-api/ {
            proxy_pass http://user_backend/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        location /forum-api/ {
            proxy_pass http://forum_backend/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        location /admin-api/ {
            proxy_pass http://admin_backend/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

"""
needle = "        # 鈹€鈹€鈹€ User Service 鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€"
if "/user-api/" not in text:
    text = text.replace(needle, block + needle)
path.write_text(text)
PY

echo "==> Create .env"
JWT_SECRET="$(openssl rand -hex 32)"
DB_PASS="$(openssl rand -hex 16)"
cp .env.production.example .env
sed -i "s/replace-with-openssl-rand-hex-32/${JWT_SECRET}/" .env
sed -i "s/replace-with-strong-password/${DB_PASS}/g" .env
echo "DEMO_LOGIN_ALLOWLIST=0.0.0.0/0" >> .env

echo "==> Patch compose for 1Panel port 80 conflict"
sed -i 's/POSTGRES_PASSWORD: ai_forum_dev/POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}/' docker-compose.yml
sed -i 's/- "80:80"/- "8888:80"/' docker-compose.yml

cat > docker-compose.host.yml <<'EOF'
services:
  user-service:
    ports: []
  forum-service:
    ports: []
  admin-service:
    ports: []
  postgres:
    ports: []
  redis:
    ports: []
  frontend:
    ports: []
EOF

echo 'allow all;' > nginx/includes/demo-login-allow.conf

echo "==> Init database schemas"
set -a
source .env
set +a
docker compose -f docker-compose.yml up -d postgres redis
for i in $(seq 1 60); do
  if docker compose exec -T postgres pg_isready -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" >/dev/null 2>&1; then
    break
  fi
  sleep 1
done
docker compose exec -T postgres psql -v ON_ERROR_STOP=1 -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" \
  < migrations/init/000_init_schemas.up.sql

echo "==> Build and start full stack (may take several minutes)"
docker compose -f docker-compose.yml -f docker-compose.host.yml up -d --build

echo "==> Wait for migrations"
sleep 25

bash scripts/seed-content.sh || true

echo "==> Verify"
curl -fsS http://127.0.0.1:8888/ >/dev/null && echo "frontend OK"
curl -fsS http://127.0.0.1:8888/api/v1/boards >/dev/null && echo "boards OK"
curl -fsS -X POST http://127.0.0.1:8888/user-api/api/v1/demo-login \
  -H "Content-Type: application/json" -d '{"role":"student"}' | grep -q token && echo "demo-login OK"
curl -fsS -X POST http://127.0.0.1:8888/user-api/api/v1/login \
  -H "Content-Type: application/json" -d '{"username":"demo_student","password":"demo123456"}' | grep -q token && echo "password login OK"

echo ""
echo "Upgrade complete: http://107.172.138.10/"
echo "Demo accounts: demo_student / demo_admin / demo_platform_admin  password: demo123456"
