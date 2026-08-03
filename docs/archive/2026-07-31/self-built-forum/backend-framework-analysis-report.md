# Gitee 后端框架分析报告

> **分析对象**：`reference/ai-forum-gitee/`（[hebaonaodai/ai-forum](https://gitee.com/hebaonaodai/ai-forum)）  
> **分析日期**：2026-05-24  
> **范围**：仅后端（`services/`、`migrations/`、`nginx/`、`docker-compose.yml`、`.env.example`），不含前端实现细节。

---

## 1. 执行摘要

Gitee 仓库提供了一套 **生产向的三微服务后端**：Go 1.23 + Gin，共享 PostgreSQL 16（按 Schema 隔离），Redis 7 做刷新令牌与部分缓存，Nginx 作统一入口。业务按 **用户 / 论坛 / 中台** 拆分，层内采用经典的 **Handler → Service → Model**，配套各服务独立 SQL 迁移与服务间 HTTP 调用。

与当前主工程 MVP（`shared/mock-data/` + 轻量 Go 单文件服务或演示接口）相比，该框架在 **数据持久化、邀请注册、JWT 鉴权、敏感词审核、统计与中台编排** 上更完整，但存在 **网关路由遗漏、Schema 初始化未自动化、管理员 JWT 与角色表不一致** 等需在对接前处理的问题。

---

## 2. 总体架构

### 2.1 架构图

```
                    ┌─────────────┐
  客户端 :80 ──────►│   Nginx     │  API 网关 + 静态前端
                    └──────┬──────┘
           ┌───────────────┼───────────────┐
           ▼               ▼               ▼
    user-service      forum-service    admin-service
        :8001             :8002            :8003
           │               │               │
           └───────┬───────┴───────┬───────┘
                   ▼               ▼
              PostgreSQL 16      Redis 7
              DB: ai_forum
              schema_auth | schema_forum | schema_admin
```

### 2.2 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| 语言/框架 | Go 1.23 + Gin | 三个独立 `go.mod` 微服务 |
| 数据库 | PostgreSQL 16 | 单库多 Schema |
| 缓存 | Redis 7 | 刷新令牌、角色缓存等 |
| 迁移 | golang-migrate | 各服务独立版本表 |
| 网关 | Nginx | `nginx/nginx.conf` |
| 编排 | Docker Compose | 根目录 `docker-compose.yml` |

### 2.3 服务职责

| 服务 | 端口 | Schema | 核心职责 |
|------|------|--------|----------|
| user-service | 8001 | `schema_auth` | 注册/登录/JWT、用户资料、邀请码、用户级统计 |
| forum-service | 8002 | `schema_forum` | 板块、帖子、评论、点赞收藏、附件、发帖审核联动 |
| admin-service | 8003 | `schema_admin` | 系统配置、审核队列、敏感词、角色、中台统计 |

---

## 3. 代码结构与分层

三个服务目录结构一致，体现统一的 **垂直分层 + 横向 pkg** 约定：

```text
services/<name>/
├── cmd/main.go              # 入口：连库、跑迁移、注册路由、监听端口
├── Dockerfile
├── migrations/              # 本服务表结构（up/down）
├── internal/
│   ├── handler/             # HTTP 入参校验、响应封装
│   ├── service/             # 业务逻辑、事务、跨表/跨服务编排
│   ├── model/               # 领域模型与 DTO
│   ├── middleware/          # JWT、CORS、等级校验
│   └── client/              # 调用其他微服务（forum/admin）
└── pkg/
    ├── database/            # pgx 连接池 + RunMigrations
    ├── jwt/                 # HS256 签发与解析
    └── redis/               # 连接与键操作
```

**设计特点**：

- 无共享 `go.work` 公共库，各服务复制 `pkg/database`、`pkg/jwt`（便于独立部署，但存在重复维护成本）。
- 对外 API 统一前缀 `/api/v1`；服务间能力暴露在 `/internal/v1`（不经过 Nginx）。
- 审核流：发帖/评论时 forum 通过 `AdminClient` 调 admin 的 `POST /internal/v1/moderation/check`。

---

## 4. 各服务 API 概览

### 4.1 用户服务（user-service）

**入口**：`services/user/cmd/main.go`

| 类型 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 公开 | POST | `/api/v1/register` | 邀请码注册 |
| 公开 | POST | `/api/v1/login` | 登录，返回 JWT |
| 公开 | POST | `/api/v1/auth/refresh` | 刷新令牌 |
| 需 JWT | GET/PUT | `/api/v1/users/:id` | 资料读写 |
| 需 JWT | POST | `/api/v1/users/:id/avatar` | 头像 |
| 内部 | * | `/internal/v1/users/...` | 封禁、邀请码、用户列表、统计 |
| 健康 | GET | `/health` | 健康检查 |

### 4.2 论坛服务（forum-service）

**入口**：`services/forum/cmd/main.go`

| 类型 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 公开 | GET | `/api/v1/boards` | 板块列表 |
| 公开 | GET | `/api/v1/posts` | 帖子列表 |
| 需 JWT | POST | `/api/v1/posts` | 发帖（level≥1） |
| 需 JWT | POST | `/api/v1/posts/:id/comments` | 评论 |
| 需 JWT | POST | `/api/v1/posts/:id/like` | 点赞 |
| 需 JWT | POST | `/api/v1/attachments/upload` | 附件（level≥2） |
| 内部 | * | `/internal/v1/...` | 中台删帖、置顶、待审列表、统计 |
| 出站 | POST | `{ADMIN_SERVICE_URL}/internal/v1/moderation/check` | 敏感词/审核 |

**注意**：`forum_service` 在 SQL 中 **JOIN `schema_auth.users`** 取作者名，属于同库跨 Schema 访问，而非 HTTP 调 user-service。

### 4.3 中台服务（admin-service）

**入口**：`services/admin/cmd/main.go`

| 类型 | 方法 | 路径前缀 | 说明 |
|------|------|----------|------|
| 需管理员 JWT | * | `/api/v1/admin/` | 配置、审核、用户、角色、敏感词、统计 |
| 内部 | POST | `/internal/v1/moderation/check` | 供 forum 调用 |
| 出站 | HTTP | user/forum internal API | `UserClient`、`ForumClient` |

**鉴权要求**：JWT 中 `role` 为 `admin` 或 `platform_admin`（见第 6 节问题）。

---

## 5. 数据层设计

### 5.1 Schema 初始化

全局脚本：`migrations/init/000_init_schemas.up.sql`

- 创建 `schema_auth`、`schema_forum`、`schema_admin`
- **不会**被 `docker-compose` 自动执行，需人工或 init 容器先跑一遍

### 5.2 各服务迁移

| 服务 | 版本表 | 主要表 |
|------|--------|--------|
| user | `schema_auth_migrations` | `users`, `invite_codes`, `operation_logs` |
| forum | `schema_forum_migrations` | `boards`, `posts`, `comments`, `likes`, `collections`, `attachments` |
| admin | `schema_admin_migrations` | `sensitive_words`, `system_config`, `roles`, `user_roles`, `statistics_daily` |

服务启动时在 `pkg/database/db.go` 中调用 **golang-migrate**，从各自 `migrations/` 目录升级。

### 5.3 审核模型

- **无** 独立 `audit_records` 表（与当前 MVP 的 `audit-records.json` 不同）
- 待审内容体现为 forum 帖子 `status = 'pending_review'`
- 中台通过 admin API 拉取待审列表并 approve/reject

### 5.4 种子数据

- forum `001` 迁移内置 3 个板块：`ai-learning`、`announcements`、`tech-help`
- admin `001` 迁移内置敏感词种子

---

## 6. 认证与安全

| 项目 | 实现 |
|------|------|
| 算法 | HS256 JWT（`JWT_SECRET` 环境变量） |
| Claims | `user_id`, `username`, `role`, `level` |
| 刷新 | Redis 键 `refresh:{userId}` |
| 论坛发帖 | `RequireLevel(1)` 等中间件 |
| 中台 API | 要求 `role ∈ {admin, platform_admin}` |

### 6.1 已识别的框架层问题

1. **Nginx 未代理用户登录/注册**  
   `nginx.conf` 仅配置 `/api/v1/auth/`，而注册/登录在 `/api/v1/register`、`/api/v1/login`。经 `:80` 网关无法完成完整认证流程，需直连 `:8001` 或改配置。

2. **管理员 JWT 与角色表脱节**  
   登录签发时 `role` 固定为 `"user"`；管理员角色在 `schema_admin.user_roles`。不额外发 token 则中台 API 可能全部 403。

3. **Schema 初始化非自动化**  
   首次部署若未执行 `000_init_schemas.up.sql`，各服务迁移会失败。

4. **内部 API 无鉴权**  
   `/internal/v1/*` 依赖 Docker 网络隔离，无 mTLS 或内部密钥。

5. **forum 未接 Redis**  
   `pkg/redis` 存在但 `main.go` 未连接（与 compose 中 redis 依赖声明不完全一致）。

---

## 7. 部署与运维

### 7.1 Docker Compose 服务清单

| 容器 | 镜像/构建 | 宿主机端口 |
|------|-----------|------------|
| nginx | nginx:alpine | 80 |
| user-service | `services/user/Dockerfile` | 8001 |
| forum-service | `services/forum/Dockerfile` | 8002 |
| admin-service | `services/admin/Dockerfile` | 8003 |
| postgres | postgres:16 | 5432 |
| redis | redis:7-alpine | 6379 |
| frontend | `frontend/Dockerfile` | 3000 |

环境变量模板：`.env.example`（DB、Redis、JWT、各服务 URL）。

### 7.2 与仓库内设计文档的关系

`reference/ai-forum-gitee/AI智联论坛-需求文档+开发边界+MVP范围+微服务架构设计.md` 描述了完整产品愿景（多角色、中台、远期社会级能力）。**当前 Gitee 代码实现的是 MVP 子集**：三服务 + 基础中台，尚未覆盖文档中的全部远期功能（积分、私信、第三方登录等）。

`.planning/` 目录记录了分阶段实施（基础设施 → 用户功能 → 中台 → 统计 → 集成测试），可作为能力完成度的对照清单。

---

## 8. 与主工程 MVP 的差异对照

主工程现状（`frontend/` + `shared/mock-data/`）：

| 维度 | 主工程 MVP | Gitee 参考框架 |
|------|------------|----------------|
| 存储 | JSON 文件 | PostgreSQL + 迁移 |
| 用户 ID | 字符串 `user-student` | 自增整数 |
| 登录 | `POST /api/v1/demo-login` 按角色切换 | 邀请码注册 + 密码登录 |
| 板块 slug | `board-study` 等 | `ai-learning` 等（种子不同） |
| 审核 | `audit-records.json` | 帖子 `pending_review` 状态 |
| 配置 | 单 JSON 对象 | `system_config` 键值表 |
| API 响应 | `{ success, message, data }` | 多为 REST 风格 `gin.H` / 错误字段 |
| 网关 | Vite 代理 `/user-api` 等 | Nginx 按路径分流 |

**对接影响**：不能简单「替换 services 目录」即跑通现有前端，需要 **API 适配层** 或分阶段改前端契约。

---

## 9. 结论与建议

### 9.1 框架评价

| 优点 | 风险/不足 |
|------|-----------|
| 清晰的微服务边界与分层 | Nginx 与用户 API 不完整 |
| Schema 隔离，便于日后拆库 | 同库 JOIN `schema_auth` 削弱边界 |
| 完整迁移与 Docker 化 | Schema 初始化需人工步骤 |
| 中台与用户/论坛编排较全 | 管理员 JWT 逻辑需补齐 |
| 与产品需求文档方向一致 | 与当前演示前端契约差异大 |

### 9.2 建议使用方式

1. **保持** `reference/ai-forum-gitee/` 为只读参考，不在此目录改业务代码。  
2. **采纳** 分层目录、三服务拆分、Schema 命名、迁移策略作为目标架构。  
3. **勿一次性替换** 主工程后端；按《后端对接实施计划》分阶段迁移。  
4. **优先修复** 网关与 JWT 角色问题后再对接前端。

---

## 10. 关键文件索引

| 用途 | 路径（相对 `reference/ai-forum-gitee/`） |
|------|------------------------------------------|
| 编排 | `docker-compose.yml` |
| 网关 | `nginx/nginx.conf` |
| 环境变量 | `.env.example` |
| Schema 初始化 | `migrations/init/000_init_schemas.up.sql` |
| 用户服务入口 | `services/user/cmd/main.go` |
| 论坛服务入口 | `services/forum/cmd/main.go` |
| 中台服务入口 | `services/admin/cmd/main.go` |
| Forum→Admin 客户端 | `services/forum/internal/client/admin_client.go` |
| 产品/架构说明 | `AI智联论坛-需求文档+开发边界+MVP范围+微服务架构设计.md` |

---

*本报告仅基于 `reference/ai-forum-gitee` 静态代码与文档分析，未执行完整联调测试。*
