# Ubuntu 云主机部署指南

> 决策基线：单机 Docker Compose · 保留学院 UI · 核心闭环 · `demo-login` 仅内网/白名单

## 1. 服务器要求

| 项目 | 建议 |
|------|------|
| 系统 | Ubuntu 22.04 LTS |
| 配置 | 2 vCPU / 4 GB RAM 起 |
| 软件 | Docker Engine 24+、Docker Compose v2、Git |
| 端口 | 开放 **80**（HTTP）；443 待 Phase 4 配置 HTTPS |

```bash
sudo apt update && sudo apt install -y docker.io docker-compose-v2 git
sudo usermod -aG docker "$USER"
# 重新登录 shell 后生效
```

## 2. 首次部署（推荐脚本）

```bash
git clone <your-repo-url> ai-forum
cd ai-forum

cp .env.production.example .env
# 编辑 .env：JWT_SECRET、POSTGRES_PASSWORD（必改）

chmod +x scripts/*.sh
bash scripts/deploy-ubuntu.sh
```

脚本将依次：启动 Postgres/Redis → 初始化 Schema → `compose` 构建全栈 → 种子数据 → 冒烟测试。

## 3. 手动分步（与脚本等价）

```bash
# 1) 基础设施 + Schema
bash scripts/init-db.sh

# 2) 生产 Compose（仅暴露 Nginx 80）
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

# 3) 等待迁移后写入示例帖
sleep 15
bash scripts/seed-content.sh

# 4) 健康检查（需临时暴露 8001-8003 或在本机网络执行）
bash scripts/smoke-test.sh
```

## 4. 访问

| 入口 | 地址 |
|------|------|
| **学生正式入口（推荐）** | `https://<你的固定域名>/` — 仅向师生分发此地址 |
| 源站运维（可选） | `http://<云主机公网IP>:8888/` — 勿对外推广；Tunnel 正常后可关公网 8888 |
| **大陆手机 Tunnel** | [docs/cloudflare-tunnel-setup.md](docs/cloudflare-tunnel-setup.md) |
| 大陆手机备选 | [docs/cross-border-access.md](docs/cross-border-access.md) |
| 演示登录 | 学号单框登录；**公网默认拒绝** `demo-login`，见下文 |

校内试运行运维详见 [docs/pilot-runbook.md](docs/pilot-runbook.md)。

### 演示账号（密码均为 `demo123456`）

| 用户名 | 角色 |
|--------|------|
| demo_student | 学生 |
| demo_admin | 协会管理员 |
| demo_platform_admin | 中台管理员 |

生产环境请使用 **注册 + 邀请码**（种子邀请码 `DEMO2026`）。

## 5. demo-login 白名单（选项 B）

两层限制：

1. **Nginx**（`nginx/includes/demo-login-allow.conf`）  
   默认仅 RFC1918 私网 + 127.0.0.1。公网访问返回 403。

2. **user-service**（`DemoLoginGuard`）  
   - `APP_ENV=development`：不限制（本地开发）  
   - `APP_ENV=production`：私网 IP 或通过 `DEMO_LOGIN_ALLOWLIST` 配置的 IP/CIDR

在云主机 `.env` 中增加办公室出口 IP 示例：

```env
DEMO_LOGIN_ALLOWLIST=203.0.113.10
```

修改后：

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build user-service nginx
```

## 6. 日常运维

```bash
# 查看状态
docker compose -f docker-compose.yml -f docker-compose.server.yml ps

# 生产冒烟（在 .env.smoke 中配置 BASE_URL=https://你的域名）
bash scripts/smoke-test.sh

# 每日备份
sudo bash scripts/backup-postgres.sh

# 强密钥检查（部署前）
bash scripts/check-production-secrets.sh .env

# 更新版本
git pull
bash scripts/deploy-full-platform.sh
```

备份与 cron 示例见 [docs/pilot-runbook.md](docs/pilot-runbook.md)。放量与验收见 [docs/pilot-launch-guide.md](docs/pilot-launch-guide.md)。

## 7. 本地开发（Windows / macOS）

```powershell
docker compose -f infra/docker-compose.yml up -d
.\scripts\init-db.ps1
$env:APP_ENV="development"
.\start-dev.ps1
.\scripts\seed-content.ps1
```

## 8. 故障排查

| 现象 | 处理 |
|------|------|
| 迁移失败 | 确认已执行 `scripts/init-db.sh` |
| 论坛为空 | 执行 `bash scripts/seed-content.sh` |
| demo-login 403 | 从校园网/VPN 访问，或配置 `DEMO_LOGIN_ALLOWLIST` |
| 中台 403 | 使用 demo_admin / demo_platform_admin 登录 |
| 大陆手机打不开 / 仅首页能开 | [Cloudflare Tunnel](docs/cloudflare-tunnel-setup.md)；备选 [跨境指南](docs/cross-border-access.md) |
| API 502 但容器 healthy | `docker compose ... restart nginx`（上游 IP 变更后） |

## 9. 后续（Phase 2+）

- HTTPS（Let's Encrypt）
- CI/CD 自动部署
- 关闭公网数据库端口（`docker-compose.prod.yml` 已默认不暴露）

详见 [docs/production-launch-plan.md](docs/production-launch-plan.md)。
