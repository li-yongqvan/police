---
phase: 03-admin-system
plan: 01
subsystem: admin
tags: [admin, rbac, audit, moderation, user-management, config-management]
requires: []
provides:
  - admin_migrations: system_config, operation_logs, roles, user_roles tables
  - forum_admin_api: 12 internal endpoints for post/board management
  - user_admin_api: 8 internal endpoints for user/invite management
  - admin_service_handlers: audit, post, user, config, role, board, invite handlers
  - admin_frontend: 9 Vue pages with Element Plus
affects:
  - services/admin/migrations/002_admin_tables.up.sql
  - services/admin/internal/model/
  - services/admin/internal/client/forum_client.go
  - services/admin/internal/service/audit_service.go
  - services/admin/internal/service/config_service.go
  - services/admin/internal/service/user_admin_service.go
  - services/admin/internal/service/role_service.go
  - services/admin/internal/service/post_admin_service.go
  - services/admin/internal/handler/
  - services/admin/cmd/main.go
  - services/forum/internal/service/forum_admin_service.go
  - services/forum/internal/handler/forum_admin_handler.go
  - services/forum/cmd/main.go
  - services/user/internal/service/user_admin_service.go
  - services/user/internal/handler/user_admin_handler.go
  - services/user/cmd/main.go
  - services/user/migrations/003_operation_logs.up.sql
  - frontend/src/api/admin.js
  - frontend/src/stores/admin.js
  - frontend/src/router/index.js
  - frontend/src/views/admin/
tech-stack:
  added: [Element Plus, @element-plus/icons-vue, Pinia admin store, JWT role guard]
  patterns: [service-proxy, HTTP client delegation, JSONB detail logging, Redis role caching]
key-files:
  created:
    - services/admin/migrations/002_admin_tables.up.sql
    - services/admin/migrations/002_admin_tables.down.sql
    - services/admin/internal/model/system_config.go
    - services/admin/internal/model/operation_log.go
    - services/admin/internal/model/role.go
    - services/admin/internal/model/audit_record.go
    - services/admin/internal/client/forum_client.go
    - services/admin/internal/service/audit_service.go
    - services/admin/internal/service/config_service.go
    - services/admin/internal/service/user_admin_service.go
    - services/admin/internal/service/role_service.go
    - services/admin/internal/service/post_admin_service.go
    - services/admin/internal/handler/audit_handler.go
    - services/admin/internal/handler/config_handler.go
    - services/admin/internal/handler/user_admin_handler.go
    - services/admin/internal/handler/role_handler.go
    - services/admin/internal/handler/post_admin_handler.go
    - services/admin/internal/handler/invite_admin_handler.go
    - services/forum/internal/service/forum_admin_service.go
    - services/forum/internal/handler/forum_admin_handler.go
    - services/user/internal/service/user_admin_service.go
    - services/user/internal/handler/user_admin_handler.go
    - services/user/migrations/003_operation_logs.up.sql
    - frontend/src/api/admin.js
    - frontend/src/stores/admin.js
    - frontend/src/views/admin/AdminLayout.vue
    - frontend/src/views/admin/AdminLogin.vue
    - frontend/src/views/admin/AdminDashboard.vue
    - frontend/src/views/admin/AdminAudit.vue
    - frontend/src/views/admin/AdminPostManagement.vue
    - frontend/src/views/admin/AdminUserManagement.vue
    - frontend/src/views/admin/AdminInviteManagement.vue
    - frontend/src/views/admin/AdminBoardManagement.vue
    - frontend/src/views/admin/AdminConfig.vue
    - frontend/src/views/admin/AdminSensitiveWords.vue
    - frontend/src/views/admin/AdminRoleManagement.vue
  modified:
    - services/admin/internal/model/admin.go
    - services/admin/internal/client/user_client.go
    - services/admin/cmd/main.go
    - services/forum/cmd/main.go
    - services/user/cmd/main.go
    - frontend/package.json
    - frontend/vite.config.js
    - frontend/src/main.js
    - frontend/src/router/index.js
decisions:
  - Ban user deletes Redis refresh token for instant revocation (D-04, D-05)
  - Admin actions logged to JSONB detail field for repudiation mitigation (T-03-03)
  - ForumClient defines local Board/PendingPost structs to avoid cross-module dependency
  - PostAdminService intermediate layer between handler and ForumClient for consistency
  - Role caching via Redis SET with user_id key prefix
  - Navigation guard checks JWT role claim (admin/platform_admin) for admin routes
  - Sensitive words API called directly via dynamic axios import instead of api/admin.js wrapper
metrics:
  duration: ~45min
  tasks_completed: 12/12
  files_created: 35
  files_modified: 9
  total_commits: 7
---

# Phase 03 Plan 01: Admin System Summary

**One-liner:** Complete admin management system with RBAC, audit workflow, user ban/level control, config management, board CRUD, invite tracking, and 9-page Vue admin dashboard.

## Commits

| # | Hash | Message |
|---|------|---------|
| 1 | 19b2efb | feat(03-01): create admin-service migrations for system_config, operation_logs, roles, user_roles |
| 2 | f261321 | feat(03-01): create admin-service model types for SystemConfig, OperationLog, Role, AuditRecord |
| 3 | c15a596 | feat(03-01): create ForumClient in admin-service for calling forum-service admin endpoints |
| 4 | 7ce4b67 | feat(03-01): add admin API endpoints to forum-service for post status management and board CRUD |
| 5 | 0506d46 | feat(03-01): add admin API endpoints to user-service for ban, list users, level change, operation logs, invite code management |
| 6 | 911a290 | feat(03-01): implement admin-service handlers for audit workflow, post management, user admin, config, roles, boards, invites |
| 7 | 11445da | feat(03-01): create admin Vue frontend with Element Plus, 9 pages, routing, API layer |

## Migration Results

### admin-service: 002_admin_tables

| Table | Schema | Status | Seed Data |
|-------|--------|--------|-----------|
| system_config | schema_admin | Created | 11 entries (post.require_level, post.daily_limit, post.interval, post.title_min/max, post.content_min, board.activity/enabled, audit.auto_reject_level, audit.sensitive_action) |
| operation_logs | schema_admin | Created | None (indexes on operator_id, action, created_at) |
| roles | schema_admin | Created | 2 entries (admin with permissions: ["audit_post","delete_post","manage_user","manage_config","manage_board"], platform_admin with ["*"]) |
| user_roles | schema_admin | Created | None |

### user-service: 003_operation_logs

| Table | Schema | Status | Seed Data |
|-------|--------|--------|-----------|
| operation_logs | schema_auth | Created | None (indexes on target_id, operator_id, created_at) |

## Build Results

| Service | Command | Status |
|---------|---------|--------|
| admin-service | `go build ./...` | PASS |
| forum-service | `go build ./...` | PASS |
| user-service | `go build ./...` | PASS |
| frontend | `npm run build` | PASS (1702 modules, ~1.1MB) |

## API Endpoint Listing

### forum-service Internal Endpoints (`/internal/v1/`)

| Method | Path | Handler | Status |
|--------|------|---------|--------|
| GET | `/internal/v1/posts/pending` | ListPendingPosts | Implemented |
| POST | `/internal/v1/posts/:id/admin-approve` | ApprovePost | Implemented |
| POST | `/internal/v1/posts/:id/admin-reject` | RejectPost | Implemented |
| POST | `/internal/v1/posts/:id/admin-delete` | AdminDeletePost | Implemented |
| POST | `/internal/v1/posts/:id/admin-featured` | SetPostFeatured | Implemented |
| POST | `/internal/v1/posts/:id/admin-pinned` | SetPostPinned | Implemented |
| POST | `/internal/v1/posts/:id/admin-status` | ChangePostStatus | Implemented |
| POST | `/internal/v1/posts/batch-delete` | BatchDeletePosts | Implemented |
| POST | `/internal/v1/boards/admin` | CreateBoard | Implemented |
| PUT | `/internal/v1/boards/admin` | UpdateBoard | Implemented |
| DELETE | `/internal/v1/boards/admin` | DeleteBoard | Implemented |
| GET | `/internal/v1/boards` | ListAllBoards | Implemented |

### user-service Internal Endpoints (`/internal/v1/`)

| Method | Path | Handler | Status |
|--------|------|---------|--------|
| POST | `/internal/v1/users/:id/ban` | BanUser | Implemented |
| POST | `/internal/v1/users/:id/unban` | UnbanUser | Implemented |
| GET | `/internal/v1/users` | ListUsers | Implemented |
| PUT | `/internal/v1/users/:id/level` | UpdateUserLevel | Implemented |
| GET | `/internal/v1/users/:id/logs` | GetUserLogs | Implemented |
| GET | `/internal/v1/invite-codes` | ListInviteCodes | Implemented |
| GET | `/external/v1/invite-codes/:code/status` | GetInviteCodeStatus | Implemented |
| PUT | `/internal/v1/invite-codes/:code/void` | VoidInviteCode | Implemented |

### admin-service External Endpoints (`/api/v1/admin/`, AuthMiddleware protected)

| Method | Path | Handler | Service Delegate | Status |
|--------|------|---------|-----------------|--------|
| GET | `/api/v1/admin/audit/pending` | AuditHandler.ListPendingAudit | AuditService -> ForumClient | Implemented |
| POST | `/api/v1/admin/audit/:id/approve` | AuditHandler.ApprovePost | AuditService -> ForumClient | Implemented |
| POST | `/api/v1/admin/audit/:id/reject` | AuditHandler.RejectPost | AuditService -> ForumClient | Implemented |
| POST | `/api/v1/admin/audit/batch-delete` | AuditHandler.BatchDeletePosts | AuditService -> ForumClient | Implemented |
| GET | `/api/v1/admin/posts` | PostAdminHandler.ListAllBoards | ForumClient | Implemented |
| DELETE | `/api/v1/admin/posts/:id` | PostAdminHandler.DeletePost | PostAdminService | Implemented |
| POST | `/api/v1/admin/posts/:id/featured` | PostAdminHandler.SetPostFeatured | PostAdminService | Implemented |
| POST | `/api/v1/admin/posts/:id/pinned` | PostAdminHandler.SetPostPinned | PostAdminService | Implemented |
| GET | `/api/v1/admin/users` | UserAdminHandler.ListUsers | UserAdminService -> UserClient | Implemented |
| POST | `/api/v1/admin/users/:id/ban` | UserAdminHandler.BanUser | UserAdminService -> UserClient | Implemented |
| POST | `/api/v1/admin/users/:id/unban` | UserAdminHandler.UnbanUser | UserAdminService -> UserClient | Implemented |
| PUT | `/api/v1/admin/users/:id/level` | UserAdminHandler.UpdateUserLevel | UserAdminService -> UserClient | Implemented |
| GET | `/api/v1/admin/users/:id/logs` | UserAdminHandler.GetUserLogs | UserAdminService -> UserClient | Implemented |
| GET | `/api/v1/admin/invites` | InviteAdminHandler.ListInviteCodes | UserClient | Implemented |
| GET | `/api/v1/admin/invites/:code/status` | InviteAdminHandler.GetInviteCodeStatus | UserClient | Implemented |
| PUT | `/api/v1/admin/invites/:code/void` | InviteAdminHandler.VoidInviteCode | UserClient | Implemented |
| GET | `/api/v1/admin/config` | ConfigHandler.GetAllConfig | ConfigService | Implemented |
| GET | `/api/v1/admin/config/:key` | ConfigHandler.GetConfigByKey | ConfigService | Implemented |
| PUT | `/api/v1/admin/config/:key` | ConfigHandler.UpdateConfig | ConfigService | Implemented |
| GET | `/api/v1/admin/boards` | BoardAdminHandler.ListAllBoards | ForumClient | Implemented |
| POST | `/api/v1/admin/boards` | BoardAdminHandler.CreateBoard | ForumClient | Implemented |
| PUT | `/api/v1/admin/boards/:id` | BoardAdminHandler.UpdateBoard | ForumClient | Implemented |
| DELETE | `/api/v1/admin/boards/:id` | BoardAdminHandler.DeleteBoard | ForumClient | Implemented |
| GET | `/api/v1/admin/roles` | RoleHandler.ListRoles | RoleService | Implemented |
| POST | `/api/v1/admin/roles/assign` | RoleHandler.AssignRole | RoleService | Implemented |
| GET | `/api/v1/admin/roles/:userId` | RoleHandler.GetUserRoles | RoleService | Implemented |
| GET | `/api/v1/admin/sensitive-words` | (frontend direct) | - | Implemented |
| POST | `/api/v1/admin/sensitive-words` | (frontend direct) | - | Implemented |
| DELETE | `/api/v1/admin/sensitive-words/:id` | (frontend direct) | - | Implemented |

### Admin Frontend Routes

| Path | Component | Auth Guard | Status |
|------|-----------|------------|--------|
| `/admin/login` | AdminLogin | None | Implemented |
| `/admin` | AdminLayout + AdminDashboard | requiresAdmin | Implemented |
| `/admin/audit` | AdminAudit | requiresAdmin | Implemented |
| `/admin/posts` | AdminPostManagement | requiresAdmin | Implemented |
| `/admin/users` | AdminUserManagement | requiresAdmin | Implemented |
| `/admin/invites` | AdminInviteManagement | requiresAdmin | Implemented |
| `/admin/boards` | AdminBoardManagement | requiresAdmin | Implemented |
| `/admin/config` | AdminConfig | requiresAdmin | Implemented |
| `/admin/sensitive-words` | AdminSensitiveWords | requiresAdmin | Implemented |
| `/admin/roles` | AdminRoleManagement | requiresAdmin | Implemented |

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed Board type undefined in ForumClient**
- **Found during:** Task 3 (ForumClient creation)
- **Issue:** ForumClient imported `ai-forum/admin-service/internal/model` and tried to use `model.Board`, but Board type is defined in forum-service module, not admin-service
- **Fix:** Defined local `Board` and `PendingPost` structs directly in `forum_client.go` and removed the model import
- **Files modified:** `services/admin/internal/client/forum_client.go`

**2. [Rule 2 - Missing] Added resolve.alias to Vite config**
- **Found during:** Task 11 (frontend build)
- **Issue:** `@/` import paths in admin.js, router, and components failed to resolve because vite.config.js lacked `resolve.alias` configuration
- **Fix:** Added `resolve: { alias: { '@': path.resolve(__dirname, 'src') } }` to vite.config.js with `import path from 'path'`
- **Files modified:** `frontend/vite.config.js`

**3. [Rule 3 - Blocking] Split BoardAdminHandler and PostAdminHandler into separate files**
- **Found during:** Task 6 (admin-service handler wiring)
- **Issue:** `post_admin_handler.go` contained both BoardAdminHandler (using `client.ForumClient`) and PostAdminHandler (using `service.PostAdminService`), creating conflicting type references
- **Fix:** Created separate `post_admin_handler.go` for PostAdminHandler and `invite_admin_handler.go` for InviteAdminHandler. Created `post_admin_service.go` as intermediate service layer
- **Files created:** `services/admin/internal/service/post_admin_service.go`, `services/admin/internal/handler/invite_admin_handler.go`

**4. [Rule 2 - Missing] Added Redis client to RoleService for role caching**
- **Found during:** Task 8 (role service implementation)
- **Issue:** RoleService.AssignRole needed to write role to Redis for caching per D-06, but no Redis client was wired
- **Fix:** Added `SetRedis(rdb *redis.Client)` method to RoleService and called it from `main.go` during service initialization
- **Files modified:** `services/admin/internal/service/role_service.go`, `services/admin/cmd/main.go`

### Unused Import Cleanup

- Removed unused `model` import from `services/admin/internal/model/audit_service.go`
- Removed unused `strconv` import from `services/admin/internal/handler/config_handler.go`

## Security Notes

- Ban user workflow implements instant token revocation by deleting Redis refresh token (D-04, D-05 decision)
- All admin actions logged to `operation_logs` with JSONB detail for audit trail (T-03-03 repudiation mitigation)
- Admin routes protected by AuthMiddleware + role claim check in JWT (admin or platform_admin)
- Batch delete limited to 100 posts per request
- Void invite code only works on unused codes (status='unused')
- Level change validates range 0-5

## Known Stubs

None. All admin pages are fully wired to their respective API endpoints.

## Self-Check: PASSED

- All 7 commits verified in git log
- All created files confirmed on disk
- All 3 Go services build successfully
- Frontend builds successfully (1702 modules)
- All 12 plan tasks completed
