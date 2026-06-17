# 国内单机迁移指南（无 Cloudflare）

适用于将论坛从海外 VPS 迁到国内云主机，直接 **HTTP 80** 对外。

## 迁移范围（必要文件）

| 目录/文件 | 说明 |
|-----------|------|
| `docker-compose.yml` + `docker-compose.prod.yml` | 全栈 Compose，Nginx 暴露 80 |
| `services/` | user / forum / admin 三个 Go 服务 |
| `frontend/` | Vue 前端 |
| `nginx/` | 网关 |
| `scripts/` | 部署、种子、冒烟 |
| `migrations/` | 数据库初始化 |
| `.env.production.example` | 生产环境变量模板 |

**不需要**：`docker-compose.cloudflare.yml`、`infra/cloudflared/`、Cloudflare 相关脚本。

## 一键迁移（从开发机）

```powershell
cd "项目根目录"
powershell -File scripts/pack-pilot-deploy.ps1
python scripts/remote-deploy-domestic.py
```

或手动：

```bash
scp pilot-deploy.tgz liyongquan@122.51.233.225:/tmp/
ssh liyongquan@122.51.233.225
sudo mkdir -p /opt/ai-forum && sudo tar -xzf /tmp/pilot-deploy.tgz -C /opt/ai-forum
sudo chown -R liyongquan:liyongquan /opt/ai-forum
cd /opt/ai-forum && bash scripts/deploy-domestic.sh
```

## 从海外迁数据（可选）

若需保留海外 **用户/帖子** 而非仅种子数据：

```bash
# 在海外源站
bash /opt/ai-forum/scripts/backup-postgres.sh
# 将 /opt/backups/ai-forum/*.sql.gz 拷到国内机后：
gunzip -c ai_forum_*.sql.gz | docker compose -f docker-compose.yml -f docker-compose.prod.yml exec -T postgres psql -U ai_forum -d ai_forum
```

## 安全组

云厂商控制台放行：**80/TCP**（HTTP）。暂不用 8888。

## 演示账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| demo_platform_admin | demo123456 | 中台 |
| demo_admin | demo123456 | 协会管理员 |
| demo_student | demo123456 | 学生 |

生产请改密码并走邀请码注册。
