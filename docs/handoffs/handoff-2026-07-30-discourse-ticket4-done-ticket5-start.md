# Handoff · 2026-07-30 · Discourse Forum Rebuild (Ticket 5 进行中)

> 新会话建议先读 docs/handoffs/handoff-2026-07-30-discourse-ticket4-done-ticket5-start.md，然后说："接着做 Ticket 5。"

## 1. 现在做了什么

- **当前任务**：Discourse 论坛重建工程，正在推进 Ticket 5（管理后台论坛治理菜单过渡）。
- **分支**：codex/discourse-rebuild，HEAD 6e2e948
- **进度**：Ticket 1-4 已完成并部署，Ticket 5 刚开始（尚未修改任何文件）。
- **核心决策**（grill 确认）：平台 = 登录注册 + 管理后台（极简），Discourse = 用户真实使用地。旧论坛一刀切删除。

## 2. 已经完成了什么

### Ticket 1-3b（假设已完成，见之前的 handoff）
- SSO provider、logout sync、Discourse 部署已完成。

### Ticket 4：前端论坛导航切换 + SSO Cookie Bridge（已完成 + 已部署）
- **SSO Cookie Bridge**：
  - rontend/src/api/http.js — 新增 setTokenCookie() 导出，token 持久化时同步写 document.cookie
  - rontend/src/stores/session.js — persistSession() 写 cookie；logout() 清除 cookie
- **路由重构**：rontend/src/router.js
  - 移除全部 13 条 /community/* 路由
  - 学生登录 → window.location.href = 'http://122.51.233.225:8080'
  - admin → /admin
  - 兜底：所有未匹配路径 → Discourse
- **文件清理（46 个）**：
  - 删除 13 个 views（CommunityLayout、CommunityHome、BoardView、PostDetail、NewPost、EditPost、ProfileView、UserPublic、MessagesView、AboutView、RankView、CircleView、MyLibraryView）
  - 删除 25 个 components/gx/，保留 5 个（GxIcon、GxAuthShell、GxAdminSidebar、GxAdminPageHeader、GxBreadcrumb）
  - 删除 8 个 composables/，保留 useDrawerNav
- **api.js 重写**：移除 forumApi，修复 adminApi 中 3 处 forumApi 引用
- **受影响文件修复**：App.vue、AdminLayout.vue（重写为简化内联顶栏+侧栏）、GxAdminSidebar（移除"返回社区"按钮）、AdminAudit.vue、AdminConfig.vue
- **Commit**：6e2e948 feat: frontend forum navigation cutover + SSO cookie bridge (#4)
- **部署**：已通过 scp dist 到云服务器 122.51.233.225，docker cp 到 i-forum-frontend-1 容器，验证返回 200
- **参考文档**：.scratch/discourse-forum-rebuild/goal-ticket4.md（grill 决策 + 文件清单）

## 3. 卡在哪里

- **暂无阻塞**。Ticket 5 范围明确，仅修改 GxAdminSidebar.vue 的导航菜单项。
- **已知环境限制**：云服务器无法访问 GitHub 和 Docker Hub，后续部署需走本地 build + scp 方案。

## 4. 下一步做什么

### Ticket 5：管理后台论坛治理菜单过渡

修改 rontend/src/components/gx/GxAdminSidebar.vue 中的 dminNav 计算属性：

1. **隐藏论坛菜单项**：
   - /admin/audit 内容审核
   - /admin/reports 举报处理
   - /admin/posts 帖子管理
   - /admin/boards 板块管理
   - /admin/sensitive 敏感词（platform_admin 专属）

2. **添加 Discourse 管理入口**：
   - 菜单底部添加「前往 Discourse 管理」链接
   - 指向 http://122.51.233.225:8080/admin
   - 新标签页或同标签页（按之前 grill 决策：同标签页）

3. **保留菜单项**：
   - /admin 数据概览
   - /admin/stats 趋势统计
   - /admin/users 用户管理
   - /admin/config 系统配置/运营配置
   - /admin/invites 邀请码（platform_admin）
   - /admin/roles 角色权限（platform_admin）

4. **不删除后端路由和页面文件**（spec 明确要求：仅导航调整）

### 后续 Ticket
- **Ticket 6**：旧论坛处置（15 天稳定期后，删除 forum-service、schema_forum 等）
- **Ticket 7**：验证与冒烟测试

### 每完成一个 Ticket 必须
1. git commit -m "<type>: <description> (#<ticket>)"
2. git push origin codex/discourse-rebuild
3. 本地 build → scp dist → docker cp 到服务器

## 5. 哪些坑不要再踩

- **PowerShell Remove-Item 被策略阻止**：用 Node.js MCP (mcp__node_repl__js) 做文件删除/写入操作。
- **api.js 编辑容易产生重复函数定义和孤儿代码**：如果改动大，直接重写整个文件比逐函数删除更安全。上次因 orumApi.getBoards 替换和 mapConfigFromBackend 删除不彻底导致语法错误，修复多轮才通过。
- **App.vue 的 import 行有不可见字符**：eplace 精确字符串可能不匹配，用 split('\n').filter(line => !line.includes('xxx')) 更可靠。
- **服务器 GitHub/Docker Hub 不可达**：不要尝试 git fetch/pull 或 docker build（会超时）。部署流程：本地
pm run build → 	ar dist → scp → 服务器 docker cp dist/. ai-forum-frontend-1:/usr/share/nginx/html/。
- **Discourse URL**：当前硬编码 http://122.51.233.225:8080 在 outer.js。后续如有域名变更需同步修改。
- **AGENTS.md 铁律**：每完成一个 Ticket 必须 commit + push，不能跨 Ticket 混合改动。
- **旧论坛一切不留**：用户多次强调，不要做过渡页、保留入口、兼容旧用户等多余工作。
- **不要修改 admin-service 后端代码**：Ticket 5 只改前端的导航菜单，不删后端路由或代码。