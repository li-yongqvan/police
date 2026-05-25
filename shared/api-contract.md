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
| `level` | `level` | |
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
