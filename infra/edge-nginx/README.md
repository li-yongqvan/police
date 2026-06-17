# 香港 / 国内边缘反向代理

大陆手机访问海外源站时，建议**不要**直连美国 IP，而是在边缘机部署本配置。

## 拓扑

```
手机(大陆) --HTTPS:443--> 边缘 VPS(香港) --HTTP:8888--> 美国源站 /opt/ai-forum
```

## 快速部署（Ubuntu 22.04 边缘机）

```bash
sudo apt update && sudo apt install -y nginx certbot python3-certbot-nginx
sudo mkdir -p /var/cache/nginx/ai-forum

# 编辑 nginx.conf：替换 ORIGIN_HOST、YOUR_DOMAIN
sudo cp nginx.conf /etc/nginx/conf.d/ai-forum-edge.conf
sudo sed -i 's/ORIGIN_HOST/107.172.138.10/g' /etc/nginx/conf.d/ai-forum-edge.conf
sudo sed -i 's/YOUR_DOMAIN/forum.example.com/g' /etc/nginx/conf.d/ai-forum-edge.conf

sudo nginx -t && sudo systemctl reload nginx
sudo certbot --nginx -d forum.example.com
```

或使用仓库脚本（在边缘机上 clone 后）：

```bash
export ORIGIN_HOST=107.172.138.10
export EDGE_DOMAIN=forum.example.com
bash scripts/deploy-edge-hk.sh
```

## 源站防火墙（推荐）

在美国源站仅允许边缘机 IP 访问 8888：

```bash
ufw allow from <EDGE_IP> to any port 8888
ufw deny 8888
```

大陆用户只推广：`https://forum.example.com/`
