# Handoff · Discourse 重建 — Ticket 6 完成 + 前端构建修复 + Ticket 7 计划

> 日期：2026-07-30 | 服务器：122.51.233.225 | 分支：codex/discourse-rebuild | 状态：? Ticket 6 完成，?? Ticket 7 待执行

---

## 1. 当前进度总览

| Ticket | 描述 | 状态 |
|--------|------|------|
| Ticket 1 | Production Baseline Lock | ? 完成 |
| Ticket 2 | Discourse PoC Deployment | ? 完成 |
| Ticket 3a | Discourse SSO Provider | ? 完成 |
| Ticket 3b | Logout Sync | ? 完成 |
| Ticket 4 | 前端论坛导航切换 + SSO Cookie Bridge | ? 完成 |
| Ticket 5 | 管理后台论坛治理菜单过渡 + 旧页面删除 | ? 完成 |
| **Ticket 6** | **旧论坛彻底处置** | ? 完成 |
| **Ticket 6 补充** | **前端 Docker 构建修复 + 部署** | ? 完成 |
| **Ticket 7** | **验证与冒烟测试** | ?? **待执行** |

---

## 2. 本次线程执行摘要（Ticket 6 补充工作）

### 2.1 新增 Git 提交

| Commit | 内容 |
|--------|------|
| `15d172a` | `fix: upgrade Node to 22-alpine and set UTF-8 locale in frontend Dockerfile (#6)` |

改动：`frontend/Dockerfile` — `NODE_IMAGE` 从 `node:20-alpine` 升级到 `node:22-alpine`，添加 `ENV LANG=C.UTF-8 LC_ALL=C.UTF-8`。仅此一步即解决了 Vite 8 Oxide CSS 提取器的 UTF-8 panic 崩溃。

### 2.2 服务器部署中修复的额外问题

部署过程中因 `docker compose up -d frontend` 触发级联重建，暴露了 3 个此前未发现的配置问题，均已修复：

| # | 问题 | 根因 | 修复 |
|---|------|------|------|
| 1 | admin-service/user-service 重启循环，报 `dial tcp 127.0.0.1:5432: connection refused` | `.env` 中 `DB_HOST=127.0.0.1`，Docker 容器内无法通过 loopback 访问 postgres 容器 | 改为 `DB_HOST=postgres`，`REDIS_HOST=redis` |
| 2 | nginx 启动崩溃，报 `host not found in upstream "forum-service:8002"` | 服务器 `nginx/nginx.conf` 未同步 Git——仍含 `upstream forum_backend`（Ticket 6 的 commit `63f78ad` 已在 Git 中移除，但服务器文件是 volume mount 未更新） | 用 sed 移除 `upstream forum_backend` + 5 个 forum location blocks |
| 3 | nginx 返回 502，上游 IP 不可达 | 容器重建后 IP 变更，nginx 缓存了旧 IP（`proxy_pass` 用 Docker DNS 名称但 nginx 在启动时解析并缓存 IP） | `docker compose restart nginx` |

**重要教训**：
- 服务器 `docker-compose.yml` 对 admin/user-service 使用 `build:` 指令，`docker compose up -d frontend` 会连带重建依赖服务。以后应使用 `--no-deps` 标志。
- nginx 的 `upstream` 块在启动时解析 DNS，容器重建后必须重启 nginx 才能生效。

---

## 3. 服务器当前状态（122.51.233.225）

| 组件 | 状态 | 说明 |
|------|------|------|
| nginx | ? running | 端口 8091/8888 → 80，forum 路由已全部移除 |
| frontend | ? running（新镜像） | AdminOverview 显示「统计功能即将上线」占位页 |
| admin-service | ? healthy | 端口 8003，`/admin-api/health` 返回 200 |
| user-service | ? healthy | 端口 8001 |
| postgres | ? healthy | schema_forum 已 DROP |
| redis | ? healthy | — |
| forum-service | ? 已删除 | 容器已 docker stop + rm |
| Discourse | ? 运行中 | `http://122.51.233.225:8080` |

---

## 4. Ticket 7 详细验证计划

### 4.1 验证目标

确认系统进入「平台 = 登录注册 + 管理后台，Discourse = 社区」的最终形态，所有 Ticket 1-6 的改动完整生效，无回归。

### 4.2 验证清单

#### A. 基础设施健康检查

| # | 检查项 | 方法 | 预期结果 |
|---|--------|------|----------|
| A1 | 6 个容器全部 healthy | `docker compose ps` | nginx / frontend / admin-service / user-service / postgres / redis 全部 Up + healthy |
| A2 | admin-service `/health` | `curl localhost:8888/admin-api/health` | `{"service":"admin-service","status":"ok"}` |
| A3 | user-service `/health` | `curl localhost:8888/user-api/health` | 200 |
| A4 | PostgreSQL 可连接 | `docker exec ai-forum-postgres-1 pg_isready` | accepting connections |
| A5 | Redis 可连接 | `docker exec ai-forum-redis-1 redis-cli ping` | PONG |
| A6 | schema_forum 已删除 | `docker exec ai-forum-postgres-1 psql -U ai_forum -c "SELECT schema_name FROM information_schema.schemata WHERE schema_name='schema_forum'"` | 0 rows |
| A7 | forum-service 容器不存在 | `docker ps -a --filter name=forum-service` | 空 |

#### B. 网关路由验证

| # | 检查项 | 方法 | 预期结果 |
|---|--------|------|----------|
| B1 | 首页正常加载 | 浏览器访问 `http://122.51.233.225:8888/` | Vue SPA 正常渲染（非 502） |
| B2 | 旧 forum-api 路由返回 SPA（非 502） | `curl -I http://122.51.233.225:8888/forum-api/` | 200（前端 SPA 的 index.html 兜底） |
| B3 | 旧 forum 帖子路由不报 502 | `curl http://122.51.233.225:8888/api/v1/posts` | 前端 SPA 兜底或 404 JSON |
| B4 | 旧 forum 板块路由不报 502 | `curl http://122.51.233.225:8888/api/v1/boards` | 同上 |
| B5 | admin-api 路由正常 | `curl localhost:8888/admin-api/api/v1/admin/config` | 401（需 token，非 502/404） |
| B6 | user-api 路由正常 | `curl localhost:8888/user-api/api/v1/auth/login` | 正常响应（非 502） |

#### C. 认证与 SSO 流程

| # | 检查项 | 方法 | 预期结果 |
|---|--------|------|----------|
| C1 | 演示登录 | 浏览器 → 点击演示登录 | 成功登录，跳转首页 |
| C2 | 用户名密码登录 | 用已知账号登录 | 成功登录 |
| C3 | SSO 跳转 Discourse | 已登录状态下点击「论坛」导航 | 无感跳转 Discourse，已登录状态 |
| C4 | Discourse 同步登出 | 平台登出 → 访问 Discourse | Discourse 也处于登出状态 |
| C5 | Token 刷新 | 等待 token 过期或手动测试 | 401 后自动刷新 token 并重试 |

#### D. 管理后台 UI 验证

| # | 检查项 | 方法 | 预期结果 |
|---|--------|------|----------|
| D1 | AdminOverview 占位页 | 管理员登录 → 管理概览 | 显示 ?? 图标 +「统计功能即将上线」+ 说明文字 |
| D2 | AdminConfig 精简 | 管理员 → 系统配置 | 仅显示 posting + moderation 配置，无 boards 相关字段 |
| D3 | AdminStats 页面不存在 | 直接访问 `/admin/stats` | SPA 路由兜底（跳转首页或 404） |
| D4 | 侧边栏无旧论坛菜单 | 查看管理后台侧边栏 | 无 audit / reports / posts / boards / sensitive 菜单项 |
| D5 | 侧边栏「前往 Discourse 管理」 | 点击该链接 | 新窗口打开 `http://122.51.233.225:8080/admin` |
| D6 | 用户管理正常 | 管理员 → 用户管理 | 列表正常加载，角色编辑可用 |
| D7 | 邀请码管理正常 | 管理员 → 邀请码 | 生成/列表/删除正常 |

#### E. 前端代码清理验证

| # | 检查项 | 方法 | 预期结果 |
|---|--------|------|----------|
| E1 | api.js 无 forum 残留 | `rg -i "forum" frontend/src/api.js` | 0 结果（或仅注释中的 forum 字样） |
| E2 | router.js 无 AdminStats 路由 | `rg "AdminStats" frontend/src/router.js` | 0 结果 |
| E3 | 无残留 forum 页面文件 | `ls frontend/src/views/AdminAudit.vue` 等 | 所有 5 个文件不存在 |
| E4 | vite.config.js 无 forum-api 代理 | `rg "forum-api" frontend/vite.config.js` | 0 结果 |

#### F. 后端代码清理验证

| # | 检查项 | 方法 | 预期结果 |
|---|--------|------|----------|
| F1 | admin-service 无 forum_client | `rg "forum_client" services/admin/` | 0 结果 |
| F2 | admin-service 无 forum 相关 handler | `ls services/admin/internal/handler/post_admin_handler.go` 等 | 文件不存在 |
| F3 | services/forum/ 目录不存在 | `ls services/forum/` | 目录不存在或为空 |
| F4 | migrations 无 schema_forum | `rg "schema_forum" migrations/` | 0 结果 |
| F5 | 环境变量无 FORUM_SERVICE_URL | `rg "FORUM_SERVICE_URL" .env .env.example services/admin/.env` | 0 结果 |
| F6 | Docker Compose 无 forum-service | `rg "forum-service" docker-compose*.yml` | 0 结果 |

#### G. 移动端窄屏检查

| # | 检查项 | 方法 | 预期结果 |
|---|--------|------|----------|
| G1 | 管理概览页窄屏 | Chrome DevTools → iPhone SE 模拟 → 访问 AdminOverview | 布局不错乱，占位内容居中可读 |
| G2 | 系统配置页窄屏 | 同上 → AdminConfig | 表单字段不溢出，保存按钮可点击 |
| G3 | 用户管理窄屏 | 同上 → AdminUsers | 表格可横向滚动或自适应 |
| G4 | 侧边栏窄屏 | 同上 → 打开/关闭侧边栏 | 遮罩正常，菜单项可点击，底部「前往 Discourse 管理」可见 |

#### H. 配置文件一致性

| # | 检查项 | 方法 | 预期结果 |
|---|--------|------|----------|
| H1 | 6 个 docker-compose 文件均已移除 forum-service | 逐文件检查 | 全部不含 forum-service |
| H2 | .env / .env.example 等 7 个文件无 FORUM_SERVICE_URL | 逐文件检查 | 全部不含 |
| H3 | 服务器 .env 与 Git 版本一致 | diff 比较 | 仅 DB_HOST/REDIS_HOST 差异（有意为之） |

### 4.3 执行顺序

```
1. 基础设施健康检查（A1-A7）                          ← 快速扫一遍，5 分钟
2. 网关路由验证（B1-B6）                               ← curl 即可，3 分钟
3. 后端代码清理验证（F1-F6）                           ← 本地 rg 扫描，3 分钟
4. 前端代码清理验证（E1-E4）                           ← 本地 rg 扫描，2 分钟
5. 配置文件一致性（H1-H3）                             ← 本地扫描，3 分钟
6. 认证与 SSO 流程（C1-C5）                            ← 浏览器手动，10 分钟
7. 管理后台 UI 验证（D1-D7）                           ← 浏览器手动，10 分钟
8. 移动端窄屏检查（G1-G4）                             ← DevTools 模拟，5 分钟
```

### 4.4 预期总耗时

约 40 分钟（代码扫描 10 分钟 + 浏览器手动验证 20 分钟 + 移动端 5 分钟 + 文档 5 分钟）。

---

## 5. Git 状态

- **分支**：`codex/discourse-rebuild`
- **HEAD**：`15d172a`（fix: upgrade Node to 22-alpine and set UTF-8 locale in frontend Dockerfile (#6)）
- **远程**：已 push
- **最近两次改动**：
  1. `da3af64`：前端文件 UTF-8 编码修正
  2. `15d172a`：Dockerfile Node 22 + UTF-8 locale

---

## 6. 已知约束

- **服务器无法访问 GitHub/Docker Hub**：部署走本地 build → SCP tar → docker load → docker compose up
- **docker compose up 会级联重建**：后续部署用 `--no-deps` 避免重建依赖服务
- **nginx 缓存 DNS**：容器重建后必须 `docker compose restart nginx`
- **SSH**：`liyongquan@122.51.233.225`，密码 `Liyongquan@123`
- **Discourse**：`http://122.51.233.225:8080`
- **项目路径**：服务器 `/home/liyongquan/projects/ai-forum`
- **每次 commit 后必须 push**，Ticket 7 如产生改动同样遵守

---

## 7. 新线程启动指令

> 读取 `docs/handoffs/handoff-2026-07-30-discourse-ticket6-done-ticket7-plan.md`，了解现状后执行 Ticket 7 验证与冒烟测试。优先跑 A 组（基础设施）→ B 组（网关）→ E/F/H（代码扫描），发现问题即时修复；C/D/G 组（浏览器手动）可在代码扫描通过后集中进行。