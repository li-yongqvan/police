# AI 智联论坛 MVP

面向学院级人工智能协会展示的可运行 MVP：Vue 3 前端 + Go 三微服务后端（PostgreSQL + Redis）。

**功能说明与使用指南（领域驱动设计视角）** → [docs/项目介绍与使用指南-领域驱动设计视角.md](docs/项目介绍与使用指南-领域驱动设计视角.md)

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

可选：启用 QQ 一键登录（QQ 互联网站应用）。

- 在 QQ 互联创建「网站应用」，并把回调地址设置为 `.env` 里的 `QQ_REDIRECT_URI`（必须完全一致）。
- 本地开发默认回调：`http://127.0.0.1:8001/api/v1/auth/qq/callback`
- 前端回跳页：`.env` 里的 `FRONTEND_OAUTH_REDIRECT_URL`（默认 `http://127.0.0.1:8091/oauth/qq`）
- 若希望 QQ 首次登录必须先邀请码注册：设置 `QQ_OAUTH_REQUIRE_INVITE=1`
- QQ 首次登录会自动创建本地用户（默认角色 student），并将 `profile_completed=false`，进入社区后可继续补全资料

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
npm run dev -- --host 127.0.0.1 --port 8091
```

### 4. 演示账号

完整列表见 **[docs/demo-accounts.md](docs/demo-accounts.md)**。

| 角色 | 用户名（示例） | 密码 |
|------|----------------|------|
| 学生 | `demo_student` 或 `demo01`…`demo06` | `demo123456` 或 **与用户名相同** |
| 协会管理员 | `demo_admin` 或 `admin01`…`admin06` | 同上 |
| 中台管理员 | `demo_platform_admin` 或 `plat01`…`plat06` | 同上 |

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
| frontend (dev) | 8091 |
| PostgreSQL | 5432 |
| Redis | 6379 |
| Nginx (compose) | 80 |

## 文档

- **[傻瓜式全功能操作清单](docs/step-by-step-uat-guide.md)** — 62 步逐步点击说明（推荐测试同学使用）
- **[人工全功能 UAT 方案](docs/human-full-experience-uat.md)** — 六大测试方法论融合；含体验反馈单（提交给 AI 改代码）
- **[项目介绍与使用指南（DDD 领域视角）](docs/项目介绍与使用指南-领域驱动设计视角.md)** — 限界上下文、角色能力、演示与本地用法
- [后端框架分析报告](docs/backend-framework-analysis-report.md)
- [后端对接实施计划](docs/backend-framework-integration-plan.md)

## 仅前端演示

无需启动 Go 服务时，使用 `frontend-standalone/`（数据存于浏览器 localStorage）。
