# Cloudflare Tunnel 部署指南（0 元跨境）

> 美国源站 Docker 网关 `127.0.0.1:8888` → cloudflared → Cloudflare 边缘 HTTPS → 大陆手机浏览器

## 前置

- 源站论坛已运行：`curl -fsS http://127.0.0.1:8888/` 返回 200
- Cloudflare 免费账号：[dash.cloudflare.com](https://dash.cloudflare.com)
- **无需 Worker**（全站 Tunnel 即可；Worker 易触达免费额度且限制上传）

## 方案 A：Quick Tunnel（无域名，当天测试）

URL 每次重启会变，**不要**长期发给学生。

```bash
cd /opt/ai-forum
chmod +x scripts/run-cloudflare-quick-tunnel.sh
sudo bash scripts/run-cloudflare-quick-tunnel.sh
```

输出示例：`https://xxxx.trycloudflare.com` — 手机打开该链接测试登录。

查看日志：`tail -f /var/log/cloudflared-quick.log`

停止：

```bash
sudo kill "$(cat /var/run/cloudflared-quick.pid)"
```

## 方案 B：Named Tunnel（常驻，推荐）

### B1 无自定义域名（控制台配 Public Hostname）

```bash
cd /opt/ai-forum
chmod +x scripts/setup-cloudflare-tunnel.sh
bash scripts/setup-cloudflare-tunnel.sh
```

1. 浏览器完成 `cloudflared tunnel login`
2. 脚本创建隧道 `ai-forum` 并启动 Docker 容器 `cloudflared`
3. 打开 [Cloudflare Zero Trust](https://one.dash.cloudflare.com/) → **Networks** → **Tunnels** → `ai-forum` → **Public Hostname**
4. 添加：`Subdomain`（或系统分配的 `*.cfargotunnel.com`）→ Service URL `http://127.0.0.1:8888`

### B2 有自己的域名（固定链接）

1. 域名 DNS 托管到 Cloudflare（橙云）
2. 执行：

```bash
export TUNNEL_HOSTNAME=forum.example.com
export TUNNEL_ZONE=example.com
bash scripts/setup-cloudflare-tunnel.sh
```

3. SSL/TLS 模式建议 **Full**（源站仍可用 HTTP 8888）

Compose 启动示例：

```bash
docker compose -f docker-compose.yml -f docker-compose.server.yml \
  -f docker-compose.cloudflare.yml up -d cloudflared
```

## 文件说明

| 路径 | 说明 |
|------|------|
| [infra/cloudflared/config.yml.example](../infra/cloudflared/config.yml.example) | Named Tunnel ingress 模板 |
| [infra/cloudflared/config.yml](../infra/cloudflared/config.yml) | 本地生成，勿提交 Git |
| [infra/cloudflared/credentials.json](../infra/cloudflared/credentials.json) | 隧道密钥，勿提交 Git |
| [docker-compose.cloudflare.yml](../docker-compose.cloudflare.yml) | cloudflared 容器（host 网络） |

## Quick Tunnel 开机自启（无域名）

```bash
sudo bash scripts/install-cloudflared-quick-service.sh
sudo systemctl start cloudflared-quick
grep trycloudflare /var/log/cloudflared-quick.log
```

当前入口也可能写在源站 `/opt/ai-forum/PUBLIC_TUNNEL_URL.txt`。

## 安全（可选）

Tunnel 稳定后关闭公网直连 8888（**先确认手机能打开 Tunnel URL**）：

```bash
sudo bash scripts/cloudflare-harden-firewall.sh
# 非交互: FORCE=1 sudo bash scripts/cloudflare-harden-firewall.sh
```

## 验收

- 手机 4G/Wi‑Fi 打开 Tunnel URL，3s 内首页可交互
- `demo_student` / `demo123456` 登录成功
- 不再向学生分发 `http://107.172.138.10:8888`

## 故障排查

源站一键诊断：

```bash
cd /opt/ai-forum
bash scripts/check-tunnel-health.sh
```

| 现象 | 处理 |
|------|------|
| **522 Connection timed out** | Cloudflare 连不上隧道：源站 `cloudflared` 未运行或隧道未 Healthy。执行 `bash scripts/restore-named-tunnel.sh`（需 `CLOUDFLARE_TUNNEL_TOKEN`），并在 Zero Trust 确认 Public Hostname `api.shgenren.dpdns.org` → `http://127.0.0.1:8888` |
| 502 / error 1033 | 源站未起：`docker compose ... ps`；本机 `curl 127.0.0.1:8888` |
| 登录超时 | 源站已加固重试；检查 cloudflared 日志：`docker logs ai-forum-cloudflared-1` 或 `journalctl -u cloudflared-named -n 50` |
| Quick URL 失效 | 进程退出后需重新运行 quick 脚本 |
| API 502 仅源站 | `docker compose ... restart nginx` |

## 与香港边缘方案

香港 VPS 反代仍可用（[infra/edge-nginx](../infra/edge-nginx/README.md)）。**优先 Cloudflare Tunnel** 可省第二台机器；若 CF 在部分运营商仍慢，再叠加香港入口。
