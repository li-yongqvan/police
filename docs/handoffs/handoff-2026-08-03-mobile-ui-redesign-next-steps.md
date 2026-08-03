# Handoff - 2026-08-03 - Mobile UI Redesign Next Steps

> 新会话建议先读本文，然后说：继续按移动端 UI 修改意见实施。

## 1. 现在做了什么

- 当前目标：基于一次只读移动端实际体验，整理下一步移动端 UI 修改意见；本阶段不改业务代码、不提交、不部署。
- 当前分支：`codex/discourse-rebuild`。
- 当前项目状态：工作区已有大量历史/用户改动，本次只新增 handoff 文档；不要回滚或清理无关改动。
- 最新验收口径：移动端优先，覆盖本地 Vue 登录/管理端、线上平台入口、外部 Discourse 论坛；线上只读，不发帖、不回复、不点赞、不封禁、不删除、不改配置。

## 2. 已经完成了什么

- 已完成移动端只读体验，覆盖视口：
  - `375x667`
  - `390x844`
  - `412x915`
  - `667x375` 横屏
- 已确认本地与线上平台登录页可访问：
  - 本地：`http://127.0.0.1:8091/`
  - 线上平台：`http://122.51.233.225:8888/`
  - Discourse：`http://122.51.233.225:8080/`
- 学生链路：使用项目迁移数据中识别出的 QQ 号账号做了一次只读登录验证，成功进入 Discourse；完整账号和密码不写入 handoff。
- 管理员链路：线上演示管理员账号单次登录失败，页面提示“用户名或密码错误”；已按规则停止试错。
- 管理端 UI：用本地 mock 登录态只读打开 `/admin`、`/admin/users`、`/admin/invites`，验证移动端布局。
- 证据文件：
  - `.scratch/mobile-ux/report.json`
  - `.scratch/mobile-ux/admin-mock-report.json`
  - `.scratch/mobile-ux/online-login-iphone-14.png`
  - `.scratch/mobile-ux/student-login-result.png`
  - `.scratch/mobile-ux/local-admin-mock-admin-iphone-14.png`
  - `.scratch/mobile-ux/local-admin-mock-admin-users-iphone-14.png`
  - `.scratch/mobile-ux/local-admin-mock-admin-invites-iphone-14.png`
  - `.scratch/mobile-ux/local-admin-mock-drawer-open.png`

## 3. 卡在哪里

- 线上管理员账号口径不确定：`demo_admin` 登录失败。下一轮如需真实线上管理端验证，先向用户确认有效管理员账号，不要继续猜密码。
- Discourse 主题属于外部论坛侧，不完全在 `frontend/` Vue 代码内；要区分“平台前端可改”和“Discourse 主题/配置可改”。
- 自动报告里的中文文本在 PowerShell/JSON 输出中有编码污染，但截图显示页面中文正常；不要把报告文本乱码误判为真实页面乱码。
- 暂无代码层阻塞；主要是需要按优先级实施 UI 修复并重新做移动端截图验收。

## 4. 下一步做什么

1. 修复管理端移动抽屉错位，优先级 P1。
   - 重点检查 `frontend/src/styles/gx-theme.css` 中重复的 `.gx-admin-shell .gx-admin-sidebar` 移动端规则。
   - 当前证据：抽屉打开后 `drawerRect.left = -24`，内容被裁切；截图为 `.scratch/mobile-ux/local-admin-mock-drawer-open.png`。
   - 建议让移动抽屉显式从 `left: 0` 开始，打开态 `transform: translateX(0)`，避免后续规则、阴影或旧兼容 CSS 把抽屉视觉起点推到屏幕外。
   - 验收：390x844 下抽屉左边界为 0，导航文字完整可读，body 保持 `mw-drawer-open` 滚动锁。

2. 补齐管理端顶部触控尺寸，优先级 P2。
   - 目标文件：`frontend/src/views/AdminLayout.vue` scoped style，必要时同步全局移动端 token。
   - 当前指标：菜单按钮约 `36x30`，退出按钮约 `56x34`，低于 44px 触控基线。
   - 建议把菜单按钮和退出按钮统一为 `min-width: 44px; min-height: 44px;`，并保持品牌文字不挤压。
   - 验收：390x844 和 375x667 下按钮不换行、不遮挡，点击区域不低于 44px。

3. 优化登录页横屏首屏可操作性，优先级 P2。
   - 目标文件：`frontend/src/components/gx/GxAuthShell.vue`。
   - 当前证据：`667x375` 横屏时表单输入区在首屏下半部分，需要滚动才能开始登录。
   - 建议新增短高视口媒体查询，例如 `@media (max-height: 480px) and (orientation: landscape)`，压缩品牌图标、标题间距、面板 padding，确保账号输入框在首屏可见。
   - 验收：667x375 下至少账号输入框和密码输入框首屏可见；不破坏 390x844 竖屏视觉。

4. 整理 Discourse 移动主题问题，优先级 P1/P2，单独 ticket。
   - 当前证据：Discourse 首页多处按钮/链接高度低于 40px，顶部 logo/头像区域显示异常，帖子列表首屏信息层级弱。
   - 如果项目继续使用 Discourse，则应在 `discourse-themes/ai-forum-premium-preview/` 或对应主题包中处理，不要在 Vue 登录页里硬改论坛体验。
   - 建议方向：移动端 header 固定高度、登录按钮 44px、搜索/菜单图标 44px、话题列表标题与摘要行数约束、首屏品牌/分类入口可读。

5. 验证与交付。
   - 本地启动：`cd frontend && npm run dev -- --host 127.0.0.1 --port 8091`
   - 跑移动截图验收，至少覆盖 `/`、`/admin`、`/admin/users`、`/admin/invites` 的 390x844 和 667x375。
   - 如果改 Discourse 主题，单独打包主题 zip 并在线上 Discourse 预览后再部署。
   - 完成后按项目约定提交、推送和部署；但只有在用户确认进入实施阶段后再做。

## 5. 哪些坑不要再踩

- 不要继续尝试多个管理员密码。线上管理员账号当前阻塞，需要用户确认有效口径。
- 不要把 `.scratch/mobile-ux/*.json` 中的中文乱码当成页面真实乱码；以截图和浏览器可见文本为准。
- 不要在未确认前对线上数据执行写操作，包括发帖、回复、点赞、收藏、封禁、删除、改等级、生成或作废邀请码。
- 不要把 Discourse 论坛 UI 问题塞进 Vue 登录页修改；论坛主题和平台前端是两个改动边界。
- 不要回滚当前工作区已有的大量历史改动；本 handoff 只新增文档。
- 不要把截图和 `.scratch/` 临时产物提交，除非用户明确要求把体验证据纳入仓库。
