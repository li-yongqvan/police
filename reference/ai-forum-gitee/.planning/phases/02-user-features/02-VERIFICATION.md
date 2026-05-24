---
phase: 02-user-features
verified: 2026-05-21T12:00:00Z
status: gaps_found
score: 13/15 must-haves verified
overrides_applied: 0
gaps:
  - truth: "Level 2用户可发帖、评论、上传图片/附件"
    status: partial
    reason: "Attachment upload route POST /attachments/upload is under auth middleware but missing RequireLevel(2) guard. Level 0 and Level 1 users can upload attachments directly, bypassing the planned level restriction."
    artifacts:
      - path: "services/forum/cmd/main.go"
        issue: "Line 107: auth.POST('/attachments/upload', ...) lacks middleware.RequireLevel(2). CreatePost/CreateComment have RequireLevel(1) but attachment upload has no level check at all."
    missing:
      - "Add middleware.RequireLevel(2) to the attachment upload route in forum-service cmd/main.go"
  - truth: "发帖/评论提交时自动校验敏感词（调用admin-service）"
    status: partial
    reason: "forum-service calls admin_client.CheckSensitiveWords on both post and comment creation, but the CreateComment handler always saves the comment regardless of the sensitive word result (clean variable is discarded with _ = clean). Clean posts correctly get status='pending_review', but comments with sensitive words are always published."
    artifacts:
      - path: "services/forum/internal/service/forum_service.go"
        issue: "Line 329: '_ = clean' discards sensitive word result; comment always saved as published. Compare with CreatePost (line 181) which sets status='pending_review' when !clean."
    missing:
      - "Handle comment sensitive word result: either set comment status to pending_review or reject the comment when sensitive words are detected"
human_verification:
  - test: "Test frontend responsive design on mobile (375px) and desktop (1440px) viewports"
    expected: "All pages (BoardList, BoardDetail, PostDetail, PostCreate, Profile, Login, Register) render correctly with hamburger menu on mobile, grid layout on desktop"
    why_human: "CSS media queries exist in all views, but actual rendering quality, touch target sizes, and visual overlap can only be verified by visual inspection in a browser"
  - test: "End-to-end user flow: register with invite code, login, create post, comment, like, collect, upload attachment"
    expected: "Full flow completes without errors across services"
    why_human: "Requires running Docker Compose and making actual HTTP requests through the Nginx gateway. Automated checks cannot verify the full inter-service data flow without a live environment"
---

# Phase 02: 用户端核心功能 Verification Report

**Phase Goal:** 完成用户邀请码注册/登录系统、用户等级体系（Level 0~5+）、三大板块、发帖/评论/点赞/收藏/附件上传功能，以及用户端Vue前端界面，支持PC+移动端H5自适应。
**Verified:** 2026-05-21T12:00:00Z
**Status:** gaps_found
**Re-verification:** No — initial verification

## Goal Achievement

### Observable Truths

| # | Truth | Status | Evidence |
|---|-------|--------|----------|
| 1 | 用户凭有效邀请码可完成注册，注册后level=0 | VERIFIED | `user_service.go:34-112` validates invite code (status='unused'), bcrypt(cost=12) password, inserts user with level=0, marks invite code used in atomic transaction. Handler at `user_handler.go:28-44` binds request and calls service. |
| 2 | 用户可用用户名+密码登录，获得JWT access token + refresh token | VERIFIED | `user_service.go:116-156` authenticates with bcrypt, generates JWT access token (30min), stores refresh token in Redis (30 days), returns LoginResponse. Banned user check included. |
| 3 | JWT access token 30分钟后过期，需通过refresh token刷新 | VERIFIED | `user_service.go:17` accessTokenExpiry=30min, refreshTokenExpiry=30days. Refresh endpoint at `user_service.go:159-193` scans Redis keys for matching token, generates new access token, handles banned user deletion. |
| 4 | refresh token存储在Redis中，支持用户封禁后即时失效 | VERIFIED | `user_service.go:110,148` stores `refresh:{user_id}` keys in Redis. `user_service.go:178-181` deletes token when user banned. `user_handler.go:65-81` exposes POST /auth/refresh endpoint. |
| 5 | Level 0用户仅可浏览帖子和板块，不可发帖/评论/上传 | PARTIAL | `forum/cmd/main.go:95,100` — POST /posts and POST /comments have `RequireLevel(1)`. BUT POST /attachments/upload at line 107 has only `AuthMiddleware()` without `RequireLevel(2)`. Level 0/1 users can upload attachments. See gap #1. |
| 6 | Level 1用户可发帖、评论 | VERIFIED | `forum/cmd/main.go:95` `auth.POST("/posts", middleware.RequireLevel(1), ...)` and line 100 `auth.POST("/posts/:id/comments", middleware.RequireLevel(1), ...)`. `middleware/auth.go:54-66` implements RequireLevel with proper level comparison from JWT claims. |
| 7 | Level 2用户可发帖、评论、上传图片/附件 | PARTIAL | Upload handler `attachment_handler.go:28-138` has full implementation: MIME validation (images: jpeg/png/gif/webp, 10MB), document validation (pdf/doc/docx/xlsx/txt/md, 20MB), link type support, file saving to `/data/uploads/YYYY/MM/`. BUT the route at `forum/cmd/main.go:107` lacks `RequireLevel(2)`. See gap #1. |
| 8 | 三大板块（AI学习交流区、协会公告&活动区、技术问答求助区）存在于数据库 | VERIFIED | `forum/migrations/001_boards_posts.up.sql:32-35` seeds three boards: AI学习交流区 (ai-learning), 协会公告&活动区 (announcements), 技术问答求助区 (tech-help). `forum_service.go:26-52` queries boards with post count subquery. |
| 9 | 用户可创建/编辑/删除（自己的）帖子 | VERIFIED | `forum_service.go:171-211` CreatePost with sensitive word check, attachment linking. `forum_service.go:214-263` UpdatePost with ownership verification and re-check sensitive words. `forum_service.go:266-284` DeletePost with soft delete and ownership check. Handlers wired in `post_handler.go`. |
| 10 | 用户可对帖子评论、点赞、收藏 | VERIFIED | `forum_service.go:287-320` ListComments (pagination, 50/page). `forum_service.go:322-347` CreateComment (sensitive word check, auto-increment count). `forum_service.go:349-377` LikePost (toggle, update like_count). `forum_service.go:379-401` CollectPost (toggle). Handlers in `comment_handler.go` and `interaction_handler.go`. |
| 11 | 用户可在发帖/评论中上传图片、分享网盘链接、上传文档 | VERIFIED | `attachment_handler.go:28-138` handles three types: image (MIME validation, 10MB), document (extension validation, 20MB), link (URL stored). File saved to `/data/uploads/YYYY/MM/` with unique name using `mimetype` library for detection. |
| 12 | 用户可下载/查看已上传的附件 | VERIFIED | `attachment_handler.go:141-172` serves images inline (`c.File`), documents with Content-Disposition attachment, links redirect to external URL. `forum_service.go:422-432` GetAttachment from DB. `forum_service.go:435-454` GetPostAttachments. |
| 13 | 发帖/评论提交时自动校验敏感词（调用admin-service） | PARTIAL | `forum_service.go:174` CreatePost calls `AdminClient.CheckSensitiveWords`, sets status='pending_review' when dirty. `forum_service.go:240` UpdatePost re-checks. `forum_service.go:325-329` CreateComment calls check but **discards result** (`_ = clean`), always saves as published. `admin_client.go:28-57` makes HTTP POST to admin-service. See gap #2. |
| 14 | 管理员可生成单个/批量邀请码 | VERIFIED | `admin/cmd/main.go:108-109` routes POST /api/v1/admin/invite-codes and /batch. `invite_handler.go:22-51` validates count (1-1000), calls admin service. `admin_service.go:129-144` delegates to UserClient. `user/cmd/main.go:137-166` internal endpoints create codes in DB transaction. |
| 15 | 前端PC+移动端H5自适应浏览和操作 | VERIFIED | All 9 Vue views have `@media (max-width: 768px)` responsive rules. `App.vue:182-199` implements hamburger menu toggle. CSS custom properties in App.vue:55-64. Register (161px), Login (159px), Profile (202px), PostCreate (197px), BoardList (52px), BoardDetail (116px), PostDetail (255px), PostEdit (135px), PostCard (94px). Vite build produced dist/ output (111 modules). |

**Score:** 13/15 truths fully verified, 2 truths partial

### Required Artifacts

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `services/user/migrations/001_users_table.up.sql` | users表（含level/password_hash/nickname/bio/avatar/status） | VERIFIED | 15 lines, CREATE TABLE with all required columns, indexes on username and status |
| `services/user/migrations/002_invite_codes_table.up.sql` | invite_codes表（code/used_by/used_at/status） | VERIFIED | 12 lines, CREATE TABLE with code/created_by/used_by/used_at/status, indexes |
| `services/forum/migrations/001_boards_posts.up.sql` | boards表+posts表，含三大板块seed data | VERIFIED | 35 lines, boards+posts tables, 3 seed rows for the three boards |
| `services/forum/migrations/002_comments_likes_collections.up.sql` | comments/likes/collections表 | VERIFIED | 25 lines, all three tables with proper FK constraints and unique constraints |
| `services/forum/migrations/003_attachments.up.sql` | attachments表 | VERIFIED | 14 lines, attachments table with file_type, file_path, file_size columns |
| `services/admin/migrations/001_sensitive_words_seed.up.sql` | sensitive_words表+初始敏感词seed data | VERIFIED | 22 lines, CREATE TABLE + 9 seed words (赌博, 代写, 作弊, etc.) |
| `services/admin/migrations/001_invite_management.up.sql` | admin侧邀请码批量生成所需表 | MISSING | Plan lists this artifact but invite_codes table lives in user-service schema. Admin calls user-service internal API — no separate admin table needed. This is an acceptable plan deviation (documented in plan task 3). |

### Key Link Verification

| From | To | Via | Status | Details |
|------|-----|-----|--------|---------|
| `services/user/internal/handler/user_handler.go` | `services/user/internal/service/user_service.go` | service.(Register\|Login) | WIRED | `user_handler.go:35` calls `h.Service.Register`, line 55 calls `h.Service.Login` |
| `services/forum/internal/handler/post_handler.go` | `services/forum/internal/client/admin_client.go` | adminClient.CheckSensitiveWords | WIRED | `forum_service.go:174,240,325` all call `s.AdminClient.CheckSensitiveWords`. AdminClient created in `forum/cmd/main.go:56` and passed to ForumService |
| `services/user/internal/middleware/auth.go` | `services/user/pkg/database/db.go` | rdb.(Get\|Set\|Del) for refresh token | WIRED | `user_service.go:110,148,161-165,180` use `s.RDB.Set`, `s.RDB.Scan`, `s.RDB.Get`, `s.RDB.Del`. Redis client initialized in `user/cmd/main.go:55-58` |
| `frontend/src/views/PostCreate.vue` | `/api/v1/posts` | axios.post for post creation | WIRED | `PostCreate.vue:109` calls `forumStore.createPost(postData)`. `forum.js:18-21` `api.post('/posts', data)`. `stores/forum.js:49` calls `createPost(data)`. Axios interceptor in `api/index.js:12-20` attaches Bearer token |
| `services/forum/internal/service/forum_service.go` | `services/forum/internal/client/admin_client.go` | CheckSensitiveWords via HTTP | WIRED | `admin_client.go:35-42` POST to `{ADMIN_SERVICE_URL}/internal/v1/moderation/check`. `forum/cmd/main.go:52-56` sets ADMIN_SERVICE_URL from env with default `http://localhost:8003` |

### Data-Flow Trace (Level 4)

| Artifact | Data Variable | Source | Produces Real Data | Status |
|----------|--------------|--------|-------------------|--------|
| BoardList.vue | `forumStore.boards` | `GET /api/v1/boards` -> DB query | Yes — `forum_service.go:26-38` SELECT with post count subquery | FLOWING |
| BoardDetail.vue | `forumStore.posts` | `GET /api/v1/posts` -> DB query | Yes — `forum_service.go:100-129` SELECT posts with JOIN to users/boards | FLOWING |
| PostDetail.vue | `forumStore.currentPost` | `GET /api/v1/posts/:id` -> DB query | Yes — `forum_service.go:133-168` SELECT post + attachments | FLOWING |
| PostCreate.vue | form data | Multipart POST -> DB INSERT | Yes — `forum_service.go:186-194` INSERT with RETURNING | FLOWING |
| Profile.vue | `authStore.user` | `GET /api/v1/users/:id` -> DB query | Yes — `user_service.go:196-208` SELECT without password_hash | FLOWING |
| Register.vue | registration form | POST /register -> DB INSERT | Yes — `user_service.go:77-84` INSERT with RETURNING | FLOWING |

### Behavioral Spot-Checks

No runnable entry points available for spot-checks without Docker Compose environment. Services compile but database and Redis are not available for direct testing.

**Step 7b: SKIPPED (no runnable entry points — requires Docker Compose environment)**

### Requirements Coverage

| Requirement | Source Plan | Description | Status | Evidence |
|-------------|-------------|-------------|--------|----------|
| AUTH-01 | 02-01-PLAN | 用户可通过邀请码注册账号 | SATISFIED | `user_service.go:34-112` validates invite code, creates user with level=0 |
| AUTH-02 | 02-01-PLAN | 用户可通过用户名+密码登录 | SATISFIED | `user_service.go:116-156` bcrypt auth, JWT + Redis refresh token |
| AUTH-03 | 02-01-PLAN | 用户可查看和编辑个人资料 | SATISFIED | `user_service.go:196-246` GetProfile/UpdateProfile, `Profile.vue` frontend |
| AUTH-04 | 02-01-PLAN | 用户可上传和更换头像 | SATISFIED | `user_handler.go:135-205` MIME validation, save to /data/uploads/avatars |
| AUTH-05 | 02-01-PLAN | 用户登录后会话保持有效 | SATISFIED | JWT access token + Redis refresh token, axios interceptor auto-refreshes on 401 |
| LEVEL-01 | 02-01-PLAN | 注册后默认Level 0，仅可浏览 | SATISFIED | `user_service.go:79` level=0 on insert, read endpoints have no level restriction |
| LEVEL-02 | 02-01-PLAN | Level 1可发帖、评论 | SATISFIED | `forum/cmd/main.go:95,100` RequireLevel(1) on POST /posts, POST /comments |
| LEVEL-03 | 02-01-PLAN | Level 2可上传附件 | BLOCKED | Upload route missing RequireLevel(2) middleware |
| LEVEL-06 | 02-01-PLAN | 操作根据用户等级进行权限校验 | SATISFIED | `middleware/auth.go:54-66` RequireLevel middleware in both user-service and forum-service |
| BOARD-01 | 02-01-PLAN | 系统初始化三大核心板块 | SATISFIED | Migration seeds three boards |
| BOARD-02 | 02-01-PLAN | 用户可浏览板块列表 | SATISFIED | `forum_service.go:26-52` ListBoards, `BoardList.vue` frontend |
| POST-01~05 | 02-01-PLAN | 帖子CRUD | SATISFIED | Full CRUD in `forum_service.go:171-284`, handlers in `post_handler.go` |
| COMM-01~02 | 02-01-PLAN | 评论CRUD | SATISFIED | `forum_service.go:287-347`, `comment_handler.go` |
| INTER-01~02 | 02-01-PLAN | 点赞/收藏 | SATISFIED | `forum_service.go:349-401`, `interaction_handler.go` |
| ATTACH-01~04 | 02-01-PLAN | 附件上传下载 | PARTIAL | Full upload/download implementation exists, but level check missing on upload route |
| INVITE-01~02 | 02-01-PLAN | 管理员生成邀请码 | SATISFIED | `admin/cmd/main.go:108-109`, `invite_handler.go`, user-service internal endpoints |
| FE-01~02 | 02-01-PLAN | 前端PC+移动端H5自适应 | SATISFIED | Responsive CSS in all 9 views, hamburger menu in App.vue, @media queries at 768px |
| MOD-01 | 02-01-PLAN | 发帖/评论提交时自动校验敏感词 | PARTIAL | Post creation sets pending_review on dirty content; comment creation ignores result |

### Anti-Patterns Found

| File | Line | Pattern | Severity | Impact |
|------|------|---------|----------|--------|
| `services/forum/cmd/main.go` | 107 | Missing RequireLevel(2) on attachment upload route | Blocker | Level 0/1 users can upload files bypassing permission system |
| `services/forum/internal/service/forum_service.go` | 329 | `_ = clean` discards sensitive word check result | Warning | Comments with sensitive words are always published |
| `services/admin/internal/handler/config_handler.go` | 11,17 | TODO stubs | Info | Phase 3 admin config — expected |
| `services/admin/internal/handler/audit_handler.go` | 11,17,23 | TODO stubs | Info | Phase 3 audit workflow — expected |
| `services/admin/internal/handler/user_admin_handler.go` | 11,17,23,29 | TODO stubs | Info | Phase 3 user management — expected |
| `services/admin/internal/handler/stats_handler.go` | 11,17 | TODO stubs | Info | Phase 4 analytics — expected |
| `services/admin/internal/service/admin_service.go` | 112,118,124 | TODO stubs | Info | Phase 3 admin config/audit — expected |

Phase 3/4 stubs are expected and documented in SUMMARY.md Known Stubs section — not counted as gaps.

### Human Verification Required

1. **Responsive Design on Real Devices**
   - Test: Open frontend in browser at 375px (mobile) and 1440px (desktop) viewports
   - Expected: All pages render correctly, hamburger menu works on mobile, grid layouts adapt, text readable (min 14px), touch targets >= 44px
   - Why human: CSS media queries exist in code but visual rendering quality can only be verified by visual inspection

2. **End-to-End User Flow**
   - Test: Register with invite code -> Login -> Create post with image -> Comment -> Like -> Collect -> Download attachment
   - Expected: Full flow completes without errors across all three services via Nginx gateway
   - Why human: Requires live Docker Compose environment with PostgreSQL, Redis, and Nginx routing. Cannot verify without running the full stack

### Gaps Summary

**2 gaps found, both functional:**

1. **Attachment upload missing level restriction (affects LEVEL-03)**
   - The route `POST /api/v1/attachments/upload` in `services/forum/cmd/main.go:107` is protected by `AuthMiddleware()` but lacks `RequireLevel(2)`. This means Level 0 and Level 1 users can upload image and document attachments, bypassing the planned permission system.
   - Fix: Add `middleware.RequireLevel(2)` before `attachmentHandler.UploadAttachment` on line 107. Note: link-type attachments should remain at level 1 per the plan.

2. **Comment sensitive word check result is ignored (affects MOD-01)**
   - In `services/forum/internal/service/forum_service.go:325-329`, `CreateComment` calls `AdminClient.CheckSensitiveWords` but discards the result with `_ = clean`. The comment is always saved as published, even when it contains sensitive words. By contrast, `CreatePost` correctly sets `status='pending_review'` when content is dirty.
   - Fix: Handle the `clean` variable — either set a pending_review status on comments (would require adding status column to comments table) or reject the comment when sensitive words are detected.

_Verified: 2026-05-21T12:00:00Z_
_Verifier: Claude (gsd-verifier)_
