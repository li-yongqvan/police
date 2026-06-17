# 大陆手机访问海外源站指南

> 适用：源站在美国（如 `107.172.138.10`），用户在境内用手机浏览器访问 AI智联平台。

## 0. 推荐入口（0 元）

**首选：Cloudflare Tunnel** — 无需香港 VPS，手机访问 `https://*.trycloudflare.com`（测试）或你的 CF 域名（正式）。

完整步骤：[cloudflare-tunnel-setup.md](cloudflare-tunnel-setup.md)

```bash
# 源站上 Quick Tunnel（无域名）
sudo bash /opt/ai-forum/scripts/run-cloudflare-quick-tunnel.sh
```

## 1. 现象与层次定位

| 手机上的表现 | 优先怀疑层次 |
|--------------|--------------|
| 一直加载后失败 | 跨境链路 RTT/丢包；浏览器无超时 |
| 部分手机正常、部分不行 | 运营商国际出口路径不同 |
| 首页能开、登录失败 | POST `/user-api/api/v1/login` 更敏感；或 nginx 上游 502 |

**不要**只改 CSS；先区分「页面」与「API」是否都成功。

## 2. 十分钟诊断清单

### 2.1 客户端

1. 确认访问的是 **Tunnel HTTPS 链接**，而非过期 Quick URL。
2. 开发者工具 / 抓包：`GET /assets/*` 与 `POST /user-api/api/v1/login` 状态码。
3. 切换 Wi‑Fi 与 4G 各试一次。

### 2.2 源站（SSH）

```bash
cd /opt/ai-forum
curl -sS http://127.0.0.1:8888/forum-api/api/v1/stats/community
docker compose -f docker-compose.yml -f docker-compose.server.yml restart nginx
docker logs ai-forum-cloudflared-1 --tail 30 2>/dev/null || true
```

### 2.3 判定表

| 结果 | 结论 |
|------|------|
| Tunnel URL 正常、直连 IP 失败 | 预期；推广 Tunnel，可关闭公网 8888 |
| 静态 200、login 502 | restart nginx；见 [deploy-full-platform.sh](../scripts/deploy-full-platform.sh) |
| Tunnel 也超时 | 换运营商测试；或备选香港边缘 |

## 3. 备选：香港边缘反代

```
手机 --HTTPS:443--> 香港 VPS --HTTP:8888--> 美国源站
```

- [infra/edge-nginx/nginx.conf](../infra/edge-nginx/nginx.conf)
- [scripts/deploy-edge-hk.sh](../scripts/deploy-edge-hk.sh)

## 4. 应用层加固（源站）

- [frontend/src/api/http.js](../frontend/src/api/http.js)：超时与登录重试
- [nginx/nginx.conf](../nginx/nginx.conf)：proxy 超时
- [scripts/deploy-full-platform.sh](../scripts/deploy-full-platform.sh)：部署后 restart nginx

## 5. 验收

- 移动 / 联通 / 电信手机：Tunnel URL 3s 内可交互
- `demo_student` / `demo123456` 连续登录 10 次 ≥ 9 次成功
- 不再推广 `http://107.172.138.10:8888` 直连
