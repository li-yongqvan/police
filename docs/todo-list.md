# AI 智联论坛 — TODO List（2026-06-04）

> 基础：已产出测试方案 `test-plan.md`、全功能脚本 `full-experience-test.sh`、安全脚本 `security-smoke.sh`、人工清单 `manual-walkthrough-checklist.md`。

---

## 🔴 P0 — 阻塞上线

- [ ] **BUG-001 修复验证**：在远程服务器 rebuild forum-service → 确认 `/forum-api/api/v1/me/collections` 返回 200 → 前端「我的收藏」正常。详见 [bug-001-collections-404-归因报告.md](./bug-001-collections-404-归因报告.md)

---

## 🟠 P1 — 上线前必须完成

- [ ] **本地 Docker 完整测试**：启动 Docker Desktop → `infra/docker-compose.yml up -d` → `start-dev.ps1` → 运行 `bash scripts/full-experience-test.sh`（含 27 组写入操作），目标 ≥95% PASS
- [ ] **本地安全全量扫描**：运行 `bash scripts/security-smoke.sh`（XSS/SQLi/越权/限流/信息泄露），全部 PASS
- [ ] **人工走查完整执行**：按 `manual-walkthrough-checklist.md` 或 `step-by-step-uat-guide.md` 逐步走完学生→协会管理→中台管理三条路径，填写反馈
- [ ] **手机端三网验收**：移动/联通/电信 × Wi‑Fi/4G 各测一轮（共 6 个组合），覆盖登录→首页→发帖→消息→管理

---

## 🟡 P2 — 质量增强

- [ ] **CI 自动化冒烟**：将 `smoke-test.sh` 接入 GitHub Actions，每次 push 自动执行
- [ ] **版本号暴露**：在 `/health` 或新增 `/version` 端点暴露 `git rev-parse HEAD`，方便判断部署版本
- [ ] **demos 登录加固**：生产环境确认 `demo-login` 白名单生效（从公网访问应 403）
- [ ] **HTTPS 配置**：按 `production-launch-plan.md` 接入 Let's Encrypt
- [ ] **数据库备份验证**：执行一次 `backup-postgres.sh` → 验证备份文件完整 → 尝试还原到测试库
- [ ] **邀请码发放清单**：确认生产环境已作废测试码 `DEMO2026`，批量生成正式邀请码

---

## 🟢 P3 — 体验优化

- [ ] **响应头收敛**：Nginx 当前暴露 `Server: openresty`，建议配置 `server_tokens off`
- [ ] **Harmony 设计自查**：按 `pilot-acceptance-checklist.md` 逐一核对 UI（日期格式、主色一致、无轮播图、管理侧栏）
- [ ] **种子内容检查**：确保每板块 ≥3 条帖子，内容不重复、分类合理
- [ ] **运维文档可查性**：确认值班人能快速找到 `pilot-runbook.md`（502 重启 nginx、回滚命令）

---

## 📋 补丁/脚本待合并（如未合并）

- [ ] `frontend-refresh.tgz` → 前端增量更新是否已推到远程
- [ ] `nav-fix-deploy.tgz` → 导航修复是否已生效
- [ ] `reddit-feed-deploy.tgz` → Reddit 风格改版是否已上线
- [ ] 清理 `scripts/_*.py` 临时部署脚本（完成使命后可归档）

---

## 🔁 灰度节奏参考

| 阶段 | 人数 | 触发条件 |
|------|------|---------|
| 内测 | 5～10（教师+骨干） | 所有 P0/P1 关闭 |
| 小班 | 30～50（单班/区队） | 内测 3 天无新增 P0/P1 |
| 全院 | 按需 | 小班稳定 + 审核人力到位 |

---

**上次更新**：2026-06-04
**关联**：[test-plan.md](./test-plan.md) · [bug-001 报告](./bug-001-collections-404-归因报告.md) · [pilot-runbook.md](./pilot-runbook.md)
