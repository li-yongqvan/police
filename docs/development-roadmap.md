# AI 智联论坛 · 四阶段开发计划

> 基于 XenForo 产品分层（内容模型 → 身份权限 → 治理中台 → 持久化运维）拆解。  
> 原则：**演示线不断档**（`frontend-standalone` 保留），**API 契约稳定**（`/api/v1` 不大改）。

---

## 当前基线

| 维度 | 现状 |
|------|------|
| 后端 | 3 个单文件 Go 服务，JSON 存储，内存 Session |
| 前端 | Vue 3 + Vite，PC 演示 + Android WebView |
| 路由规模 | user 6 / forum 11 / admin 10 |
| 缺口 | 无真实注册、无 DB、无 Mod Log、无附件上传、点赞不可取消 |

---

## Week 1 · Foundation（数据库 + JWT + 权限）

**目标**：后端换「心脏」，演示动线不变。

### Issue 清单

- [ ] **W1-01** 添加 `infra/docker-compose.yml`（PostgreSQL 16）
- [ ] **W1-02** 编写 `infra/postgres/init.sql`（users / boards / posts / comments / audit_records / system_config）
- [ ] **W1-03** 种子数据迁移（3 角色账号 + 3 板块 + 现有帖子/评论 JSON 导入）
- [ ] **W1-04** user-service 分层：`handler → service → repository`
- [ ] **W1-05** 实现 `POST /api/v1/register`（邮箱 + 密码 + 昵称）
- [ ] **W1-06** 实现 `POST /api/v1/login`（返回 access_token + user）
- [ ] **W1-07** JWT 中间件替换内存 Session；保留 `POST /api/v1/demo-login` 兼容演示
- [ ] **W1-08** `GET /api/v1/users/me` 从 DB 读取
- [ ] **W1-09** 权限中间件 `RequireAuth` / `RequireStaff`（admin + platform_admin）
- [ ] **W1-10** forum-service 帖子/评论读写切 PostgreSQL
- [ ] **W1-11** admin-service 配置/审核切 PostgreSQL
- [ ] **W1-12** 更新 README：Docker 启动 + 「做了什么 · 放在哪 · 怎么验证」

### 数据库 Schema（Week 1 最小集）

```text
users          id, email, password_hash, name, avatar, role, department, bio, status, created_at
boards         id, slug, name, description, enabled, sort_order
posts          id, board_id, author_id, title, content, tags, status, is_featured, like_count, comment_count, created_at
comments       id, post_id, author_id, parent_id, content, created_at
audit_records  id, post_id, reason, status, reviewer_id, created_at
system_config  id (singleton), posting_enabled, board_switches, moderation_mode
```

### 验收标准

1. `docker compose up -d postgres` 后，三个 Go 服务读写 DB，重启不丢数据
2. 可用邮箱注册新账号并登录；演示登录三角色仍可用
3. 学生发帖 → 审核 → 管理员通过，全流程与 MVP 一致

---

## Week 2 · Content & Moderation（状态机 + 治理链）

**目标**：补齐 XenForo Approval Queue 核心。

### Issue 清单

- [ ] **W2-01** 帖子状态机：`draft | pending_review | published | rejected | deleted`
- [ ] **W2-02** 发帖前敏感词检测（admin-service `/internal/v1/moderation/check-post` 读 DB 词库）
- [ ] **W2-03** 命中敏感词自动 `pending_review`，记录 reason
- [ ] **W2-04** 新增 `moderation_logs` 表（actor_id, action, target_type, target_id, detail, created_at）
- [ ] **W2-05** approve / reject 写入 Mod Log
- [ ] **W2-06** 评论支持 `parent_id`（一级嵌套回复）
- [ ] **W2-07** 点赞 toggle（`post_likes` 表，防重复 + 可取消）
- [ ] **W2-08** 前台评论树渲染

### 验收标准

1. 中台可看到「为何进审核」
2. 每条审核操作有 moderation_log 可查
3. 点赞可取消，计数正确

---

## Week 3 · Engagement & Admin（运营能力）

**目标**：从「公告板」升级为「可继续阅读

```text
favorites      user_id, post_id, created_at
attachments    id, post_id, type, name, url, mime, size_bytes, created_at
```

### Issue 清单

- [ ] **W3-01** 收藏 / 取消收藏 + 个人页「我的收藏」
- [ ] **W3-02** 附件上传（本地 `uploads/` 或 MinIO），MIME + 大小校验
- [ ] **W3-03** 板块 CRUD（中台增删改，不只开关）
- [ ] **W3-04** 置顶帖 `is_pinned`，列表排序优先
- [ ] **W3-05** 统计增强：今日发帖 / 评论 / 新用户（SQL 聚合）
- [ ] **W3-06** 前台分页（posts 列表 `?page=&limit=`）

### 验收标准

1. 协会可在中台自行调整板块
2. 技术问答区可上传 PDF / 图片
3. 首页展示今日活跃数字

---

## Week 4 · Production Shell（可交付）

**目标**：新人 10 分钟内跑通，材料可交接。

### Issue 清单

- [ ] **W4-01** 三个服务统一分层目录（handler / service / repository / model）
- [ ] **W4-02** `docker-compose.yml` 扩展：postgres + 3 服务 + frontend
- [ ] **W4-03** 各服务 `GET /health`（含 DB ping）
- [ ] **W4-04** `.env.example` 统一环境变量
- [ ] **W4-05** 部署文档 + 验证清单（对齐例会「三件套」）
- [ ] **W4-06** 打包 Android APK 演示包

### 验收标准

1. `docker compose up` 一键启动全栈
2. 新同学按 README 10 分钟内完成验证
3. 汇报不再依赖口头解释端口

---

## Phase 2  backlog（本月不做）

| 能力 | 说明 |
|------|------|
| Redis 缓存 | 热门帖、配置缓存 |
| Nginx 网关 | 公网部署时再上 |
| 邀请码系统 | 协会封闭注册 |
| 私信 / 通知 | XenForo Alert / Conversation |
| 搜索 | 全文检索（PostgreSQL tsvector 或 Meilisearch） |
| 积分等级 | 用户成长体系 |

---

## 汇报三件套模板

每次例会按此结构汇报：

1. **做了什么** — 勾选上表 Issue 编号 + 一句话
2. **放在哪** — 分支 / 目录 / Docker 命令
3. **怎么验证** — 3 步以内可复现的操作

---

## 分支建议

```text
main                 稳定演示版（含 frontend-standalone）
feat/week1-foundation
feat/week2-moderation
feat/week3-engagement
feat/week4-production
```
