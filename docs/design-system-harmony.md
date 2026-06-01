# Harmony 设计系统执行规范（v1）

> 依据 [harmony-design-audit-vlad-perspective.md](./harmony-design-audit-vlad-perspective.md)

## 色彩分层

| 层级 | Token | 用途 |
|------|-------|------|
| L0 | `--color-bg` / `--color-surface` / `--color-muted` | 背景、边框、次要文字 |
| L1 | `--color-primary` / `--color-brand` | 导航 active、主按钮 |
| L2 | `--color-danger` | 举报、删除、未读（Badge `danger` / `accent`） |
| L2 | `--color-warning` | 审核待办、警告 |
| L2 | `--color-success` | 成功提示 |
| L3 | `--color-gold` + 板块 Badge | 置顶、公告板块标签 |

**禁止**：`accent` 红用于板块 Badge、置顶条；置顶使用 `gold` 或左边框 `.gx-post-row--pinned`。

## 排版五级

| 类名 | 用途 |
|------|------|
| `text-display` | 帖子标题、板块 H1 |
| `text-title` | 区块标题 |
| `text-body` | 正文、表单 |
| `text-meta` | 日期、统计、面包屑 |
| `text-caption` | 页脚说明、辅助 |

## 时间格式

- 列表/侧栏：`formatDisplayDate` → `YYYY-MM-DD`
- 评论/消息：`formatDisplayTime` → 短日期+时分
- API 层保留 `createdAtIso` 供 `<time datetime>`

## 禁止事项

- 页面内随意 `max-w-2xl` 破坏 `--gx-content-max`
- UI 直接展示原始 ISO（含 `T17:53`）
- 首页首屏轮播（已移除）
- Card 套 Card（管理端/表单页单层 Card）
