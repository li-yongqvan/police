# AI 智联论坛 — 全栈对接审视与集成计划书

> **审视视角**：Addy Osmani 式全栈工程 — 先度量真实用户路径，再谈功能清单；API 契约清晰、按需加载、避免「后端有、前端假」的演示债。  
> **范围**：学院版 `frontend/` ↔ Go 微服务（user / forum / admin）↔ `nginx` 网关。  
> **结论摘要**：**MVP 主路径已接通**（登录 → 社区浏览/发帖/评论/点赞 → 中台审核/配置/封禁）；**约 60% 后端公开能力未暴露到 UI**，属产品裁剪与待集成 backlog，非全部「断线」。

---

## 1. 审视方法（Measure → Map → Fix）

| 步骤 | 做法 |
|------|------|
| **Measure** | 以用户故事为纲：演示登录、社区、中台四条旅程 |
| **Map** | 逐条对照 `services/*/cmd/main.go` 路由与 `frontend/src/api.js` + 视图 |
| **Fix** | 先修「已接线但损坏」（bug），再补「高价值缺口」，最后扩展 Gitee 级能力 |

**性能原则（Osmani）**

- 社区列表：优先 `GET /posts?board_id=`，避免拉全量再前端 filter（当前 Board 页可优化）。
- 管理概览：并行 `stats/overview` + `audit/pending` + `boards`，控制瀑布请求（见 Phase 1）。
- 静态前端：保持 route-level code split（Vite 已拆 chunk）；大图/附件走 CDN 或懒加载（Phase 3）。

---

## 2. 架构与流量路径

```
浏览器
  → :80  1Panel OpenResty
  → :8888  ai-forum-nginx
       ├─ /user-api/*   → user-service:8001
       ├─ /forum-api/*  → forum-service:8002
       ├─ /admin-api/*  → admin-service:8003
       └─ /             → frontend (Vue SPA)
```

学院前端 **仅** 通过 `api.js` 三前缀访问后端；`/internal/v1/*` 仅供服务间调用，**不应**由浏览器直连。

---

## 3. 对接矩阵（用户可见能力）

### 3.1 已接通（Yes）— MVP 可演示

| 用户故事 | 前端 | 后端 API |
|----------|------|----------|
| 账号密码登录 | `DemoLogin` → `userApi.login` | `POST /api/v1/login`, `GET /users/me` |
| 社区板块导航 | `CommunityLayout` | `GET /boards` |
| 首页信息流 | `CommunityHome` | `GET /posts`, `GET /boards` |
| 板块帖子列表 | `BoardView` | `GET /boards`, `GET /posts`（客户端按 board 过滤） |
| 帖子详情 / 评论 / 点赞 | `PostDetail` | `GET /posts/:id`, `GET/POST comments`, `POST like` |
| 发帖 | `NewPost` | `POST /posts` |
| 中台审核 | `AdminAudit` | `GET audit/pending`, `POST approve/reject` |
| 系统配置（审核模式、板块开关） | `AdminConfig` | `GET/PUT /admin/config/:key` |
| 用户列表 / 封禁 | `AdminUsers` | `GET /admin/users`, `POST ban` |
| 中台概览（部分指标） | `AdminOverview` | `GET /admin/stats/overview` |

### 3.2 部分接通（Partial）

| 能力 | 现状 | 缺口 |
|------|------|------|
| 中台概览指标 | UI 有今日发帖、待审、板块热度 | 映射不全；`board_activity` 未从 overview 返回；待审需 audit 列表计数 |
| 发帖开关 | 配置页可勾选 | `updateConfig` 未写 `post_requires_level` |
| 用户解封 | 按钮文案「恢复账号」 | 管理端无 public `unban`，仅 user internal |
| 社区待审提示 | `CommunityHome` 尝试 overview | `pendingAuditCount` 写死为 0 |
| 帖子标签 / 附件 | `NewPost` 有表单字段 | `createPost` 未传 tags；附件未走 `attachments/upload` |
| 个人主页 | `ProfileView` | 仅读 localStorage，未 `PUT /users/:id`、未上传头像 |
| 精华帖 | API 占位 throw | 后端有 `POST /admin/posts/:id/featured`，无 UI |

### 3.3 未接通（No）— 后端有、学院 UI 无

按 **业务域** 归纳（公开 `/api/v1` 与中台 API）：

| 域 | 后端能力（示例） | 建议 |
|----|------------------|------|
| **认证扩展** | `register`, `auth/refresh`, `demo-login` | 注册页 / Token 静默刷新（Phase 2） |
| **用户** | `PUT users/:id`, `POST avatar` | 个人资料编辑（Phase 2） |
| **论坛写操作** | `PUT/DELETE posts`, `collect` | 编辑帖、删帖、收藏（Phase 3） |
| **附件** | `upload`, `GET attachments/:id` | 发帖上传与详情展示（Phase 2） |
| **中台内容** | `admin/posts` 列表、featured/pinned、delete | 帖子管理页（Phase 3） |
| **中台用户** | `level`, `logs`, roles CRUD | 用户详情抽屉（Phase 3） |
| **邀请码** | invite-codes 全套 | 邀请码管理页（Phase 4，偏运营） |
| **板块 CRUD** | `admin/boards` POST/PUT/DELETE | 板块管理页（Phase 3） |
| **敏感词** | `admin/sensitive-words` | 敏感词配置页（Phase 4） |
| **统计** | `admin/stats/daily` | 趋势图 / 仪表盘（Phase 3） |
| **角色权限** | `admin/roles`, assign/remove | 权限矩阵（Phase 4） |

**Internal API**（`/internal/v1/*`）不计划直接暴露给浏览器；保持经 admin-service 聚合。

---

## 4. 已发现缺陷（Phase 0 — 立即修复）

| ID | 问题 | 影响 | 处理 |
|----|------|------|------|
| P0-1 | `AdminOverview.vue` 使用 `session` 未 import | 中台首页白屏 / 报错 | 修复 store 引用 |
| P0-2 | `getOverview` 未映射 `posts_today`、待审数、板块热度 | 指标恒为 0 | 完善 `api.js` 聚合 |
| P0-3 | `updateConfig` 忽略 `postingEnabled` | 发帖开关保存无效 | 写入 `post_requires_level` |
| P0-4 | `AdminUsers` 解封无 API | 点击恢复无效果 | UI 标注 + Phase 2 补 `unban` 路由 |

---

## 5. 分阶段实施计划

### Phase 0 — 修复接线 ✅

- [x] 修复 `AdminOverview`、配置发帖开关、`BoardView` 按板块拉帖

### Phase 1 — 契约与性能 ✅（2026-05-25）

- [x] `api/http.js`：401 自动 refresh / 跳转登录
- [x] `BoardView` / `getOverview` 并行请求

### Phase 2 — 社区闭环 ✅

- [x] 注册页 `/register`
- [x] 个人资料 `PUT` + 头像上传
- [x] 发帖附件 upload + `attachment_ids`
- [x] Token refresh
- [x] 帖子编辑/删除/收藏

### Phase 3 — 中台运营 ✅

- [x] 帖子管理（精华/置顶/删除）
- [x] 板块 CRUD
- [x] 趋势统计表 `stats/daily`
- [x] 用户解封 `POST /admin/users/:id/unban`

### Phase 4 — 平台能力 ✅

- [x] 邀请码生成/列表/作废
- [x] 敏感词增删
- [x] 角色分配（`platform_admin`）

### Phase 4 — 产品级扩展（按需）

- 邀请码、敏感词、角色权限 — 对齐 Gitee 参考后端，**仍用学院 UI 设计语言**
- E2E：`scripts/smoke-test.sh` + Playwright 三条旅程
- 可观测：结构化日志、网关 request-id、慢查询告警

---

## 6. 契约与字段对照（防漂移）

| 前端字段 | 后端来源 |
|----------|----------|
| `userCount` | `stats/overview.data.total_users` |
| `todayPostCount` | `posts_today` |
| `pendingAuditCount` | `len(audit/pending.posts)` |
| `boardActivity[].count` | `boards[].post_count` 或聚合 posts |
| `postingEnabled` | `post_requires_level !== '99'` |
| `moderationMode` | `sensitive_word_action` |

---

## 7. 不建议做的事（Osmani：Less but faster）

1. **不要把 internal 路由暴露给浏览器** — 安全与耦合双输。
2. **不要为对接而恢复 Gitee 整套 UI** — 学院版是刻意产品选择；只补 API 与页面模块。
3. **不要一次性接满 40+ admin 路由** — 按演示剧本优先级迭代。
4. **不要在首屏拉 admin stats** — 仅 `canAccessAdmin` 时请求。

---

## 8. 演示剧本验收清单（Definition of Done）

| # | 剧本 | 通过标准 |
|---|------|----------|
| 1 | 学生登录 | 表单登录 → 进入 `/community` |
| 2 | 浏览发帖 | 看板 → 发帖 → 审核模式下进入待审提示 |
| 3 | 互动 | 评论、点赞成功并刷新详情 |
| 4 | 协会管理员 | 审核通过/驳回生效 |
| 5 | 中台管理员 | 配置保存、概览数字合理、封禁生效 |
| 6 | 性能 | 4G 下首屏 &lt; 3s（缓存后）、API 并行无串联瀑布 |

---

## 9. 相关文档

- `docs/backend-framework-integration-plan.md` — 后端域模型
- `docs/frontend-ui-role-recovery-plan.md` — 学院 UI 与角色
- `docs/production-launch-plan.md` — 生产与 demo-login 策略

---

*文档版本：2026-05-25 · 审视基线：`frontend/` 学院 MVP + `feature/ai-forum-production` 后端*
