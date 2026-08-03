# Handoff - 2026-07-24 - 帖子推荐与论坛细节优化

> 新会话建议先读本文档，然后说：接着做。

## 1. 现在做了什么

- 当前目标：在现有论坛视觉风格已定型的前提下，参考 Discourse 的成熟论坛交互习惯，继续增加更符合人类使用习惯的小细节。
- 当前活跃分支任务：先从帖子详情页 `/community/posts/:id` 开始做细节优化。
- 最新用户决策：推荐帖子补位逻辑“不加区分”，即同板块推荐不足 5 条时，用全站热门帖子补齐，但 UI 不显示“补位来源”标签。
- 已确定参考源：优先参考 `discourse/discourse` 里的内置 `plugins/styleguide`，不要直接套小众 forum design system。

## 2. 已经完成了什么

- 后端新增推荐帖子 API：
  - `services/forum/cmd/main.go` 注册 `GET /api/v1/posts/:id/recommended`。
  - `services/forum/internal/handler/post_handler.go` 新增 `GetRecommendedPosts` handler。
  - `services/forum/internal/service/forum_service.go` 新增 `GetRecommendedPosts` service 方法。
- 推荐逻辑已升级：
  - 先取当前帖同板块热门帖。
  - 排除当前帖。
  - 不足 5 条时，从全站热门帖补齐。
  - 排除已选帖子，避免重复。
  - 不在前端区分同板块推荐和补位推荐。
- 前端新增推荐模块：
  - `frontend/src/api.js` 新增 `forumApi.getRecommendedPosts(id)`。
  - `frontend/src/views/PostDetail.vue` 新增 `recommendedPosts`、加载推荐数据，并在文章和评论区之间渲染推荐模块。
  - `frontend/src/styles/gx-theme.css` 新增 `.gx-post-recommended` 相关样式。
- 数据与本地环境：
  - 本地 PostgreSQL 中已手动插入 20 条帖子数据。
  - `services/forum/migrations/013_seed_more_posts.up.sql` 和 `.down.sql` 已创建，但本地容器迁移曾出过 dirty/BOM 问题，实际本地数据是手动 SQL 插入验证的。
  - `docker-compose.yml` 已补充容器内服务环境变量，例如 `DB_HOST=postgres`、`REDIS_HOST=redis`、内部服务 URL。
- 验证结果：
  - `go build ./cmd/...` 在 `services/forum` 下通过。
  - 本地 forum-service 容器已通过 `docker cp` 替换 Linux 二进制并重启。
  - `GET http://127.0.0.1:8002/api/v1/posts/2/recommended` 返回 5 条：前 4 条为公告区同板块，第 5 条为全站热门补位。
  - `GET http://127.0.0.1:8092/forum-api/api/v1/posts/2/recommended` 通过 Vite 代理返回同样 5 条。
  - Playwright 已用完整 Chromium 指定路径完成页面验证：`http://127.0.0.1:8092/community/posts/2` 显示正文，推荐区 5 条正常。
- 文档：
  - `docs/dev-records/2026-07-23-post-recommend.md` 记录技术过程。
  - `docs/dev-records/2026-07-23-chat-log.md` 记录对话摘要式日志。

## 3. 卡在哪里

- 暂无明确阻塞。
- 但有几个未收尾风险：
  - `services/forum/forum-service` 是本地交叉编译产物，当前为 untracked 文件，不应提交到 Git。
  - 013 迁移文件需要重新审查编码和实际可执行性，不能直接假设生产迁移安全。
  - 本地容器曾出现 `schema_forum_migrations` dirty 状态，后续如果重建数据库，需要重新确认迁移链完整。
  - Playwright 的 `chromium-headless-shell` 下载失败，但完整 Chromium 已下载，可通过显式 `executablePath` 运行自动化。
  - 前端仍有已有路由警告：`/rank`、`/messages`、`/library`、`/about` 无匹配。与推荐功能无关，但后续应单独处理。

## 4. 下一步做什么

1. 先清理工作区：确认不提交 `services/forum/forum-service` 编译产物。
2. 审查 `services/forum/migrations/013_seed_more_posts.up.sql` 和 `.down.sql`，确保 UTF-8 无 BOM、schema 字段真实存在、迁移可重复执行。
3. 用 `go build ./cmd/...`、前端构建和 Playwright 页面检查重新跑一遍基线。
4. 继续参考 Discourse，从帖子详情页做低风险细节优化，建议优先级：
   - 评论加载失败不阻断正文显示。
   - 推荐区移动端排版更可读，避免文字挤压。
   - 互动按钮增加明确 loading/成功反馈。
   - 返回按钮行为统一：有历史则返回，无历史则去社区首页。
   - 评论输入框在移动端更稳定，避免遮挡和布局跳动。
5. 做完每一组细节后，必须用浏览器自动化打开真实页面验证，而不只测 API。

## 5. 哪些坑不要再踩

- 不要用 PowerShell `Set-Content -Encoding UTF8` 随意改含中文的源码或 SQL；之前导致过 Vue 模板中文损坏、Go 字符串乱码和 SQL BOM 错误。
- 不要用 `git checkout --` 恢复文件，除非用户明确同意；当前规则要求不能回退可能由用户或工具产生的改动。
- 不要只测帖子详情 API 和推荐 API；帖子详情页还依赖评论 API，之前就是评论表缺列导致页面只显示背景。
- 不要把评论请求失败设计成整个帖子详情页失败；这是后续应该修的关键 UX 问题。
- 不要假设 Docker rebuild 可用；Docker Hub 网络可能超时，必要时继续使用宿主机交叉编译加 `docker cp` 的临时验证方式。
- 不要把本地手动数据库补丁当成正式迁移完成；手动补过 `comments.like_count`、`comments.dislike_count`、`comment_likes`、`comment_dislikes`，正式迁移仍需整理。
- 不要把 Discourse 视觉风格直接套进本项目；用户明确说整体设计已定型，只需要交互、功能和 UI 小细节。
- 使用 `grill-me` 时一次只问一个关键决策。已经完成的决策：推荐补位不加 UI 区分。
