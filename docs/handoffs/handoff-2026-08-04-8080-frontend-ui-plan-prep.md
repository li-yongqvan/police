# Handoff - 2026-08-04 - 8080 端口前端 UI 计划准备

> 新会话建议先读本文档，然后说：继续制定 8080 端口前端 UI 修改计划书。

## 1. 现在做了什么

- 当前主线任务已经从 8888 平台前端移动端 UI 修复，转入准备制定 8080 端口前端 UI 修改计划书。
- 当前本地分支：`codex/discourse-rebuild`。
- 用户最新要求：生成 handoff 文档，为后续“修改 8080 端口的前端 UI 计划书”做准备。
- 8080 端口对应的是 Discourse 侧界面，前面已部署过 Discourse preview theme 相关内容；后续计划书应先只做需求统一和方案制定，不要直接改代码。
- 用户偏好仍然有效：
  - 不确定时先问，不要继续执行。
  - 遇到突发情况先停止并汇报。
  - 如果工作区产生文件，必须说明。
  - 登录凭据不写入计划书/截图；需要验证时使用既有安全登录方式。
  - QQ 号账号只作为普通账号；管理员账号使用 `demo_admin`。

## 2. 已经完成了什么

- 8888 平台前端 MUI-01 已完成并部署：
  - 本地提交：`3ad6817 fix: improve mobile platform UI (#MUI-01)`
  - 服务器部署提交：`e84627e fix: deploy mobile platform UI MUI-01`
  - 已部署到：`http://122.51.233.225:8888/` 和 `/admin`
- 管理员密码登录角色问题已修复并云端部署：
  - 本地提交：`13c3b26 fix: include login user role (#AUTH-ROLE-01)`
  - 服务器部署提交：`45be065 fix: deploy login user role response (#AUTH-ROLE-01)`
  - 只重建并重启了 `user-service`
  - 验证结果：
    - `demo_admin` 登录响应 `role=admin`
    - QQ 普通账号 `2075773070` 登录响应 `role=student`
    - `user-service` healthy
- 开发记录日志已迁出到专用分支：
  - 分支：`codex/dev-records`
  - 已推送到 GitHub
  - 后续提交 dev-record 自身时，必须使用 `DEV_RECORD_SKIP=1` 防止 post-commit 递归生成日志。
- 已生成 MUI-01 前后截图：
  - 截图目录：`work/screenshots/mobile-ui-comparison-2026-08-04/`
  - 对比图包括登录竖屏、登录横屏、管理端顶栏、管理端抽屉。
- 为生成 before 截图创建了临时 worktree：
  - `work/mui-before`
  - 指向 MUI-01 前一版 `3ad6817^`。

## 3. 卡在哪里

- 暂无明确阻塞。
- 当前工作区有未跟踪的 `work/` 目录，内容包括：
  - `work/mui-before`：临时 Git worktree，用于 before 截图。
  - `work/screenshots/mobile-ui-comparison-2026-08-04`：截图、对比图、Vite 日志和 PID 记录。
- 这些文件尚未删除，因为属于批量生成内容，需用户确认后才能清理。
- 8080 端口 UI 修改尚未开始制定计划书，当前只完成了交接准备。

## 4. 下一步做什么

1. 先向用户确认 8080 端口 UI 计划书的范围：
   - 只改 Discourse theme，还是包含 8888 平台到 8080 的跳转一致性？
   - 优先移动端，还是移动端 + 桌面端同时规划？
   - 是否沿用 MUI-01 的原则：图标优先、单字兜底、保留 `aria-label`？
   - 是否需要先截图审计当前 8080 页面作为基线？
2. 读取与 8080/Discourse 相关的本地材料：
   - `discourse-themes/ai-forum-premium-preview/`
   - `docs/handoffs/handoff-2026-07-30-discourse-ticket7-done-feedback.md`
   - `docs/handoffs/handoff-2026-07-30-discourse-sso-ticket3b.md`
   - `docs/reddit-inspired-redesign-plan.md`
   - `docs/harmony-design-audit-vlad-perspective.md`
3. 对 8080 线上页面做只读截图审计：
   - `http://122.51.233.225:8080/`
   - 重点看移动端首页、话题列表、话题详情、登录/SSO入口、管理相关入口是否存在横向溢出、文字拥挤、按钮识别不清。
4. 形成计划书，不直接修改代码：
   - 按 Ticket 拆分。
   - 每个 Ticket 写清楚文件范围、验收截图、风险和是否需要部署。
   - 明确 8080 的 Discourse theme 部署路径和回滚方式。
5. 如果用户确认后再实施：
   - 先创建/使用 `codex/discourse-rebuild` 上的新 commit。
   - 每个 Ticket 单独提交并推送。
   - 云端部署前必须说明会产生哪些文件、如何推送、如何部署、如何回滚。

## 5. 哪些坑不要再踩

- 不要把本地整文件直接 `scp` 覆盖服务器文件，尤其是含中文或 CRLF/LF 差异的文件；优先用最小补丁。
- Windows PowerShell 执行远端 SSH 命令时，不要直接复制复杂 Linux shell；优先用 `@'...'@ | ssh host 'python3 -'` 或远端单行 Python，避免本地解析 `$()`、`curl`、引号和 CRLF。
- `git commit --no-verify` 不会跳过 `post-commit`；提交 dev-record 自身时用 `DEV_RECORD_SKIP=1`。
- Docker Compose 部署单服务时要注意 `depends_on` 和 `build:` 可能级联影响其它服务；只更新单服务时优先考虑 `--no-deps`，并评估是否需要重启 nginx。
- 服务器 `git fetch` GitHub 可能失败，之前遇到过 TLS 中断；如果再次失败，先汇报，不要自行扩大操作范围。
- 8080 是 Discourse 侧，不要把 8888 Vue SPA 的组件结构假设套过去；计划书要先确认 Discourse theme 的可改边界。
- 用户强调空间紧张时可保留代表性汉字或图标；但视觉只显示图标/单字时必须保留语义标签。
