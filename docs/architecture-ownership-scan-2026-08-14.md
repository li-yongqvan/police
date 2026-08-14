# 架构"无主共享物"扫描报告（2026-08-14）

> 目标：找出项目中**没有明确所有者、但多个模块都在引用**的"东西"——即跨模块共享却无人单一负责的概念/契约/数据，避免其成为漂移与事故的温床。
> 方法：静态扫描引用面（代码 grep + 迁移文件 + compose 配置 + 契约文档），并对照所有权声明（schema 归属、包归属、注释/ADR）。词汇借用 `improve-codebase-architecture` skill：**owner（所有者）、seam（接缝）、locality（局部性）、leverage（杠杆）**。
> 判定标准：引用方 ≥ 2 个模块，且仓库内不存在唯一权威定义/所有者说明。

---

## 1. 发现总览

| # | 共享物 | 引用方 | 无主表现 | 建议所有者 | 优先级 |
|---|--------|--------|----------|------------|--------|
| O-1 | JWT 令牌契约 | user-service（签发）、admin-service（校验）、前端 session | `pkg/jwt` 双份复制且已分叉（user 有 `level` claim，admin 没有）；无契约文档/ADR | user-service（BC-Identity，唯一签发方） | 高 |
| O-2 | 角色与等级授权语义 | user-service、admin-service、JWT claims、前端 Admin*、`system_config`、迁移 | 角色数据在 `schema_admin`，但 user-service 登录时**跨 schema 直连**读取；`role:*` Redis 缓存只写不读；`level` 语义在 `api-contract.md` 中为空白 | admin-service 为角色唯一权威源，user-service 经内部 API 消费 | 高 |
| O-3 | 演示账号供给 | user 迁移 004/007、admin 迁移 004/005 | 账号与角色分属两个服务的迁移链，靠"碰巧的顺序"工作，依赖无声明 | 单一种子/供给模块（见建议） | 中 |
| O-4 | 共享 `.env`（厨房水槽） | 两个服务、compose、前端构建 | 一个 `.env` 整份注入所有容器：`FORUM_SERVICE_URL`、`ADMIN_SERVICE_URL` 无人引用；`DISCOURSE_*` 只有 user 用；`DB_PASSWORD` 是超级用户口令 | 每个服务独立 env 文件 + 显式共享清单 | 中 |
| O-5 | 审计日志（operation_logs）概念 | user-service、admin-service、两套 schema | 同构双表（`schema_auth` 66 行 / `schema_admin` 76 行），各写各的，同一用户目标的审计动作被劈成两半 | admin-service（治理）统一收口 | 中 |
| O-6 | 敏感词治理能力 | admin-service（CRUD + 缓存 + 内部检测端点）、`system_config` | forum-service 移除后**没有任何消费者**（Discourse 独立治理），能力悬空 | 决策：下线 or 明确 Discourse 侧消费者 | 中 |
| O-7 | `/internal/v1` 无鉴权边界 | user-service（暴露）、admin-service（消费） | 代码注释引用不存在的 ADR "D-09"；保护靠"nginx 不路由 + 端口不发布"的配置缺席 | 主控/部署侧，补 ADR + 网络隔离 | 中 |
| O-8 | 字段契约（api-contract.md ↔ api.js） | 前端 8 个文件、契约文档、历史 agent | `api.js` 是前端唯一 seam（好局部性），但契约文档无人校验、已滞后（`level` 行空白、`post_requires_level` 是 forum 时代映射） | 集成契约 agent / 生成器或测试 | 低 |

---

## 2. 逐项证据

### O-1 JWT 令牌契约（格式 + 密钥）

- 双份实现已分叉：`services/user/pkg/jwt/jwt.go:15` 的 Claims 含 `Level int`，`services/admin/pkg/jwt/jwt.go:12` 的 Claims **没有** Level；两者各自 `ValidateToken`。
- user-service 签发（`user_service.go:150,203,348`），admin-service 校验并信任 `claims["role"]`（`admin/internal/middleware/auth.go:37-41`）。
- 密钥靠一份 `.env` 的 `JWT_SECRET` 注入两个容器，轮换时无单一负责方；claim 变更会静默影响另一个服务。
- 仓库中没有任何文档/ADR 描述 token 格式与 claim 语义（`shared/api-contract.md` 未覆盖）。

### O-2 角色与等级授权语义

- **跨 schema 直连**（违反项目自己的 6.3 规则）：`services/user/internal/service/user_service.go:44-63` 的 `resolveJWTReole`（函数名还拼错了 Reole）直接 `SELECT ... FROM schema_admin.user_roles JOIN schema_admin.roles`。之所以能跑通，是因为两个服务共用超级用户 `ai_forum`（见数据库审查报告 P1-4）。
- **死缓存**：`services/admin/internal/service/role_service.go:73,100` 写/删 `role:%d` Redis 键（TTL=0 永不过期），全仓库**无任何读者**。
- **双轨授权**：`users.level`（int，值域无定义：迁移种 1/2、库中实际只有 2/5、前端 AdminUsers 输入 0–5、`system_config.upload_requires_level=2`）与 roles（admin/platform_admin 名称）并存；`shared/api-contract.md:14` 的 level 映射值为空——语义没有唯一权威源。
- 引用面：JWT claim（user）、admin 权限校验、前端 `AdminUsers.vue:112`、`api.js:323`、`system_config`。

### O-3 演示账号供给

- 账号创建在 user 迁移：`004_demo_accounts.up.sql`、`007_demo_batch_accounts.up.sql`（INSERT users）。
- 角色授予在 admin 迁移：`004_demo_user_roles.up.sql`、`005_demo_batch_roles.up.sql` —— 从 **schema_admin 的迁移里跨 schema 读 `schema_auth.users`** 按 username 找 user_id。
- 两个服务的迁移链没有任何顺序/依赖声明；新环境能否正确长出演示账号取决于部署脚本的执行顺序（碰巧对）。

### O-4 共享 `.env`（厨房水槽）

- `docker-compose.yml:22-25,42-45`：两个服务 `env_file: .env`，所有变量整份注入。
- 无人引用的死变量：`FORUM_SERVICE_URL`（全仓库 0 处 Go 引用）、`ADMIN_SERVICE_URL`（0 处）。
- 只被单方使用的变量也注入了双方：`DISCOURSE_URL/DISCOURSE_API_KEY/DISCOURSE_CONNECT_SECRET`（仅 user 用）、`USER_SERVICE_URL`（仅 admin 用）、`PLATFORM_BASE_URL`（仅 user 用）。
- 最敏感的是 `DB_PASSWORD`（超级用户口令）同时存在于两个服务的内存环境里，任一服务被攻破即得整库。

### O-5 审计日志概念

- `schema_auth.operation_logs`（66 行）与 `schema_admin.operation_logs`（76 行）结构同构、无外键、各自服务各写各的（`user_admin_service.go:87,116,194,216` vs `admin` 的 `config_service.go:72`、`role_service.go:80`、`user_admin_service.go:34`）。
- 例如"封禁用户"这类同一目标的治理动作，一部分落在 auth、一部分落在 admin，审计视图必须拼接两个 schema——"审计"这个概念没有所有者。

### O-6 敏感词治理能力（有生产者无消费者）

- admin-service 维护完整能力：CRUD（`main.go:106-108`）、内存缓存（`admin_service.go:27-47`）、内部检测端点 `/internal/v1/moderation/check`（`main.go:113`、`moderation_handler.go:43`）。
- 曾经的消费者 forum-service 已移除；Discourse 有独立治理。当前**无人调用**检测端点 → 能力悬空，`sensitive_words` 表成为无消费者的数据。
- 同类：`statistics_daily` 空表无代码写入、`system_config` 11 个 forum 时代键（见数据库审查报告 P3-1/P3-2）。

### O-7 `/internal/v1` 无鉴权边界

- user-service 暴露内部端点（`cmd/main.go:124` 起：users/status、invite-codes、ban/unban、level、stats……），注释写着 "no external auth per D-09"，但 `docs/adr/` 里只有 0001/0002（Discourse CSS），**D-09 并不存在**。
- admin-service 通过 `user_client.go` 消费（invite-codes/ban/unban/users/level/logs）。
- 保护机制 = nginx 不路由 `/internal/v1`（`nginx.conf:136-138` 注释）+ 容器端口不发布——边界由"配置缺席"构成，无人声明、无人测试；将来任何代理规则调整都可能瞬间把它暴露到公网。

### O-8 字段契约（api-contract.md ↔ api.js）

- `frontend/src/api.js` 是前端唯一 API seam（8 个消费文件：main.js、session store、AdminConfig/AdminInvites/AdminUsers/AdminRoles/Register 等），snake_case→camelCase 映射集中于此——这是**好局部性**的正面案例。
- 但契约文档 `shared/api-contract.md` 与实现无任何自动校验：`level` 行空白、`post_requires_level` 是 forum 时代映射（`api.js:280,289` 还在用它）、路由表仍含 forum 端点。文档自身也是"无主共享物"，靠人工自觉。

---

## 3. 修复建议（按优先级）

1. **O-2 角色语义收口**：删除 `resolveJWTReole` 的跨 schema 直连，改为 user-service 经内部 API 向 admin 查询角色（复用现有 `user_client.go` 模式）；删除 `role:*` 死缓存；为 `level` 建立唯一权威定义（枚举 + api-contract.md 填值）。
2. **O-1 令牌契约锁定**：确定 user-service 为唯一签发方，admin 仅消费；把 Claims 结构收敛（建议共享类型或生成契约），补一份 token 契约说明（含 claim 清单、轮换流程）。
3. **O-4 拆 env**：user.env / admin.env 分离，共享密钥显式列举（JWT_SECRET、REDIS/DB 口令按需最小化）；删除死变量 `FORUM_SERVICE_URL`、`ADMIN_SERVICE_URL`。
4. **O-7 内部边界显式化**：补 ADR（记录"内部 API 无鉴权 + 靠网络层隔离"的决策与假设）；nginx 显式 `location /internal/v1 { return 404; }` 化被动为主动。
5. **O-3 演示账号供给统一**：合并为单一 seed 模块（如 `migrations/init/` 或独立脚本），声明唯一 owner。
6. **O-5/O-6 治理能力盘点**：审计日志二选一收口（admin 统一 + user 上报）；敏感词能力与 forum 时代 config 键明确下线或指定新消费者。
7. **O-8 契约可校验**：api-contract.md 改为生成或由集成测试校验，指派"集成契约 agent"为 owner。

> 以上均为计划草案，按 AGENTS.md 规则需你确认后再执行；其中 O-2 与数据库审查的 DB-FIX-03（降权）强相关，建议合并规划。

---

## 4. 后续选项

- 本报告是 `improve-codebase-architecture` skill 扫描阶段的书面版。如需，我可以把 8 个候选做成该 skill 的可视化 HTML 报告（before/after 图），或对其中任意一项（建议 O-2）进入它的 grilling 决策流程。