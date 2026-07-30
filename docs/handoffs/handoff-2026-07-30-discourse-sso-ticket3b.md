# Handoff · 2026-07-30 · Discourse SSO 集成 — Ticket 3b 完成，准备 Ticket 4

> 新会话建议先读本文件，然后说："读取 .scratch/discourse-forum-rebuild/spec.md 的 Ticket 4，按 goal.md 的 Execution order 逐步执行。每完成一个 Step 报告并等待确认。"

---

## 1. 现在做了什么

- 正在执行 **Discourse Forum Rebuild** 系列 tickets。
- **Ticket 3a**（SSO Provider Core）和 **Ticket 3b**（Logout Synchronization）已全部完成并通过验证。
- 当前分支：`codex/discourse-rebuild`，HEAD：`8806550`。
- 代码已提交并 push 到 GitHub，已部署到云服务器 `122.51.233.225`。
- **用户指令**：生成 handoff 后切换线程推进 Ticket 4。

---

## 2. 已经完成了什么

### Ticket 3a: SSO Provider Core + Discourse Configuration ✅

- `services/user/internal/handler/discourse_sso_handler.go` — 完整 DiscourseConnect 协议实现
- `services/user/cmd/main.go` — 注册 `GET /api/v1/auth/discourse/sso`
- Discourse 配置：`enable_discourse_connect`、SSO URL、secret、overrides

### Ticket 3b: Logout Synchronization ✅

| 文件 | 操作 | 说明 |
|------|------|------|
| `services/user/internal/handler/auth_handler.go` | **新建** | Logout handler，通过 `/u/{username}.json` 公开接口查找 Discourse 用户，调 `/admin/users/{id}/log_out` 同步登出。goroutine fire-and-forget 模式 |
| `services/user/cmd/main.go` | 修改 | 注册 `POST /api/v1/auth/logout`（AuthMiddleware 保护） |
| `services/user/.env.example` | 追加 | `DISCOURSE_URL` + `DISCOURSE_API_KEY` |
| `frontend/src/api.js` | 修改 | `userApi` 新增 `logout()` 方法 |
| `frontend/src/stores/session.js` | 修改 | `logout()` 插入 `userApi.logout().catch(() => {})` fire-and-forget |
| 服务器 `.env` | 追加 | `DISCOURSE_URL=http://122.51.233.225:8080` + `DISCOURSE_API_KEY`

**核心设计决策**：
- 前后端均 fire-and-forget，不阻塞用户体验
- 后端 handler 先查用户 username → 调 Discourse `/u/{username}.json`（无需 API Key）→ 获取 Discourse 用户 ID → 调 `/admin/users/{id}/log_out`（需 API Key）
- API Key 需要「Users: log out」+「Users: read」两个 scope
- 用户身份通过 `c.Get("user_id")` 从 AuthMiddleware 获取

---

## 3. 部署状态

- 所有代码已通过 scp 同步至 `122.51.233.225:/home/liyongquan/projects/ai-forum/`
- `ai-forum-user-service-1` 和 `ai-forum-frontend-1` 容器已重建并健康运行
- 分支 `codex/discourse-rebuild` 已 push 到 GitHub

### 关键 URL

| 服务 | 地址 |
|------|------|
| 平台前端 | `http://122.51.233.225:8888` |
| Discourse 论坛 | `http://122.51.233.225:8080` |
| user-service（内部） | `http://127.0.0.1:8001` |
| SSO 端点 | `GET /api/v1/auth/discourse/sso`（需 sso+sig query params） |
| Logout 端点 | `POST /api/v1/auth/logout`（需 Bearer token） |

### 环境变量（服务器 .env）

```
DISCOURSE_CONNECT_SECRET=shared-secret-change-in-production
PLATFORM_BASE_URL=http://122.51.233.225:8888
DISCOURSE_URL=http://122.51.233.225:8080
DISCOURSE_API_KEY=27014a794b4f4157ca8791ab92df42612d44b37b843cd7bacef645cc11c8a912
JWT_SECRET=dev-secret-change-in-production
```

---

## 4. 卡在哪里

无阻塞。Ticket 3a 和 3b 全部完成。

---

## 5. 下一步做什么

### 立即：Ticket 4 — Frontend Forum Navigation Cutover

按 spec.md 定义：

- 所有论坛入口（community links、board links、发帖按钮、首页 hero CTA）指向 Discourse
- 旧论坛页面（BoardView、NewPost、EditPost、PostDetail 等）从导航和路由中移除
- 旧论坛写操作（发帖、编辑、删除、评论、点赞、收藏、举报、附件上传）在 forum-service 服务端禁用
- 现有非论坛页面（portal、login、profile、about、AI pages、admin）保持功能
- 旧论坛不可被任何用户访问
- 论坛相关通知徽章和"消息"隐藏或用 Discourse 指示器替换

**涉及文件（参考）：**
- `frontend/src/router.js`
- `frontend/src/components/gx/GxSiteHeader.vue`
- `frontend/src/components/gx/GxSiteSidebar.vue`
- `frontend/src/views/` 下各论坛相关页面
- `services/forum/` 相关 handler（禁用写操作）

**依赖**：Ticket 3a + 3b（已完成）

### 后续 Tickets

| 顺序 | Ticket | 说明 |
|------|--------|------|
| 4 | Frontend Nav Cutover | 前端路由/导航切到 Discourse |
| 5 | Admin Menu Transition | 管理后台菜单调整 |
| 6 | Old Forum Disposal | 15 天稳定期后清理旧论坛 |
| 7 | Smoke Tests | 全面验收 |

---

## 6. 哪些坑不要再踩

### PowerShell 坑
- **反引号在双引号字符串中是转义符**，写 Go 代码的 struct tag 时如果用双引号 heredoc（`@"..."@`），反引号会被吃掉。解决方案：用 `Add-Content` + 单引号字符串（单引号中反引号是字面量）。
- **`&` 和 `(` 在 SSH 命令中需转义**，尽量把脚本写成文件上传再执行。
- **`%{http_code}` 在 curl 中**，如果外层用了双引号会被 PowerShell 展开。用单引号包裹或写文件执行。

### Go 代码坑
- **未使用的 import 导致编译失败**。新增 handler 时注意检查。
- **中间件 `c.Get("user_id")` 返回 `interface{}`**，需要 `.(uint)` 类型断言。

### Docker 部署坑
- **Docker Hub 被墙**，使用 DaoCloud 镜像源。构建命令：`docker compose -f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml build <service>`
- **`env_file` 路径相对于 docker-compose.yml**。环境变量加到项目根目录 `.env`。
- 修改 `.env` 后必须 `docker compose up -d` 重建容器。

### Discourse API 坑
- `/u/by-external/{id}.json` 和 `/users/by-external/{id}.json` 均需要额外权限（API Key scope 不足时会返回 403）
- **正确做法**：用公开接口 `/u/{username}.json` 获取 Discourse 用户 ID（无需 API Key），拿到 ID 后再用 API Key 调 `/admin/users/{id}/log_out`
- API Key 需要 **Users: log out** + **Users: read** 两个 scope
- Discourse 管理员用户名为空时用 `Api-Username: system`

### 前端构建坑
- 修改 `api.js` 时注意**中文字符编码**。`.NET WriteAllText` 默认写 BOM 可能破坏编码。用 `new UTF8Encoding(false)` 避免 BOM。
- 文件替换时确认编码一致，否则会破坏多字节字符。