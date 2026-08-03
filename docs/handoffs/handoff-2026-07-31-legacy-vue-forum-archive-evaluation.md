# Handoff · 2026-07-31 · Vue 自研论坛归档完成，等待方案评估

> 新会话建议先读本文件，然后说：“接着评估我提供的方案。”

## 1. 现在做了什么

- 当前项目是 AI 智联论坛 MVP。
- 论坛本体已经迁移到 Discourse：
  - `8888`：Vue 登录入口和最小管理后台；
  - `8080`：Discourse 论坛，负责浏览、发帖、回帖、搜索、通知和内容治理。
- 本轮目标是隔离已经放弃的 Vue 自研论坛代码，避免后续 Discourse 论坛开发受到旧代码干扰。
- 当前分支：`codex/discourse-rebuild`。
- 下一阶段任务：用户会提供一个方案，需要先进行技术和产品评估，再决定是否执行。

## 2. 已经完成了什么

### 2.1 旧文档归档

旧自研论坛相关文档已归档到：

`docs/archive/2026-07-31/self-built-forum/`

代码归档目录为：

`docs/archive/2026-07-31/legacy-vue-forum/`

代码归档包含：

- 旧 Vue 论坛页面和隐藏页面；
- 旧论坛移动端社区/帖子样式；
- 旧论坛发帖、评论、排序测试；
- 已停用的 `services/forum/` 残留；
- 旧论坛工具组件和日期/帖子数据映射；
- 清理前的 `api.js`、`router.js`、`session.js`、`style.css` 和移动端样式入口快照。

### 2.2 活动 Vue 工程清理

已从活动工程移除或清理：

- 旧帖子、板块、评论、通知 API 映射；
- 旧注册、关注、用户资料接口；
- 旧系统配置、敏感词、角色权限 API；
- 注册、QQ 登录、系统配置、角色权限等已隐藏页面；
- 旧 `/community` 登录后跳转；
- 旧社区移动端样式和测试引用。

当前学生登录后使用：

```text
8888 登录
-> /forum
-> Discourse SSO
-> 8080 论坛
```

当前 `/register` 和 `/oauth/qq` 继续保持隐藏，不重新开放。

### 2.3 有意保留的共享文件

以下文件没有整文件移动，因为当前 8888 登录或管理后台仍然使用：

- `frontend/src/styles/gx-theme.css`
- `frontend/src/api/http.js`
- `frontend/src/composables/useDrawerNav.js`
- `mobile-web/styles/tokens.css`
- `mobile-web/styles/drawer.css`
- `mobile-web/styles/touch-forms.css`
- `mobile-web/styles/pages/login.css`
- `mobile-web/styles/pages/admin.css`

其中 `gx-theme.css` 同时包含当前平台样式和旧社区选择器。旧社区选择器目前没有路由或组件入口，但不应在没有视觉回归验证的情况下整文件拆除。

### 2.4 验证结果

已执行并通过：

```bash
cd frontend
npm run build
npm test -- --run
```

结果：

- Vite production build：通过；
- Vitest：1 个测试文件、1 个测试用例通过；
- 未修改 Discourse 源码和 8080 论坛配置；
- 未进行线上部署。

## 3. 卡在哪里

- 暂无明确阻塞。
- 本次归档和清理目前仍处于工作区未提交状态。
- 需要评估的方案尚未提供，因此暂不继续修改代码。
- `frontend/src/styles/gx-theme.css` 中仍有旧社区选择器，这是为了避免影响当前登录和管理后台而保留的共享文件，不代表旧论坛仍然可用。

当前工作区存在此前遗留的无关变更和未跟踪文件，例如：

- `docs/dev-records/`
- 镜像压缩包；
- 截图；
- 其他历史交接文档。

这些文件不是本轮清理目标，不得擅自删除或回滚。

## 4. 下一步做什么

1. 用户提供待评估的方案。
2. 先判断方案是否与当前 Discourse 架构、MVP 目标和“不能影响现有论坛开发”的要求相容。
3. 输出评估结论：
   - 可行性；
   - 对 `8888`、`8080`、SSO、登录和管理后台的影响；
   - 是否属于大改动；
   - 实施成本；
   - 回滚难度；
   - 潜在风险；
   - 建议保留、修改或否决的部分。
4. 如果方案合理，先生成修改计划，再等待用户确认后执行。
5. 如果方案不合理，明确说明原因，不直接修改代码。
6. 若方案涉及前端视觉改版，优先评估：
   - `8888` 登录页和管理首页是否需要改；
   - `8080` Discourse 主题是否应通过后台主题/CSS实现；
   - 是否误把旧 Vue 自研论坛重新带回当前工程；
   - 是否会影响移动端和登录输入焦点。

## 5. 哪些坑不要再踩

- 不要恢复 `/community` 路由、旧 Vue 论坛页面或旧 `forum-service`。
- 不要把整个 `frontend` 目录移动到归档目录，里面仍有当前登录和管理后台。
- 不要整文件移动 `frontend/src/styles/gx-theme.css`，它仍被当前页面使用。
- 不要修改 `frontend/src/router.js`、`frontend/src/stores/session.js`、`frontend/src/api/http.js` 或 SSO handler，除非方案确实需要且先说明影响。
- 不要重新开放注册、QQ 登录、系统配置、角色权限、统计等已隐藏功能。
- 不要把 `8080` Discourse 的用户界面当成 Vue 页面修改；论坛用户真正看到的主体界面属于 Discourse。
- 不要为了通过测试恢复已归档的旧注册或旧发帖测试；应以当前 MVP 测试范围为准。
- 不要执行批量删除，也不要删除工作区中与本轮无关的用户文件、镜像、截图或开发记录。
- PowerShell 不使用 Bash 风格的 `||` 连接语法。
- 任何方案都必须先验证不会破坏：

```text
8888 登录
-> SSO
-> 8080 Discourse
-> 浏览 / 发帖 / 回帖
```
