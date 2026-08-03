# Handoff · Discourse 重建 — Ticket 5 完成 + Ticket 6 准备

> 日期：2026-07-30 | 服务器：122.51.233.225 | 分支：codex/discourse-rebuild | 状态：✅ Ticket 5 完成，🔄 Ticket 6 待执行

---

## 1. 当前进度总览

| Ticket | 描述 | 状态 |
|--------|------|------|
| Ticket 1 | Production Baseline Lock | ✅ 完成 |
| Ticket 2 | Discourse PoC Deployment | ✅ 完成 |
| Ticket 3a | Discourse SSO Provider | ✅ 完成 |
| Ticket 3b | Logout Sync | ✅ 完成 |
| Ticket 4 | 前端论坛导航切换 + SSO Cookie Bridge | ✅ 完成 |
| Ticket 5 | 管理后台论坛治理菜单过渡 + 旧页面删除 | ✅ 完成 |
| **Ticket 6** | **旧论坛彻底处置** | 🔄 **待执行** |
| Ticket 7 | 验证与冒烟测试 | ⬜ 后续 |

---

## 2. Ticket 5 执行摘要（本次线程完成）

### 提交记录

| Commit | 内容 |
|--------|------|
| `0d3a667` | `GxAdminSidebar.vue`：隐藏 audit/reports/posts/boards/sensitive 五个菜单项；侧边栏底部新增「外部管理」分区 + 「前往 Discourse 管理」链接 → `http://122.51.233.225:8080/admin` |
| `f950fd3` | `router.js`：移除 5 条路由定义 + 5 个懒加载 import + `admin-sensitive` 的 platform_only 守卫；删除 5 个页面文件（AdminAudit/AdminReports/AdminPosts/AdminBoards/AdminSensitiveWords.vue） |

### 部署

- 已通过 Python paramiko SFTP + SSH 部署到 `122.51.233.225` 的 `ai-forum-frontend-1` 容器
- 服务端验证通过：`index.html` 时间戳已更新，`AdminLayout` bundle 含 Discourse 链接

---

## 3. Ticket 6 目标

**彻底删除旧论坛（forum-service），清理所有相关配置、代理、死代码，使系统进入「平台 = 登录注册 + 管理后台，Discourse = 社区」的最终形态。**

---

## 4. Ticket 6 改动范围分析

### 4.1 核心删除：forum-service 服务（51 个文件）

```
services/forum/
├── cmd/main.go
├── internal/
│   ├── handler/       (board, comment, interaction, post, report, attachment, etc.)
│   ├── service/
│   ├── model/
│   ├── middleware/
│   ├── repository/
│   └── client/
├── pkg/
│   ├── database/
│   └── redis/
├── .env / .env.example
├── Dockerfile
└── go.mod / go.sum
```

### 4.2 Nginx 配置（nginx/nginx.conf + includes/）

需移除：
- `upstream forum_backend { server forum-service:8002; }`
- `/forum-api/` location block（含 `/api/v1/posts/*`、`/api/v1/boards/*`、`/api/v1/comments/*`、`/api/v1/attachments/*` 等四组 proxy_pass）

### 4.3 Docker Compose（6 个文件）

需移除 forum-service 段 + 相关依赖：
- `docker-compose.yml`
- `docker-compose.dev.yml`
- `docker-compose.prod.yml`
- `docker-compose.server.yml`
- `docker-compose.server-full.yml`
- `docker-compose.cn-mirror.yml`

注意：`server-full.yml` 中 postgres 的 `depends_on: forum-service` 需要移除。

### 4.4 Vite 前端代理（frontend/vite.config.js）

移除 `/forum-api` 代理规则。

### 4.5 admin-service 死代码清理（影响面大）

admin-service 通过 `internal/client/forum_client.go` 深度依赖 forum-service。

**需删除的文件**：
- `services/admin/internal/client/forum_client.go`（约 230 行，含所有 forum API 的 HTTP 封装）

**需修改的文件**：

| 文件 | 改动 |
|------|------|
| `cmd/main.go` | 移除 forumClient 创建 + 所有通过 forumClient 初始化的 handler/service 替换或移除 |
| `internal/handler/post_admin_handler.go` | 移除 ForumClient 依赖，`ListAllPosts`/`CreateBoard`/`UpdateBoard`/`DeleteBoard`/`ListAllBoards` 等方法需要处理 |
| `internal/handler/report_handler.go` | 移除 ForumClient 依赖 |
| `internal/service/audit_service.go` | 移除 ForumClient 依赖（ListPendingPosts/ChangePostStatus/BatchDeletePosts） |
| `internal/service/post_admin_service.go` | 移除 ForumClient 依赖（DeletePost/SetPostFeatured/SetPostPinned） |
| `internal/service/stats_service.go` | 移除 ForumClient 依赖（GetForumOverview/GetDailyPosts/GetBoardActivity） |

**关键决策点**：这些 handler 对应的前端页面已在 Ticket 5 中删除（AdminAudit/AdminReports/AdminPosts/AdminBoards），路由也已移除。但后端 handler 和 service 代码仍注册在 main.go 中。直接删除 handler 文件会导致编译失败，需要策划合理的删除路径：

- **方案 A（推荐）**：删除所有依赖 forumClient 的 handler、service、route 注册，集中清理 main.go → handler → service → client 整条链。stats_service 需要保留但移除 forum 统计部分。AuditService/PostAdminService 随 handler 一并删除。
- **方案 B（保守）**：保留 handler/service 文件但移除 ForumClient 依赖，让这些接口返回空数据或友好的过渡提示。不推荐——增加维护负担。

### 4.6 数据库 schema_forum

- `migrations/init/000_init_schemas.up.sql`：移除 `CREATE SCHEMA IF NOT EXISTS schema_forum;`
- `migrations/init/000_init_schemas.down.sql`：移除 `DROP SCHEMA IF EXISTS schema_forum CASCADE;`
- 运行中的 PostgreSQL：执行 `DROP SCHEMA IF EXISTS schema_forum CASCADE;`（⚠️ 不可逆，建议备份）

### 4.7 环境变量清理

以下文件中移除 `FORUM_SERVICE_URL`：
- `.env`
- `.env.example`
- `.env.local.example`
- `.env.production.example`
- `services/admin/.env`
- `services/user/.env`
- `services/user/.env.example`

### 4.8 前端 api.js（验证）

Ticket 4 已移除 `forumApi`，但需验证 `frontend/src/api.js` 无残留 forum 相关引用。

---

## 5. 不改动的范围

- **user-service**：不依赖 forum，无需修改
- **Discourse 容器/配置**：不动
- **schema_auth / schema_admin**：不删（schema_admin 是 admin-service 自身数据，与 forum 无关）
- **前端路由、页面、组件**：Ticket 4/5 已完成
- **API 服务间调用**：admin-service 对 user-service 的 UserClient 保留不动

---

## 6. 执行顺序建议

```
1. nginx.conf → 移除 forum 路由
2. Docker Compose → 移除 forum-service（6 个文件）
3. vite.config.js → 移除 /forum-api 代理
4. admin-service → 方案 A 清理死代码（handler + service + client + main.go 路由）
5. 环境变量 → 移除 FORUM_SERVICE_URL（7 个文件）
6. 迁移脚本 → 移除 schema_forum
7. services/forum/ → 删除整个目录
8. 构建验证 → docker compose build admin-service nginx frontend
9. 服务器 DROP SCHEMA → 在运行数据库中执行
10. 部署 → build → SCP → docker cp
```

---

## 7. Git 状态

- **分支**：`codex/discourse-rebuild`
- **HEAD**：`f950fd3`（feat: remove forum governance routes and page files, keep only admin core pages (#5)）
- **远程**：已 push

### 最近两次代码改动

1. `0d3a667`：GxAdminSidebar 菜单调整（Ticket 5 前半）
2. `f950fd3`：router + 5 页面删除（Ticket 5 后半）

---

## 8. 已知约束

- **服务器无法访问 GitHub/Docker Hub**：部署走本地 build → SCP → docker cp 方案
- **SSH 账号**：`liyongquan@122.51.233.225`，密码 `Liyongquan@123`
- **Discourse URL**：`http://122.51.233.225:8080`
- **项目路径**：服务器 `/home/liyongquan/projects/ai-forum`
- **每次 commit 后必须 push**，禁止跨 Ticket 混合改动
- **admin-service 的 UserClient 保留不动**（用于查询用户信息）

---

## 9. 新线程启动指令

> 读取 `docs/handoffs/handoff-2026-07-30-discourse-ticket5-done-ticket6-start.md`，先向用户列出修改计划对齐需求，确认后执行 Ticket 6。
