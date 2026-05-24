---
wave: 1
depends_on: []
files_modified:
  - docker-compose.yml
  - .env.example
  - nginx/nginx.conf
  - services/user/**
  - services/forum/**
  - services/admin/**
  - frontend/**
  - migrations/**
autonomous: true
requirements: [INFRA-01, INFRA-02, INFRA-03, FE-04]
---

# Phase 1 Plan: 基础设施与项目骨架

**Goal:** 搭建完整的开发基础设施，包括Docker Compose编排、数据库schema、Nginx网关、以及三个Go微服务的项目骨架，确保服务间可以互相通信。

**Mode:** standard

## must_haves

- Docker Compose一键启动所有服务（PostgreSQL 16, Redis 7, Nginx, 3 Go服务, 前端）
- 三个Go服务可通过各自端口响应健康检查
- 通过Nginx网关可访问各服务API
- PostgreSQL三个独立schema已创建（schema_auth, schema_forum, schema_admin）
- golang-migrate集成到各服务启动流程

## Wave 1: 基础设施编排

### Task 1.1: Create docker-compose.yml

<objective>
Create the Docker Compose orchestration file for all services.
</objective>

<read_first>
- AI智联论坛-需求文档+开发边界+MVP范围+微服务架构设计.md (section 6.1 for Docker Compose template)
</read_first>

<action>
Create `docker-compose.yml` at project root with the following services:

1. **nginx** — `nginx:alpine`, port `80:80`, depends on all 3 services + frontend, mounts `./nginx/nginx.conf:/etc/nginx/nginx.conf`
2. **user-service** — build from `./services/user`, port `8001:8001`, depends on postgres + redis, env vars: `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_SCHEMA`, `DB_USER`, `DB_PASSWORD`, `REDIS_HOST`, `REDIS_PORT`, `JWT_SECRET`, `JWT_EXPIRY`
3. **forum-service** — build from `./services/forum`, port `8002:8002`, depends on postgres + redis, env vars: `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_SCHEMA`, `DB_USER`, `DB_PASSWORD`, `REDIS_HOST`, `REDIS_PORT`, `ADMIN_SERVICE_URL`
4. **admin-service** — build from `./services/admin`, port `8003:8003`, depends on postgres + redis, env vars: `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_SCHEMA`, `DB_USER`, `DB_PASSWORD`, `REDIS_HOST`, `REDIS_PORT`
5. **postgres** — `postgres:16`, env `POSTGRES_DB=ai_forum`, `POSTGRES_USER=ai_forum`, `POSTGRES_PASSWORD=ai_forum_dev`, volume `pgdata:/var/lib/postgresql/data`
6. **redis** — `redis:7-alpine`, port `6379`
7. **frontend** — build from `./frontend`, port `3000:80` (or dev server port)

Add shared volume `pgdata` at the bottom.

All services on the same default Docker Compose network.
</action>

<acceptance_criteria>
- docker-compose.yml exists at project root
- Contains 7 services: nginx, user-service, forum-service, admin-service, postgres, redis, frontend
- postgres uses image postgres:16
- redis uses image redis:7-alpine
- user-service exposes port 8001
- forum-service exposes port 8002
- admin-service exposes port 8003
- nginx exposes port 80
- frontend exposes port 3000
- Volume pgdata defined for postgres
- All backend services depend_on postgres and redis
</acceptance_criteria>

---

### Task 1.2: Create .env.example

<objective>
Create environment variable template for all services.
</objective>

<read_first>
- docker-compose.yml (task 1.1 output)
</read_first>

<action>
Create `.env.example` at project root with these variables:

```
# PostgreSQL
POSTGRES_DB=ai_forum
POSTGRES_USER=ai_forum
POSTGRES_PASSWORD=ai_forum_dev

# Database connection (services)
DB_HOST=postgres
DB_PORT=5432
DB_NAME=ai_forum
DB_USER=ai_forum
DB_PASSWORD=ai_forum_dev

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# JWT
JWT_SECRET=change-me-in-production
JWT_EXPIRY=30m

# Service URLs
ADMIN_SERVICE_URL=http://admin-service:8003
USER_SERVICE_URL=http://user-service:8001
FORUM_SERVICE_URL=http://forum-service:8002

# Frontend
VITE_API_BASE_URL=/api/v1
```
</action>

<acceptance_criteria>
- .env.example exists at project root
- Contains POSTGRES_DB=ai_forum
- Contains JWT_EXPIRY=30m
- Contains REDIS_HOST=redis
- Contains ADMIN_SERVICE_URL=http://admin-service:8003
- Contains VITE_API_BASE_URL=/api/v1
- Contains DB_SCHEMA comment noting each service uses its own schema
</acceptance_criteria>

---

### Task 1.3: Create Nginx reverse proxy configuration

<objective>
Create Nginx config routing frontend and API requests to correct services.
</objective>

<read_first>
- AI智联论坛-需求文档+开发边界+MVP范围+微服务架构设计.md (section 6.1, port assignments)
- docker-compose.yml
</read_first>

<action>
Create `nginx/nginx.conf` with:

1. Upstream blocks:
   - `user_backend` → `user-service:8001`
   - `forum_backend` → `forum-service:8002`
   - `admin_backend` → `admin-service:8003`
   - `frontend_backend` → `frontend:80` (or 3000)

2. Location blocks:
   - `/` → proxy to `frontend_backend` (serve Vue SPA)
   - `/api/v1/auth/` → proxy to `user_backend` (strip `/api/v1/auth/` prefix or pass through)
   - `/api/v1/posts`, `/api/v1/boards`, `/api/v1/comments`, `/api/v1/attachments` → proxy to `forum_backend`
   - `/api/v1/admin/` → proxy to `admin_backend`
   - `/internal/v1/` → internal only, not routed from external

3. Add headers: `proxy_set_header Host $host`, `proxy_set_header X-Real-IP $remote_addr`, `proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for`

4. CORS headers for development
</action>

<acceptance_criteria>
- nginx/nginx.conf exists
- Contains upstream block for user_backend pointing to user-service:8001
- Contains upstream block for forum_backend pointing to forum-service:8002
- Contains upstream block for admin_backend pointing to admin-service:8003
- Contains location / proxying to frontend
- Contains location /api/v1/auth/ proxying to user_backend
- Contains location /api/v1/posts proxying to forum_backend
- Contains location /api/v1/admin/ proxying to admin_backend
- Contains proxy_set_header directives
- Contains CORS headers (Access-Control-Allow-Origin)
</acceptance_criteria>

---

## Wave 2: Database schema setup (parallel with Wave 3)

### Task 2.1: Create PostgreSQL schema migration init script

<objective>
Create SQL migration to initialize three independent PostgreSQL schemas.
</objective>

<read_first>
- .planning/phases/01-infrastructure/01-CONTEXT.md (D-06, D-07, D-08 decisions)
- AI智联论坛-需求文档+开发边界+MVP范围+微服务架构设计.md (section 5 for schema names)
</read_first>

<action>
Create `migrations/init/000_init_schemas.up.sql` with:

```sql
-- Create schema_auth for user-service
CREATE SCHEMA IF NOT EXISTS schema_auth;

-- Create schema_forum for forum-service
CREATE SCHEMA IF NOT EXISTS schema_forum;

-- Create schema_admin for admin-service
CREATE SCHEMA IF NOT EXISTS schema_admin;
```

Create corresponding `migrations/init/000_init_schemas.down.sql`:

```sql
DROP SCHEMA IF EXISTS schema_admin CASCADE;
DROP SCHEMA IF EXISTS schema_forum CASCADE;
DROP SCHEMA IF EXISTS schema_auth CASCADE;
```

Note: Each service will have its own migrations/ directory for service-specific tables (created in later tasks).
</action>

<acceptance_criteria>
- migrations/init/000_init_schemas.up.sql exists
- Contains CREATE SCHEMA IF NOT EXISTS schema_auth
- Contains CREATE SCHEMA IF NOT EXISTS schema_forum
- Contains CREATE SCHEMA IF NOT EXISTS schema_admin
- migrations/init/000_init_schemas.down.sql exists
- Contains DROP SCHEMA for all three schemas with CASCADE
</acceptance_criteria>

---

## Wave 3: Go service skeletons (all 3 parallel)

### Task 3.1: Create user-service skeleton

<objective>
Create the user-service Go + Gin project skeleton with health check endpoint.
</objective>

<read_first>
- AI智联论坛-需求文档+开发边界+MVP范围+微服务架构设计.md (section 3.1, 7 for directory structure and routes)
- .planning/phases/01-infrastructure/01-CONTEXT.md (D-01 to D-10 decisions)
</read_first>

<action>
Create the following directory structure and files under `services/user/`:

**go.mod:**
```
module ai-forum/user-service
go 1.22
require (
    github.com/gin-gonic/gin v1.9.1
    github.com/jackc/pgx/v5 v5.5.1
    github.com/redis/go-redis/v9 v9.3.1
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/joho/godotenv v1.5.1
    github.com/golang-migrate/migrate/v4 v4.17.0
)
```

**cmd/main.go** — Entry point:
- Load .env via godotenv
- Connect to PostgreSQL (env: DB_HOST, DB_PORT, DB_NAME, DB_SCHEMA=schema_auth, DB_USER, DB_PASSWORD)
- Connect to Redis (env: REDIS_HOST, REDIS_PORT)
- Run migrations from `./migrations/` using golang-migrate
- Initialize Gin router
- Register routes
- Start server on `:8001`

**internal/handler/health.go** — Health check handler:
- `GET /health` returns `{"status": "ok", "service": "user-service"}`

**internal/handler/user_handler.go** — Stub handlers (empty returns with TODO comments):
- Register, Login, GetProfile, UpdateProfile, UploadAvatar

**internal/service/user_service.go** — Stub service layer (empty struct with methods)

**internal/model/user.go** — User model struct:
```go
type User struct {
    ID        uint      `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Avatar    string    `json:"avatar"`
    Level     int       `json:"level"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}
```

**internal/middleware/auth.go** — JWT middleware stub (function that extracts and validates token, returns 401 on failure)

**internal/middleware/cors.go** — CORS middleware stub

**pkg/database/db.go** — PostgreSQL connection helper (returns *pgxpool.Pool from env vars)

**pkg/redis/redis.go** — Redis connection helper (returns *redis.Client from env vars)

**pkg/jwt/jwt.go** — JWT helper (GenerateToken, ValidateToken, with 30min expiry)

**migrations/000_init_schemas.up.sql** — Symlink or copy of root migrations/init/000_init_schemas.up.sql (or reference)
</action>

<acceptance_criteria>
- services/user/go.mod exists with module ai-forum/user-service
- services/user/cmd/main.go exists
- cmd/main.go contains gin.Default() or gin.New()
- cmd/main.go contains ":8001" for server port
- services/user/internal/handler/health.go exists
- health.go contains GET /health returning JSON with "status": "ok"
- services/user/internal/model/user.go exists with User struct
- services/user/internal/middleware/auth.go exists
- services/user/pkg/database/db.go exists
- services/user/pkg/redis/redis.go exists
- services/user/pkg/jwt/jwt.go exists
- go.mod contains github.com/gin-gonic/gin dependency
- go.mod contains github.com/golang-migrate/migrate/v4 dependency
</acceptance_criteria>

---

### Task 3.2: Create forum-service skeleton

<objective>
Create the forum-service Go + Gin project skeleton with health check endpoint.
</objective>

<read_first>
- AI智联论坛-需求文档+开发边界+MVP范围+微服务架构设计.md (section 3.2, 7 for directory structure and routes)
- .planning/phases/01-infrastructure/01-CONTEXT.md
</read_first>

<action>
Create the following under `services/forum/` (same structure as user-service):

**go.mod:**
```
module ai-forum/forum-service
go 1.22
require (
    github.com/gin-gonic/gin v1.9.1
    github.com/jackc/pgx/v5 v5.5.1
    github.com/redis/go-redis/v9 v9.3.1
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/joho/godotenv v1.5.1
    github.com/golang-migrate/migrate/v4 v4.17.0
)
```

**cmd/main.go:**
- Same pattern as user-service but schema = `schema_forum`, port = `:8002`
- Includes `ADMIN_SERVICE_URL` env var for internal calls to admin-service

**internal/handler/health.go** — `GET /health` returns `{"status": "ok", "service": "forum-service"}`

**internal/handler/post_handler.go** — Stub: ListPosts, GetPost, CreatePost, UpdatePost, DeletePost
**internal/handler/board_handler.go** — Stub: ListBoards, GetBoard
**internal/handler/comment_handler.go** — Stub: ListComments, CreateComment
**internal/handler/interaction_handler.go** — Stub: LikePost, CollectPost
**internal/handler/attachment_handler.go** — Stub: UploadAttachment, DownloadAttachment

**internal/service/forum_service.go** — Stub with methods for post/board/comment operations

**internal/model/post.go** — Post, Board, Comment, Like, Collection, Attachment structs

**internal/middleware/auth.go** — JWT validation middleware

**pkg/database/db.go**, **pkg/redis/redis.go** — Same as user-service

**internal/client/admin_client.go** — HTTP client stub for calling admin-service at `ADMIN_SERVICE_URL`:
- Method `CheckSensitiveWords(text string) (bool, error)` — calls `POST /internal/v1/moderation/check`
</action>

<acceptance_criteria>
- services/forum/go.mod exists with module ai-forum/forum-service
- services/forum/cmd/main.go exists with port ":8002"
- services/forum/internal/handler/health.go exists returning "forum-service"
- services/forum/internal/model/post.go exists with Post struct
- services/forum/internal/client/admin_client.go exists
- admin_client.go contains ADMIN_SERVICE_URL reference
- admin_client.go has method CheckSensitiveWords or similar
- services/forum/internal/handler/post_handler.go exists
- services/forum/internal/handler/board_handler.go exists
- services/forum/internal/handler/comment_handler.go exists
- go.mod contains github.com/gin-gonic/gin
- go.mod contains github.com/golang-migrate/migrate/v4
</acceptance_criteria>

---

### Task 3.3: Create admin-service skeleton

<objective>
Create the admin-service Go + Gin project skeleton with health check endpoint.
</objective>

<read_first>
- AI智联论坛-需求文档+开发边界+MVP范围+微服务架构设计.md (section 3.3, 7 for directory structure and routes)
- .planning/phases/01-infrastructure/01-CONTEXT.md
</read_first>

<action>
Create the following under `services/admin/`:

**go.mod:**
```
module ai-forum/admin-service
go 1.22
require (
    github.com/gin-gonic/gin v1.9.1
    github.com/jackc/pgx/v5 v5.5.1
    github.com/redis/go-redis/v9 v9.3.1
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/joho/godotenv v1.5.1
    github.com/golang-migrate/migrate/v4 v4.17.0
)
```

**cmd/main.go:**
- Same pattern but schema = `schema_admin`, port = `:8003`
- Includes `USER_SERVICE_URL` and `FORUM_SERVICE_URL` env vars

**internal/handler/health.go** — `GET /health` returns `{"status": "ok", "service": "admin-service"}`

**internal/handler/config_handler.go** — Stub: GetConfig, UpdateConfig
**internal/handler/audit_handler.go** — Stub: ListPendingAudit, ApprovePost, RejectPost
**internal/handler/user_admin_handler.go** — Stub: BanUser, ListUsers, UpdateUserLevel, GetUserLogs
**internal/handler/stats_handler.go** — Stub: GetOverview, GetDailyStats
**internal/handler/moderation_handler.go** — Stub: AddSensitiveWord, CheckSensitiveWords (internal endpoint)

**internal/service/admin_service.go** — Stub with admin methods

**internal/model/admin.go** — SensitiveWord, AuditRecord, SystemConfig, ViolationRecord, DailyStats structs

**internal/middleware/auth.go** — JWT validation + admin role check middleware

**pkg/database/db.go**, **pkg/redis/redis.go** — Same pattern

**internal/client/user_client.go** — HTTP client stub for user-service
**internal/client/forum_client.go** — HTTP client stub for forum-service
</action>

<acceptance_criteria>
- services/admin/go.mod exists with module ai-forum/admin-service
- services/admin/cmd/main.go exists with port ":8003"
- services/admin/internal/handler/health.go exists returning "admin-service"
- services/admin/internal/model/admin.go exists with SensitiveWord struct
- services/admin/internal/handler/audit_handler.go exists
- services/admin/internal/handler/config_handler.go exists
- services/admin/internal/handler/stats_handler.go exists
- services/admin/internal/handler/moderation_handler.go exists
- services/admin/internal/client/user_client.go exists
- services/admin/internal/client/forum_client.go exists
- go.mod contains github.com/gin-gonic/gin
- go.mod contains github.com/golang-migrate/migrate/v4
</acceptance_criteria>

---

## Wave 4: Frontend skeleton (depends on Wave 1)

### Task 4.1: Create Vue 3 + Vite frontend skeleton

<objective>
Create minimal Vue 3 + Vite project structure with basic routing.
</objective>

<read_first>
- AI智联论坛-需求文档+开发边界+MVP范围+微服务架构设计.md (section 7 for directory structure)
- docker-compose.yml
</read_first>

<action>
Create the following under `frontend/`:

**package.json:**
```json
{
  "name": "ai-forum-frontend",
  "version": "0.1.0",
  "type": "module",
  "scripts": {
    "dev": "vite --host",
    "build": "vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "vue": "^3.4.0",
    "vue-router": "^4.2.0",
    "pinia": "^2.1.0",
    "axios": "^1.6.0"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.0.0",
    "vite": "^5.0.0"
  }
}
```

**vite.config.js:**
```js
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3000,
    proxy: {
      '/api': 'http://localhost:80'
    }
  }
})
```

**index.html:** Standard Vite HTML entry point with `<div id="app"></div>`

**src/main.js:**
```js
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
```

**src/App.vue:** Basic shell with `<router-view />`

**src/router/index.js:** Basic routes:
- `/` → Home view
- `/login` → Login view
- `/register` → Register view

**src/views/Home.vue** — Placeholder "AI智联论坛" landing page
**src/views/Login.vue** — Placeholder login form
**src/views/Register.vue** — Placeholder registration form

**src/api/index.js** — Axios instance with base URL from `import.meta.env.VITE_API_BASE_URL`

**Dockerfile** (for docker-compose):
```Dockerfile
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
```

**.env.example:**
```
VITE_API_BASE_URL=/api/v1
```
</action>

<acceptance_criteria>
- frontend/package.json exists with vue, vue-router, pinia, axios dependencies
- frontend/vite.config.js exists with port 3000 and /api proxy
- frontend/index.html exists
- frontend/src/main.js exists with createApp, createPinia, router setup
- frontend/src/App.vue exists with router-view
- frontend/src/router/index.js exists with /, /login, /register routes
- frontend/src/views/Home.vue exists
- frontend/src/views/Login.vue exists
- frontend/src/views/Register.vue exists
- frontend/src/api/index.js exists with axios instance
- frontend/Dockerfile exists with multi-stage build (node:20-alpine build, nginx:alpine serve)
- frontend/.env.example exists with VITE_API_BASE_URL
</acceptance_criteria>

---

## Wave 5: Verification & documentation

### Task 5.1: Create Dockerfile for each Go service

<objective>
Create Dockerfiles for all three Go services.
</objective>

<read_first>
- services/user/cmd/main.go
- services/forum/cmd/main.go
- services/admin/cmd/main.go
</read_first>

<action>
Create identical-pattern Dockerfiles for each service:

**services/user/Dockerfile:**
```Dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service ./cmd/main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /user-service .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8001
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8001/health || exit 1
CMD ["./user-service"]
```

**services/forum/Dockerfile:** Same pattern, EXPOSE 8002, HEALTHCHECK on :8002/health
**services/admin/Dockerfile:** Same pattern, EXPOSE 8003, HEALTHCHECK on :8003/health
</action>

<acceptance_criteria>
- services/user/Dockerfile exists
- Contains golang:1.22-alpine build stage
- Contains HEALTHCHECK on http://localhost:8001/health
- services/forum/Dockerfile exists with HEALTHCHECK on :8002
- services/admin/Dockerfile exists with HEALTHCHECK on :8003
- All Dockerfiles have multi-stage build pattern
</acceptance_criteria>

---

### Task 5.2: Update docker-compose.yml with env_file and healthchecks

<objective>
Add env_file references and health checks to docker-compose.yml.
</objective>

<read_first>
- docker-compose.yml
- .env.example
</read_first>

<action>
Update docker-compose.yml:

1. Add `env_file: .env` to all three Go services
2. Add healthcheck to postgres: `pg_isready -U ai_forum`
3. Add healthcheck to redis: `redis-cli ping`
4. Ensure depends_on uses health check condition pattern for postgres and redis
</action>

<acceptance_criteria>
- docker-compose.yml contains env_file: .env for user-service, forum-service, admin-service
- docker-compose.yml contains healthcheck for postgres with pg_isready
- docker-compose.yml contains healthcheck for redis with redis-cli ping
</acceptance_criteria>

---

### Task 5.3: Create project README.md

<objective>
Create project README with setup instructions.
</objective>

<read_first>
- AI智联论坛-需求文档+开发边界+MVP范围+微服务架构设计.md
- .planning/ROADMAP.md
</read_first>

<action>
Create `README.md` with:
- Project name: AI智联论坛 (AI Forum)
- Tech stack summary: Go + Gin, Vue 3 + Vite, PostgreSQL 16, Redis 7, Nginx
- Quick start: `cp .env.example .env && docker-compose up`
- Service ports table (8001, 8002, 8003, 80, 3000)
- Directory structure overview
- Links to API route documentation
</action>

<acceptance_criteria>
- README.md exists at project root
- Contains "docker-compose up" instruction
- Contains port 8001, 8002, 8003
- Contains tech stack: Go, Gin, Vue 3, PostgreSQL, Redis, Nginx
</acceptance_criteria>

---

## Verification Criteria

1. `docker-compose config` validates the compose file syntax
2. `docker-compose up -d` starts all 7 services without errors
3. `curl http://localhost:8001/health` returns `{"status":"ok","service":"user-service"}`
4. `curl http://localhost:8002/health` returns `{"status":"ok","service":"forum-service"}`
5. `curl http://localhost:8003/health` returns `{"status":"ok","service":"admin-service"}`
6. `curl http://localhost/` returns the Vue frontend
7. PostgreSQL has three schemas: schema_auth, schema_forum, schema_admin (verified via `psql -U ai_forum -c "\dn"`)
8. All services respond to health checks through Nginx proxy

## Schema Push

No schema push required for this phase — migrations run automatically on service startup via golang-migrate integration (D-08).

---

*Phase: 01-infrastructure*
*Plan created: 2026-05-21*
