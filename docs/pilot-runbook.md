# AI智联平台 — 校内试运行运维手册

## 正式入口

- 用户访问：**https://api.shgenren.dpdns.org**（Cloudflare Tunnel / 边缘反代），勿向学生分发 `http://IP:8888`
- 运维 SSH：`107.172.138.10`（示例，以实际为准）
- 项目目录：`/opt/ai-forum`

## 日常巡检（每日）

```bash
cd /opt/ai-forum
docker compose -f docker-compose.yml -f docker-compose.server.yml ps
bash scripts/smoke-test.sh   # 需配置 .env.smoke 中 BASE_URL
ls -lt /opt/backups/ai-forum | head
```

## 部署 / 更新

```bash
cd /opt/ai-forum
git pull   # 或 scp 同步
bash scripts/deploy-full-platform.sh
```

**重要**：`user-service` / `forum-service` 容器重建后，脚本会自动 `restart nginx`。若手工 `compose up`，务必执行：

```bash
docker compose -f docker-compose.yml -f docker-compose.server.yml restart nginx
```

## 数据库备份

```bash
sudo bash /opt/ai-forum/scripts/backup-postgres.sh
```

建议 cron（每日 3:00）：

```cron
0 3 * * * root /opt/ai-forum/scripts/backup-postgres.sh >> /var/log/ai-forum-backup.log 2>&1
```

### 恢复演练（简化）

```bash
gunzip -c /opt/backups/ai-forum/ai_forum_YYYYMMDD-HHMMSS.sql.gz | \
  docker compose -f docker-compose.yml -f docker-compose.server.yml exec -T postgres \
  psql -U ai_forum -d ai_forum
```

恢复前请停止写入或维护窗口内进行。

## 冒烟监控 cron

1. 创建 `/opt/ai-forum/.env.smoke`：

```bash
BASE_URL=https://api.shgenren.dpdns.org
SMOKE_USER=demo_student
SMOKE_PASS=demo123456
```

2. 添加 cron（每 15 分钟）：

```cron
*/15 * * * * root /opt/ai-forum/scripts/cron-smoke.sh
```

日志：`/opt/ai-forum/logs/smoke-cron.log`

## 邀请码

- 中台 → **邀请码**：批量生成，一码一人或一班
- 生产环境勿公开 `DEMO2026`；演示账号仅教师内部使用

## 试运行值班

| 项 | SLA |
|----|-----|
| 待审核帖子 | 24h 内处理 |
| 用户举报 | 24h 内处理（中台 → 举报处理） |
| P0（全站不可用） | 立即 |

## 回滚

```bash
cd /opt/ai-forum
docker compose -f docker-compose.yml -f docker-compose.server.yml pull
docker compose -f docker-compose.yml -f docker-compose.server.yml up -d
docker compose -f docker-compose.yml -f docker-compose.server.yml restart nginx
```

## 相关文档

- [DEPLOY.md](../DEPLOY.md)
- [cloudflare-tunnel-setup.md](cloudflare-tunnel-setup.md)
- [pilot-acceptance-checklist.md](pilot-acceptance-checklist.md)
