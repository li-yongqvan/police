# Handoff · Discourse 重建 — Ticket 2 完成 + Ticket 3 盘问结束

> 日期：2026-07-29 | 服务器：122.51.233.225 | 状态：✅ Ticket 2 完成，🔄 Ticket 3 待进入 to-spec

---

## 1. 已完成工作总览

| Ticket | 状态 | Handoff |
|--------|------|---------|
| Ticket 1: Production Baseline Lock | ✅ 完成 | docs/handoffs/handoff-2026-07-29-discourse-baseline-lock.md |
| Ticket 2: Discourse PoC Deployment | ✅ 完成 | 本文档 |
| Ticket 3: DiscourseConnect SSO Integration | 🔄 盘问结束 | 本文档（待新线程进入 to-spec → to-goal → 执行） |

---

## 2. Ticket 2 执行摘要

### 部署结果

- **Discourse URL**：http://122.51.233.225:8080
- **管理员账号**：dmin / AIForum2026!Admin / Lyq1841777060@outlook.com
- **5 个默认分类**：AI 前沿、技术讨论、学习资源、项目展示、闲聊灌水
- **站点设置**：公开浏览已启用、发帖需登录、force_https 已禁用
- **数据库**：Discourse 使用容器内独立 PostgreSQL，未连接现有 schema_*
- **旧平台**：7 个容器（user-service, forum-service, admin-service, frontend, nginx, postgres, redis）正常运行

### 站点配置详情

| 设置项 | 值 |
|--------|-----|
| login_required | disabled（公开浏览） |
| must approve users | disabled |
| min trust to create topic | 0 |
| 	itle | AI 智联论坛 |
| site description | AI 主题交流社区 |
| orce https | disabled |
| contact email | Lyq1841777060@outlook.com |

### 服务器访问

- SSH：ssh liyongquan@122.51.233.225（key 已配置）
- Discourse 容器操作：cd /home/liyongquan/discourse && ./launcher enter app
- Rails console：docker exec app rails runner "..."（批量配置用此方式最快）
- 配置文件：/home/liyongquan/discourse/containers/app.yml
- AI 论坛：http://122.51.233.225:8888
- 项目路径：/home/liyongquan/projects/ai-forum

### 经验教训：Docker 镜像快速上传

> 以后需要用 docker save + scp 管道直传，避免先 save 到磁盘再 scp（磁盘 I/O 瓶颈）：
> `ash
> docker save discourse/base:release | ssh liyongquan@122.51.233.225 'docker load'
> `
> 或使用 skopeo 在服务器上直接从 registry 拉取（不占本地磁盘）。

---

## 3. Ticket 3 盘问结果（3 问全部解答）

盘问基于 spec.md 中已锁定的 10 项决策进行深化，针对 SSO 实现细节。

| # | 问题 | 决定 |
|---|------|------|
| 1 | 未登录用户访问 Discourse SSO 端点时如何处理？ | **A**：重定向到平台登录页，登录后回跳 Discourse |
| 2 | 平台退出时 Discourse 会话如何同步？ | **A**：调用 Discourse API POST /admin/users/{id}/log_out 同步登出 |
| 3 | SSO 共享密钥存放位置？ | **A**：环境变量 DISCOURSE_CONNECT_SECRET（不与 system_config 表耦合） |

### 决策理由（Q3 选 A）

- 环境变量是 12-factor app 标准做法，Go 服务已有 os.Getenv() 模式
- system_config 表属于 admin-service 的数据主权，user-service 不应反向依赖
- 密钥变更只需重启 user-service 容器，比改数据库记录更可控
- Discourse 侧 discourse_connect_secret 本身就是容器环境变量设置，两边对称

---

## 4. 关键代码上下文（Ticket 3 涉及文件）

### 文件清单

| 文件 | 作用 | 改动类型 |
|------|------|----------|
| services/user/cmd/main.go | 路由注册 | 新增 SSO 路由（/api/v1/auth/discourse/sso） |
| services/user/internal/handler/discourse_sso_handler.go | SSO handler | **新建文件** |
| services/user/internal/service/user_service.go:44-62 |esolveJWTReole() — 角色映射 | 读取参考，不改动 |
| services/user/internal/middleware/auth.go | JWT 认证中间件 | 读取参考（SSO 路由不走 AuthMiddleware，但需从 cookie 读 JWT） |
| services/user/internal/model/user.go | User 模型 | 读取参考（ID, Username, Nickname, Bio, Avatar, Status） |
| services/user/pkg/jwt/jwt.go | JWT 工具 | 读取参考（HMAC-SHA256 签名） |
| services/user/.env | 环境变量 | 新增 DISCOURSE_CONNECT_SECRET |

### 角色映射逻辑（已存在，不改动）

`
resolveJWTReole() 在 user_service.go:44-62：
  schema_admin.user_roles JOIN schema_admin.roles
  → platform_admin → Discourse admin
  → admin          → Discourse moderator
  → 无记录         → "student" → Discourse regular user
`

### 用户模型关键字段

`
User.ID       uint    → Discourse external_id
User.Username string  → Discourse username
User.Nickname string  → Discourse name
User.Avatar   string  → 相对路径，需构造完整 URL
User.Bio      string  → Discourse bio
User.Status   string  → "active" / "banned"（banned 阻止 SSO）
`

### 认证中间件

- AuthMiddleware() 从 Authorization: Bearer <token> 头解析 JWT
- SSO handler 不走这个中间件，而是从 cookie i-forum-token 读取 JWT（因为 Discourse 重定向过来时无法携带自定义 Header）
- JWT 解析后通过 c.Set("user_id", ...) 注入上下文

---

## 5. Ticket 3 实现要点预览（仅供参考，由 to-spec 阶段确认）

### 需新增的端点

`
GET /api/v1/auth/discourse/sso
  - 查询参数：sso（Base64 编码的 payload）、sig（HMAC-SHA256 签名）
  - 认证方式：从 cookie 读 JWT（不经过 AuthMiddleware）
  - 已验证：验证 HMAC-SHA256 签名 → 解析 payload → 查用户 → 构建返回 payload → 签名 → 重定向回 Discourse
  - 未验证：重定向到平台登录页（/login?redirect=/api/v1/auth/discourse/sso?...）
  - Banned 用户：返回错误，不生成 SSO payload
`

### 需配置的 Discourse 设置

`
enable_discourse_connect        → true
discourse_connect_url           → http://122.51.233.225:8888/api/v1/auth/discourse/sso
discourse_connect_secret        → <与 user-service 环境变量一致>
`

### 登出同步

`
平台登出时 → POST /admin/users/{discourse_user_id}/log_out
通过 Discourse API key 认证
`

---

## 6. 待完成：Ticket 3 后续流程

按铁律三，盘问（阶段 1）已结束。下一线程应走：

1. **阶段 2（to-spec）**：将盘问结果沉淀到 .scratch/discourse-forum-rebuild/spec.md 的 Ticket 3 章节
   - 更新 spec.md，补充 Ticket 3 的实现细节（路由、handler 逻辑、SSO payload 字段、错误处理）
   - 明确约束：不改动esolveJWTReole()、不改动 schema_admin 表、不改动 AuthMiddleware
2. **阶段 3（to-tickets）**：如需可进一步拆分 Ticket 3 为子 ticket（如 SSO 端点 + 登出同步 + 前端登录流程调整）
3. **阶段 4（to-goal）**：编译 .scratch/discourse-forum-rebuild/goal.md（Ticket 3 版本）
4. **阶段 5（新会话执行）**：读 goal.md 逐步执行

---

## 7. 新线程启动指令

> 读取 .scratch/discourse-forum-rebuild/spec.md 和 docs/handoffs/handoff-2026-07-29-discourse-ticket2-done-ticket3-grill.md，盘问已结束（3 问全部解答），进入 to-spec 阶段，将 Ticket 3 的 SSO 实现细节沉淀到 spec.md 中。

### 关键约束提醒

- 不改动 schema_forum、schema_admin 现有表结构
- 不删除 user-service 现有路由
- 不修改前端 auth store 的 token 存储逻辑
- DiscourseConnect SSO 协议：HMAC-SHA256 签名
- 不改动 login_required/orce_https 等已配置的 Discourse 设置
- 不 push、PR、merge、关闭 issue
- 不提前实现下游 ticket（Ticket 4/5/6/7）

---

## 8. 项目 Git 状态

- **分支**：master
- **HEAD**：cf7f860（docs: add production baseline snapshot 2026-07-29）
- **Tag**：pre-discourse-baseline → 05a575
- **工作区**：干净（服务器侧）
