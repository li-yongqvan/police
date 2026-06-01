# Cloudflare Tunnel

| 模式 | 命令 | 适用 |
|------|------|------|
| Quick（无域名） | `sudo bash scripts/run-cloudflare-quick-tunnel.sh` | 测试、演示 |
| Quick 开机自启 | `sudo bash scripts/install-cloudflared-quick-service.sh` | 无域名但需重启后仍可用 |
| Named（固定） | `bash scripts/setup-cloudflare-tunnel.sh` | 有 CF 账号；有域名时设 `TUNNEL_HOSTNAME` |

Secrets（勿提交 Git）：`config.yml`, `credentials.json` — 从 `config.yml.example` 复制。
