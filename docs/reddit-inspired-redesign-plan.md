# Reddit 风格借鉴 · AI智联平台改版计划书

> 概念稿效果图见 [`mockups/`](mockups/)：`reddit-style-home-desktop.png`、`reddit-style-board-desktop.png`、`reddit-style-post-detail-desktop.png`  
> 设计约束：Harmony 内容优先 + 警院校色（`#0f2b5b` / `#d4b86a` / 警示红仅用于举报）  
> 预览端口：**8091**

## 1. 调研结论（reddit.com）

| 维度 | Reddit 做法 | 本项目映射 |
|------|-------------|------------|
| 布局 | 左导航 + 居中窄 Feed + 右信息栏 | `GxSiteSidebar` + `GxFeedLayout`（720px）+ `gx-feed-layout__aside` |
| 帖子 | 左投票轨 + 标题/摘要/操作条 | `GxVoteRail` + `GxFeedPostCard` |
| 排序 | Hot / New / Top / Rising Pills | 热门 / 最新 / 精华 / 今日 → `sort=hot\|new\|featured\|today` |
| 评论 | 树状、可折叠 | `parent_id` + `GxCommentTree` |
| 品牌 | 橙红 + Snoo | **不采用**；沿用 `--color-primary` 警院蓝 |

## 2. 设计 Token（`gx-theme.css`）

| Token | 值 | 用途 |
|-------|-----|------|
| `--gx-feed-max` | `720px` | Feed 主列最大宽度 |
| `--gx-feed-bg` | `#f6f7f8` | Feed 区背景 |
| `--gx-vote-rail-w` | `40px` | 投票轨宽度 |
| `--gx-feed-card-radius` | `8px` | 帖子卡片圆角 |

## 3. 组件清单

| 组件 | 职责 |
|------|------|
| `GxVoteRail` | 赞/数字；列表与详情共用 |
| `GxFeedPostCard` | Feed 帖子单元（投票+meta+摘要+操作） |
| `GxFeedSortTabs` | 排序 Pill，`v-model` + URL `?sort=` |
| `GxFeedLayout` | 主列 + 右侧栏栅格 |
| `GxCommentTree` | 嵌套评论、折叠、回复 |

## 4. API

```
GET /api/v1/posts?sort=hot|new|featured|today&board_id=&page=&limit=&q=
POST /api/v1/posts/:id/comments  { "content", "parent_id"? }
```

评论列表返回 `parent_id`、`depth`、`author_name`。

## 5. 分期与状态

| 阶段 | 内容 | 状态 |
|------|------|------|
| A | 本文档 + `docs/mockups/*.png` | 已交付 |
| B | 前端 Feed 组件与首页/板块重构 | 已交付 |
| C | 后端 `sort` 参数 | 已交付 |
| D | 评论树迁移与 UI | 已交付 |

## 6. 验收标准

- [ ] 桌面 Feed 居中，卡片含投票轨与操作条
- [ ] 四种排序 Tab 切换且分页正确（URL 可分享 `?sort=hot`）
- [ ] 详情页 VoteRail 与列表点赞一致
- [ ] 评论支持回复与折叠（至少 2 层可见）
- [ ] 移动端单栏无横向溢出

## 7. 风险边界

- 不做踩帖（downvote）
- 不使用 Reddit 商标、Snoo、`r/` 前缀
- 首页不再以轮播为主视觉，数据进右侧栏
