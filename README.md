# AI 智联论坛 MVP

面向学院级人工智能协会展示的可运行 MVP：Vue 3 前端 + Go 三微服务后端（PostgreSQL + Redis）。

## 目录结构

```text
frontend/                 Vue 3 + Vite 前端（学院演示 UI，唯一线上前端）
frontend-standalone/      纯前端 Mock 演示（本地离线用，不部署）
services/user/            用户服务 (:8001)
services/forum/           论坛服务 (:8002)
services/admin/           中台服务 (:8003)
migrations/init/          数据库 Schema 初始化
nginx/                    API 网关配置（Docker 全栈）
infra/                    本地 Postgres + Redis
reference/ai-forum-gitee/ Gitee 后端参考（只读，不含前端）
docs/                     分析与计划文档
```

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.23 + Gin |
| 前端 | Vue 3 + Vite |
| 数据库 | PostgreSQL 16（三 Schema） |
| 缓存 | Redis 7 |
| 部署 | Docker Compose + Nginx |

## 本地开发

### 1. 启动基础设施

```powershell
docker compose -f infra/docker-compose.yml up -d
```

首次需初始化 Schema（仅需一次）：

```powershell
# 需本机安装 psql，或使用 Docker 执行
psql -h 127.0.0.1 -U ai_forum -d ai_forum -f migrations/init/000_init_schemas.up.sql
```

### 2. 环境变量

```powershell
cp .env.local.example .env
```

### 3. 启动后端与前端

```powershell
.\start-dev.ps1
```

或分别启动：

```powershell
cd services/user
$env:DB_SCHEMA="schema_auth"
go run ./cmd/main.go
```

```powershell
cd services/forum
$env:DB_SCHEMA="schema_forum"
go run ./cmd/main.go
```

```powershell
cd services/admin
$env:DB_SCHEMA="schema_admin"
go run ./cmd/main.go
```

```powershell
cd frontend
npm install
npm run dev -- --host 127.0.0.1 --port 4173
```

### 4. 演示账号

| 角色 | 演示登录按钮 | 密码 |
|------|--------------|------|
| 学生 | 学生用户 | `demo123456` |
| 协会管理员 | 协会管理员 | `demo123456` |
| 中台管理员 | 中台管理员 | `demo123456` |

## Ubuntu 云主机部署

见 **[DEPLOY.md](DEPLOY.md)**。快速命令：

```bash
cp .env.production.example .env
# 编辑 JWT_SECRET、数据库密码
bash scripts/deploy-ubuntu.sh
```

## Docker 全栈（开发）

```powershell
cp .env.example .env
docker compose up -d --build
```

生产建议叠加：`docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build`

- 前端入口：`http://localhost/`
- 用户服务健康检查：`http://localhost:8001/health`

## 默认端口

| 服务 | 端口 |
|------|------|
| user-service | 8001 |
| forum-service | 8002 |
| admin-service | 8003 |
| frontend (dev) | 4173 |
| PostgreSQL | 5432 |
| Redis | 6379 |
| Nginx (compose) | 80 |

## 文档

- [后端框架分析报告](docs/backend-framework-analysis-report.md)
- [后端对接实施计划](docs/backend-framework-integration-plan.md)

## 仅前端演示

无需启动 Go 服务时，使用 `frontend-standalone/`（数据存于浏览器 localStorage）。
