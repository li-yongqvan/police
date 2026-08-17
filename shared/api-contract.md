# API 字段契约（前端 mapper ↔ Go JSON）

> 维护原则：后端 JSON 使用 **snake_case**；前端 `api.js` 映射为 **camelCase**。修改 Go struct 时同步更新此表。

## User

| 后端字段 | 前端字段 | 说明 |
|----------|----------|------|
| `id` | `id` | string 化 |
| `username` | `username` | |
| `nickname` / `name` | `name` | 展示名 |
| `avatar` | `avatar` | |
| `role` | `role` | student / admin / platform_admin |
| `level` | `level` | 0-5（纯展示属性，不作为权限依据） |
| `bio` | `bio` | |
| `status` | `status` | active / banned |
| `access_token` / `token` | localStorage `ai-forum-token` | 登录响应 |

## Forum — Board

| 后端 | 前端 |
|------|------|
| `sort_order` | `sortOrder` |
| `post_count` | `postCount` |
| `enabled` | `enabled` |

## Forum — Post

| 后端 | 前端 |
|------|------|
| `board_id` | `boardId` |
| `board_name` | `boardName` |
| `author_id` | `authorId` |
| `author_name` | `authorName` |
| `is_featured` | `isFeatured` |
| `is_pinned` | `isPinned` |
| `like_count` | `likeCount` |
| `comment_count` | `commentCount` |
| `created_at` | `createdAt` |
| `status` | `status` | published / pending_review |

## Forum — 列表分页

| 后端 | 前端 |
|------|------|
| `posts` | `posts` |
| `total` | `total` |
| `page` | `page` |
| `limit` | `limit` |

## Admin — Config keys

| 配置键 | 前端 `getConfig` |
|--------|------------------|
| `post_requires_level` | `postingEnabled`（≠99 为 true） |
| `sensitive_word_action` | `moderationMode` manual/auto |
| `board_{slug}_enabled` | `boardSwitches[boardId]` |

## Admin — Audit

| 后端 | 前端 |
|------|------|
| `post_ids` (batch-delete body) | `batchDeleteAudit(ids, reason)` |
| `reason` (reject body) | `rejectAudit(id, reason)` |

## Admin — Invite status

`GET /api/v1/admin/invite-codes/:code/status` → 原样展示 `status` / `message` 字段。

## Session Token Contract（JWT 会话令牌）

> 由 user-service 唯一签发（HS256），admin-service 仅校验消费，任何服务不得自行签发。

| Claim | 类型 | 说明 |
|-------|------|------|
| `user_id` | uint | 用户 ID |
| `username` | string | 用户名 |
| `role` | string | student / admin / platform_admin（签发时校验域，域外拒绝签发） |
| `exp` / `iat` | int | 过期 / 签发时间戳 |
| ~~`level`~~ | — | 已移除：等级是展示属性，不得作为权限依据 |

校验行为：admin 路由要求 role ∈ {admin, platform_admin}，否则 403；user-service 中间件仅设置 user_id / username，不做角色判定。

## 内部端点 — 角色权威（admin-service → user-service）

`GET /internal/v1/users/:id/role` — user-service 登录时读取权威角色名（Redis 60s 缓存读穿），不再直连 schema_admin。

| 状态码 | 语义 |
|--------|------|
| `200` | `{"user_id": 42, "role": "admin"}`（无角色分配返回 student，永不 404） |
| `400` | :id 非法 |
| `500` | 权威查询失败（user-service 降级为 student 并记日志，降级值不入缓存） |

