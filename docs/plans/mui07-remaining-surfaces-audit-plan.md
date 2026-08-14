# MUI-07 继续审计 · 未覆盖页面元素收尾扫描

> 日期：2026-08-14
> 状态：已执行完成（2026-08-14，20/20 场景通过，无新缺陷）

## 1. 背景

MUI-07 修复（/top period-chooser 收窄 + 15px 字号）已上线并全绿。但最终验证脚本与计划书 §6 对照仍有遗漏：signup/aq/	os 与 /c/ai/5 在修复后未复测，另有若干从未扫过的表面。本次补齐并收尾。

## 2. 审计范围（修复后基线）

### 2.1 补测计划书承诺场景
- 登出态：/signup 390/320、/faq（跳 /guidelines）390/320；/tos 已 404，跳过并在结论中说明。
- 登录态：/c/ai/5 390/375。

### 2.2 新增未覆盖表面
- /tags 390/320（实例暂无真实标签，只扫列表页）。
- /g 分组列表 390/320 + /g/moderators/members 390。
- /u/123/messages 私信 inbox 390/320（仅测指标；若收件箱非空则跳过截图，避免泄露私信内容）。
- 菜单展开态 320/375：/latest 汉堡菜单与搜索下拉（此前仅 390 测过）。
- 用户主页 320：/u/123/summary、/u/123/activity。
- 登录页 320（登出态；装饰球 gx-auth-spotify__orb 为既有白名单）。

## 3. 硬性门禁

- 每场景 docOverflow=0、odyOverflow=0、ad=0（除白名单装饰球）。
- 发现缺陷 → CSS-only 修复（MUI-07 区块）→ 部署 → 全量回归（MUI-07 31 项 + MUI-04/05/06 17 项）→ 前后截图 → commit/push。
- 全部通过 → 输出各场景截图归档 work/mui07-validation/，提交审计计划书并推送。

## 4. 交付

- 一次 commit：<type>: <description> (#MUI-07)，立即 push codex/discourse-rebuild。
- 若无需 CSS 修复：commit 类型 chore，只含计划书与审计结论。
## 5. 执行结果（2026-08-14）

- 扫描 20 个场景全部通过：docOverflow=0、odyOverflow=0、ad=0（登录/注册页装饰球 gx-auth-spotify__orb 为既有白名单，仅 3 场景出现且不计失败）。
- 场景明细：signup/guidelines/login 320、/c/ai/5 390/375、/tags 390/320、/g 390/320、/g/moderators/members 390、/u/123/messages 390/320（收件箱为空，截图安全）、latest 汉堡/搜索菜单 320/375、用户主页 summary/activity 320。
- /tos 已 404（实例未配置服务条款路由），按计划跳过。
- 结论：MUI-07 剩余表面无移动端溢出缺陷，无需新增 CSS；审计截图与指标归档 work/mui07-validation/（mui07-sweep-remaining.mjs、metrics-remaining.json、40-64 截图）。