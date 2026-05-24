---
phase: 01-infrastructure
plan: 01
subsystem: infrastructure
tags: [docker, postgresql, redis, nginx, go, gin, vue3, vite]
dependency_graph:
  requires: []
  provides: [infra-scaffold, db-schemas, go-services, frontend-skeleton]
  affects: []
tech-stack:
  added:
    - Go 1.22 + Gin
    - Vue 3 + Vite
    - PostgreSQL 16
    - Redis 7
    - Nginx Alpine
    - golang-migrate v4
    - pgx/v5
    - go-redis/v9
    - golang-jwt/v5
    - axios
    - pinia
    - vue-router
  patterns:
    - Multi-stage Docker builds
    - Shared PostgreSQL + schema isolation
    - JWT authentication middleware
    - REST API gateway (Nginx)
key-files:
  created:
    - docker-compose.yml
    - .env.example
    - nginx/nginx.conf
    - migrations/init/000_init_schemas.up.sql
    - migrations/init/000_init_schemas.down.sql
    - services/user/ (16 files)
    - services/forum/ (16 files)
    - services/admin/ (18 files)
    - frontend/ (12 files)
    - README.md
  modified: []
decisions:
  - Added nginx to user-service Dockerfile healthcheck deps for wget support
  - Inlined admin_client handler in forum main.go to avoid cross-package gin import
  - Added CreatePostRequest type to forum model (was referenced in service layer)
  - docker-compose.yml created with env_file and healthchecks from the start (combined Tasks 1.1 + 5.2)
metrics:
  duration_minutes: 45
  completed_date: "2026-05-21"
  tasks_completed: 10
  files_created: 64
  commits: 5
---

# Phase 01 Plan 01: 基础设施与项目骨架 Summary

**One-liner:** Complete Docker Compose orchestration with 7 services, PostgreSQL schema migrations, three Go microservice skeletons (user:8001, forum:8002, admin:8003) with JWT auth and DB/Redis integration, Vue 3 frontend with routing, multi-stage Dockerfiles for all services, and project README.

## Commits

| # | Hash    | Message |
|---|---------|---------|
| 1 | 85c7ac5 | create docker-compose, env template, nginx config |
| 2 | 28b4bad | add PostgreSQL schema init migrations |
| 3 | 1d35b0c | create three Go service skeletons |
| 4 | 71ab1f5 | create Vue 3 frontend skeleton |
| 5 | 42dd62a | add Go service Dockerfiles and project README |

## Tasks Completed

### Wave 1: 基础设施编排

| Task | Status | Files |
|------|--------|-------|
| 1.1 docker-compose.yml | Done | `docker-compose.yml` |
| 1.2 .env.example | Done | `.env.example` |
| 1.3 nginx/nginx.conf | Done | `nginx/nginx.conf` |

- 7 services: nginx (:80), user-service (:8001), forum-service (:8002), admin-service (:8003), postgres (16), redis (7-alpine), frontend (:3000)
- postgres/redis healthchecks with `pg_isready` and `redis-cli ping`
- Backend services depend_on postgres/redis with `condition: service_healthy`
- `env_file: .env` on all three Go services
- pgdata volume for PostgreSQL

### Wave 2: Database schema setup

| Task | Status | Files |
|------|--------|-------|
| 2.1 Schema migrations | Done | `migrations/init/000_init_schemas.{up,down}.sql` |

- Three schemas: `schema_auth`, `schema_forum`, `schema_admin`
- Down migration with CASCADE

### Wave 3: Go service skeletons

| Task | Status | Files |
|------|--------|-------|
| 3.1 user-service | Done | 16 files under `services/user/` |
| 3.2 forum-service | Done | 16 files under `services/forum/` |
| 3.3 admin-service | Done | 18 files under `services/admin/` |

Each service includes:
- `go.mod` with gin, pgx, go-redis, jwt, godotenv, golang-migrate
- `cmd/main.go` with env loading, DB/Redis connection, migration run, route registration
- `internal/handler/health.go` returning `{"status":"ok","service":"<name>"}`
- `internal/handler/<type>_handler.go` stub handlers with TODO comments
- `internal/service/<name>_service.go` stub service layer
- `internal/model/<type>.go` with typed structs
- `internal/middleware/auth.go` JWT validation middleware (admin has role check)
- `internal/middleware/cors.go` (user-service)
- `internal/client/admin_client.go` (forum-service, calls ADMIN_SERVICE_URL)
- `internal/client/user_client.go`, `forum_client.go` (admin-service)
- `pkg/database/db.go` pgxpool connection + golang-migrate integration
- `pkg/redis/redis.go` Redis client with ping check
- `pkg/jwt/jwt.go` GenerateToken + ValidateToken with HS256
- `migrations/` directory for service-specific migrations

### Wave 4: Frontend skeleton

| Task | Status | Files |
|------|--------|-------|
| 4.1 Vue 3 + Vite | Done | 12 files under `frontend/` |

- package.json: vue 3.4, vue-router 4.2, pinia 2.1, axios 1.6
- vite.config.js: port 3000, /api proxy
- Router: /, /login, /register routes
- Views: Home (hero + board cards), Login (email/password), Register (username/email/password/invitationCode)
- api/index.js: axios instance with JWT interceptor, 401 auto-redirect
- Dockerfile: multi-stage (node:20-alpine build -> nginx:alpine serve)
- .env.example: VITE_API_BASE_URL=/api/v1

### Wave 5: Verification & documentation

| Task | Status | Files |
|------|--------|-------|
| 5.1 Go service Dockerfiles | Done | `services/{user,forum,admin}/Dockerfile` |
| 5.2 docker-compose env_file + healthchecks | Done | Integrated into initial docker-compose.yml |
| 5.3 README.md | Done | `README.md` |

- Multi-stage Go builds (golang:1.22-alpine builder -> alpine:3.19 runtime)
- Healthchecks: wget -qO- http://localhost:<port>/health
- README: tech stack, quick start, port table, API routes, directory structure

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Removed unused `io` import in admin_client.go**
- **Found during:** Wave 5 Dockerfile creation (reviewing all service files)
- **Issue:** `io` import in `services/forum/internal/client/admin_client.go` was unused, causing compilation failure
- **Fix:** Removed the unused import
- **Files modified:** `services/forum/internal/client/admin_client.go`

**2. [Rule 2 - Missing functionality] Added missing `CreatePostRequest` type**
- **Found during:** Wave 3 service file creation
- **Issue:** `ForumService.CreatePost()` referenced `model.CreatePostRequest` which was not defined in `post.go`
- **Fix:** Added `CreatePostRequest` struct with Title, Content, BoardID fields
- **Files modified:** `services/forum/internal/model/post.go`

**3. [Rule 1 - Bug] Fixed gin import in admin_client**
- **Found during:** Wave 5
- **Issue:** `admin_client.go` imported `github.com/gin-gonic/gin` but only used `gin.Context` in a handler method that was not needed externally
- **Fix:** Removed gin import and inlined the handler logic in `cmd/main.go` instead
- **Files modified:** `services/forum/internal/client/admin_client.go`, `services/forum/cmd/main.go`

**4. [Rule 3 - Blocker] Added `wget` to Go service Dockerfiles**
- **Found during:** Wave 5
- **Issue:** Alpine base image does not include wget; healthcheck CMD references wget
- **Fix:** Added `apk --no-cache add ca-certificates wget` to each Dockerfile runtime stage
- **Files modified:** `services/{user,forum,admin}/Dockerfile`

## Known Stubs

All handler endpoints return `http.StatusNotImplemented` with `{"error": "not implemented"}` and TODO comments. These are intentional for this infrastructure phase -- the next phase (02: user endpoint core features) will wire up actual business logic.

Specific stubs:
- All user-service handlers: Register, Login, GetProfile, UpdateProfile, UploadAvatar
- All forum-service handlers: ListPosts, GetPost, CreatePost, UpdatePost, DeletePost, ListComments, CreateComment, LikePost, CollectPost, UploadAttachment, DownloadAttachment, ListBoards, GetBoard
- All admin-service handlers: GetConfig, UpdateConfig, ListPendingAudit, ApprovePost, RejectPost, BanUser, ListUsers, UpdateUserLevel, GetUserLogs, GetOverview, GetDailyStats, AddSensitiveWord, CheckSensitiveWords (returns hardcoded `{"clean": true}`)

## Threat Flags

| Flag | File | Description |
|------|------|-------------|
| threat_flag:jwt-secret | services/*/cmd/main.go | JWT secret loaded from env; default must be changed before production |
| threat_flag:internal-routes | services/*/cmd/main.go | /internal/v1/ endpoints exposed but Nginx does NOT proxy them externally |
| threat_flag:schema-isolation | migrations/init/000_init_schemas.up.sql | Shared PostgreSQL with schema isolation -- schema bypass possible if search_path not set |

## Self-Check: PASSED

All created files verified:
- docker-compose.yml: 7 services, postgres:16, redis:7-alpine, ports 8001/8002/8003, pgdata volume, healthchecks, env_file
- .env.example: all required variables present
- nginx/nginx.conf: 4 upstreams, 6+ location blocks, CORS, proxy headers
- migrations/init/: up.sql (3 schemas), down.sql (3 DROP CASCADE)
- services/user/: go.mod (ai-forum/user-service), cmd/main.go (:8001, gin.Default, health), all handler/model/service/middleware/pkg files
- services/forum/: go.mod (ai-forum/forum-service), cmd/main.go (:8002, ADMIN_SERVICE_URL), admin_client with CheckSensitiveWords, all handlers, model with Post struct
- services/admin/: go.mod (ai-forum/admin-service), cmd/main.go (:8003), SensitiveWord struct, audit/config/stats/moderation handlers, user_client + forum_client
- frontend/: package.json (vue/vue-router/pinia/axios), vite.config.js (port 3000, /api proxy), router with /login/register, views, api/index.js, Dockerfile (node:20-alpine), .env.example
- services/*/Dockerfile: golang:1.22-alpine, healthcheck on respective ports
- README.md: docker-compose up, ports 8001/8002/8003, tech stack
