# 校内试运行 — 第 3 周放量与值守

## 真机验收矩阵

参考 [mobile-web-adaptation-plan-chris-coyier.md](mobile-web-adaptation-plan-chris-coyier.md) 精简执行：

| 维度 | 要求 |
|------|------|
| 运营商 | 移动 / 联通 / 电信各 1 台 |
| 网络 | Wi‑Fi + 4G 各测一轮 |
| 路径 | 登录 → 首页 → 发帖 → 消息 → 管理端审核 |

记录表见 [pilot-acceptance-checklist.md](pilot-acceptance-checklist.md)。

## 灰度节奏

| 阶段 | 人数 | 动作 |
|------|------|------|
| 内测 | 5～10 | 教师 + 骨干学生，修 P0/P1 |
| 小班 | 30～50 | 单班/单区队，观察审核与举报量 |
| 全院试点 | 按需 | 盯 CPU、磁盘、备份、限流日志 |

## 每日值守清单

1. `docker compose ... ps` — 容器均为 Up
2. `tail /opt/ai-forum/logs/smoke-cron.log` — 最近一条为 PASS
3. `ls -lt /opt/backups/ai-forum | head` — 昨夜备份存在
4. 中台：待审核数、待处理举报数（目标 24h 内清零）

## 邀请码批量发放

1. 使用 `demo_platform_admin` 登录 → **邀请码**
2. 「批量生成」→ 导出列表给辅导员（一码一人）
3. 生产环境作废测试码 `DEMO2026`

## 出问题怎么办

- API 502：`restart nginx`（见 [pilot-runbook.md](pilot-runbook.md)）
- 全站不可用：回滚上一版镜像 tag，恢复备份库（维护窗口）
- 限流误伤：查 Redis key `rl:*`，必要时临时调高阈值后重新部署

## 自动化

```bash
# 生产冒烟（写入 .env.smoke）
BASE_URL=https://api.shgenren.dpdns.org bash scripts/backup-postgres.sh
bash scripts/check-production-secrets.sh /opt/ai-forum/.env
```

可选：将 `scripts/smoke-test.sh` 接入 GitHub Actions，对 staging 域名只读探测。
