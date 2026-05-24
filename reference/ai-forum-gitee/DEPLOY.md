# AI智联论坛 MVP 部署指南

## 快速启动

```bash
# 一键启动所有服务
docker compose up -d

# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f
```

## 前置条件

- Docker Desktop (Windows/Mac) 或 Docker Engine (Linux)
- Docker Compose v2
- 至少 2GB 可用内存
- 端口 80, 8001, 8002, 8003, 5432, 6379 未被占用

## 服务说明

| 服务 | 端口 | 说明 |
|------|------|------|
| nginx | 80 | 网关，前端和API路由 |
| user-service | 8001 | 用户认证、资料管理 |
| forum-service | 8002 | 帖子、评论、板块、附件 |
| admin-service | 8003 | 中台管理、审核、统计 |
| postgres | 5432 | PostgreSQL 16 数据库 |
| redis | 6379 | Redis 7 缓存 |
| frontend | 5173 | Vue 3 开发服务器（dev模式） |

## 环境变量

复制 `.env.example` 为 `.env` 并修改：

```env
DB_HOST=postgres
DB_PORT=5432
DB_NAME=ai_forum
DB_USER=postgres
DB_PASSWORD=your_password
REDIS_HOST=redis
REDIS_PORT=6379
```

## 数据库迁移

服务启动时会自动执行迁移。手动执行：

```bash
# 每个服务内部会运行 pkg/database 的 RunMigrations
# 使用 golang-migrate CLI 手动迁移：
migrate -path services/admin/migrations -database "postgres://..." up
migrate -path services/forum/migrations -database "postgres://..." up
migrate -path services/user/migrations -database "postgres://..." up
```

## 开发模式

```bash
# 只启动数据库和Redis
docker compose up -d postgres redis

# 分别启动各服务（从各自目录）
cd services/user && go run ./cmd/main.go
cd services/forum && go run ./cmd/main.go
cd services/admin && go run ./cmd/main.go
cd frontend && npm run dev
```

## E2E 集成测试

服务启动后运行：

```bash
bash scripts/e2e-test.sh
```

## 常见问题

### 1. 服务启动失败
检查 `docker compose logs <service-name>` 查看详细错误。

### 2. 数据库连接失败
确保 PostgreSQL 已就绪，检查 healthcheck 状态。

### 3. 端口被占用
修改 `.env` 中的端口或关闭占用端口的进程。

### 4. 前端页面空白
确认后端服务已启动，检查浏览器控制台 API 错误。

### 5. 迁移重复执行
golang-migrate 使用 schema_migrations 表追踪版本，重复执行不会报错。

## 停止服务

```bash
# 停止所有服务（保留数据）
docker compose down

# 停止并删除数据卷（数据清空）
docker compose down -v
```
