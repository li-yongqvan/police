# Handoff · 2026-07-30 · Discourse SSO 集成 — Ticket 3a 完成，准备 Ticket 3b

> 新会话建议先读本文件，然后说："读取 .scratch/discourse-forum-rebuild/spec.md 的 Ticket 3b，按 goal.md 的 Execution order 逐步执行。每完成一个 Step 报告并等待确认。"

---

## 1. 现在做了什么

- 正在执行 **Ticket 3a: SSO Provider Core + Discourse Configuration**，已全部完成并通过验证。
- 用户已验证浏览器端 SSO 流程可用。
- 当前分支：`master`，HEAD：`cf7f860`。
- **用户指令**：继续 Ticket 3b（Logout Synchronization），但先生成此 handoff 后切换对话。

---

## 2. 已经完成了什么

### Ticket 3a 完整执行（5 个 Step 全部完成）

| Step | 内容 | 文件 | 状态 |
|---|---|---|---|
| 1 | 添加环境变量 | `services/user/.env`（追加 `DISCOURSE_CONNECT_SECRET`、`PLATFORM_BASE_URL`），新建 `services/user/.env.example` | ✅ |
| 2 | 创建 SSO handler | `services/user/internal/handler/discourse_sso_handler.go`（NEW）— 完整 DiscourseConnect 协议：HMAC-SHA256 签名校验、JWT cookie 认证、用户查询与封禁检查、角色映射（platform_admin → admin、admin → moderator）、302 重定向 | ✅ |
| 3 | 注册路由 | `services/user/cmd/main.go`：第 72 行添加 `discourseSSOHandler` 实例化，第 99 行添加 `v1.GET("/auth/discourse/sso", ...)` | ✅ |
| 4 | 配置 Discourse | 云服务器 `122.51.233.225` Discourse 容器内执行 Rails runner，设置 `enable_discourse_connect`、`discourse_connect_url=http://122.51.233.225:8888/api/v1/auth/discourse/sso`、`discourse_connect_secret`、overrides（groups/bio/avatar） | ✅ |
| 5 | 验证 | 无效签名 → 403 ✅；无 cookie → 302 `/login` ✅；有效 JWT → 302 Discourse SSO login ✅；platform_admin → admin=true ✅ | ✅ |

### 部署状态

- 代码已通过 scp 同步至 `122.51.233.225:/home/liyongquan/projects/ai-forum/`
- `ai-forum-user-service-1` 容器已重建并健康运行
- Discourse `app` 容器 SSO 配置已生效
- SSO URL：`http://122.51.233.225:8888/api/v1/auth/discourse/sso`（通过 Nginx → user-service:8001）

### 关键 URL

| 服务 | 地址 |
|---|---|
| 平台前端 | `http://122.51.233.225:8888` |
| Discourse 论坛 | `http://122.51.233.225:8080` |
| user-service（内部） | `http://127.0.0.1:8001` |
| SSO 端点 | `GET /api/v1/auth/discourse/sso`（需 sso+sig query params） |

### 共享密钥

- `DISCOURSE_CONNECT_SECRET` = `shared-secret-change-in-production`（服务器 `.env` 和 Discourse 配置中一致）
- JWT Secret fallback：`default-secret-change-in-production`

---

## 3. 卡在哪里

暂无明确阻塞。Ticket 3a 全部完成。

---

## 4. 下一步做什么

### 立即：Ticket 3b — Logout Synchronization

按 spec.md 定义：

> 当用户从平台登出时，调用 Discourse API `POST /admin/users/{discourse_user_id}/log_out`，使 Discourse 端也登出。

**涉及文件：**
- `services/user/internal/handler/auth_handler.go`（**READ FIRST** — 需确认现有登出 handler 位置和逻辑）
- `services/user/.env`（追加 `DISCOURSE_API_KEY`）
- 云服务器根目录 `.env`（同步追加 `DISCOURSE_API_KEY`）

**完成标准：**
- 平台登出时调用 Discourse API 同步登出
- `DISCOURSE_API_KEY` 从环境变量读取
- Discourse API 调用失败时记录日志但不影响本地登出
- 本地 cookie 正常清除

**依赖：** Ticket 3a（已完成）。需要先在 Discourse 管理后台生成一个 Admin API Key。

### 后续 Tickets

| 顺序 | Ticket | 说明 |
|---|---|---|
| 3b | Logout Sync | 当前即将执行 |
| 4 | Frontend Nav Cutover | 前端路由/导航切到 Discourse |
| 5 | Admin Menu Transition | 管理后台菜单调整 |
| 6 | Old Forum Disposal | 15 天稳定期后清理旧论坛 |
| 7 | Smoke Tests | 全面验收 |

---

## 5. 哪些坑不要再踩

### PowerShell 坑
- **`&` 符号在 PowerShell 中会被解析为后台操作符**，传给 SSH 的 bash 命令中 `&` 必须做特殊处理。解决方案：把脚本写成文件上传到服务器再执行，不要内联在 SSH 命令中。
- **heredoc（<<）在 PowerShell 中不兼容**，会报 "Missing file specification after redirection operator"。

### Docker 部署坑
- **Docker Hub 被墙**（服务器在中国大陆），需要使用 DaoCloud 镜像源。构建命令必须加 `-f docker-compose.cn-mirror.yml`。
- **`env_file` 路径是相对于 docker-compose.yml**。环境变量应加到项目根目录 `C:\Users\liyongquan\Documents\New project 6\.env`（不是 `services/user/.env`），因为 `docker-compose.yml` 引用 `env_file: - .env`。
- 修改 `.env` 后必须 `docker compose up -d` 重建容器，否则环境变量不生效。
- 正确的构建+重启命令：
  ```bash
  cd /home/liyongquan/projects/ai-forum
  docker compose -f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml build user-service
  docker compose -f docker-compose.yml -f docker-compose.server.yml up -d user-service
  ```

### Discourse 配置坑
- Discourse `enable_discourse_connect` 必须在 `discourse_connect_url` 设置之后才能启用，否则报错 "You must set a 'discourse connect url' before enabling this setting"。
- 配置脚本通过 `docker exec app rails runner` 执行，需要先把 Ruby 脚本 cp 到容器内再执行。
- 如果 `docker exec` 报 unicode/font 问题，用 `-i` 传 stdin 可能会因 heredoc 不兼容而失败，最可靠的方式是 scp 脚本 → docker cp → docker exec。

### Go 代码坑
- **未使用的 import 会导致编译失败**。`discourse_sso_handler.go` 最初 import 了 `model` 包但因未显式引用被移除。Go 编译器只检查显式引用的 import，通过函数返回值类型隐式使用的不算。
- JWT secret 的 fallback 值必须与 `AuthMiddleware` 一致：`"default-secret-change-in-production"`。

### SSO 协议坑
- Discourse SSO payload 需要 URL-encode → Base64 编码。返回的 Base64 字符串在 302 redirect URL 中需 `url.QueryEscape`（避免 `+`/`/`/`=` 被误解析）。
- HMAC 比较必须用 `hmac.Equal`（常数时间比较防时序攻击），不能直接 `==` 字符串比较。

### 验证测试
- 生成测试用 JWT 可用服务器上的 Python + PyJWT 库（`pip3 install pyjwt`）。
- 测试 SSO 的 curl 脚本模板（上传到服务器执行，不要内联在 PowerShell）：
  ```bash
  #!/bin/bash
  JWT='your_jwt_token'
  SSO='base64_encoded_payload'
  SIG='hex_signature'
  curl -s -v --cookie "ai-forum-token=${JWT}" "http://127.0.0.1:8001/api/v1/auth/discourse/sso?sso=${SSO}&sig=${SIG}"
  ```