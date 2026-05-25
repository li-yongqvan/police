# R1–R3 执行计划书（AI 智联论坛 · 学院 MVP）

> **基线**：Phase 0–4 全栈集成已完成；生产入口 `http://107.172.138.10/`  
> **执行日期**：2026-05-24  
> **范围**：R1 工程微调 + R2 补齐 API/UI 缺口 + R3 体验抛光（**不含 R4 迁栈/SSR**）

---

## 1. 目标摘要

| 阶段 | 目标 | 状态 |
|------|------|------|
| **R1** | 并行加载、去掉 token 穿透、配置保存单次拉板块 | ✅ 已完成 |
| **R2** | 板块编辑、驳回理由、批量审核删除、邀请码状态、他人主页 | ✅ 已完成 |
| **R3** | 分页、点赞/收藏反馈、友好错误文案、API 契约文档 | ✅ 已完成 |

---

## 2. R1 — 数据加载与 API 层

### 2.1 任务清单

- [x] `frontend/src/composables/usePageLoad.js` — `loadPage({ key: loader })` 并行 `Promise.all`
- [x] `api.js` 统一 `authHeaders()` 从 `localStorage` 读 token，视图不再传 `session.token`
- [x] `CommunityHome`：`loadPage` 并行 boards + posts；仅 `canAccessAdmin` 时请求 `getOverview`
- [x] `adminApi.updateConfig`：循环前 **一次** `getBoards(true)`
- [x] `api/http.js` + `api/errors.js`：`formatApiError` 统一错误文案
- [ ] 拆分 `api.js` 为多文件（**延后**，当前单文件可维护）

### 2.2 验收

1. 学生登录打开 `/community`：Network 中 boards、posts **并行**，**无** `/admin/stats/overview`
2. 管理员首页：boards、posts、overview **并行**
3. 保存系统配置：boards 接口 **仅 1 次**（非每板块一次）

---

## 3. R2 — 有 API 无 UI 的补齐

| ID | 任务 | 实现位置 | 状态 |
|----|------|----------|------|
| R2-1 | 板块行内编辑 | `AdminBoards.vue` → `adminApi.updateBoard` | ✅ |
| R2-2 | 审核驳回自定义理由 | `AdminAudit.vue` → `prompt` + `rejectAudit(id, reason)` | ✅ |
| R2-3 | 待审批量删除 | `AdminAudit.vue` 勾选 → `batchDeleteAudit` | ✅ |
| R2-4 | 邀请码状态查询 | `AdminInvites.vue` → `getInviteCodeStatus` | ✅ |
| R2-5 | 他人公开主页 | `UserPublic.vue` + `/community/users/:id` | ✅ |

**未做（刻意）**：`GET /admin/config/:key` 单键读取 — 列表页已用聚合 `getConfig`。

---

## 4. R3 — 列表与反馈抛光

- [x] `forumApi.getPosts({ page, limit, boardId })` + 社区首页/板块「加载更多」
- [x] `adminApi.listPosts` / `userApi.listUsers` 分页 + 中台上一页/下一页
- [x] `likePost` 返回计数后 toast；`collectPost` 展示后端 message
- [x] `formatApiError`：等级 / 敏感词 / 审核 / 登录过期
- [x] `shared/api-contract.md` 字段对照表

---

## 5. 关键文件索引

| 文件 | 变更 |
|------|------|
| `frontend/src/api.js` | 无 token 参数；新 batch/reject/invite/profile API |
| `frontend/src/composables/usePageLoad.js` | 新增 |
| `frontend/src/api/errors.js` | 新增 |
| `frontend/src/views/CommunityHome.vue` | 并行 + 分页 |
| `frontend/src/views/AdminAudit.vue` | 驳回理由 + 批量删除 |
| `frontend/src/views/AdminBoards.vue` | 编辑表单 |
| `frontend/src/views/AdminInvites.vue` | 状态查询 |
| `frontend/src/views/UserPublic.vue` | 他人主页 |

---

## 6. 演示剧本（R1–R3 DoD）

| # | 旅程 | 通过标准 |
|---|------|----------|
| 1 | 学生刷首页 | 无 admin overview 请求；帖子可「加载更多」 |
| 2 | 管理员审核驳回 | 弹出理由，非写死「演示驳回」 |
| 3 | 批量删除待审帖 | 勾选 ≥2 篇，调用 batch-delete 成功 |
| 4 | 编辑板块名称 | 保存后列表刷新 |
| 5 | 邀请码状态查询 | 输入码显示 status |
| 6 | 帖子详情 → 作者链接 | 打开 `/community/users/:id` |
| 7 | Token 过期 | refresh 或友好「请重新登录」 |

---

## 7. 部署提示

```bash
cd frontend && npm run build
# 服务器 /opt/ai-forum
docker compose -f docker-compose.yml -f docker-compose.server.yml build frontend
docker compose -f docker-compose.yml -f docker-compose.server.yml up -d frontend
```

---

## 8. 后续（R4，非本次）

- Vitest + MSW 契约测试
- 可选 SvelteKit 子应用 / 首页 SSR
- OpenAPI 生成类型替代手写 mapper

---

*关联：[rich-harris-fullstack-audit-plan.md](./rich-harris-fullstack-audit-plan.md) · [fullstack-integration-audit-plan.md](./fullstack-integration-audit-plan.md)*
