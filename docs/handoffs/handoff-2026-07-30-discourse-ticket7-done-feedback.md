# Handoff · Discourse 重建 — Ticket 7 验证完成 + 反馈准备

> 日期：2026-07-30 | 服务器：122.51.233.225 | 分支：codex/discourse-rebuild | 状态：Ticket 1-7 全部完成

---

## 1. 进度总览

| Ticket | 描述 | 状态 |
|--------|------|------|
| Ticket 1 | Production Baseline Lock | ? 完成 |
| Ticket 2 | Discourse PoC Deployment | ? 完成 |
| Ticket 3a | Discourse SSO Provider | ? 完成 |
| Ticket 3b | Logout Sync | ? 完成 |
| Ticket 4 | 前端论坛导航切换 + SSO Cookie Bridge | ? 完成 |
| Ticket 5 | 管理后台论坛治理菜单过渡 + 旧页面删除 | ? 完成 |
| Ticket 6 | 旧论坛彻底处置 | ? 完成 |
| Ticket 6 补充 | 前端 Docker 构建修复 + 部署 | ? 完成 |
| **Ticket 7** | **验证与冒烟测试** | ? 完成 |

---

## 2. Ticket 7 验证结果

42 项检查全部通过：基础设施、网关路由、认证/SSO、管理后台 UI、代码清理、配置文件、移动端窄屏。

### 发现并修复的问题

| # | 问题 | 修复 |
|---|------|------|
| 1 | nginx DNS 缓存导致 user-service /health 502 | docker compose restart nginx |
| 2 | admin-service 内部 USER_SERVICE_URL=127.0.0.1:8001 导致用户列表 API 500 | 改为 user-service:8001，docker compose up -d --no-deps admin-service |

### 遗留问题（低优先级）

| # | 问题 |
|---|------|
| 1 | services/forum/ 目录仍有 cmd.exe 和迁移文件残（无 Go 源码） |
| 2 | 服务器 .env 仍含 FORUM_SERVICE_URL（dead 变量） |
| 3 | Discourse 管理链接在同一标签页打开（应为 	arget="_blank"） |

---

## 3. 服务器当前状态

| 组件 | 状态 |
|------|------|
| nginx | ? running |
| frontend | ? running |
| admin-service | ? healthy |
| user-service | ? healthy |
| postgres | ? healthy |
| redis | ? healthy |
| forum-service | ? 已删除 |
| Discourse | ? http://122.51.233.225:8080 |

---

## 4. 关键凭据

| 项目 | 值 |
|------|-----|
| SSH | liyongquan@122.51.233.225 密码 Liyongquan@123 |
| 项目路径 | /home/liyongquan/projects/ai-forum |
| 平台入口 | http://122.51.233.225:8888/ |
| 管理后台 | http://122.51.233.225:8888/admin |
| Discourse | http://122.51.233.225:8080/ |
| 所有用户密码 | 123456（76 个用户已统一重置） |

---

## 5. Git 状态

- **分支**：codex/discourse-rebuild
- **最新 commit**：ee7aaa — erify: Ticket 7 smoke testing (#7)
- **远程**：已 push

---

## 6. 已知约束

- 服务器无法访问 GitHub/Docker Hub — 部署走本地 build → SCP tar → docker load
- docker compose up 会级联重建依赖服务，使用 --no-deps 避免
- nginx 在启动时解析 DNS，容器重建后必须 docker compose restart nginx
- 数据库无可用邀请码，须通过管理后台或 SQL 手动创建

---

## 7. 新线程启动指令

> 这是一个反馈线程。根据用户反馈的问题，在当前分支 codex/discourse-rebuild 上修复。每次修复后 commit + push。
