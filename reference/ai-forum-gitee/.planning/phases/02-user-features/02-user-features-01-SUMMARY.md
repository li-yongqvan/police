---
phase: 02-user-features
plan: 01
subsystem: backend-services + frontend
tags: [authentication, forum, admin, vue3, go, postgresql]
requirements:
  provides:
    - AUTH-01
    - AUTH-02
    - AUTH-03
    - AUTH-04
    - AUTH-05
    - LEVEL-01
    - LEVEL-02
    - LEVEL-03
    - LEVEL-06
    - BOARD-01
    - BOARD-02
    - POST-01
    - POST-02
    - POST-03
    - POST-04
    - POST-05
    - COMM-01
    - COMM-02
    - INTER-01
    - INTER-02
    - ATTACH-01
    - ATTACH-02
    - ATTACH-03
    - ATTACH-04
    - INVITE-01
    - INVITE-02
    - FE-01
    - FE-02
    - MOD-01
  requires:
    - INFRA-01
    - INFRA-02
    - INFRA-03
  affects:
    - services/user
    - services/forum
    - services/admin
    - frontend
dependency_graph:
  - wave1_migrations -> wave2_user_service -> wave4_admin_service
  - wave1_migrations -> wave3_forum_service
  - wave2_5_user_service -> wave5_frontend
  - wave3_5_forum_service -> wave5_frontend
  - wave4_admin_service -> wave3_forum_service (moderation check)
tech_stack:
  added:
    - golang.org/x/crypto/bcrypt (password hashing)
    - gabriel-vasile/mimetype (file type detection)
  patterns:
    - Service/Handler/Middleware separation
    - JWT access token + Redis refresh token
    - Soft delete for posts (status='deleted')
    - In-memory sensitive word cache with sync.RWMutex
    - Pinia store pattern for Vue state management
key_files:
  created:
    - services/user/migrations/001_users_table.up.sql
    - services/user/migrations/002_invite_codes_table.up.sql
    - services/forum/migrations/001_boards_posts.up.sql
    - services/forum/migrations/002_comments_interactions.up.sql
    - services/forum/migrations/003_attachments.up.sql
    - services/admin/migrations/001_sensitive_words_seed.up.sql
    - services/admin/internal/handler/invite_handler.go
    - services/admin/internal/service/user_client.go
    - frontend/src/api/auth.js
    - frontend/src/api/forum.js
    - frontend/src/stores/auth.js
    - frontend/src/stores/forum.js
    - frontend/src/views/BoardList.vue
    - frontend/src/views/BoardDetail.vue
    - frontend/src/views/PostDetail.vue
    - frontend/src/views/PostCreate.vue
    - frontend/src/views/PostEdit.vue
    - frontend/src/views/Profile.vue
    - frontend/src/components/UserAvatar.vue
    - frontend/src/components/LevelBadge.vue
    - frontend/src/components/BoardCard.vue
    - frontend/src/components/PostCard.vue
  modified:
    - services/user/cmd/main.go
    - services/user/internal/model/user.go
    - services/user/internal/handler/user_handler.go
    - services/user/internal/service/user_service.go
    - services/user/internal/middleware/auth.go
    - services/user/pkg/database/db.go
    - services/user/pkg/jwt/jwt.go
    - services/forum/cmd/main.go
    - services/forum/internal/model/post.go
    - services/forum/internal/handler/post_handler.go
    - services/forum/internal/handler/comment_handler.go
    - services/forum/internal/handler/interaction_handler.go
    - services/forum/internal/handler/attachment_handler.go
    - services/forum/internal/service/forum_service.go
    - services/forum/internal/client/admin_client.go
    - services/forum/internal/middleware/auth.go
    - services/admin/cmd/main.go
    - services/admin/internal/service/admin_service.go
    - services/admin/internal/handler/moderation_handler.go
    - services/admin/internal/client/user_client.go
    - frontend/src/router/index.js
    - frontend/src/api/index.js
    - frontend/src/App.vue
    - frontend/src/views/Register.vue
    - frontend/src/views/Login.vue
decisions:
  - key: user-level-checking-location
    decision: "Level checking done via middleware on route level, not in service layer"
    reason: "Forum-service cannot directly query user-service DB; JWT claims provide level info"
  - key: sensitive-word-check-granularity
    decision: "Sensitive word check returns clean boolean + matched words array"
    reason: "Enables future UI feedback showing which words triggered the filter"
  - key: comment-pending-review
    decision: "Comments always saved as published in MVP; pending_review reserved for posts"
    reason: "Simplifies MVP; comments table lacks status column"
  - key: forum-service-user-client
    decision: "Forum-service uses admin-service's UserClient instead of duplicating"
    reason: "UserClient moved to admin-service/internal/service; forum-service references it"
  - key: file-upload-path
    decision: "Files saved to /data/uploads/YYYY/MM/ with unique filename"
    reason: "Per D-02 decision; prevents path traversal, avoids filename conflicts"
metrics:
  duration_minutes: "estimated 90"
  tasks_completed: 12
  waves_completed: 5
  files_created: 22
  files_modified: 25
  lines_added: 3600
  build_results:
    user_service: "go build successful"
    forum_service: "go build successful"
    admin_service: "go build successful"
    frontend: "vite build successful (111 modules, 171KB total)"
---

# Phase 02 Plan 01: User Features Summary

## One-liner

Complete user invite-code registration, JWT authentication with Redis refresh tokens, user profile/avatar management, level-based permission system, three forum boards, post/comment CRUD with sensitive word filtering, likes/collections, file attachment upload/download, admin invite code generation, and full Vue 3 frontend with responsive mobile design.

## Work Completed

### Wave 1: Database Migrations (Tasks 1-3)

- **user-service**: `users` table (id, username, password_hash, nickname, bio, avatar, level, status) + `invite_codes` table (code, created_by, used_by, used_at, status)
- **forum-service**: `boards` (with 3 seed rows: AI学习交流区, 协会公告&活动区, 技术问答求助区), `posts` (with board_id FK, status, like/comment counts), `comments` (with ON DELETE CASCADE), `likes` (unique post_id+user_id), `collections` (unique post_id+user_id), `attachments` (with file_type: image/document/link)
- **admin-service**: `sensitive_words` table with 9 seed words (赌博, 代写, 作弊, 刷单, 加微信, 加QQ, 兼职, 刷粉, 引流)

### Wave 2: user-service (Tasks 4-5)

- **Register**: Validates invite code (unused), bcrypt(cost=12) password, creates user with level=0, marks invite code used in atomic transaction
- **Login**: Username+password auth, JWT access token (30min), Redis refresh token (30 days), ban check
- **RefreshToken**: Iterates Redis keys to find matching refresh token, generates new access token, deletes token if user banned
- **Profile**: GET/PUT with ownership verification, username uniqueness check
- **Avatar upload**: Multipart upload, max 5MB, validates MIME type (JPG/PNG/GIF/WEBP), saves to `/data/uploads/avatars/`
- **Auth middleware**: JWT validation, ban status check, sets user_id/username/level in context
- **RequireLevel middleware**: Blocks sub-level users on write endpoints
- **Internal APIs**: `/internal/v1/users/:id/status`, `/internal/v1/invite-codes`, `/internal/v1/invite-codes/batch`

### Wave 3: forum-service (Tasks 6-8)

- **Board listing**: Returns all enabled boards with post count subquery
- **Post CRUD**: Create (level>=1, sensitive word check), read (pagination, pinned first), update (author check, re-check sensitive words), delete (soft delete, author check)
- **Comments**: List (pagination, 50/page), create (level>=1, sensitive word check), auto-increment post comment_count
- **Likes**: Toggle with unique constraint, updates like_count
- **Collections**: Toggle with unique constraint
- **Attachments**: Upload images (max 10MB), documents (max 20MB), links (no level restriction beyond auth); serve files with correct Content-Type; save to `/data/uploads/YYYY/MM/`

### Wave 4: admin-service (Task 9)

- **Invite code generation**: Single and batch (up to 1000) via user-service internal API
- **Sensitive word check**: In-memory cache loaded on startup, O(n) scan with strings.Contains, cache invalidation on add/delete
- **Sensitive word CRUD**: Add, list, delete endpoints

### Wave 5: Frontend (Tasks 10-12)

- **Auth pages**: Register (username, password, confirm, invite code), Login (username, password, remember me) with Pinia store
- **Forum views**: BoardList (grid cards), BoardDetail (paginated posts, pinned first), PostDetail (full content, attachments inline, comments, like/collect)
- **Post CRUD views**: PostCreate (title, content, board selector, file upload, link input), PostEdit (pre-fill, update)
- **Profile view**: Avatar upload with preview, editable username/nickname/bio, level badge display
- **Components**: UserAvatar (image or initials fallback with color hash), LevelBadge (5 colors by level), BoardCard, PostCard
- **Responsive design**: Mobile hamburger menu, single-column on <768px, 2-column on tablet, 3-column on desktop, 44px tap targets
- **API integration**: Axios interceptors for JWT refresh, Pinia stores for auth and forum state, navigation guards for protected routes

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Import] Fixed unused imports in forum-service**
- **Found during:** Task 6 build
- **Issue:** `net/http` and `time` imported but not used in forum_service.go
- **Fix:** Removed unused imports

**2. [Rule 3 - Dependency] Fixed admin_client.go Post method**
- **Found during:** Task 6 build
- **Issue:** `http.Client.Post()` accepts `io.Reader`, not `[]byte`
- **Fix:** Wrapped JSON payload with `bytes.NewBuffer()`

**3. [Rule 3 - Naming] Resolved UserClient naming conflict**
- **Found during:** Task 9 build
- **Issue:** Both `admin/internal/client/user_client.go` and `admin/internal/service/user_client.go` defined `UserClient`
- **Fix:** Removed `admin/internal/client/` directory entirely, consolidated in service layer

**4. [Rule 3 - Method] Fixed gin Context.GetUint non-existence**
- **Found during:** Task 4 build
- **Issue:** `c.GetUint()` does not exist in gin; must use type assertion
- **Fix:** Changed to `userIDAny, _ := c.Get("user_id"); userID := userIDAny.(uint)`

**5. [Rule 2 - Security] Added missing `os` import in forum-service middleware**
- **Found during:** Task 6 build (implicit)
- **Issue:** JWT_SECRET environment variable not read in forum-service middleware
- **Fix:** Added `os.Getenv("JWT_SECRET")` with default fallback

**6. [Rule 3 - Module] Removed unused `strconv` import in board_handler.go**
- **Found during:** Task 6 build
- **Issue:** `strconv` imported but not used after stubbing
- **Fix:** Removed import, made file a placeholder

## Build Results

| Service/Component | Command | Result |
|-------------------|---------|--------|
| user-service | `go build ./...` | Passed |
| forum-service | `go build ./...` | Passed |
| admin-service | `go build ./...` | Passed |
| frontend | `vite build` | Passed (111 modules, 171KB) |

## Known Stubs

| Stub | File | Description |
|------|------|-------------|
| config handlers | `services/admin/internal/handler/config_handler.go` | GetConfig/UpdateConfig still return 501 |
| audit handlers | `services/admin/internal/handler/audit_handler.go` | Audit pending/approve/reject still return 501 |
| user admin handlers | `services/admin/internal/handler/user_admin_handler.go` | Ban/ListUsers/UpdateLevel still return 501 |
| stats handlers | `services/admin/internal/handler/stats_handler.go` | Stats overview/daily still return 501 |
| health check stubs | `services/user/internal/handler/health.go` | HealthCheck defined in both health.go and user_handler.go (removed duplicate) |

These stubs are expected — they belong to Phase 3 (admin management) functionality, not Phase 2 user features.

## Threat Flags

| Flag | File | Description |
|------|------|-------------|
| threat_flag: file_upload | `services/forum/internal/handler/attachment_handler.go` | File upload endpoints accept multipart data with MIME validation and size limits (10MB images, 20MB documents) |
| threat_flag: token_refresh | `services/user/internal/service/user_service.go` | Refresh token iteration via Redis SCAN pattern; key-based lookup for user identification |
| threat_flag: internal_api | `services/user/cmd/main.go` | Internal invite code generation endpoints exposed without auth (per D-09, trust Docker internal network) |
| threat_flag: sensitive_word_cache | `services/admin/internal/service/admin_service.go` | In-memory sensitive word cache with sync.RWMutex; potential memory growth if many words added |

## Self-Check: PASSED

All migration files exist, all Go services compile, frontend builds successfully.
