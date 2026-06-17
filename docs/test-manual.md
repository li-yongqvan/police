# AI 智联论坛 — 全功能测试手册

> **版本**：v1.0（2026-06-03）  
> **适用对象**：产品验收、答辩演示、样式走查、向 AI 助手提交修改需求  
> **关联文档**：**[step-by-step-uat-guide.md](./step-by-step-uat-guide.md)**（傻瓜式 62 步，推荐）· [human-full-experience-uat.md](./human-full-experience-uat.md) · [test-plan.md](./test-plan.md)  
> **Cursor Skill**：`.cursor/skills/forum-web-uat/`

---

## 目录

1. [测试前准备](#1-测试前准备)
2. [功能总览与路由地图](#2-功能总览与路由地图)
3. [分角色逐步测试（含预期结果）](#3-分角色逐步测试)
4. [样式与多端专项](#4-样式与多端专项)
5. [自动化与 API 冒烟](#5-自动化与-api-冒烟)
6. [向 AI 提交修改建议的标准格式](#6-向-ai-提交修改建议的标准格式)
7. [系统自检报告（静态 + 构建）](#7-系统自检报告)

---

## 1. 测试前准备

### 1.1 环境二选一

| 模式 | 访问地址 | 说明 |
|------|----------|------|
| 本地开发 | `http://127.0.0.1:8091` | `.\start-dev.ps1` 或分别启动三 Go 服务 + `npm run dev` |
| Docker / 云主机 | `http://localhost/` 或部署域名 | `docker compose up` 或见 [DEPLOY.md](../DEPLOY.md) |

### 1.2 健康检查（必须先过）

在终端执行（本地直连）：

```powershell
Invoke-WebRequest http://127.0.0.1:8001/health -UseBasicParsing
Invoke-WebRequest http://127.0.0.1:8002/health -UseBasicParsing
Invoke-WebRequest http://127.0.0.1:8003/health -UseBasicParsing
```

三者均应返回 **200**。若失败，先启动服务，不要开始功能测试。

### 1.3 数据库与种子

| 步骤 | 命令/文件 | 说明 |
|------|-----------|------|
| Schema 初始化 | `migrations/init/000_init_schemas.up.sql` | 首次一次 |
| 各服务迁移 | `services/*/migrations/*.up.sql` | 含 QQ OAuth `006_oauth_identities` |
| GX 板块种子 | `scripts/seed/002_pilot_gx_content.sql` | **强烈建议执行**，否则顶栏「学业研讨」等 slug 与 DB 不一致 |

### 1.4 演示账号

完整列表见 [demo-accounts.md](./demo-accounts.md)。

| 角色 | 用户名 | 密码 | 登录方式 |
|------|--------|------|----------|
| 学生 | `demo_student` | `demo123456` | 登录页输入学号/用户名 |
| 学生（并行） | `demo01` … `demo06` | 与用户名相同 | 同上 |
| 协会管理员 | `demo_admin` | `demo123456` | 同上，自动进 `/admin` |
| 协会管理员（并行） | `admin01` … `admin06` | 与用户名相同 | 同上 |
| 中台管理员 | `demo_platform_admin` | `demo123456` | 同上 |
| 中台管理员（并行） | `plat01` … `plat06` | 与用户名相同 | 同上 |
| 邀请码注册 | 新用户名自填 | 8 位+大小写+数字 | 邀请码 `DEMO2026`（演示环境） |

**注意**：生产环境 `demo-login` API 默认对公网 IP 拒绝；演示按钮仅 `import.meta.env.DEV` 前端可用，生产请用账号密码或 QQ 登录。

### 1.5 QQ 登录（可选）

仅在 `.env` 配置 `QQ_APP_ID`、`QQ_APP_KEY`、`QQ_REDIRECT_URI` 后测试：

1. QQ 互联回调地址与 `QQ_REDIRECT_URI` **完全一致**
2. 执行迁移 `006_oauth_identities.up.sql`
3. 登录页点「QQ 登录」→ 授权 → 回跳 `/oauth/qq` → 进入社区

---

## 2. 功能总览与路由地图

### 2.1 公开 / 认证入口

| 路由 | 页面 | 必测 |
|------|------|------|
| `/` | 登录 | ✓ |
| `/register` | 邀请码注册 | ✓ |
| `/oauth/qq` | QQ 回调落 token | 配置后测 |

### 2.2 社区（需登录，`demo_student`）

| 路由 | 功能 | 必测 |
|------|------|------|
| `/community` | 首页帖子流、热榜排序 | ✓ |
| `/community/boards/:slug` | 板块帖子列表 | ✓ |
| `/community/posts/new` | 发帖（含附件） | ✓ |
| `/community/posts/:id` | 详情、评论、点赞、收藏、举报 | ✓ |
| `/community/posts/:id/edit` | 编辑自己的帖 | ✓ |
| `/community/profile` | 个人资料、头像 | ✓ |
| `/community/users/:id` | 他人公开主页 | ✓ |
| `/community/messages` | 通知中心 | ✓（建议双账号） |
| `/community/my/posts` | 我的帖子 | ✓ |
| `/community/my/favorites` | 我的收藏 | ✓ |
| `/community/my/history` | 浏览历史 | ✓ |
| `/community/circle` | 校园圈（`campus-circle` 板块） | ✓ |
| `/community/rank` | 排行榜 | ✓ |
| `/community/about` | 帮助/关于 | 样式 |

### 2.3 管理端

| 路由 | 功能 | 角色 |
|------|------|------|
| `/admin` | 概览统计 | admin / platform_admin |
| `/admin/audit` | 待审核 | admin+ |
| `/admin/reports` | 举报处理 | admin+ |
| `/admin/posts` | 精华/置顶/删帖 | admin+ |
| `/admin/users` | 用户封禁、等级 | admin+ |
| `/admin/boards` | 板块 CRUD | admin+ |
| `/admin/config` | 系统配置 | admin+ |
| `/admin/invites` | 邀请码 | **仅 platform_admin** |
| `/admin/sensitive` | 敏感词 | **仅 platform_admin** |
| `/admin/roles` | 角色分配 | **仅 platform_admin** |
| `/admin/stats` | 统计图表 | admin+ |

**权限快测**：`demo_student` 访问 `/admin` → 应回到社区首页，不能看中台。

---

## 3. 分角色逐步测试

> 每一步格式：**操作 → 预期结果 → 失败时记录为缺陷**（用第 6 节模板）。

### 3.1 学生 — 注册与登录（约 10 分钟）

| # | 操作 | 预期结果 |
|---|------|----------|
| S-01 | 打开 `/register`，邀请码 `DEMO2026`，填写新用户名与强密码 | 注册成功，可登录 |
| S-02 | 邀请码填 `INVALID` | 提示邀请码无效 |
| S-03 | 用户名与已有账号重复 | 提示用户名已存在 |
| S-04 | 密码 `123` | 前端或后端拒绝，有强度提示 |
| S-05 | `/` 用 `demo_student` / `demo123456` 登录 | 进入 `/community`，顶栏显示用户 |
| S-06 | 错误密码连续 5 次 | 触发限流提示（约 1 分钟） |
| S-07 | 退出登录（个人中心或菜单） | 回到 `/`，受保护路由需重新登录 |
| S-08 | （可选）QQ 登录完整流程 | 回跳成功，`/community` 可浏览 |

### 3.2 学生 — 资料与浏览（约 10 分钟）

| # | 操作 | 预期结果 |
|---|------|----------|
| S-10 | `/community/profile` 改昵称、院系、年级、简介并保存 | 刷新后仍保留 |
| S-11 | 上传 jpg/png 头像 | 头像更新（或刷新后可见） |
| S-12 | 点击帖子作者进入 `/community/users/:id` | 只读，无编辑按钮 |
| S-13 | 首页列表 | 有标题、作者、时间、点赞/评论数；日期为 `YYYY-MM-DD` 非 ISO 原串 |
| S-14 | 切换顶栏板块（学业研讨等） | 列表随板块变化 |
| S-15 | 分页到第 2 页 | 内容与第 1 页不同 |
| S-16 | 进入帖子详情 | 正文、评论、附件链接可用 |

### 3.3 学生 — 发帖与互动（约 15 分钟）

| # | 操作 | 预期结果 |
|---|------|----------|
| S-20 | `/community/posts/new` 选板块，填标题+正文提交 | 列表可见新帖（或待审核提示） |
| S-21 | 标题或内容为空提交 | 拦截并提示 |
| S-22 | 正文含 `<script>alert(1)</script>` | **不弹窗**，内容被转义/过滤 |
| S-23 | 编辑自己的帖 | 保存后详情更新 |
| S-24 | 删除自己的帖 | 列表不再出现 |
| S-25 | 他人帖子 | 无编辑/删除入口 |
| S-26 | 发表评论 | 立即出现在列表 |
| S-27 | 空评论 | 被拦截 |
| S-28 | 点赞 → 再点取消 | 计数 +1 再 -1 |
| S-29 | 收藏 → 「我的收藏」可见 → 取消 | 状态与列表一致 |
| S-30 | 举报他人帖子并填原因 | 成功提示 |
| S-31 | 上传 jpg/pdf 附件发帖 | 详情可下载；上传 exe 应失败 |

### 3.4 学生 — 通知与扩展页（约 10 分钟）

| # | 操作 | 预期结果 |
|---|------|----------|
| S-40 | 账号 B 评论账号 A 的帖 | A 的 `/community/messages` 有新通知 |
| S-41 | 点开通知 | 标记已读，未读数减少 |
| S-42 | `/community/circle` | 校园圈内容加载正常 |
| S-43 | `/community/rank` | 排行榜有数据或空态文案合理 |
| S-44 | `/community/my/history` | 浏览过的帖有记录 |

### 3.5 协会管理员 — `demo_admin`（约 20 分钟）

| # | 操作 | 预期结果 |
|---|------|----------|
| A-01 | 登录后进入 `/admin` | 概览有用户数、帖子数等 |
| A-02 | `demo_student` 访问 `/admin` | 拒绝，回社区 |
| A-03 | `/admin/audit` 通过一条待审 | 帖变为公开 |
| A-04 | 驳回并填理由 | 状态为已驳回 |
| A-05 | `/admin/posts` 设精华、置顶 | 社区列表有标识/置顶 |
| A-06 | 管理员删除任意帖 | 帖消失 |
| A-07 | `/admin/users` 封禁测试用户 | 无法登录或发帖 |
| A-08 | 解封 | 恢复正常 |
| A-09 | `/admin/boards` 禁用板块 | 社区侧栏/列表不显示 |
| A-10 | `/admin/reports` 处理举报 | 状态变为已处理 |

**制造待审数据**：中台将 `sensitive_word_action` 设为 `pending_review`，添加敏感词后用学生号发帖。

### 3.6 中台管理员 — `demo_platform_admin`（约 15 分钟）

| # | 操作 | 预期结果 |
|---|------|----------|
| P-01 | `/admin/invites` 生成 1 个码 | 列表出现新记录 |
| P-02 | 批量生成 10 个 | 数量增加 |
| P-03 | 作废未使用码后用其注册 | 注册失败 |
| P-04 | `/admin/sensitive` 增删敏感词 | 发帖行为随之变化 |
| P-05 | `/admin/roles` 给用户加 `admin` | 重新登录可进中台 |
| P-06 | 移除角色 | 无法再进中台 |
| P-07 | `/admin/config` 提高发帖等级门槛 | 低等级学生发帖被拒 |
| P-08 | `demo_admin` 访问 `/admin/invites` | 被重定向到概览（非 platform_admin） |

### 3.7 安全抽检（约 10 分钟）

与 [test-plan.md §5–6](./test-plan.md) 一致，至少做：

- XSS：帖标题、评论、个人简介各测 1 次
- 越权：学生 token 调 `POST /admin-api/...` → 403
- 恶意文件：`.exe` 上传 → 拒绝

自动化：`bash scripts/security-smoke.sh`（Git Bash / WSL）。

---

## 4. 样式与多端专项

### 4.1 桌面（Chrome / Edge）

| 检查项 | 通过标准 |
|--------|----------|
| 登录页 | 双栏布局不溢出，错误提示可见 |
| 社区首页 | 卡片对齐、标签颜色与板块语义一致 |
| 帖子详情 | 评论区宽度、按钮可点 |
| 管理端表格 | 小屏横向可滚动，操作列不被裁切 |

### 4.2 手机（Safari iOS + Chrome Android）

| 检查项 | 通过标准 |
|--------|----------|
| 底部 Tab | 首页/板块/发帖/消息/我的可切换 |
| 抽屉菜单 | 可打开并导航 |
| 发帖页 | 键盘不遮挡提交按钮 |
| 管理端 | 核心列表可读（可接受横向滚动） |

### 4.3 样式问题记录方式

样式类问题统一标 **类型=UI**，严重度多为 **P3**；若导致按钮无法点击则为 **P1**。

---

## 5. 自动化与 API 冒烟

在 **三服务已启动** 的前提下（Git Bash / WSL / Linux）：

```bash
# 30 秒冒烟
bash scripts/smoke-test.sh

# 5 分钟全量 API（写操作会创建测试数据）
bash scripts/full-experience-test.sh

# 网关模式（只读 + 登录）
BASE_URL=http://127.0.0.1:8091 bash scripts/full-experience-test.sh
```

**当前缺口**：脚本未覆盖 QQ OAuth；需按 §1.5 手工测。

通过标准建议：

- `smoke-test.sh`：全部 PASS  
- `full-experience-test.sh`：≥ 95% PASS（本地全量模式）  
- 所有 **P0/P1** 手工项已关闭

---

## 6. 向 AI 提交修改建议的标准格式

复制下面模板填写（可一次多条）。**信息越完整，AI 越能一次改对。**

```markdown
## 反馈单 #<序号>

### 元信息
- **类型**：BUG | UI | 文案 | 功能建议 | 安全
- **严重度**：P0 | P1 | P2 | P3
- **环境**：本地 8091 | Docker 80 | 云主机 <域名>
- **角色/账号**：demo_student | demo_admin | 新注册 xxx
- **浏览器**：Chrome 136 / iOS Safari 18 / …

### 现象（必填）
用一句话说明「实际看到了什么」。

### 复现步骤（必填，编号列表）
1. 打开 …
2. 点击 …
3. 输入 …
4. 观察到 …

### 预期结果（必填）
应该发生什么。

### 实际结果（必填）
实际发生了什么（含 HTTP 状态码、报错原文更佳）。

### 证据（强烈建议）
- 截图路径或描述
- 浏览器 Network：请求 URL、Method、Status、Response 片段（勿贴 Token）
- 控制台 Console 报错（如有）

### 修改诉求（必填，告诉 AI 要做什么）
- [ ] 修复逻辑  （例：学生不应看到编辑按钮）
- [ ] 调整样式  （例：移动端发帖按钮被键盘挡住）
- [ ] 改文案      （例：错误提示改为中文「邀请码无效」）
- [ ] 新增能力  （例：QQ 登录失败时显示具体原因）

### 影响范围（可选）
社区首页 / 管理端审核 / 仅 QQ 登录 / …

### 相关文件猜测（可选，不懂可留空）
frontend/src/views/PostDetail.vue
```

### 6.1 严重度定义（与 AI 对齐）

| 等级 | 含义 | AI 处理优先级 |
|------|------|----------------|
| P0 | 不可用或安全漏洞可利用 | 立即修 |
| P1 | 主流程受损，有绕行 | 本轮必修 |
| P2 | 次要功能异常 | 排期修 |
| P3 | 体验/样式/文案 | 批量优化 |

### 6.2 示例（合格反馈）

```markdown
## 反馈单 #1
- **类型**：BUG
- **严重度**：P1
- **环境**：本地 http://127.0.0.1:8091
- **角色**：demo_student

### 现象
收藏后「我的收藏」列表为空。

### 复现步骤
1. 登录 demo_student
2. 打开任意帖子，点收藏，按钮变为已收藏
3. 侧栏进入「我的收藏」

### 预期
刚收藏的帖子出现在列表中。

### 实际
列表为空，Network 中 GET /forum-api/api/v1/me/collections 返回 200 但 posts: []。

### 修改诉求
- [x] 修复逻辑：核对 collect API 与 collections 列表查询是否同一 user_id。
```

---

## 7. 系统自检报告

> **执行时间**：2026-06-03  
> **范围**：本机静态构建 + 代码审查；**未**在会话中完成对你方运行中服务的全量 API 回归（需你启动服务后跑 §5 脚本）。

### 7.1 构建与编译 — 通过

| 检查项 | 结果 |
|--------|------|
| `services/user` `go build ./cmd` | ✅ 通过 |
| `services/forum` `go build ./cmd` | ✅ 通过 |
| `services/admin` `go build ./cmd` | ✅ 通过 |
| `frontend` `npm run build` | ✅ 通过（Vite 无 TS/打包错误） |

### 7.2 已发现问题（建议优先处理）

| ID | 严重度 | 问题 | 说明与建议 |
|----|--------|------|------------|
| SC-01 | P2 | **GX 板块种子未执行时导航错位** | 仅 `001` 迁移时 DB 为 `ai-learning` 等，顶栏链接为 `study`/`training` 等；`resolveBoardByKey` 会 fallback 到第一个板块，体验混乱。**验收前执行** `scripts/seed/002_pilot_gx_content.sql`。 |
| SC-02 | P2 | **中台概览待审数可能恒为 0** | `frontend/src/api.js` 中 `getOverview` 对 `audit/pending` 使用 `.catch(() => ({ posts: [] }))`，接口失败时静默为 0，不利于发现审核服务问题。 |
| SC-03 | P3 | **死代码：AdminService 内 TODO 桩** | `services/admin/internal/service/admin_service.go` 中 `GetConfig`/`UpdateConfig`/`ListPendingAudit` 为未实现桩；实际路由使用 `ConfigService`/`AuditService`，**不影响运行**，可后续删除避免误导。 |
| SC-04 | — | **QQ OAuth 未纳入自动化脚本** | `scripts/full-experience-test.sh` 无 QQ 流程；需手工 + 配置 `.env`。 |
| SC-05 | — | **生产 demo-login 限制** | 非 `APP_ENV=development` 且非内网 IP 时，`POST /demo-login` 返回 403；生产演示需账号密码或白名单 `DEMO_LOGIN_ALLOWLIST`。 |
| SC-06 | — | **迁移 006 需手动应用** | 启用 QQ 登录前必须跑 `006_oauth_identities.up.sql`，否则 OAuth 绑号失败。 |

### 7.3 未发现阻断级编译问题

- 路由守卫：`/admin` 对非管理员重定向、`platform_admin` 独占 invites/sensitive/roles 与代码一致。  
- QQ 回调页 `OAuthQQ.vue`：从 query 取 token 写入 session 逻辑完整。  
- 三服务与前端构建均无错误。

### 7.4 建议你本地补跑的验证

```bash
bash scripts/smoke-test.sh
bash scripts/full-experience-test.sh
bash scripts/security-smoke.sh
```

将终端 **FAIL** 行按 §6 模板反馈，我可据此做针对性修复。

---

## 附录：与现有文档的关系

| 文档 | 用途 |
|------|------|
| **本文 test-manual.md** | 操作步骤 + AI 反馈格式 + 自检结论 |
| test-plan.md | 测试方法论、风险地图、安全矩阵 |
| manual-walkthrough-checklist.md | 可打印的勾选清单（45–60 分钟） |

---

**测试人**：________  
**日期**：________  
**环境**：□ 本地 □ Docker □ 云主机  
**自动化结果**：smoke ___/___ · full ___/___ · security ___/___
