# AI 智联论坛 — Rich Harris 视角全栈审视与优化计划书

> **审视视角**：Rich Harris（Svelte / Rollup / SvelteKit）— 少运行时、端到端数据流、编译期能确定的就不要留给浏览器猜。  
> **技术栈现实**：学院 MVP 使用 **Vue 3 + Vite + Pinia**（非 Svelte），本计划 **不建议为对接而整体迁栈**，而是把 SvelteKit 里「load 函数、薄客户端、路由级代码分割」的思想移植过来。  
> **基线日期**：2026-05-24（Phase 0–4 集成已完成；**R1–R3 已落地**，详见 [r1-r3-execution-plan.md](./r1-r3-execution-plan.md)）

---

## 1. 核心判断（一句话）

**浏览器可调用的后端能力，约 90% 已在 `api.js` 接线；约 75% 有对应页面；约 65% 达到可演示的完整交互。**  
剩余缺口主要是 **运营细节**（分页、批量审核、板块编辑、驳回理由）和 **工程卫生**（并行加载、类型契约、去掉 token 穿透），不是「大模块没做」。

---

## 2. Harris 式审视原则（映射到本项目）

| 原则 | 在本项目的含义 |
|------|----------------|
| **Less work in the browser** | 列表用 `board_id` 查询，不在前端过滤全站帖子；统计用后端聚合，不在前端重算 |
| **Colocated data loading** | 每个路由的 `onMounted` 应像 SvelteKit 的 `load`：声明依赖、并行拉取、错误边界清晰 |
| **Thin client** | `api.js` 只做 HTTP + 最小字段映射；业务规则留在 Go 服务 |
| **Explicit over magic** | 401 → refresh → 登出 已在 `api/http.js`；继续避免隐式全局副作用 |
| **Ship small chunks** | Vite 已按路由 lazy；主包 ~120KB gzip 尚可，避免再塞图表库 |
| **End-to-end types** | 当前最大技术债：手写 `mapUser` / `mapPost`，与 Go struct 易漂移 |

---

## 3. 前后端对接总表（浏览器可访问 API）

### 3.1 User 服务

| 后端路由 | 前端 API | 页面/流程 | 状态 |
|----------|----------|-----------|------|
| `POST /register` | `userApi.register` | `/register` | ✅ 完整 |
| `POST /login` | `userApi.login` | `DemoLogin` | ✅ 完整 |
| `POST /auth/refresh` | `api/http.js` 自动 | 全局 | ✅ 完整 |
| `POST /demo-login` | `userApi.demoLogin` | 无（`session.loginAs` 遗留） | ⚠️ 仅开发/备用 |
| `GET /users/me` | `userApi.me` | 全局 session | ✅ 完整 |
| `GET /users/:id` | `userApi.getProfile` | `/community/users/:id` | ✅ 公开主页 |
| `PUT /users/:id` | `userApi.updateProfile` | `/community/profile` | ✅ 完整 |
| `POST /users/:id/avatar` | `userApi.uploadAvatar` | profile | ✅ 完整 |

### 3.2 Forum 服务

| 后端路由 | 前端 API | 页面/流程 | 状态 |
|----------|----------|-----------|------|
| `GET /boards` | `forumApi.getBoards` | 社区/中台 | ✅ 完整 |
| `GET /boards/:id` | —（用 list + slug） | `BoardView` | ⚠️ 可用但多一次数据 |
| `GET /posts` | `forumApi.getPosts` | 首页/板块 | ✅ 完整 |
| `GET /posts/:id` | `forumApi.getPost` | 详情/审核预览 | ✅ 完整 |
| `GET /posts/:id/comments` | 合并在 `getPost` | 详情 | ✅ 完整 |
| `POST /posts` | `forumApi.createPost` | 发帖 | ✅ 含 `attachment_ids` |
| `PUT /posts/:id` | `forumApi.updatePost` | `/posts/:id/edit` | ✅ 完整 |
| `DELETE /posts/:id` | `forumApi.deletePost` | 详情（作者） | ✅ 完整 |
| `POST .../comments` | `forumApi.createComment` | 详情 | ✅ 完整 |
| `POST .../like` | `forumApi.likePost` | 详情 | ✅ 完整 |
| `POST .../collect` | `forumApi.collectPost` | 详情 | ✅ toast + 后端 message |
| `POST /attachments/upload` | `forumApi.uploadAttachment` | 发帖 | ✅ 完整 |
| `GET /attachments/:id` | `attachmentUrl()` | 详情链接 | ✅ 完整 |
| `GET /posts/:id/attachments` | —（详情内嵌） | — | ⚠️ 未单独拉取 |

### 3.3 Admin 服务（均需 JWT + 管理员角色）

| 后端路由 | 前端 API | 页面 | 状态 |
|----------|----------|------|------|
| `GET /stats/overview` | `adminApi.getOverview` | 概览 | ✅ 完整 |
| `GET /stats/daily` | `adminApi.getDailyStats` | 趋势统计 | ✅ 表格 |
| `GET/PUT /config` | get/update Config | 系统配置 | ✅ 完整 |
| `GET /config/:key` | — | — | ❌ 未用 |
| `GET /audit/pending` | `getPendingAudit` | 审核 | ✅ 完整 |
| `POST approve/reject` | approve/reject | 审核 | ✅ 自定义驳回理由 |
| `POST /audit/batch-delete` | `batchDeleteAudit` | 审核 | ✅ 勾选批量删除 |
| `GET /posts` | `listPosts` | 帖子管理 | ✅ 分页 UI |
| `POST delete/featured/pinned` | 有 | 帖子管理 | ✅ 完整 |
| `GET/POST/PUT/DELETE /boards` | list/create/**update**/delete | 板块管理 | ✅ 行内编辑 |
| `GET/POST ban/unban` | `setUserStatus` | 用户管理 | ✅ 完整 |
| `PUT /users/:id/level` | `updateUserLevel` | 用户管理 | ✅ 完整 |
| `GET /users/:id/logs` | `getUserLogs` | 用户管理 | ⚠️ 展示简陋 |
| 邀请码全套 | list/generate/batch/void | 邀请码 | ⚠️ 无单码状态查询页 |
| `GET /invite-codes/:code/status` | `getInviteCodeStatus` | 邀请码页 | ✅ |
| 敏感词 CRUD | 有 | 敏感词 | ✅ 完整 |
| 角色 CRUD | list/assign/remove | 角色权限 | ✅ 基础 |

### 3.4 刻意不对接（正确）

- 全部 `/internal/v1/*`：仅服务间调用，**不应**暴露给浏览器（Harris 也会赞成边界清晰）。

---

## 4. 代码可优化点（按收益排序）

### R-A — 数据加载（高，SvelteKit `load` 等价物）

| 问题 | 位置 | 建议 |
|------|------|------|
| 串行 `await` | `CommunityHome`：boards → posts → overview | 改为 `Promise.all`；学生角色不请求 overview |
| 重复 `getBoards(true)` | `updateConfig` 循环内多次拉板块 | 一次拉取传入 |
| 全量 `limit=100` | 帖子/用户列表 | 分页参数 +「加载更多」 |
| token 层层传递 | 每个 API 方法 `(token)` | `api/client.js` 从 Pinia 读 token，视图只调 `forumApi.getPosts()` |

**建议新增**：`frontend/src/composables/usePageLoad.js`

```js
// 用法类似 SvelteKit load: usePageLoad({ boards: () => forumApi.getBoards(), posts: () => forumApi.getPosts() })
```

### R-B — 模块与契约（中高）

| 问题 | 建议 |
|------|------|
| 单文件 `api.js` ~550 行 | 拆为 `api/user.js`、`api/forum.js`、`api/admin.js` + `api/mappers.js` |
| 手写 mapper 与 Go 漂移 | 增加 `shared/api-contract.md` 或 OpenAPI 生成 TS 类型（JSDoc 亦可） |
| `mapPost` 中 `tags: []` 写死 | 后端若返回 `tags` 则映射；否则 UI 隐藏标签表单项 |

### R-C — UI 与状态（中）

| 问题 | 建议 |
|------|------|
| 收藏无反馈 | `collectPost` 后 toast；若后端返回 `liked/collected` 则回显 |
| 发帖等级限制 | 上传附件需 level≥2：在 `NewPost` 根据 `currentUser.level` 提示 |
| 配置页 `postingEnabled` | 已与后端同步；保存后 `CommunityLayout` 应 `refresh` config |
| 审核驳回理由 | `AdminAudit` 增加 reason 输入框传入 `rejectAudit` |
| 板块管理无编辑 | 增加行内编辑或弹层，调用已有 `adminApi.updateBoard` |

### R-D — 构建与运行时（中低）

| 指标 | 当前 | Harris 目标 |
|------|------|-------------|
| 主包 gzip | ~45KB（`index-*.js`） | 保持 <50KB；admin 路由已 lazy |
| CSS | 单文件 ~14KB | 可接受 |
| 依赖 | 仅 vue + pinia + router | **保持零重型 UI 库**（正确） |
| Hydration | 纯 CSR | 学院演示足够；生产可考虑 SSR 仅首页（可选） |

### R-E — 死代码与一致性（低）

- 删除或 `import.meta.env.DEV` 包裹 `demoLogin` / `session.loginAs`
- `PostDetail` 附件字段统一用 `item.title`（已修）
- `AdminInvites` 列表字段与后端 `codes[]` 结构对齐（兼容 `code` 字符串 / 对象）

---

## 5. 覆盖率量化（Harris：用数字驱动，不凭感觉）

| 维度 | 数量 | 说明 |
|------|------|------|
| 浏览器可达 Admin API | 28 | `admin/cmd/main.go` 公开组 |
| 浏览器可达 User+Forum 写/读 API | 17 | 不含 internal |
| `api.js` 已封装方法 | ~42 | 含重复语义（deletePost 社区/中台） |
| 有路由的 Vue 页面 | 20 | 含 login/register |
| **API 有封装但无 UI** | **5** | batch-delete、config/:key、invite status、users/:id 公开、demo-login |
| **API 有 UI 但体验不完整** | **8** | 分页、板块编辑、驳回理由、收藏态、标签、日志展示、批量邀请展示、帖子管理筛选、他人主页 |

**结论**：相对 Osmani 审视（~60% 未接），本轮集成后 **核心演示路径已闭合**；剩余是 **产品完成度** 与 **工程整洁度**，不是架构性断线。

---

## 6. 分阶段优化计划（R1–R4）

### R1 — 「SvelteKit load」层 ✅（2026-05-24）

- [x] 新增 `usePageLoad` composable
- [x] API 层去掉 `token` 参数，统一从 localStorage 读取
- [x] `CommunityHome` 并行加载；学生不请求 overview
- [x] `updateConfig` 单次拉 boards
- [ ] 拆分 `api.js` 为按域模块（延后）

**验收**：Network 中首页 boards/posts 并行；学生无 overview 请求。

### R2 — 补齐「有 API 无 UI」✅

| 任务 | 状态 |
|------|------|
| R2-1 板块编辑 | ✅ `AdminBoards.vue` |
| R2-2 驳回理由 | ✅ `AdminAudit.vue` |
| R2-3 批量删除 | ✅ `batchDeleteAudit` |
| R2-4 邀请码状态 | ✅ `AdminInvites.vue` |
| R2-5 他人主页 | ✅ `UserPublic.vue` |

### R3 — 列表与反馈抛光 ✅

- [x] 帖子/用户/中台帖子分页
- [x] 点赞计数 toast；收藏 message
- [x] `formatApiError` 友好文案
- [x] `shared/api-contract.md`

### R4 — 可选演进（非必须）

| 选项 | 说明 |
|------|------|
| **保持 Vue** | 推荐；已投入学院 UI，迁 Svelte 成本高、收益低 |
| **SvelteKit 子应用** | 仅当团队熟悉 Svelte 时，新建 `apps/web-svelte` 复用同一 nginx API |
| **SSR 首页** | 仅 SEO/首屏；学院内网演示可跳过 |
| **Vitest + MSW** | 契约测试，防止 mapper 漂移 |

---

## 7. 与上一份计划（Osmani）的对照

| 项目 | Osmani 计划 Phase 0–4 | 当前状态 |
|------|----------------------|----------|
| 登录表单 | Phase 0 | ✅ |
| 中台概览/配置 | Phase 0–1 | ✅ |
| 注册/资料/附件/refresh | Phase 2 | ✅ |
| 帖子/板块/统计/邀请/敏感词/角色 | Phase 3–4 | ✅ |
| Lighthouse / OpenAPI | Phase 1 未做 | ⏳ 归入 **R1/R3** |
| 解封 API | Phase 2 | ✅ 已实现 |

---

## 8. 演示剧本（Definition of Done — Harris 版「一条用户旅程算一个 feature」）

| # | 旅程 | 通过标准 |
|---|------|----------|
| 1 | 注册 → 发帖 → 过审 → 列表可见 | 端到端无手动改库 |
| 2 | 学生编辑/删自己的帖 | 作者校验生效 |
| 3 | 管理员审核 + **自定义驳回理由** | R2 完成后 |
| 4 | 平台管理员改板块名 | R2 完成后 |
| 5 | 邀请码生成 → 注册消耗 | 已有 |
| 6 | Token 过期自动 refresh | 已有 |
| 7 | 首屏 3 请求并行 | R1 完成后 |

---

## 9. 不建议做的事

1. **为「Harris 审美」全面重写为 Svelte** — 对接已完成，迁栈只增加风险。  
2. **引入 Redux / TanStack Query** — 当前规模 Pinia + 显式 load 更轻。  
3. **在前端实现敏感词检测** — 应继续走后端 moderation。  
4. **暴露 internal API 给浏览器** — 破坏服务边界。

---

## 10. 相关文档

- [fullstack-integration-audit-plan.md](./fullstack-integration-audit-plan.md) — Osmani 首轮审视与 Phase 0–4 记录  
- [production-launch-plan.md](./production-launch-plan.md) — 生产与 demo-login 策略  

---

*审视人视角：Rich Harris — 强调编译期与端到端数据流；实施仍基于 Vue 3 学院 MVP，吸收思想而非换框架。*
