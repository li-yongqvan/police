# Handoff · 暖色流动 UI 改版部署 — 规划审查完成

> 日期：2026-07-27 | 分支：`experiment/bg-animation` | HEAD：`7fce935` | 状态：规划阶段完成，等待用户确认后进入执行

---

## 1. 功能概述

基于 `Kimi_Agent_暖色流动UI设计/` 原型，将 AI 智联论坛 Vue 3 SPA 全面改版为暖色流动风格，部署到云服务器 `http://122.51.233.225/`。

核心改造：
- **Canvas 2D 双层背景引擎**：光晕粒子 + 织线纹理，替代现有 CSS blur orb
- **双主题 CSS Token**：warm（学生端暖色）/ cool（管理端冷色），通过 `body[data-theme]` 切换
- **6 个 Ticket**：Token 重构 → 基础组件 → 登录/布局壳 → 学生端页面 → 管理后台 → Docker 部署

---

## 2. 规划产物（只读输入，执行线程不读聊天记录）

| 文件 | 大小 | 用途 |
|------|------|------|
| `.scratch/warm-ui-redesign/spec.md` | 10KB | 子问题、假设、约束、验收标准、6 个 Ticket |
| `.scratch/warm-ui-redesign/goal.md` | 6KB | Ticket 1 执行契约（CSS Token + Canvas 2D 引擎） |

---

## 3. 本轮审查修复清单

在本线程中，对 spec.md 和 goal.md 进行了**两轮审查**，共修复/确认以下问题：

### 第一轮修复（6 个矛盾 + 5 个技术补充）

| # | 问题 | 修复 |
|---|---|---|
| 1 | Canvas 引擎行号矛盾（goal 引 877-900 vs 实际 3246-3600） | 统一为 3246-3600 |
| 2 | `gx-circle-theme.css` 三步说法打架 | 删除 Step 1 中「不改」条款 |
| 3 | `data-theme` 挂载位置不一致 | 明确 `document.body.dataset.theme` |
| 4 | 「不改 `<script setup>`」vs Step 5 冲突 | Constraints 增加例外条款 |
| 5 | 「不 push」vs Ticket 6 部署 | Ticket 6 标注需用户单独授权 push |
| 6 | mode 推断规则不完整 | 补路由→mode 映射表 + 兜底 warm |
| 7 | 硬编码色值预检 | Step 0 新增 grep 命令 |
| 8 | Canvas composable 生命周期 | 补 RAF 清理、非响应式粒子数组、离屏 canvas、prefers-reduced-motion |
| 9 | 验收标准过严 | 「无错误」→「无新增错误」 |
| 10 | Docker 构建 OOM | 优先本地 build → push registry |
| 11 | FOUC 白屏闪烁 | Step 2 补 `index.html` `<body data-theme="warm">` |

### 第二轮修复（本轮，5 个残留问题）

| # | 问题 | 结果 |
|---|---|---|
| 1 | 登录页路由 `/` vs `/login` 打架 | ✅ 确认已修复（goal.md L82 为 `/login`） |
| 2 | spec 部署路径二选一未决 | ✅ 确认已修复（统一为 registry 路径） |
| 3 | FOUC 白屏 | ✅ 确认已修复 |
| 4 | `rg --type vue` 跑不通 | ✅ 确认已修复 |
| 5a | Step 5 残留「角色」 | ✅ 确认已修复 |
| **5b** | **映射表缺 markdown 表头** | ✅ **本轮修复**（goal.md L63-64 插入表头） |

---

## 4. 如何启动执行线程

用户在新 Codex 对话中，第一句话：

> 读取 `.scratch/warm-ui-redesign/goal.md`，按 Execution order 逐步执行。不要重新问需求，不要修改 Completion criteria 和 Constraints 里没提到的东西。每完成一个 Step 报告并等待确认。

执行顺序：Step 0（预检）→ Step 1（gx-theme.css cool Token）→ Step 2（body[data-theme]）→ Step 3（gx-circle-theme.css）→ Step 4（useCanvasBackground.js）→ Step 5（App.vue 挂载）→ Step 6（GxAnimatedBg.vue 保留不用）→ Step 7（验证）

---

## 5. 已知注意事项

- **index.html 的 `<body>`**：Step 2 需在 `<body>` 上预写 `data-theme="warm"`，防止 FOUC
- **粒子数组**：Step 4 中 `motes[]` 必须用模块级普通变量，禁止 `ref/reactive`
- **离屏 canvas 缓存**：织线纹理层只在 resize/主题切换时重绘，禁止每帧全屏画线
- **prefers-reduced-motion**：匹配时粒子静止或数量减半
- **rg 命令**：Step 0 的命令已修正为 `rg -n "#[0-9a-fA-F]{3,6}" -g '*.vue' -g '*.css'`
- **编码**：spec.md 和 goal.md 均为 UTF-8，PowerShell cat 会显示乱码，用 Node REPL 读
- **原型参考行号**：Canvas 引擎逻辑在原型的第 3246-3600 行（已用 grep 验证）

---

## 6. 未修改范围

- 后端 Go 服务（user/forum/admin）
- 数据库 schema、API 路由、API 字段契约
- `router.js`、`stores/`、`api.js`、`api/http.js`
- `composables/` 中已有文件
- 不引入新 npm 依赖（Canvas 2D 用浏览器原生 API）

---

## 7. 部署关键决策

- **构建方式**：本地 `docker build` → push registry → 服务器 `docker compose pull`
- **版本标签**：`ai-forum-frontend:v1.0`（旧版）、`ai-forum-frontend:v2.0`（新版）
- **回滚**：`scripts/rollback-frontend.sh` 切回 v1.0 标签重启
- **docker-compose.yml**：前端服务从 `build:` 改为 `image: registry/ai-forum-frontend:v2.0`
- **push 授权**：Ticket 6 执行 push 和部署前需用户单独口头授权
