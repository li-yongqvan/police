# 校内小规模试运行 — 启动清单

**学生入口（只发这个）：** https://api.shgenren.dpdns.org

## 第一天（管理员）

1. 使用 **中台账号**（`demo_platform_admin` / 密码见内部文档）登录 → **邀请码** → 批量生成 50～100 个，导出给辅导员。
2. 使用 **协会管理员**（`demo_admin`）登录 → **内容审核**、**举报处理**。
3. 打开 **关于本站**（社区内 `/community/about`）确认注册说明正确。

## 学生侧流程

邀请码注册 → 登录 → 完善院系/区队 → 浏览四板块 → 发帖/评论 → 消息中心看通知。

## 运维（服务器）

| 项 | 命令/位置 |
|----|-----------|
| 冒烟日志 | `/opt/ai-forum/logs/smoke-cron.log` |
| 备份目录 | `/opt/backups/ai-forum/` |
| 本机冒烟 | `bash /opt/ai-forum/scripts/smoke-vps-remote.sh` |
| 手册 | [pilot-runbook.md](pilot-runbook.md) |

**注意：** VPS 上 cron 使用 `BASE_URL=http://127.0.0.1:8888`，不要用公网域名（会超时）。

**若学生打开域名出现 Cloudflare 522：** 在 VPS 上执行 `bash scripts/check-tunnel-health.sh`；隧道挂了则 `bash scripts/restore-named-tunnel.sh`（见 [cloudflare-tunnel-setup.md](cloudflare-tunnel-setup.md)）。

## 灰度建议

| 阶段 | 人数 |
|------|------|
| 内测 | 5～10 人（教师+骨干） |
| 小班 | 30～50 人（一个班） |
| 扩大 | 按监控与审核能力再加 |

## 验收

打印或共享 [pilot-acceptance-checklist.md](pilot-acceptance-checklist.md)，内测结束后勾选。
