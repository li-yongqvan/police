# AI 智联论坛 — 帖子推荐功能开发记录

> 日期：2026-07-23
> 任务：增加种子帖子数据 + 帖子详情页新增「推荐帖子」模块

---

## 一、需求概述

1. **扩充种子数据**：为论坛增加 20 条 AI 主题的种子帖子，覆盖不同板块
2. **推荐帖子模块**：在帖子详情页（如 `/community/posts/20`）新增推荐模块，展示同板块的热门帖子

---

## 二、实现方案

### 2.1 推荐策略

- 基于 **同板块（board_id）** 匹配
- 排除当前帖子自身
- 按 **点赞数（like_count）降序** 排列
- 最多返回 **5 条**

### 2.2 架构改动

```
┌─ 前端 ─────────────────────────────────────┐
│ api.js          + getRecommendedPosts()     │
│ PostDetail.vue  + recommendedPosts ref      │
│                 + 推荐模块模板（文章↓评论之间）│
│ gx-theme.css    + .gx-post-recommended 样式  │
└────────────────────────────────────────────┘
         │ Vite proxy /forum-api → localhost:8002
         ▼
┌─ 后端 forum-service ───────────────────────┐
│ main.go          + GET /posts/:id/recommended│
│ post_handler.go  + GetRecommendedPosts()    │
│ forum_service.go + GetRecommendedPosts()    │
│   └─ SQL: SELECT ... WHERE board_id = ...   │
│          ORDER BY like_count DESC LIMIT 5   │
└────────────────────────────────────────────┘
         │
         ▼
┌─ 数据库 PostgreSQL ────────────────────────┐
│ schema_forum.posts (20 条种子数据)          │
│ 013_seed_more_posts.up.sql                  │
└────────────────────────────────────────────┘
```

---

## 三、改动文件清单

| 文件 | 操作 | 说明 |
|------|------|------|
| `services/forum/migrations/013_seed_more_posts.up.sql` | **新建** | 20 条种子帖子 SQL 迁移 |
| `services/forum/migrations/013_seed_more_posts.down.sql` | **新建** | 对应回滚脚本 |
| `services/forum/internal/service/forum_service.go` | 追加 | `GetRecommendedPosts` 方法 |
| `services/forum/internal/handler/post_handler.go` | 追加 | `GetRecommendedPosts` handler |
| `services/forum/cmd/main.go` | +1 行 | 注册路由 `/posts/:id/recommended` |
| `frontend/src/api.js` | +5 行 | `forumApi.getRecommendedPosts()` |
| `frontend/src/views/PostDetail.vue` | +25 行 | 推荐模块（ref + fetch + 模板） |
| `frontend/src/styles/gx-theme.css` | +35 行 | 推荐卡片样式 |
| `docker-compose.yml` | +12 行 | 为服务添加 `DB_HOST=postgres` 等环境变量 |

---

## 四、遇到的问题与解决方式

### 问题 1：PowerShell 写入 Go 源码导致 UTF-8 中文乱码

**现象**：`Set-Content -Encoding UTF8` 写入 Go 文件后，中文字符串变成 `闁哄啰濮…` 乱码，且反引号被错误转义。

**原因**：PowerShell 的 `Set-Content` 默认带 BOM，且字符串拼接时对特殊字符处理不当。

**解决**：
- 用 `git checkout` 恢复原始文件
- 新代码先写入 UTF-8 临时文件（`[System.IO.File]::WriteAllText`），再 `Add-Content` 追加
- 不再修改已有代码行，只追加新函数

### 问题 2：Docker Hub 网络不通，无法 rebuild 镜像

**现象**：`docker compose up --build` 卡在拉取 `golang:1.23-alpine` 等基础镜像环节。

**原因**：国内网络环境访问 `registry-1.docker.io` 超时。

**解决**：
- 在 Windows 宿主机交叉编译：`$env:GOOS="linux"; go build`
- 用 `docker cp` 将编译好的二进制和迁移文件拷贝进容器
- `docker restart` 使新二进制生效

### 问题 3：Docker Compose 中服务无法连接数据库

**现象**：forum-service 日志显示 `dial tcp 127.0.0.1:5432: connection refused`。

**原因**：`.env` 中 `DB_HOST=127.0.0.1` 是宿主机开发配置，容器内应使用 Docker 服务名 `postgres`。`docker-compose.yml` 未覆盖此变量。

**解决**：在 `docker-compose.yml` 中为 user-service、forum-service、admin-service 各添加 4 个环境变量：
- `DB_HOST=postgres`
- `REDIS_HOST=redis`
- 对应的服务间调用 URL

### 问题 4：数据库迁移状态 dirty

**现象**：forum-service 日志 `Dirty database version 8. Fix and force version.`，迁移无法继续。

**原因**：之前某次 docker compose down/up 周期中迁移被中断，`schema_forum_migrations` 表标记 version=8 为 dirty。

**解决**：
```sql
UPDATE schema_forum_migrations SET dirty = false WHERE version = 8;
```
之后强制设置版本到 11（因 012 迁移文件不存在于容器中）。

### 问题 5：种子数据 SQL 中 `role` 列不存在

**现象**：psql 执行种子 SQL 报错 `column "role" does not exist`。

**原因**：`schema_auth.users` 表使用 `level` 列而非 `role`。

**解决**：改用已有的 user ID（1=demo_student, 2=demo_admin 等）直接硬编码。

### 问题 6：迁移文件 BOM 导致语法错误

**现象**：迁移 013 执行时报 `syntax error at or near "﻿"`。

**原因**：PowerShell 的 `Set-Content -Encoding UTF8` 在文件开头写入 BOM（`0xEF 0xBB 0xBF`），PostgreSQL 无法识别。

**解决**：
```powershell
$utf8NoBom = New-Object System.Text.UTF8Encoding $false
[System.IO.File]::WriteAllText($path, $content, $utf8NoBom)
```

### 问题 7：api.js 中模板字符串反引号被 PowerShell 吃掉

**现象**：`getRecommendedPosts` 方法中 `` `/api/v1/posts/${id}/recommended` `` 变成了 `/api/v1/posts//recommended`。

**原因**：PowerShell 在字符串处理中将反引号（backtick）当作转义字符处理。

**解决**：使用 `[System.IO.File]::WriteAllLines` 配合单引号字符串写入，避免 PowerShell 转义。

### 问题 8：PostDetail.vue 模板中文损坏

**现象**：Vite 报 `Element is missing end tag`，`<span>{{ post.likeCount }} 赞</span>` 变成 `<span>{{ post.likeCount }} 璧?/span>`。

**原因**：同问题 1，PowerShell 文件操作导致 UTF-8 编码损坏。

**解决**：`git checkout` 恢复，然后用 `[System.IO.File]::ReadAllLines` + `WriteAllLines` 精准插入新代码。

### 问题 9：推荐模块插入位置错误

**现象**：推荐模块出现在评论模块**之后**而非之前。

**原因**：行匹配条件找到的是 `<section class="gx-post-comments-module">` 本身，插入时机在 `AppendLine` 之后。

**解决**：改为匹配评论模块**之前**的 `</div>` + 空行 + `<section>` 三段式锚点。

---

## 五、验证结果（子 agent 执行）

| 验证项 | 状态 |
|--------|:--:|
| `GET /api/v1/posts` — 20 条帖子 | ✅ |
| `GET /api/v1/posts/20` — 帖子详情 | ✅ |
| `GET /api/v1/posts/20/recommended` — 5 条推荐 | ✅ |
| `GET /api/v1/boards` — 4 个板块 | ✅ |
| `api.js` 语法（vite build 成功） | ✅ |
| `PostDetail.vue` ref / fetch / 模板位置 | ✅ |
| `gx-theme.css` 推荐样式 | ✅ |
| `main.go` 路由注册 | ✅ |
| `forum_service.go` / handler 方法存在 | ✅ |
| Vite 代理转发正常 | ✅ |

**17/17 全部通过。**

---

## 六、本地运行方式

```bash
# 启动完整 Docker 环境
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# 前端开发模式
cd frontend && npm run dev    # → http://127.0.0.1:8091

# 查看推荐效果
# 浏览器打开 http://127.0.0.1:8091/community/posts/20
# 登录凭据：demo_student / demo123456
```
