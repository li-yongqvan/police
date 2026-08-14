# MUI-07 补充计划 · /top 日期筛选控件字号适配

> 状态：已确认执行（用户反馈：箭头已可见，但日期文字被省略号截断，要求缩小字号完整显示）
> 日期：2026-08-14

## 问题

MUI-07 修复让 /top period-chooser 的箭头回到框内后，selected-name 仍沿用 24px 的 h2 字号，完整日期串（如 年2025 年 8月 14 日 – 2026 年 8月 14 日）宽度 380px，在 390/375/320 视口（可用 323/308/253px）被 	ext-overflow: ellipsis 截断。

## 方案（CSS-only，mobile.scss MUI-07 区块）

- .selected-name ont-size: 24px → 15px。
  - 实测（Playwright 注入）：16px 在最宽日期串（12月31日跨度）下 320px 仍差 2px 截断；15px 在 320/375/390 全部完整显示，且留出约 17px 余量。
- .angle-icon ont-size → 15px（当前约 24px），与缩小后的文字协调，仍保持 20px 可见图标（不违反移动端控件可见性约定）。
- 保留 max-width + ellipsis 作为极端情况兜底；@media (max-width: 320px) 的 select-kit-body 规则不动。
- 作用范围：仅 mobile 字段（Discourse 移动视口），桌面 1440 不受影响。

## 验证

- 320/375/390 三档注入后复测：无截断、箭头在框内、docOverflow=0、odyOverflow=0；再用最宽日期串模拟 320px 边界。
- 上线后重跑 MUI-07 全场景复测（31 项）+ MUI-04/05/06 回归（17 项），确认不回退。
- 输出 /top 390px 与 320px 修复前后对比截图。

## 交付

- 提交纪律：一次 commit，格式 <type>: <description> (#MUI-07)，随 commit 立即 push codex/discourse-rebuild。