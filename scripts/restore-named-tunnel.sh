#!/usr/bin/env bash
# Restore fixed hostname https://api.shgenren.dpdns.org (Named Tunnel).
# Requires CLOUDFLARE_TUNNEL_TOKEN from Zero Trust dashboard.
set -eu
cd /opt/ai-forum

if [ -z "${CLOUDFLARE_TUNNEL_TOKEN:-}" ]; then
  if [ -f /etc/cloudflared/tunnel.token ]; then
    export CLOUDFLARE_TUNNEL_TOKEN="$(cat /etc/cloudflared/tunnel.token)"
  fi
fi

if [ -z "${CLOUDFLARE_TUNNEL_TOKEN:-}" ]; then
  echo "请先在 Cloudflare Zero Trust 创建 Tunnel，复制 Docker 安装命令里的 token，然后："
  echo "  export CLOUDFLARE_TUNNEL_TOKEN='eyJ...'"
  echo "  bash scripts/setup-named-tunnel-token.sh"
  exit 1
fi

bash scripts/setup-named-tunnel-token.sh
echo "Test: curl -fsS https://api.shgenren.dpdns.org/forum-api/api/v1/stats/community"
