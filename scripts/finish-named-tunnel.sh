#!/usr/bin/env bash
# Run AFTER: cloudflared tunnel login (cert.pem exists)
set -eu
cd "$(dirname "$0")/.."

export TUNNEL_NAME="${TUNNEL_NAME:-ai-forum}"
export TUNNEL_HOSTNAME="${TUNNEL_HOSTNAME:-api.shgenren.dpdns.org}"
export TUNNEL_ZONE="${TUNNEL_ZONE:-shgenren.dpdns.org}"
export ORIGIN_URL="${ORIGIN_URL:-http://127.0.0.1:8888}"

# Stop quick tunnel
kill "$(cat /var/run/cloudflared-quick.pid 2>/dev/null)" 2>/dev/null || true
systemctl stop cloudflared-quick 2>/dev/null || true
pkill -f 'cloudflared tunnel --url' 2>/dev/null || true

if [ ! -f "$HOME/.cloudflared/cert.pem" ]; then
  echo "Missing $HOME/.cloudflared/cert.pem — run: cloudflared tunnel login" >&2
  exit 1
fi

bash scripts/setup-cloudflare-tunnel.sh

# Also run named tunnel via systemd (host network) if docker cloudflared fails
if ! docker ps --format '{{.Names}}' | grep -q cloudflared; then
  echo "Starting cloudflared via host..."
  nohup cloudflared tunnel --config /opt/ai-forum/infra/cloudflared/config.yml run >>/var/log/cloudflared-named.log 2>&1 &
  echo $! >/var/run/cloudflared-named.pid
fi

echo "Test: curl -sS https://${TUNNEL_HOSTNAME}/forum-api/api/v1/stats/community"
