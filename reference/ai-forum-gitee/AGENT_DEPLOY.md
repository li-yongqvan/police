# AI Agent Deployment Guide — AI智联论坛

> **Target:** Cloud server AI Agent (autonomous deployment)
> **Server requirements:** Linux x86_64, Docker 24+, Docker Compose 2.20+, Git, 2GB+ RAM, 20GB+ disk

## Step 1: Environment Setup

```bash
# 1.1 Install Docker & Docker Compose if missing
which docker || curl -fsSL https://get.docker.com | sh
which docker-compose || which docker compose || apt install -y docker-compose-plugin

# 1.2 Create project directory
mkdir -p /opt/ai-forum && cd /opt/ai-forum
```

## Step 2: Acquire Project

```bash
# Clone repository
git clone https://gitee.com/hebaonaodai/ai-forum.git .
# Or unpack from archive:
# unzip ai-forum.zip && cd ai-forum
```

## Step 3: Configure Environment

```bash
# 3.1 Generate secrets
JWT_SECRET=$(openssl rand -hex 32)
ADMIN_JWT_SECRET=$(openssl rand -hex 32)
DB_PASSWORD=$(openssl rand -hex 16)

# 3.2 Create .env file
cat > .env << EOF
# PostgreSQL
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=${DB_PASSWORD}
DB_NAME=ai_forum
DB_SSLMODE=disable

# Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=${JWT_SECRET}
ADMIN_JWT_SECRET=${ADMIN_JWT_SECRET}

# Service Ports
USER_SERVICE_PORT=8001
FORUM_SERVICE_PORT=8002
ADMIN_SERVICE_PORT=8003
EOF
```

## Step 4: Deploy with Docker Compose

```bash
# 4.1 Pull images (or build from source)
docker compose build

# 4.2 Start infrastructure (DB + Redis)
docker compose up -d db redis

# 4.3 Wait for PostgreSQL readiness
echo "Waiting for PostgreSQL..."
for i in $(seq 1 30); do
  docker compose exec -T db pg_isready -U postgres >/dev/null 2>&1 && break
  sleep 2
done

# 4.4 Initialize database (creates schemas and tables)
docker compose run --rm user-service /app/migrate
docker compose run --rm forum-service /app/migrate
docker compose run --rm admin-service /app/migrate

# 4.5 Start all services
docker compose up -d

# 4.6 Verify all services healthy
sleep 5
docker compose ps
```

## Step 5: Health Verification

```bash
# Check all service health endpoints
services=(
  "http://localhost:8001/api/v1/health"
  "http://localhost:8002/api/v1/health"
  "http://localhost:8003/api/v1/health"
)
for svc in "${services[@]}"; do
  status=$(curl -s -o /dev/null -w "%{http_code}" "$svc")
  echo "$svc → $status"
  [ "$status" = "200" ] || echo "WARNING: $svc not healthy"
done

# Check Nginx gateway
curl -s -o /dev/null -w "%{http_code}" "http://localhost/api/v1/health"
```

## Step 6: Initial Setup (Admin & Invite Codes)

```bash
# 6.1 Get an invite code (via admin-service internal API)
INVITE_CODE=$(curl -s -X POST "http://localhost:8003/internal/v1/invite-codes" \
  -H "Content-Type: application/json" \
  -d '{"count":1,"created_by":"system"}' | grep -o '"code":"[^"]*"' | cut -d'"' -f4)

echo "Invite code: ${INVITE_CODE}"

# 6.2 Register admin user
curl -s -X POST "http://localhost/api/v1/user/register" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"admin\",\"password\":\"Admin@123456\",\"invite_code\":\"${INVITE_CODE}\"}"

# 6.3 Login to get JWT
ADMIN_TOKEN=$(curl -s -X POST "http://localhost/api/v1/user/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"admin\",\"password\":\"Admin@123456\"}" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

echo "Admin JWT: ${ADMIN_TOKEN:0:20}..."

# 6.4 Upgrade admin to Level 99 (full access)
curl -s -X PUT "http://localhost:8003/internal/v1/users/admin/level" \
  -H "Content-Type: application/json" \
  -d '{"level":99}' > /dev/null

echo "Deployment complete."
```

## Step 7: Monitoring & Maintenance

```bash
# View logs
docker compose logs -f --tail=100 user-service
docker compose logs -f --tail=100 forum-service
docker compose logs -f --tail=100 admin-service

# Restart a service
docker compose restart user-service

# Full restart
docker compose down && docker compose up -d

# Backup database
docker compose exec db pg_dump -U postgres ai_forum > backup_$(date +%Y%m%d).sql
```

## Troubleshooting (for AI Agent)

| Symptom | Check | Fix |
|---------|-------|-----|
| Service crashes on start | `docker compose logs <service>` | Verify DB migrations ran, check .env values |
| DB connection refused | `docker compose logs db` | Wait for PG init, verify DB_HOST=db |
| 401 on API calls | Token expired/not found | Re-login via `/api/v1/user/login` |
| Nginx 502 bad gateway | `docker compose ps` | Target service not running, check port config |
| Migration failed | `docker compose run --rm <service> /app/migrate` | Table already exists (safe to re-run) |

## Service Ports (internal)

| Service | Internal Port | Nginx Route |
|---------|--------------|-------------|
| user-service | 8001 | /api/v1/user/* |
| forum-service | 8002 | /api/v1/forum/* |
| admin-service | 8003 | /api/v1/admin/* |
| PostgreSQL | 5432 | (internal only) |
| Redis | 6379 | (internal only) |
| Nginx | 80 | (public gateway) |
